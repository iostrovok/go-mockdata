package b

import (
	"github.com/iostrovok/go-mockdata/samples/2/2a"
	t "github.com/iostrovok/go-mockdata/samples/2/2a"
)

type B struct{}

func NewRefIA() a.IA {
	return &B{}
}

type C struct{}

func NewNoRefIA() t.IA {
	return C{}
}

func NewRefB() *B {
	return &B{}
}

func NewNoRefC() C {
	return C{}
}
