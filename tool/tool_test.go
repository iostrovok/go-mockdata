package tool

import (
	"context"
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

func (s *testSuite) TestFunctions1(c *C) {

	pkp := NewTool()
	c.Assert(pkp.Run(context.Background(), testHome+"/samples/1/a1/"), IsNil)
	//
	//checkFiles(c, pkp, testHome+"/samples/1/a1/a1.go")
	//checkDirs(c, pkp, testHome+"/samples/1/a1")
	//
	//file := pkp.parsedFiles[testHome+"/samples/1/a1/a1.go"]
	//c.Assert(file.pkgName, Equals, "a1")
	//c.Assert(file.file, NotNil)
}
