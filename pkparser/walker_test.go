package pkparser

import (
	. "github.com/iostrovok/check"
)

type testWalker struct{}

var _ = Suite(&testWalker{})

var WalkerTestSkip = true

func skip(c *C) {
	if WalkerTestSkip {
		c.Skip("no now")
	}
}

func (s *testWalker) TestWalker_Empty1(c *C) {
	skip(c)

	pkp := NewParser(testHome)
	c.Assert(pkp.parsePackage(testHome+"/samples/2/2b/"), IsNil)

	w := NewWalker(testHome).AddFiles(pkp.parsedFiles)

	res := w.FindFile("", "a")
	c.Assert(res.Find, Equals, false)
	c.Assert(res.ParsedFile, IsNil)
	c.Assert(res.Node, IsNil)

	res = w.FindPath("", "a")
	c.Assert(res.Find, Equals, false)
	c.Assert(res.ParsedFile, IsNil)
	c.Assert(res.Node, IsNil)
}

func (s *testWalker) TestWalker_1(c *C) {
	skip(c)

	pkp := NewParser(testHome)
	c.Assert(pkp.parsePackage(testHome+"/samples/2/2b/"), IsNil)

	w := NewWalker(testHome).AddFiles(pkp.parsedFiles)

	res := w.FindFile(testHome+"/samples/2/2b/b.go", "Check")
	c.Assert(res.Find, Equals, true)
	c.Assert(res.ParsedFile, NotNil)
	c.Assert(res.Node, NotNil)

	res = w.FindPath("github.com/iostrovok/go-mockdata/samples/2/2c", "MyError")
	c.Assert(res.Find, Equals, true)
	c.Assert(res.ParsedFile, NotNil)
	c.Assert(res.Node, NotNil)
}
