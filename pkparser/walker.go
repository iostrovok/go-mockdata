package pkparser

/*
	Simple walk through parsed package tree.
*/

import (
	"fmt"
	"go/ast"
	"go/build"
	"strings"
)

type NodeType int

const (
	UnknownType NodeType = iota + 1
	FieldType
	FieldListType
	FileType
	FuncDeclType
	FuncTypeType
	GenDeclType
	IdentType
	ImportSpecType
	InterfaceType
	PackageType
	StarExprType
	TypeSpecType
)

type Walker struct {
	srcDir      string
	parsedFiles map[string]*ParsedFile
	viewDebug   bool
}

type NodeResult struct {
	Path   string
	GoRoot bool
	Type   NodeType
	Node   interface{}
	Error  error
}

type WalkerResult struct {
	ParsedFile *ParsedFile
	Node       *NodeResult
	Find       bool
}

func NewWalker(srcDir string) *Walker {
	return &Walker{
		srcDir:      srcDir,
		parsedFiles: map[string]*ParsedFile{},
	}
}

func (w *Walker) AddFiles(new map[string]*ParsedFile) *Walker {
	for name, file := range new {
		w.parsedFiles[name] = file
	}
	return w
}

func (w *Walker) SetDebug(viewDebug bool) *Walker {
	w.viewDebug = viewDebug
	return w
}

func (w *Walker) debug(format string, d ...interface{}) {
	if w.viewDebug {
		fmt.Printf(format, d...)
	}
}

// FindPath finds  in parsed file or package by path and name of node
func (w *Walker) FindPath(path, name string) *WalkerResult {
	w.debug("FindPath path: %s, name: %s\n", path, name)
	return w.find(path, "", name)
}

// FindFile finds  in parsed file or package by srcFile and name of node
func (w *Walker) FindFile(srcFile, name string) *WalkerResult {
	w.debug("FindFile srcFile: %s, name: %s\n", srcFile, name)
	return w.find("", srcFile, name)
}

// hasPath finds  in parsed file or package by srcFile and name of node
func (w *Walker) hasPath(path string) bool {

	for _, p := range w.parsedFiles {

		w.debug("w.hasPath: path: %s, p.path: %s\n", path, p.path)
		if p.path == path {
			return true
		}
	}

	return false
}

// FindFile finds  in parsed file or package by srcFile and name of node
func (w *Walker) find(path, srcFile, name string) *WalkerResult {

	result := &WalkerResult{}

	w.debug("1. find path: %s, srcFile: %s, name: %s\n", path, srcFile, name)

	if path == "" && srcFile == "" {
		return result
	}

	if path != "" && srcFile != "" {
		return result
	}

	w.debug("2. find path: %s, srcFile: %s, name: %s\n", path, srcFile, name)
	w.debug("3. find path: w.parsedFiles: %+v\n", w.parsedFiles)

	for _, result.ParsedFile = range w.parsedFiles {

		w.debug("w.parsedFiles:  file.path: %s, file.srcFile: %s\n", result.ParsedFile.path, result.ParsedFile.srcFile)

		if srcFile != "" && result.ParsedFile.srcFile != srcFile {
			continue
		}

		if path != "" && result.ParsedFile.path != path {
			continue
		}

		fmt.Printf("file: %+v\n", result.ParsedFile)
		fmt.Printf("name: %s, pkgName: %s, path: %s\n", name, result.ParsedFile.pkgName, result.ParsedFile.path)

		if result.Node, result.Find = w.deepWalk(result.ParsedFile.file, name); result.Find {
			return result
		}
	}

	return result
}

func (walk *Walker) deepWalk(in interface{}, name string) (result *NodeResult, findResult bool) {

	walk.debug("\ndeepWalk. name: %s, in: %T => %+v\n", name, in, in)

	result = &NodeResult{
		Type:  UnknownType,
		Node:  in,
		Error: nil,
	}

	switch in.(type) {
	case *ast.Package:
		f := in.(*ast.Package)
		walk.debug("\ndeepWalk. *ast.Package\n")
		for _, d := range f.Files {
			if out, find := walk.deepWalk(d, name); find {
				return out, find
			}
		}

	case *ast.File:
		f := in.(*ast.File)
		walk.debug("\ndeepWalk. *ast.File\n")
		if f.Name.String() == name {
			result.Type = FileType
			findResult = true
			return
		}
		for _, d := range f.Decls {
			if out, find := walk.deepWalk(d, name); find {
				return out, find
			}
		}

	case *ast.GenDecl:
		f := in.(*ast.GenDecl)
		walk.debug("\ndeepWalk. *ast.GenDecl\n")
		if f.Specs != nil && len(f.Specs) > 0 {
			for _, d := range f.Specs {
				walk.debug("\n-----------deepWalk. f.Specs. d: %T, d: %+v\n", d, d)
				if out, find := walk.deepWalk(d, name); find {
					walk.debug("\n+++++++++++++deepWalk. *ast.GenDecl. out: %+v, find: %t\n", out, find)
					return out, find
				}
			}
		}

	case *ast.TypeSpec:
		f := in.(*ast.TypeSpec)
		walk.debug("\ndeepWalk. *ast.TypeSpec name: %s\n", f.Name.String())
		if f.Name.String() == name {
			result.Type = TypeSpecType
			findResult = true
			return
		}
		if out, find := walk.deepWalk(f.Type, name); find {
			return out, find
		}

	case *ast.InterfaceType:
		f := in.(*ast.InterfaceType)
		walk.debug("\ndeepWalk. *ast.InterfaceType\n")
		if f.Methods.List != nil {
			for _, method := range f.Methods.List {
				if out, find := walk.deepWalk(method, name); find {
					return out, find
				}
			}
		}

	case *ast.FuncDecl:
		f := in.(*ast.FuncDecl)
		walk.debug("\ndeepWalk. *ast.FuncDecl\n")
		if f.Name.String() == name {
			result.Type = FuncDeclType
			findResult = true
			return
		}
		if out, find := walk.deepWalk(f.Type, name); find {
			return out, find
		}

		// case for "func (s *obj) Name()"
		if f.Recv != nil {
			if _, find := walk.deepWalk(f.Recv, name); find {
				result.Type = FuncDeclType
				findResult = true
				return
			}
		}

	case *ast.FieldList:
		f := in.(*ast.FieldList)
		walk.debug("\ndeepWalk. *ast.FieldList\n")
		if f.List != nil {
			for _, method := range f.List {
				if out, find := walk.deepWalk(method, name); find {
					return out, find
				}
			}
		}

	case *ast.ImportSpec:
		f := in.(*ast.ImportSpec)
		path := strings.Trim(f.Path.Value, `"`)

		walk.debug("\n\ndeepWalk. *ast.ImportSpec\n")
		walk.debug("deepWalk. f.Path: %T, f.Name.String(): %s, %s\n", f.Path, f.Name.String(), path)

		if f.Name.String() == name || path == name {
			result.Type = ImportSpecType
			result.Path = path
			findResult = true
			return
		}

		impPkg, err := build.Import(path, walk.srcDir, build.ImportComment)
		if err != nil {
			result.Error = err
			return
		}

		walk.debug("deepWalk. *ast.ImportSpec. impPkg: %+v\n\n", impPkg)
		walk.debug("deepWalk. *ast.ImportSpec. impPkg.Dir: --%s--, impPkg.Name: %s, impPkg.Goroot: %t\n\n", impPkg.Dir, impPkg.Name, impPkg.Goroot)

		if impPkg.Name == name {
			result.GoRoot = impPkg.Goroot
			result.Path = path
			result.Type = ImportSpecType
			findResult = true
			return
		}

	case *ast.Field:
		f := in.(*ast.Field)
		walk.debug("\ndeepWalk. *ast.Field\n")
		if f.Type != nil {
			if f.Names != nil {
				// interface/func case
				for _, n := range f.Names {
					if n.String() == name {
						walk.debug("deepWalk. RETURN\n")
						result.Type = FieldType
						findResult = true
						return
					}
				}

				// object/func case
				if getObjType(f.Type) == name {
					result.Type = FieldType
					findResult = true
					return
				}
			}
		}
	}

	return
}

func getObjType(f interface{}) string {

	//fmt.Printf("getObjType IIIN . %T => %+v\n", f, f)

	switch f.(type) {
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", f.(*ast.StarExpr).X)
	case *ast.Field:
		w := f.(*ast.Field)
		if len(w.Names) > 0 {
			return w.Names[0].String()
		}
	case *ast.Ident:
		return f.(*ast.Ident).String()
	case *ast.FuncType:
		//w := f.(*ast.FuncType)
		//fmt.Printf("getObjType *ast.FuncType. %T => %+v\n", w, w)
		return ""
	default:
		panic(fmt.Sprintf("ERROR type for getObjType: %T => %+v", f, f))
	}

	return fmt.Sprintf("%s", f)
}

func (walk *Walker) FindImport(alias, fileSrc, path string) *WalkerResult {
	res := &WalkerResult{}
	if fileSrc != "" {
		res = walk.find("", fileSrc, alias)
	} else {
		res = walk.find(path, "", alias)
	}

	return res
}
