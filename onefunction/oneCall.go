package onefunction

import (
	"strings"
)

/*
	Simple package for fill mocker with data
*/

type CallParams struct {
	params []interface{}
	result []interface{}
}

type Calls struct {
	inOut        []CallParams
	functionName string
}

type StrCalls struct {
	Params, Result string
	FunctionName   string
}

func NewOneCall() *Calls {
	return &Calls{
		inOut:        []CallParams{},
		functionName: "",
	}
}

func (c *Calls) Name(functionName string) *Calls {
	c.functionName = functionName
	return c
}

func (c *Calls) Add(params, result []interface{}) *Calls {

	c.inOut = append(c.inOut, CallParams{
		params: params,
		result: result,
	})

	return c
}

func (c *Calls) ToStr(userFunc SaveStringFunc) []StrCalls {

	calls := make([]StrCalls, len(c.inOut))
	for i, p := range c.inOut {
		calls[i] = StrCalls{
			Result:       inOutCode(p.result, userFunc),
			Params:       inOutCode(p.params, userFunc),
			FunctionName: c.functionName,
		}
	}

	return calls
}

func inOutCode(in []interface{}, userFunc SaveStringFunc) string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = ToString(v, userFunc)
	}
	return strings.Join(out, ", ")
}
