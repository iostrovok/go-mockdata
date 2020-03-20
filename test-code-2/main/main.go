package m

import (
	"github.com/iostrovok/go-mockdata/test-code-2/a"
	b "github.com/iostrovok/go-mockdata/test-code-2/a"
	"github.com/iostrovok/go-mockdata/test-code-2/inter"
	inter "github.com/iostrovok/go-mockdata/test-code-2/inter"
)

type MyTest struct {
	two   pkginter.INTER
	three inter.INTER
	one   a.One
}

func New() *MyTest {
	return &MyTest{
		two:   a.NewTwo(),
		three: a.NewThree(),
		one:   b.NewOne(),
	}
}

func (m *MyTest) GoGo(a.One) (*b.Three, pkginter.INTER) {
	return &b.Three{}, m.two
}
