package onefunction

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"
)

/*
	Simple package for fill mocker with data
*/

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
	inOutValue   *Calls
	functionName string
	mockType     string
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
		inOutValue: NewOneCall(),
	}
}

func (w *MyWriter) Add(params, result []interface{}) *MyWriter {
	w.inOutValue.Add(params, result)
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
	//s = strings.ReplaceAll(s, `"`, `\"`)
	//s = strings.ReplaceAll(s, "\n", `\n`)
	return strconv.Quote(s)
}

func ToString(i interface{}) string {
	if i == nil {
		return "nil"
	}
	return _toString(reflect.ValueOf(i))
}

func _toString(v reflect.Value) string {

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

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.String:
		return Escape(v.String())
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
	}

	if v.Kind() == reflect.Struct {
		t := addPointer + typ.String()
		out := make([]string, typ.NumField())

		for i := 0; i < typ.NumField(); i++ {
			p := typ.Field(i)
			if !p.Anonymous {
				out[i] = _toString(v.Field(i))
			} else { // Anonymus structues
				out[i] = _toString(v.Field(i).Addr())
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
			out[i] = _toString(v.Index(i))
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
			out[i] = _toString(key) + `:` + _toString(v.MapIndex(key))
		}
		return t + "{" + strings.Join(out, ", ") + "}"
	}

	return "nil"
}

func (w *MyWriter) Code() string {

	data := &Data{
		MockType:         w.mockType,
		FunctionName:     w.functionName,
		FullFunctionName: w.FullFunctionName(),
		Calls:            w.inOutValue.ToStr(),
	}

	wr := &Collector{
		data: make([]byte, 0),
	}
	if err := qTmpl.ExecuteTemplate(wr, "code", data); err != nil {
		log.Panic(err)
	}

	return string(wr.data)
}
