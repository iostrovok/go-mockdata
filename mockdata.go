package mockdata

/*
	Simple package for fill mocker with data
*/

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
	"text/template"

	"github.com/golang/mock/gomock"

	"github.com/iostrovok/go-mockdata/imports"
	"github.com/iostrovok/go-mockdata/onefunction"
	"github.com/iostrovok/go-mockdata/receivers"
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

type Recorder struct {
	mMock           interface{}
	object          interface{}
	currentFunction string
	lastCall        *onefunction.MyWriter

	Imports      *imports.Imports
	Parser       *receivers.Parser
	ParserResult receivers.OneFunctionRes
	Constructor  string
	MockType     string

	maxStringLength int

	functions map[string]map[string]*onefunction.MyWriter
}

func New() *Recorder {

	im := imports.New()
	im.Add("testing", "")
	im.Add("github.com/golang/mock/gomock", "")

	par := receivers.New(im)

	return &Recorder{
		maxStringLength: -1,
		Imports:         im,
		Parser:          par,
		functions:       map[string]map[string]*onefunction.MyWriter{},
	}
}

func (m *Recorder) StringLimit(maxStringLength int) *Recorder {
	m.maxStringLength = maxStringLength
	return m
}

func (m *Recorder) SetMMock(mMock interface{}) *Recorder {
	m.mMock = mMock

	value := reflect.ValueOf(mMock)
	if value.Kind() != reflect.Func {
		panic("should be function")
	}

	ctrl := gomock.NewController(&testing.T{})
	res := value.Call([]reflect.Value{reflect.ValueOf(ctrl)})
	m.MockType = res[0].Type().String()

	ofr := m.Parser.Run(mMock)

	fmt.Printf("SetMMock .......... m.MockType: %s\n", m.MockType)
	fmt.Printf("SetMMock .......... ofr.Params: %+v\n", ofr)

	dir, file := filepath.Split(runtime.FuncForPC(value.Pointer()).Name())

	parts := strings.SplitN(file, ".", 3)
	pkgName := parts[0]

	m.Imports.Add(filepath.Join(dir, pkgName), "")
	m.Constructor = ofr.FuncName

	return m
}

func (m *Recorder) StartFunction(object interface{}) *Recorder {
	m.object = object

	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Func {
		panic("should be function")
	}

	m.ParserResult = m.Parser.Run(object)

	// set package / alias for setup local vars
	m.Imports.SetCurrentPackage(m.ParserResult.Pkg)

	fmt.Printf("StartFunction .......... ofr.Params: %+v\n", m.ParserResult)
	fmt.Printf("StartFunction .......... ofr.Params: %+v\n", m.Imports.List())
	m.currentFunction = m.ParserResult.FuncName

	return m
}

func (m *Recorder) Add(in, out []interface{}) *Recorder {

	if _, find := m.functions[m.MockType]; !find {
		m.functions[m.MockType] = map[string]*onefunction.MyWriter{}
	}

	if _, find := m.functions[m.MockType][m.currentFunction]; !find {
		m.functions[m.MockType][m.currentFunction] = onefunction.New().
			StringLimit(m.maxStringLength).
			FunctionName(m.currentFunction).
			MockType(m.MockType)

	}

	m.functions[m.MockType][m.currentFunction].SetParserResult(m.ParserResult)
	m.functions[m.MockType][m.currentFunction].Add(in, out)

	return m
}

func (m *Recorder) Save(fileName string) error {
	return ioutil.WriteFile(fileName, []byte(m.Code()), 0666)
}

// Code returns
func (m *Recorder) Code() string {

	functionBodies := make([]string, 0)
	functionCalls := make([]string, 0)
	for _, functions := range m.functions {
		for _, w := range functions {
			functionCalls = append(functionCalls, w.FullFunctionName())
			functionBodies = append(functionBodies, w.Code())
		}
	}

	data := &Data{
		FunctionBodies: functionBodies,
		FunctionCalls:  functionCalls,
		LocalPackages:  m.Imports.List(),
		Constructor:    m.Constructor,
		MockType:       m.MockType,
	}

	wr := &Collector{
		data: make([]byte, 0),
	}
	if err := qTmpl.ExecuteTemplate(wr, "code", data); err != nil {
		log.Panic(err)
	}

	src, err := format.Source(wr.data)
	if err != nil {
		log.Fatalf("Failed to format generated source code: %s\n", err)
	}

	return string(src)
}

func (m *Recorder) Clean() {
	m.Imports = imports.New()
	m.Imports.Add("testing", "")
	m.Imports.Add("github.com/golang/mock/gomock", "")

	m.Parser = receivers.New(m.Imports)
	m.ParserResult = receivers.OneFunctionRes{}

	m.mMock = nil
	m.object = nil
	m.currentFunction = ""
	m.lastCall = nil
	m.Constructor = ""
	m.MockType = ""
	m.functions = map[string]map[string]*onefunction.MyWriter{}
}

func (m *Recorder) SplitFunctionObject(i interface{}, checkResult bool) (string, string) {

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

	m.ParserResult = m.Parser.Run(i)

	dir, file := filepath.Split(runtime.FuncForPC(value.Pointer()).Name())

	parts := strings.SplitN(file, ".", 3)

	pkgName := parts[0]

	m.Imports.Add(filepath.Join(dir, pkgName), "")
	funcName := m.ParserResult.FuncName

	return funcName, outType
}

func (m *Recorder) SplitFunctionObject_old(i interface{}, checkResult bool) (string, string, string) {

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

	pkgName := parts[0]
	fmt.Printf("pkgName: %s\n", pkgName)
	fmt.Printf("runtime.FuncForPC(value.Pointer()).Name(): %s\n", runtime.FuncForPC(value.Pointer()).Name())
	fmt.Printf("filepath.Split(runtime.FuncForPC(value.Pointer()).Entry(): %+v\n", runtime.FuncForPC(value.Pointer()).Entry())

	pkg := `"` + filepath.Join(dir, pkgName) + `"`
	funcParts := strings.SplitN(parts[len(parts)-1], "-", 2)
	funcName := funcParts[0]

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
