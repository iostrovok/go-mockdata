package mockdata

/*
	Simple package for fill mocker with data
*/

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
	"text/template"

	"github.com/golang/mock/gomock"

	"github.com/iostrovok/go-mockdata/onefunction"
)

const tplCode = `
package mockdata

import ({{range .LocalPackages}}
	{{ . }}{{end}}
)

// {{.Constructor}} returns new mocker with data for {{.MockType}}
func {{.Constructor}} (t *testing.T) {{.MockType}} {
	ctrl := gomock.NewController(t)
	m := mmock.{{.Constructor}}(ctrl){{range .FunctionCalls}}
	m = {{ . }}(m){{end}}
	return m
}
{{range .FunctionBodies}}{{ . }}{{end}}
`

var qTmpl *template.Template

func init() {
	var err error
	qTmpl, err = template.New("code").Parse(tplCode)
	if err != nil {
		log.Panic(err)
	}
}

type Collector struct {
	data []byte
}

func (w *Collector) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

type Data struct {
	FunctionBodies []string
	FunctionCalls  []string
	LocalPackages  []string
	Constructor    string
	MockType       string
}

type Maker struct {
	mMock           interface{}
	object          interface{}
	currentFunction string
	lastCall        *onefunction.MyWriter

	LocalPackages []string
	Constructor   string
	MockType      string

	functions map[string]map[string]*onefunction.MyWriter
}

func New() *Maker {
	return &Maker{
		LocalPackages: []string{`"github.com/golang/mock/gomock"`, `"testing"`},
		functions:     map[string]map[string]*onefunction.MyWriter{},
	}
}

func (m *Maker) SetMMock(mMock interface{}) *Maker {
	m.mMock = mMock

	pkg, mFuncName, mOutType := SplitFunctionObject(m.mMock, true)
	m.LocalPackages = append(m.LocalPackages, pkg)
	m.Constructor = mFuncName
	m.MockType = mOutType

	fmt.Printf("\n------------------------\n\n\nm.Constructor: --%s--\n\n\n------------------------\n", m.Constructor)
	fmt.Printf("\n------------------------\n\n\nm.MockType: --%s--\n\n\n------------------------\n", m.MockType)

	return m
}

func (m *Maker) StartFunction(object interface{}) *Maker {
	m.object = object
	pkg, funcName, _ := SplitFunctionObject(m.object, false)
	m.LocalPackages = append(m.LocalPackages, pkg)
	m.currentFunction = funcName

	fmt.Printf("\n------------------------\n\n\nfuncName: --%s--\n\n\n------------------------\n", funcName)
	return m
}

func (m *Maker) InOut(in, out []interface{}) *Maker {

	if _, find := m.functions[m.MockType]; !find {
		m.functions[m.MockType] = map[string]*onefunction.MyWriter{}
	}

	if _, find := m.functions[m.MockType][m.currentFunction]; !find {
		m.functions[m.MockType][m.currentFunction] = onefunction.New().
			FunctionName(m.currentFunction).
			MockType(m.MockType)
	}

	m.functions[m.MockType][m.currentFunction].Add(in, out)

	return m
}

func (m *Maker) Code() string {

	functionBodies := make([]string, 0)
	functionCalls := make([]string, 0)
	for mockName, functions := range m.functions {
		for funcName, w := range functions {

			fmt.Printf("generate %s/%s", mockName, funcName)

			functionCalls = append(functionCalls, w.FullFunctionName())
			functionBodies = append(functionBodies, w.Code())

			fmt.Printf("functionCalls::: %+v\n", functionCalls)
			fmt.Printf("functionBodies::: %+v\n", functionBodies)
		}
	}

	//for i := range m.LocalPackages {
	//	m.LocalPackages[i] = strconv.Quote(m.LocalPackages[i])
	//}

	data := &Data{
		FunctionBodies: functionBodies,
		FunctionCalls:  functionCalls,
		LocalPackages:  uniqStrArray(m.LocalPackages),
		Constructor:    m.Constructor,
		MockType:       m.MockType,
	}

	wr := &Collector{
		data: make([]byte, 0),
	}
	if err := qTmpl.ExecuteTemplate(wr, "code", data); err != nil {
		log.Panic(err)
	}

	fmt.Printf("\n\n")
	fmt.Printf("%s", string(wr.data))
	fmt.Printf("\n\n")
	return string(wr.data)
}

type functionParts struct {
	pkg, f, outType string
}

func SplitFunctionObject(i interface{}, checkResult bool) (string, string, string) {

	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Func {
		panic("should be function")
	}

	outType := ""
	if checkResult {
		ctrl := gomock.NewController(&testing.T{})
		res := value.Call([]reflect.Value{reflect.ValueOf(ctrl)})
		outType = res[0].Type().String()
	}

	dir, file := filepath.Split(runtime.FuncForPC(value.Pointer()).Name())

	parts := strings.SplitN(file, ".", 3)

	pkg := `"` + filepath.Join(dir, parts[0]) + `"`
	funcParts := strings.SplitN(parts[len(parts)-1], "-", 2)
	funcName := funcParts[0]

	fmt.Printf("\n\n-----------------\nSplitFunctionObject >>>>>>>>\n\n")
	fmt.Printf("pkg: %s\n", pkg)
	fmt.Printf("funcName: %s\n", funcName)
	fmt.Printf("out.outType: %s\n", outType)

	return pkg, funcName, outType
}

func uniqStrArray(in []string) []string {
	check := map[string]bool{}
	out := make([]string, 0)
	for _, s := range in {
		if !check[s] {
			out = append(out, s)
			check[s] = true
		}
	}

	sort.Strings(out)

	return out
}
