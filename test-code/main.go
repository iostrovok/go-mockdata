package testdata

import (
	"errors"
	//"github.com/iostrovok/go-mockdata/test-code/otherpackage"
)

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
	//ThirdFunc(s string) ([]string, otherpackage.Error)
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

//
//func (m *One) ThirdFunc(s string) ([]string, otherpackage.Error) {
//	return []string{s}, otherpackage.Error(errors.New("one-more-error"))
//}
