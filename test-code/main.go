package testdata

import "errors"

/*
	Simple package for fill mocker with data
*/

type CountGroups struct {
	Count   int32
	GroupBy int64
}

type IOne interface {
	FirstFunc(s string) (string, error)
	SecondFunc(s string) ([]string, error)
}

type One struct {
}

func New() IOne {
	return &One{}
}

func (m *One) FirstFunc(s string) (string, error) {
	return s, errors.New("test error")
}

func (m *One) SecondFunc(s string) ([]string, error) {
	return []string{s}, nil
}
