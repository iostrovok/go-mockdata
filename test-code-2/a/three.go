package a

import (
	inter "github.com/iostrovok/go-mockdata/test-code-2/inter"
)

type Three struct {
}

func NewThree() inter.INTER {
	return &Three{}
}

func (t *Three) Plus(a, b int) (c int) {
	return a + b
}

func (t *Three) Minus(a int, b int) int {
	return a - b
}

func (t *Three) Repeat(a int, b int) (c, d int) {
	return a, b
}
