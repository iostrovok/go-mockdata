package pkparser

import (
	"fmt"
	. "github.com/iostrovok/check"
	T "github.com/iostrovok/go-mockdata/samples/2/2b"
)

type testManager struct{}

var _ = Suite(&testManager{})

func (s *testManager) TestManager_Empty1(c *C) {
	//skip(c)

	m := New(testHome)

	f := T.NewRefIA()
	res, err := m.Find(f.Check)
	c.Assert(err, IsNil)

	fmt.Printf("===> TYPE: %T\n", res)

	c.Assert(1, IsNil)
}
