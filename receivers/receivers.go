package receivers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/iostrovok/go-mockdata/imports/inter"
	"github.com/iostrovok/go-mockdata/pkparser"
)

type Parser struct {
	Files    map[string]*OneFile
	imp      inter.IImp
	pkParser *pkparser.PkParser
	srcDir   string
}

type OneFile struct {
	imp     map[string]string
	aliases map[string]string
	parsed  *ast.File
}

func New(srcDir string, imp inter.IImp) *Parser {

	srcDir, err := filepath.Abs(filepath.Dir(srcDir))
	if err != nil {
		panic(fmt.Sprintf("receivers has got wrong dir: %v", err))
	}

	return &Parser{
		imp:      imp,
		srcDir:   srcDir,
		Files:    map[string]*OneFile{},
		pkParser: pkparser.New(srcDir, imp),
	}
}

type OneFunctionRes struct {
	Err                           error
	FilePath, Pkg, Recv, FuncName string
	Params, Results               []string
}

func (pr *Parser) Run(i interface{}) OneFunctionRes {
	ofr := pr.OneFunction(i)
	if ofr.Err == nil {
		ofr.Params, ofr.Results, ofr.Err = pr.InOut(ofr)
	}

	return ofr
}

func (pr *Parser) convertImports(fileName string, f *ast.File) *OneFile {

	oneFile, find := pr.Files[fileName]
	if !find {
		oneFile = &OneFile{
			imp:     map[string]string{},
			aliases: map[string]string{},
			parsed:  f,
		}
	}

	// Store the imports from the file's AST to tmp maps
	for _, s := range f.Imports {
		value := strings.Trim(s.Path.Value, `"`)
		oneFile.imp[value] = value
		if s.Name != nil {
			oneFile.aliases[strings.Trim(s.Name.String(), `"`)] = value
		} else {
			_, file := filepath.Split(value)
			oneFile.imp[file] = value
		}
	}

	pr.Files[fileName] = oneFile
	return oneFile
}

func (pr *Parser) InOut(ofr OneFunctionRes) ([]string, []string, error) {
	f, err := pr.pkParser.ParsePackageByFile(ofr.FilePath)
	if err != nil {
		fmt.Printf("InOut ERROR: %v\n", err)
		return nil, nil, err
	}

	params, results := pr.Get(f, ofr.FilePath, ofr.Recv, ofr.FuncName)
	return params, results, nil
}

func packageNamByFile(fileName string) (string, error) {
	src, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}
	return src.Name.Name, nil
}

//
//func (pr *Parser) parsePackage(fileName string) (*OneFile, error) {
//
//	pkgName, err := packageNamByFile(fileName)
//	if err != nil {
//		return nil, err
//	}
//
//	if f, find := pr.Files[pkgName]; find {
//		return f, nil
//	}
//
//	dir, _ := filepath.Split(fileName)
//
//	fset := token.NewFileSet() // positions are relative to fset
//
//	// Parse src but stop after processing the imports.
//	parsed, err := parser.ParseDir(fset, dir, nil, parser.AllErrors)
//	if err != nil {
//		return nil, err
//	}
//
//	if _, find := parsed[pkgName]; !find {
//		return nil, fmt.Errorf("not found pacakge %s", pkgName)
//	}
//
//	if _, find := parsed[pkgName].Files[fileName]; !find {
//		return nil, fmt.Errorf("not found file %s for pacakge %s", fileName, pkgName)
//	}
//
//	for fileName, parsed := range parsed[pkgName].Files {
//		pr.convertImports(fileName, parsed)
//	}
//
//	return pr.Files[fileName], nil
//}

//func (pr *Parser) FindInFile(fileName, name string) (bool, error) {
//	f, err := pr.parsePackage(fileName)
//	if err != nil {
//		return false, err
//	}
//
//	_, find := Walker(f.parsed, name)
//	return find, nil
//}

func extractFuncType(iFn interface{}) *ast.FuncType {
	switch iFn.(type) {
	case *ast.Field:
		return extractFuncType(iFn.(*ast.Field).Type)
	case *ast.FuncDecl:
		return extractFuncType(iFn.(*ast.FuncDecl).Type)
	case *ast.FuncType:
		return iFn.(*ast.FuncType)
	}

	return nil
}

func (pr *Parser) Get(f *ast.Package, fileName, recv, funcName string) (params, results []string) {
	iFn, find := Walker(f, recv, funcName)

	fmt.Printf("\n\nGet-Walker result: %t, %T => %+v\n", find, iFn, iFn)

	if !find {
		return
	}

	fn := extractFuncType(iFn)
	if fn == nil {
		return
	}

	fmt.Printf("-1. fn: %+v, recv: %s, funcName: %s\n", fn, recv, funcName)
	packages, params, results := checkOne(fn.Params, fn.Results)

	fmt.Printf("-1. packages....: %+v\n", packages)
	fmt.Printf("-1. params....: %+v\n", params)
	fmt.Printf("-1. results....: %+v\n", results)

	for _, p := range packages {
		fmt.Printf("packages....: %+v\n", p)
		pr.imp.UseFile(fileName, p)
	}

	return params, results
}

func checkOne(paramsList, resultsList *ast.FieldList) ([]string, []string, []string) {

	fmt.Printf("checkOne. -1. %+v => %+v\n", paramsList, resultsList)

	packages := make([]string, 0)
	results := make([]string, 0)
	params := make([]string, 0)

	if paramsList != nil && paramsList.List != nil {
		for _, f := range paramsList.List {
			fmt.Printf("checkOne.-2.  %+v\n", f)

			pkg, val := readInOutType(f.Type)

			fmt.Printf("checkOne.-2. pkg: %+v, val: %s\n", pkg, val)
			fmt.Printf("checkOne.-2. f.Names: %+v\n", f.Names)
			packages = append(packages, pkg...)

			for i := 0; i < len(f.Names) || i < 1; i++ {
				params = append(params, val)
			}
		}
	}

	if resultsList != nil && resultsList.List != nil {

		for _, f := range resultsList.List {
			pkg, val := readInOutType(f.Type)
			packages = append(packages, pkg...)
			results = append(results, val)
		}
	}

	return packages, params, results
}

func readInOutType(f interface{}) ([]string, string) {

	defPkg := make([]string, 0)

	fmt.Printf("readInOutType. -1. %T=> %+v\n", f, f)

	switch f.(type) {
	case *ast.Ident:
		return defPkg, f.(*ast.Ident).String()
	case *ast.ArrayType:
		k := f.(*ast.ArrayType)
		_, pkgLoc := readInOutType(k.Elt)
		res := ""
		if k.Len != nil {
			res = fmt.Sprintf("[%s]%s", k.Len, pkgLoc)
		} else {
			res = fmt.Sprintf("[]%s", pkgLoc)
		}
		return []string{pkgLoc}, res
	case *ast.SelectorExpr:
		k := f.(*ast.SelectorExpr)
		_, pkgLoc := readInOutType(k.X)
		return []string{pkgLoc}, pkgLoc + "." + k.Sel.String()
	case *ast.StarExpr:
		pkg, val := readInOutType(f.(*ast.StarExpr).X)
		return pkg, "*" + val
	case *ast.MapType:
		pkg1, key := readInOutType(f.(*ast.MapType).Key)
		pkg2, value := readInOutType(f.(*ast.MapType).Value)
		return append(pkg1, pkg2...), fmt.Sprintf("map[%s]%s", key, value)

	}

	panic(fmt.Sprintf("\n ERRRRRRRRRRRRRRRRR : %T\n\n", f))

	return nil, ""
}

func (pr *Parser) OneFunction(i interface{}) OneFunctionRes {

	out := OneFunctionRes{}

	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Func {
		out.Err = fmt.Errorf("should be function")
		return out
	}

	// github.com/iostrovok/go-mockdata/imports.(*Imports).Add-fm
	fullName := runtime.FuncForPC(value.Pointer()).Name()
	dir, file := filepath.Split(fullName)

	//
	parts := strings.SplitN(file, ".", 3)
	if len(parts) < 3 {
		out.Err = fmt.Errorf("not object function: %s", file)
		return out
	}

	// github.com/iostrovok/go-mockdata/ + imports
	out.Pkg = filepath.Join(dir, parts[0])

	// (*Imports) => *Imports
	out.Recv = strings.TrimRight(strings.TrimLeft(parts[1], "("), ")")

	// full path to file
	out.FilePath, _ = runtime.FuncForPC(value.Pointer()).FileLine(value.Pointer())

	// return "-fm" from function name
	out.FuncName = strings.TrimRight(parts[len(parts)-1], "-fm")

	return out
}
