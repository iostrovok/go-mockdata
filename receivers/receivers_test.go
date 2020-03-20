package receivers

import (
	"log"
	"os"
	"strings"
	"testing"

	. "github.com/iostrovok/check"

	importsAlias "github.com/iostrovok/go-mockdata/imports"
	main "github.com/iostrovok/go-mockdata/test-code-2/main"
	tc "github.com/iostrovok/go-mockdata/testcode"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestServiceDiscovery(t *testing.T) { TestingT(t) }

var testHome string

func init() {
	testHome := os.Getenv("TEST_SOURCE_PATH")
	log.Printf("TEST_SOURCE_PATH: %s", testHome)
	if testHome != "" {
		err := os.Chdir(testHome)
		if err != nil {
			panic(err)
		}
	}
}

func (s *testSuite) TestFunctionsOneFunction1(c *C) {
	c.Skip("no way")
	im := importsAlias.New()
	m := New(testHome, im)

	ofr := m.OneFunction(importsAlias.New().Add)
	c.Assert(ofr.Err, IsNil)

	c.Logf("TestFunctionsOneFunction1: %+v", ofr)

	checkSubStr := "go-mockdata/imports/inter/inter.go"
	c.Assert(strings.Index(ofr.FilePath, checkSubStr), Equals, len(ofr.FilePath)-len(checkSubStr))
	c.Assert(ofr.Pkg, DeepEquals, "github.com/iostrovok/go-mockdata/imports/inter")
	c.Assert(ofr.Recv, DeepEquals, "IImp")
	c.Assert(ofr.FuncName, DeepEquals, "Add")
}

func (s *testSuite) TestFunctionsOneFunction2(c *C) {
	c.Skip("no way")
	im := importsAlias.New()
	m := New(testHome, im)

	ofr := m.OneFunction(im.Add)

	c.Logf("TestFunctionsOneFunction2: %+v", ofr)

	checkSubStr := "go-mockdata/imports/inter/inter.go"
	c.Assert(strings.Index(ofr.FilePath, checkSubStr), Equals, len(ofr.FilePath)-len(checkSubStr))
	c.Assert(ofr.Pkg, DeepEquals, "github.com/iostrovok/go-mockdata/imports/inter")
	c.Assert(ofr.Recv, DeepEquals, "IImp")
	c.Assert(ofr.FuncName, DeepEquals, "Add")
}

func (s *testSuite) TestFunctionsOneFunction3(c *C) {
	c.Skip("no way")
	im := importsAlias.New()
	m := New(testHome, im)

	ofr := m.OneFunction(tc.NewTwo().Third)

	c.Logf("TestFunctionsOneFunction3: %+v", ofr)

	checkSubStr := "go-mockdata/testcode/main.go"
	c.Assert(strings.Index(ofr.FilePath, checkSubStr), Equals, len(ofr.FilePath)-len(checkSubStr))
	c.Assert(ofr.Pkg, DeepEquals, "github.com/iostrovok/go-mockdata/testcode")
	c.Assert(ofr.Recv, DeepEquals, "Two")
	c.Assert(ofr.FuncName, DeepEquals, "Third")
}

func (s *testSuite) TestFunctionsRun1(c *C) {
	c.Skip("no way")
	im := importsAlias.New()
	m := New(testHome, im)

	ofrTMP := m.OneFunction(im.Add)

	c.Logf("TestFunctionsRun1-ofrTMP: %+v", ofrTMP)

	ofr := m.Run(im.Add)

	c.Logf("TestFunctionsRun1: %+v", ofr)

	c.Assert(ofr.Err, IsNil)
	c.Assert(ofr.Params, DeepEquals, []string{"string", "string", "string", "string"})
	c.Assert(ofr.Results, DeepEquals, []string{})
	c.Assert(im.List(), DeepEquals, []string{})
}

func (s *testSuite) TestFunctionsRun2(c *C) {
	c.Skip("no way")
	im := importsAlias.New()
	m := New(testHome, im)

	obj := tc.New()

	ofr := m.Run(obj.FourthFunc)
	c.Assert(ofr.Err, IsNil)
	//c.Assert(im.List(), DeepEquals, []string{`"github.com/iostrovok/go-mockdata/testcode/otherpackage"`})
	c.Assert(ofr.Results, DeepEquals, []string{"[]string", "otherpackage.Error"})
	c.Assert(ofr.Params, DeepEquals, []string{"string"})
}

func (s *testSuite) TestFunctionsRun3(c *C) {
	//c.Skip("no way")
	im := importsAlias.New()
	m := New(testHome, im)

	obj := tc.NewTwo()

	ofr := m.Run(obj.Fisrt)
	c.Assert(ofr.Err, IsNil)

	c.Assert(ofr.Results, DeepEquals, []string{"map[string]oe.Error", "otherpackage.Error"})
	c.Assert(ofr.Params, DeepEquals, []string{"string"})
	c.Assert(im.List(), DeepEquals, []string{`"github.com/iostrovok/go-mockdata/testcode/otherpackage"`})
}

func (s *testSuite) TestFunctionsRun4(c *C) {
	c.Skip("no way")
	im := importsAlias.New()
	m := New(testHome, im)

	obj := tc.NewTwo()

	ofr := m.Run(obj.Second)
	c.Assert(ofr.Err, IsNil)
	c.Assert(im.List(), DeepEquals, []string{`"github.com/iostrovok/go-mockdata/testcode/otherpackage"`})
	c.Assert(ofr.Results, DeepEquals, []string{"[]otherpackage.Error", "*otherpackage.Error"})
	c.Assert(ofr.Params, DeepEquals, []string{"string"})
}

func (s *testSuite) TestFunctionsRun5(c *C) {
	c.Skip("no way")
	im := importsAlias.New()
	m := New(testHome, im)

	ofr := m.Run(tc.NewTwo().Third)
	c.Assert(ofr.Err, IsNil)
	//c.Assert(im.List(), DeepEquals, []string{`"github.com/iostrovok/go-mockdata/test-code/otherpackage"`})
	c.Assert(ofr.Results, DeepEquals, []string{"*map[string]*otherpackage.Error", "*oe.Error"})
	c.Assert(ofr.Params, DeepEquals, []string{"string"})
}

func (s *testSuite) TestFunctionsRun6(c *C) {
	c.Skip("no way")
	im := importsAlias.New()
	main := main.New()

	m := New(testHome, im)
	ofr := m.Run(main.GoGo)
	c.Assert(ofr.Err, IsNil)
	c.Assert(ofr.Params, DeepEquals, []string{"a.One"})
	c.Assert(ofr.Results, DeepEquals, []string{"*b.Three", "pkginter.INTER"})
	c.Assert(im.List(), DeepEquals, []string{
		`a "github.com/iostrovok/go-mockdata/test-code-2/a"`,
		`b "github.com/iostrovok/go-mockdata/test-code-2/a"`,
		`pkginter "github.com/iostrovok/go-mockdata/test-code-2/a"`,
	})
}
