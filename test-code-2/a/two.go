package a

import (
	"github.com/iostrovok/go-mockdata/test-code-2/inter"
)

type Two struct {
}

func NewTwo() pkginter.INTER {
	return Two{}
}

func (t Two) Plus(a, b int) (c int) {
	return a + b
}

func (t Two) Minus(a int, b int) int {
	return a - b
}

func (t Two) Repeat(a int, b int) (c, d int) {
	return a, b
}
