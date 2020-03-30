package pkparser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

/*
	Simple package for fill mocker with data
*/

type ParsedFile struct {
	file    *ast.File
	pkgName string
	path    string
	srcDir  string
	srcFile string
}

type PkParser struct {
	parsedFiles map[string]*ParsedFile

	srcDir    string
	parsedDir map[string]bool
	viewDebug bool

	deep int
}

func NewParser(srcDir string) *PkParser {
	return &PkParser{
		srcDir:      srcDir,
		parsedFiles: map[string]*ParsedFile{},
		parsedDir:   map[string]bool{},
		deep:        3,
	}
}

func (pkp *PkParser) SetDeep(deep int) *PkParser {
	pkp.deep = deep
	return pkp
}

func (pkp *PkParser) SetDebug(viewDebug bool) *PkParser {
	pkp.viewDebug = viewDebug
	return pkp
}

func (pkp *PkParser) debug(format string, d ...interface{}) {
	if pkp.viewDebug {
		fmt.Printf(format, d...)
	}
}

func (pkp *PkParser) setPathByDir(src, path string) {
	for f := range pkp.parsedFiles {
		if strings.Index(f, src) == 0 {
			pkp.parsedFiles[f].path = path
		}
	}
}

func (pkp *PkParser) parsePackageByPath(path string, deep ...int) error {

	if pkp.srcDir == "" {
		return errors.New("srcDir is empty")
	}

	deep = append(deep, 3)
	path = strings.Trim(path, `"`)

	impPkg, err := build.Import(path, pkp.srcDir, build.FindOnly)
	if err != nil {
		return err
	}

	pkp.debug("ParsePackageByPath. impPkg.Dir: %+v\n", impPkg)
	pkp.debug("ParsePackageByPath. impPkg.Dir: --%s--, impPkg.Name: %s, impPkg.Goroot: %t\n", impPkg.Dir, impPkg.Name, impPkg.Goroot)

	if pkp.parsedDir[impPkg.Dir] || impPkg.Goroot {
		return nil
	}

	return pkp._parsePackage(impPkg.Dir, path, deep[0])
}

func (pkp *PkParser) parsePackage(srcIn string, path ...string) error {
	path = append(path, "")
	return pkp._parsePackage(srcIn, path[0], pkp.deep)
}

func (pkp *PkParser) _parsePackage(srcIn, path string, deep int) error {

	deep--
	if deep < 0 {
		return nil
	}

	src, err := filepath.Abs(srcIn)
	if err != nil {
		return err
	}

	if itIsFile, err := isFile(src); err != nil {
		return err
	} else if itIsFile {
		src = filepath.Dir(src)
	}

	srcDir, err := filepath.Abs(src)
	if err != nil {
		panic(fmt.Sprintf("receivers has got wrong dir: %v", err))
	}

	if pkp.parsedDir[srcDir] {
		if path != "" {
			pkp.setPathByDir(srcDir, path)
		}
		return nil
	}

	pkp.debug("ParsePackage. srcDir: %s\n", srcDir)

	filter := func(file os.FileInfo) bool {
		return !strings.HasSuffix(file.Name(), "_test.go")
	}

	parsed, err := parser.ParseDir(token.NewFileSet(), srcDir, filter, parser.AllErrors)
	if err != nil {
		return err
	}

	pkp.parsedDir[srcDir] = true

	for _, p := range parsed {

		pkp.debug("\nparsed: %s\n", p.Name)
		pkp.debug("Imports PP: %+v\n", p)
		pkp.debug("Files: %+v\n\n", p.Files)

		for srcFile, file := range p.Files {

			//ast.SortImports(token.NewFileSet(), file)

			pkp.parsedFiles[srcFile] = &ParsedFile{
				file:    file,
				pkgName: p.Name,
				path:    path,
				srcDir:  srcDir,
				srcFile: srcFile,
			}

			pkp.debug(".....fileName: %s, file.Imports: %+v\n", srcFile, file.Imports)

			for _, s := range file.Imports {
				if err := pkp.parsePackageByPath(s.Path.Value, deep); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func isFile(path string) (bool, error) {
	info, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return false, nil
	}

	return true, nil
}
