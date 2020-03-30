package pkparser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"sync"

	"github.com/iostrovok/go-mockdata/imports/inter"
)

/*
	Simple package for fill mocker with data
*/

type PkParser struct {
	sync.RWMutex

	parsed []*ast.Package
	srcDir string
	imp    inter.IImp
}

func New(srcDir string, imp inter.IImp) *PkParser {
	return &PkParser{
		srcDir: srcDir,
		parsed: make([]*ast.Package, 0),
		imp:    imp,
	}
}

func (pkp *PkParser) ByFile(fileName string) (*ast.Package, bool) {

	for _, p := range pkp.parsed {
		if _, find := p.Files[fileName]; find {
			return p, true
		}
	}
	return nil, false
}

func extractPackageNamByFile(fileName string) (string, error) {
	src, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}
	return src.Name.Name, nil
}

// parsePackage loads package specified by path, parses it and populates
// corresponding imports and importedInterfaces into the fileParser.
func (pkp *PkParser) ParsePackageByPath(path string) error {

	if pkp.srcDir == "" {
		panic("srcDir is empty")
		return errors.New("srcDir is empty")
	}

	path = strings.Trim(path, `"`)

	fmt.Printf("ParsePackageByPath. pkp.srcDir: %s, path: --%s--\n", pkp.srcDir, path)

	impPkg, err := build.Import(path, pkp.srcDir, build.FindOnly)
	if err != nil {
		return err
	}

	fmt.Printf("ParsePackageByPath. impPkg.Dir: %+v\n", impPkg)
	fmt.Printf("ParsePackageByPath. impPkg.Dir: %s, impPkg.Name: %s, impPkg.Goroot: %t\n", impPkg.Dir, impPkg.Name, impPkg.Goroot)

	if impPkg.Goroot {
		// package from std lib
		return nil
	}

	_, err = pkp.parsePackage(impPkg.Dir, impPkg.Name)
	return err
}

func (pkp *PkParser) ParsePackageByFile(fileName string) (*ast.Package, error) {
	fmt.Printf("ParsePackage.......... %s\n", fileName)

	if p, find := pkp.ByFile(fileName); find {
		return p, nil
	}

	pkgName, err := extractPackageNamByFile(fileName)
	if err != nil {
		return nil, err
	}

	fmt.Printf("pkgName: %s\n", pkgName)

	srcDir, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		return nil, fmt.Errorf("failed getting source directory: %v", err)
	}

	fmt.Printf("ParsePackageByFile. pkgName: %s, srcDir: %s\n", pkgName, srcDir)

	imports, err := pkp.parsePackage(srcDir, pkgName)
	if err != nil {
		return nil, err
	}

	for _, oneFile := range imports {
		if err := pkp.ParsePackageByPath(oneFile.path); err != nil {
			return nil, err
		}
	}

	p, find := pkp.ByFile(fileName)
	if !find {
		err = fmt.Errorf("failed getting file %s", fileName)
	}
	return p, err
}

type OneFile struct {
	file, pckName, path, alias string
}

func (pkp *PkParser) parsePackage(srcDir, pkgName string) ([]OneFile, error) {

	imports := make([]OneFile, 0)

	pkp.Lock()
	defer pkp.Unlock()

	fmt.Printf("ParsePackage. pkgName: %s, srcDir: %s\n", pkgName, srcDir)

	parsed, err := parser.ParseDir(token.NewFileSet(), srcDir, nil, parser.AllErrors)
	if err != nil {
		return imports, err
	}

	fmt.Printf("\nParsed: %+v\n\n", parsed)

	//if _, find := parsed[pkgName]; !find {
	//	return imports, fmt.Errorf("ParsePackage: not found pacakge %s in dir", pkgName)
	//}

	for _, p := range parsed {

		fmt.Printf("\nparsed: %s\n", p.Name)
		fmt.Printf("Imports PP: %+v\n", p)
		fmt.Printf("Files: %+v\n\n", p.Files)

		pkp.parsed = append(pkp.parsed, p)

		for fileName, file := range p.Files {
			fmt.Printf(".....fileName: %s, file.Imports: %+v\n", fileName, file.Imports)

			for i, s := range file.Imports {
				fmt.Printf("\n%d .....range Name: %T, Name: %+v\n", i, s.Name, s.Name)
				fmt.Printf("%d .....range Doc: %T, Doc: %+v\n", i, s.Doc, s.Doc)
				fmt.Printf("%d .....range Path: %T, Path: %+v\n", i, s.Path, s.Path)
				fmt.Printf("%d .....range Comment: %T, Comment: %+v\n", i, s.Comment, s.Comment)
				fmt.Printf("%d .....range s: %s\n", i, s)
				pkg := file.Name.String()
				alias := ""
				if s.Name != nil {
					alias = s.Name.String()
				}
				uri := s.Path.Value
				fmt.Printf("fileName: %s, pkg: %s, s.Name: %s, uri: %s\n\n", fileName, pkg, alias, uri)

				pkp.imp.Add(fileName, file.Name.String(), s.Path.Value, alias)
				if file.Name.String() != "" {
					imports = append(imports, OneFile{
						file:    fileName,
						pckName: file.Name.String(),
						path:    s.Path.Value,
						alias:   alias,
					})
				}
			}
		}
	}

	return imports, nil
}
