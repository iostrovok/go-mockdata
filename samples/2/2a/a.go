package a

import (
	"github.com/iostrovok/go-mockdata/samples/2/2c"
	e "github.com/iostrovok/go-mockdata/samples/2/2c"
	t "github.com/iostrovok/go-mockdata/samples/2/2c"
)

type IA interface {
	Plus(int, int) int
	Minus(a, b int) (c int)
	Check(a e.MyError, i t.MyError) (c c.MyError)
}
