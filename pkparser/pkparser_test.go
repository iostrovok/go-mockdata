package pkparser

import (
	"log"
	"os"
	"testing"

	. "github.com/iostrovok/check"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestServiceDiscovery(t *testing.T) { TestingT(t) }

var testHome string

func init() {
	testHome = os.Getenv("TEST_SOURCE_PATH")
	log.Printf("TEST_SOURCE_PATH: %s", testHome)
	if testHome != "" {
		err := os.Chdir(testHome)
		if err != nil {
			panic(err)
		}
	}
}

func checkFiles(c *C, pkp *PkParser, files ...string) {

	c.Assert(len(files), Equals, len(pkp.parsedFiles))

	for _, f := range files {
		c.Assert(pkp.parsedFiles[f], NotNil)
	}
}

func checkDirs(c *C, pkp *PkParser, dirs ...string) {

	c.Assert(len(dirs), Equals, len(pkp.parsedDir))

	for _, d := range dirs {
		c.Assert(pkp.parsedDir[d], Equals, true)
	}
}

func (s *testSuite) TestFunctions1(c *C) {
	skip(c)

	pkp := NewParser(testHome)
	c.Assert(pkp.parsePackage(testHome+"/samples/1/a1/"), IsNil)

	checkFiles(c, pkp, testHome+"/samples/1/a1/a1.go")
	checkDirs(c, pkp, testHome+"/samples/1/a1")

	file := pkp.parsedFiles[testHome+"/samples/1/a1/a1.go"]
	c.Assert(file.pkgName, Equals, "a1")
	c.Assert(file.file, NotNil)
}

func (s *testSuite) TestFunctions2(c *C) {
	skip(c)

	pkp := NewParser(testHome)
	c.Assert(pkp.parsePackage(testHome+"/samples/2/2b/"), IsNil)

	checkFiles(c, pkp,
		testHome+"/samples/2/2b/s.go",
		testHome+"/samples/2/2b/b.go",
		testHome+"/samples/2/2a/a.go",
		testHome+"/samples/2/2c/c.go",
	)

	checkDirs(c, pkp,
		testHome+"/samples/2/2b",
		testHome+"/samples/2/2a",
		testHome+"/samples/2/2c",
	)

	file := pkp.parsedFiles[testHome+"/samples/2/2b/b.go"]
	c.Assert(file.pkgName, Equals, "b")
	c.Assert(file.file, NotNil)
	c.Assert(file.file.Name.String(), Equals, "b")
}
