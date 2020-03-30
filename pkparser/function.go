package pkparser

/*
	Simple walk through parsed package tree.
*/

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type SFunction struct {
	Receiver string
	FileSrc  string
	Path     string
	Name     string
}

func SplitFunction(i interface{}) (*SFunction, error) {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Func {
		return nil, fmt.Errorf("not function")
	}

	ptr := value.Pointer()
	top := runtime.FuncForPC(ptr)
	fmt.Printf("\nTOP.Name: %+v\n\n", top.Name())
	fmt.Printf("\nTOP.Entry: %+v\n\n", top.Entry())
	fileSrc, line := top.FileLine(ptr)
	fmt.Printf("\nTOP.fileSrc: %s, line: %d\n\n", fileSrc, line)

	dir, file := filepath.Split(top.Name())

	parts := strings.SplitN(file, ".", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("not object function")
	}

	path := filepath.Join(dir, parts[0])
	funcParts := strings.SplitN(parts[len(parts)-1], "-", 2)
	receiver := strings.Trim(parts[1], "()")

	return &SFunction{
		FileSrc: fileSrc,

		Receiver: receiver,
		Path:     path,
		Name:     funcParts[0],
	}, nil
}
