package testdata

import (
	"errors"

	"github.com/iostrovok/go-mockdata/testcode/otherpackage"
	oe "github.com/iostrovok/go-mockdata/testcode/otherpackage"
)

/*
	Simple package for fill mocker with data
*/

type CountGroups struct {
	Count   int32
	GroupBy int64
}

type LocalError otherpackage.Error

type IOne interface {
	FirstFunc(s string) (string, error)
	SecondFunc(s string) ([]string, error)
	ThirdFunc(s string) ([]string, LocalError)
	FourthFunc(s string) ([]string, otherpackage.Error)
	FifthFunc(s string) ([]string, oe.Error)
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

func (m *One) ThirdFunc(s string) ([]string, LocalError) {
	return []string{s}, LocalError(errors.New("one-more-errorx"))
}

func (m *One) FourthFunc(s string) ([]string, otherpackage.Error) {
	return []string{s}, otherpackage.Error(errors.New("one-more-error"))
}

func (m *One) FifthFunc(s string) ([]string, oe.Error) {
	return []string{s}, otherpackage.Error(errors.New("one-more-error"))
}

type Two struct {
}

func NewTwo() Two {
	return Two{}
}

func (m Two) Fisrt(s string) (map[string]oe.Error, otherpackage.Error) {
	return nil, otherpackage.Error(errors.New(s))
}

func (m Two) Second(s string) ([]otherpackage.Error, *otherpackage.Error) {
	return nil, nil
}

func (m Two) Third(s string) (*map[string]*otherpackage.Error, *oe.Error) {
	return nil, nil
}
