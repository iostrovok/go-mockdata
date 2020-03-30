package pkparser

import (
	"fmt"
	. "github.com/iostrovok/check"
	T "github.com/iostrovok/go-mockdata/samples/2/2b"
)

type testFunction struct{}

var _ = Suite(&testFunction{})

func (s *testFunction) TestFunction_1(c *C) {
	skip(c)

	f := T.NewRefIA()
	res, err := SplitFunction(f.Plus)
	c.Assert(err, IsNil)
	fmt.Printf("\nNewRefIA - Name: %s, Path: %s, FileSrc: %s, Receiver: %s\n\n", res.Name, res.Path, res.FileSrc, res.Receiver)

	c.Assert(res.Name, Equals, "Plus")
	c.Assert(res.Path, Equals, "github.com/iostrovok/go-mockdata/samples/2/2a")
	c.Assert(res.FileSrc, Equals, testHome+"/samples/2/2a/a.go")
	c.Assert(res.Receiver, Equals, "IA")

	nf := T.NewNoRefIA()
	res, err = SplitFunction(nf.Plus)
	c.Assert(err, IsNil)
	fmt.Printf("\nNewNoRefIA - Name: %s, Path: %s, FileSrc: %s, Receiver: %s\n\n", res.Name, res.Path, res.FileSrc, res.Receiver)
	c.Assert(res.Name, Equals, "Plus")
	c.Assert(res.Path, Equals, "github.com/iostrovok/go-mockdata/samples/2/2a")
	c.Assert(res.FileSrc, Equals, testHome+"/samples/2/2a/a.go")
	c.Assert(res.Receiver, Equals, "IA")

	b := T.NewRefB()
	res, err = SplitFunction(b.Plus)
	c.Assert(err, IsNil)
	fmt.Printf("\nNewRefB - Name: %s, Path: %s, FileSrc: %s, Receiver: %s\n\n", res.Name, res.Path, res.FileSrc, res.Receiver)
	c.Assert(res.Name, Equals, "Plus")
	c.Assert(res.Path, Equals, "github.com/iostrovok/go-mockdata/samples/2/2b")
	c.Assert(res.FileSrc, Equals, testHome+"/samples/2/2b/b.go")
	c.Assert(res.Receiver, Equals, "*B")

	cc := T.NewNoRefC()
	res, err = SplitFunction(cc.Plus)
	c.Assert(err, IsNil)
	fmt.Printf("\nNewNoRefC - Name: %s, Path: %s, FileSrc: %s, Receiver: %s\n\n", res.Name, res.Path, res.FileSrc, res.Receiver)
	c.Assert(res.Name, Equals, "Plus")
	c.Assert(res.Path, Equals, "github.com/iostrovok/go-mockdata/samples/2/2b")
	c.Assert(res.FileSrc, Equals, testHome+"/samples/2/2b/b.go")
	c.Assert(res.Receiver, Equals, "C")
}
