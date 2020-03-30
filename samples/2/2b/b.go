package b

import (
	c "github.com/iostrovok/go-mockdata/samples/2/2c"
	t "github.com/iostrovok/go-mockdata/samples/2/2c"
)

func (b *B) Plus(a1, a2 int) int {
	return a1 + a2
}

func (b *B) Minus(a, d int) (c int) {
	c = a - d
	return
}

func (b *B) Check(a, r t.MyError) c.MyError {
	return nil
}

func (b C) Plus(a1, a2 int) int {
	return a1 + a2
}

func (b C) Minus(a, d int) (c int) {
	c = a - d
	return
}

func (b C) Check(a, d c.MyError) (c t.MyError) {
	return nil
}
