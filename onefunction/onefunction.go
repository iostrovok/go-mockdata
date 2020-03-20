package onefunction

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/iostrovok/go-mockdata/receivers"
)

/*
	Simple package for fill mocker with data
*/

type SaveStringFunc func(v reflect.Value, limit int) (res string, find bool)

const (
	tplCode = `
func {{.FullFunctionName}}(m {{.MockType}}) {{.MockType}} {
{{range .Calls}}	m.EXPECT().{{.FunctionName}}({{.Params}}).Return({{.Result}})
{{end}}	return m
}
`
	TimeFormat = "2006-01-02T15:04:05.999999999"
)

var qTmpl *template.Template

func init() {
	var err error
	qTmpl, err = template.New("code").Parse(tplCode)
	if err != nil {
		log.Panic(err)
	}
}

type Data struct {
	Calls            []StrCalls
	FunctionName     string
	FullFunctionName string
	MockType         string
}

type MyWriter struct {
	inOutValue      *Calls
	functionName    string
	mockType        string
	maxStringLength int
	userFunc        SaveStringFunc
	ParserResult    receivers.OneFunctionRes
}

const tmpl = `func MockIVerification{{.FunctionName}}(m {{.MMockType}}) {{.MMockType}} {
	m.EXPECT().{{.FunctionName}}().Return([]*base.CountGroups{
		{
			Count:   9,
			GroupBy: 0,
		},
		{
			Count:   10,
			GroupBy: 1,
		},
	}, nil)
	return m
}
`

type Collector struct {
	data []byte
}

func (w *Collector) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func New() *MyWriter {
	return &MyWriter{
		inOutValue:      NewOneCall(),
		maxStringLength: -1,
	}
}

func (w *MyWriter) SetParserResult(parserResult receivers.OneFunctionRes) *MyWriter {
	w.ParserResult = parserResult
	return w
}

func (w *MyWriter) SaveStringFunc(userFunc SaveStringFunc) *MyWriter {
	w.userFunc = userFunc
	return w
}

func (w *MyWriter) Add(params, result []interface{}) *MyWriter {
	w.inOutValue.Add(params, result)
	return w
}

func (w *MyWriter) StringLimit(maxStringLength int) *MyWriter {
	w.maxStringLength = maxStringLength
	return w
}

func (w *MyWriter) FunctionName(functionName string) *MyWriter {
	w.functionName = functionName
	w.inOutValue.Name(functionName)
	return w
}

func (w *MyWriter) FullFunctionName() string {
	sp := strings.Split(w.mockType, ".")
	return sp[len(sp)-1] + w.functionName
}

func (w *MyWriter) MockType(mockType string) *MyWriter {
	w.mockType = mockType
	return w
}

func Escape(s string) string {
	return strconv.Quote(s)
}

func limitString(s string, limit int) string {
	if limit == -1 || len(s) <= limit {
		return s
	}

	return s[:limit]
}

func ToString(i interface{}, userFunc SaveStringFunc, limit ...int) string {
	if i == nil {
		return "nil"
	}

	limit = append(limit, -1)
	return _toString(reflect.ValueOf(i), userFunc, limit[0])
}

func _toString(v reflect.Value, userFunc SaveStringFunc, limit int) string {

	typ := v.Type()
	addPointer := ""
	if typ.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "nil"
		}

		typ = typ.Elem()
		v = v.Elem()
		addPointer = "&"
	}

	if userFunc != nil {
		if str, find := userFunc(v, limit); find {
			fmt.Printf("v: %v, limit: %d\n", v, limit)
			fmt.Printf("str: %s, find: %t\n", str, find)
			return str
		}
	}

	if typ.String() == "interface {}" {
		return _toString(reflect.ValueOf(v.Interface()), userFunc, limit)
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.String:
		return Escape(limitString(v.String(), limit))
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Func, reflect.UnsafePointer, reflect.Chan, reflect.Uintptr:
		return "nil"
	}

	//
	switch typ.String() {
	case "time.Time":
		return Escape(v.Interface().(time.Time).Format(TimeFormat))
	case "[]uint8":
		return "[]uint8(" + Escape(limitString(string(v.Interface().([]byte)), limit)) + ")"
	}

	if v.Kind() == reflect.Struct {
		t := addPointer + typ.String()
		out := make([]string, typ.NumField())

		for i := 0; i < typ.NumField(); i++ {
			p := typ.Field(i)
			if !p.Anonymous {
				out[i] = _toString(v.Field(i), userFunc, limit)
			} else { // Anonymus structues
				out[i] = _toString(v.Field(i).Addr(), userFunc, limit)
			}
		}

		return t + "{" + strings.Join(out, ", ") + "}"
	}

	if v.Kind() == reflect.Slice {
		t := addPointer + typ.String()
		if v.Len() == 0 {
			return "{}"
		}

		out := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			out[i] = _toString(v.Index(i), userFunc, limit)
		}
		return t + "{" + strings.Join(out, ", ") + "}"
	}

	if v.Kind() == reflect.Map {
		t := addPointer + typ.String()
		keys := v.MapKeys()

		if len(keys) == 0 {
			return "{}"
		}

		out := make([]string, len(keys))
		for i, key := range keys {
			out[i] = _toString(key, nil, -1) + `:` + _toString(v.MapIndex(key), userFunc, limit)
		}
		return t + "{" + strings.Join(out, ", ") + "}"
	}

	return "nil"
}

func (w *MyWriter) Code() string {

	fmt.Printf("Code:Params: %+v\n", w.ParserResult.Params)
	fmt.Printf("Code:Results: %+v\n", w.ParserResult.Results)

	data := &Data{
		MockType:         w.mockType,
		FunctionName:     w.functionName,
		FullFunctionName: w.FullFunctionName(),
		Calls:            w.inOutValue.ToStr(w.userFunc, w.ParserResult),
	}

	wr := &Collector{
		data: make([]byte, 0),
	}
	if err := qTmpl.ExecuteTemplate(wr, "code", data); err != nil {
		log.Panic(err)
	}

	return string(wr.data)
}
