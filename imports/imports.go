package imports

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/iostrovok/go-mockdata/imports/inter"
)

/*
	Simple package for fill mocker with data
*/

type alias struct {
	value string
	used  bool
}

func setAlias(a []*alias, value string, used bool) []*alias {

	for i := range a {
		if a[i].value == value {
			if used {
				a[i].used = true
			}
			return a
		}
	}

	return append(a, &alias{value, used})
}

type OneFile struct {
	file, pkgName, path string
	aliases             []*alias
	whereUsed           map[string]bool
}

func newOneFile(file, pkgName, path, whereUsed string) *OneFile {
	return &OneFile{
		file:    file,
		pkgName: pkgName,
		path:    path,
		aliases: make([]*alias, 0),
		whereUsed: map[string]bool{
			whereUsed: true,
		},
	}
}

type Imports struct {
	files []*OneFile
}

func New() inter.IImp {
	return &Imports{
		files: make([]*OneFile, 0),
	}
}

func absFile(fileIn string) string {

	if fileIn == "" {
		return ""
	}

	if strings.HasSuffix(fileIn, ".go") {
		fileIn = filepath.Dir(fileIn)
	}

	fileOut, _ := filepath.Abs(fileIn)
	return fileOut
}

func (imp *Imports) find(file, path, whereUsed string) *OneFile {

	for i := range imp.files {

		if whereUsed

		if imp.files[i].whereUsed[whereUsed] || (path != "" && imp.files[i].path == path) || (file != "" && imp.files[i].file == file) {
			return imp.files[i]
		}
	}

	return nil
}

func (imp *Imports) Add(file, pkg, path, whereUsed, alias string) {

	fmt.Printf("*Imports.ADDDDD file: %s, pkg: %s, path: %s, alias: %s, whereUsed: %s\n", file, pkg, path, alias, whereUsed)

	path = strings.Trim(path, `"`)
	alias = strings.Trim(alias, `"`)
	pkg = strings.Trim(pkg, `"`)
	file = strings.Trim(file, `"`)
	whereUsed = strings.Trim(whereUsed, `"`)

	file = absFile(file)

	p := imp.find(file, path)
	if p == nil {
		p = newOneFile(file, pkg, path, whereUsed)
		imp.files = append(imp.files, p)
	} else {
		if p.file == "" {
			p.file = file
		}
		if p.path == "" {
			p.path = path
		}
		if p.pkgName == "" {
			p.pkgName = pkg
		}
	}

	p.aliases = setAlias(p.aliases, pkg, false)
	if alias != "" {
		p.aliases = setAlias(p.aliases, alias, false)
	}
}

func (imp *Imports) UseWhere(whereUsed, alias string) error {
	return imp.use("", "", whereUsed, alias)
}

func (imp *Imports) UseFile(file, alias string) error {
	return imp.use(file, "", "", alias)
}

func (imp *Imports) UsePath(path, alias string) error {
	return imp.use("", path, "", alias)
}

func (imp *Imports) use(file, path, whereUsed, alias string) error {

	alias = strings.Trim(alias, `"`)
	file = strings.Trim(file, `"`)
	path = strings.Trim(path, `"`)
	whereUsed = strings.Trim(whereUsed, `"`)

	file = absFile(file)

	p := imp.find(file, path)
	if p == nil {
		return fmt.Errorf("file '%s' or path '%s' not found", file, path)
	}

	if alias != "" {
		p.aliases = setAlias(p.aliases, alias, true)
	} else {
		p.aliases = setAlias(p.aliases, p.pkgName, true)
	}

	return nil
}

func (imp *Imports) Dump() string {
	out := make([]string, 0)
	for i, f := range imp.files {
		out = append(out, "")
		out = append(out, fmt.Sprintf("*Imports List() ..........  imp.files[%d]: %s, %s, %s", i, f.pkgName, f.path, f.file))
		for j, a := range f.aliases {
			out = append(out, fmt.Sprintf("				.......... aliases[%d]: %+v", j, a))
		}
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}

func (imp *Imports) List() []string {

	fmt.Println(imp.Dump())

	paths := make([]string, 0)
	for _, p := range imp.files {
		if p.path != "" {
			paths = append(paths, p.path)
		}
	}
	sort.Strings(paths)

	out := make([]string, 0)
	for _, path := range paths {

		p := imp.find("", path)
		if p == nil {
			continue
		}

		res := make([]string, 0)
		for _, a := range p.aliases {
			if a.used {
				if a.value != p.pkgName {
					res = append(res, a.value+" "+strconv.Quote(path))
				} else {
					res = append(res, strconv.Quote(path))
				}
			}
		}
		res = uniqStrArray(res)
		out = append(out, res...)
	}

	return uniqStrArray(out, false)
}

func uniqStrArray(in []string, noSort ...bool) []string {

	noSort = append(noSort, false)

	check := map[string]bool{}
	out := make([]string, 0)
	for _, s := range in {
		if !check[s] {
			out = append(out, s)
			check[s] = true
		}
	}

	if !noSort[0] {
		sort.Strings(out)
	}

	return out
}
