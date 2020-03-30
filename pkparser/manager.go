package pkparser

import (
	"crypto/md5"
	"fmt"
	"go/ast"
	"math/rand"
	"sync"
	"time"
)

/*
	Simple walk through parsed package tree.
*/

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Element struct {
	Alias, Path string
	GoRoot      bool
	Name        string
}

type Manager struct {
	sync.RWMutex

	P *PkParser
	W *Walker

	srcDir string
}

func randomAlias(x string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(x)))[:5]
}

func New(srcDir string) *Manager {
	return &Manager{
		srcDir: srcDir,
		P:      NewParser(srcDir).SetDebug(false),
		W:      NewWalker(srcDir).SetDebug(false),
	}
}

func (m *Manager) Find(in interface{}) (interface{}, error) {

	sf, err := SplitFunction(in)
	if err != nil {
		return nil, err
	}

	if sf.FileSrc != "" {
		if m.P.parsePackage(sf.FileSrc, sf.Path); err != nil {
			return nil, err
		}
	} else if sf.Path != "" {
		if m.P.parsePackageByPath(sf.Path); err != nil {
			return nil, err
		}
	}

	imp := make([][]string, 0)

	m.W.AddFiles(m.P.parsedFiles)
	alias := randomAlias(sf.Path)
	receiver := alias + "." + sf.Receiver
	res := m.W.FindPath(sf.Path, sf.Receiver)

	//receiverType := getType(res.Node)
	//if receiverType != "" {
	//	imp = append(imp, []string{alias, sf.Path})
	//}

	fmt.Printf("FFFFF res: %+v\n", res)
	fmt.Printf("FFFFF res: %+v\n", res.Node)
	fmt.Printf("FFFFF res: %T\n", res.Node.Node)
	fmt.Printf("FFFFF res: %+v\n", res.Node.Node.(*ast.TypeSpec).Type)
	fmt.Printf("FFFFF res: %+v\n", res.Node.Node.(*ast.TypeSpec).Type)

	if TypeSpecType == res.Node.Type {
		//m.W.SetDebug(true)
		f, find := m.W.deepWalk(res.Node.Node, sf.Name)
		//m.W.SetDebug(false)
		fmt.Printf("TTTTTTTTT: %T, f: %+v, find: %t\n", f, f, find)
		if find {
			typeF := getType(f)
			if typeF == "Func" {
				w := f.Node.(*ast.Field).Type.(*ast.FuncType)
				params, result := w.Params, w.Results

				fmt.Printf("params: %T\n", params)
				fmt.Printf("params: %+v\n", params)
				pp, err := m.extractFieldList(params, sf.FileSrc, sf.Path)
				if err != nil {
					return nil, err
				}
				fmt.Printf("params.pp: %+v\n", pp)

				fmt.Printf("result: %T\n", result)
				fmt.Printf("result: %+v\n", result)
				rr, err := m.extractFieldList(result, sf.FileSrc, sf.Path)
				if err != nil {
					return nil, err
				}

				fmt.Printf("\n\n TMP OT:\n")
				fmt.Printf("params.pp: %+v\n", pp)
				fmt.Printf("result.pp: %+v\n", rr)
			}
		}
	}

	fmt.Printf("sf.Path: %s, sf.Receiver: %s\n", sf.Path, sf.Receiver)
	fmt.Printf("file: %+v, node: %+v, find: %t\n", res.ParsedFile, res.Node, res.Find)

	fmt.Printf("imp: %+v\n", imp)
	fmt.Printf("receiver: %s\n", receiver)

	return nil, nil
}

func getType(res *NodeResult) string {

	fmt.Printf("getType.res: %+v\n", res)

	switch res.Type {
	case TypeSpecType:
		switch res.Node.(*ast.TypeSpec).Type.(type) {
		case *ast.InterfaceType:
			return "Interface"
		}

	case FieldType:
		switch res.Node.(*ast.Field).Type.(type) {
		case *ast.FuncType:
			return "Func"
		}
	}

	return ""
}

func (m *Manager) extractFieldList(res interface{}, fileSrc, path string) ([]Element, error) {

	out := make([]Element, 0)

	wL := make([]*ast.Field, 0)
	switch res.(type) {
	case *ast.FieldList:
		wL = res.(*ast.FieldList).List
	default:
		return out, nil
	}

	if wL == nil || len(wL) == 0 {
		return out, nil
	}

	for i, w := range wL {
		count := 1
		if w.Names != nil && len(w.Names) > 0 {
			count = len(w.Names)
		}
		for j := 0; j < count; j++ {
			fmt.Printf("extractFieldList ===> i: %d, j: %d, w: %+v\n", i, j, w)
			fmt.Printf("extractFieldList ===> i: %d, j: %d, w.Type: %T, w.Type: %+v\n", i, j, w.Type, w.Type)
			switch w.Type.(type) {
			case *ast.Ident:
				e := Element{
					Alias: randomAlias(path),
					Path:  path,
					Name:  w.Type.(*ast.Ident).String(),
				}
				out = append(out, e)
			case *ast.SelectorExpr:
				ww := w.Type.(*ast.SelectorExpr)
				alias := fmt.Sprintf("%s", ww.X)
				fmt.Printf("extractFieldList ===> w.Type.(*ast.SelectorExpr).X: %+v\n", ww.X)
				if alias != "" {
					m.W.SetDebug(true)
					fmt.Printf("\n\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n")
					fmt.Printf("extractFieldList ===> alias: %+v\n", alias)
					imp := m.W.FindImport(alias, fileSrc, path)
					fmt.Printf("\n<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n\n")
					m.W.SetDebug(false)
					if imp.Node.Error != nil {
						return out, imp.Node.Error
					}
					if imp.Find {
						fmt.Printf("extractFieldList ===> imp: %+v\n", imp.Node.Path)
						e := Element{
							Alias:  alias,
							Path:   imp.Node.Path,
							GoRoot: imp.Node.GoRoot,
							Name:   ww.Sel.String(),
						}
						out = append(out, e)
					}
				} else {
					e := Element{
						Alias: randomAlias(path),
						Path:  path,
						Name:  ww.Sel.String(),
					}
					out = append(out, e)
				}
			}
		}
	}

	return out, nil
}
