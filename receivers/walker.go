package receivers

/*
	Simple walk through parsed package tree.
*/

import (
	"fmt"
	"go/ast"
)

// Walker finds path [names, names ...] in parsed file or package
func Walker(in interface{}, names ...string) (out interface{}, find bool) {

	for _, name := range names {
		out, find = deepWalk(in, name)
		if !find {
			break
		}
	}

	return
}

func deepWalk(in interface{}, name string) (interface{}, bool) {

	//fmt.Printf("\ndeepWalk. name: %s, in: %T => %+v\n", name, in, in)

	switch in.(type) {
	case *ast.Package:
		w := in.(*ast.Package)
		for _, d := range w.Files {
			if out, find := deepWalk(d, name); find {
				return out, find
			}
		}

	case *ast.File:
		w := in.(*ast.File)
		if w.Name.String() == name {
			return w, true
		}
		for _, d := range w.Decls {
			if out, find := deepWalk(d, name); find {
				return out, find
			}
		}

	case *ast.GenDecl:
		w := in.(*ast.GenDecl)
		if w.Specs != nil && len(w.Specs) > 0 {
			for _, d := range w.Specs {
				if out, find := deepWalk(d, name); find {
					return out, find
				}
			}
		}

	case *ast.TypeSpec:
		w := in.(*ast.TypeSpec)
		if w.Name.String() == name {
			return w, true
		}
		if out, find := deepWalk(w.Type, name); find {
			return out, find
		}

	case *ast.InterfaceType:
		w := in.(*ast.InterfaceType)
		if w.Methods.List != nil {
			for _, method := range w.Methods.List {
				if out, find := deepWalk(method, name); find {
					return out, find
				}
			}
		}

	case *ast.FuncDecl:
		w := in.(*ast.FuncDecl)
		if w.Name.String() == name {
			return w, true
		}
		if out, find := deepWalk(w.Type, name); find {
			return out, find
		}

		// case for "func (s *obj) Name()"
		if w.Recv != nil {
			if _, find := deepWalk(w.Recv, name); find {
				return w, true
			}
		}

	case *ast.FieldList:
		w := in.(*ast.FieldList)
		if w.List != nil {
			for _, method := range w.List {
				if out, find := deepWalk(method, name); find {
					return out, find
				}
			}
		}

	case *ast.Field:
		w := in.(*ast.Field)
		//fmt.Printf("deepWalk. 1.ast.Field: %s, in: %T => %+v\n", name, w, w)
		if w.Type != nil {
			//fmt.Printf("deepWalk. 2.ast.Field: %s, in: %T => %+v\n", name, w.Names, w.Names)
			if w.Names != nil {
				// interface/func case
				for _, n := range w.Names {
					//fmt.Printf("deepWalk. 3.ast.Field: %s, in: %T => ---%s---\n", name, n, n.String())
					if n.String() == name {
						fmt.Printf("deepWalk. RETURN\n")
						return w, true
					}
				}

				// object/func case
				if getObjType(w.Type) == name {
					return w, true
				}
			}
		}
	}

	return nil, false
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
