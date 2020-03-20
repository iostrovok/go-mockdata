package imports

import (
	"testing"

	. "github.com/iostrovok/check"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestServiceDiscovery(t *testing.T) { TestingT(t) }

func (s *testSuite) TestFunctions1(c *C) {
	//c.Skip("no way")
	imp := New()
	imp.Add("/go-my_super_pkg/my_file1.go", "my_pkg1", "github.com/iostrovok/go-my_pkg", "")
	imp.Add("/go-my_super_pkg/my_file2.go", "my_pkg2", "github.com/iostrovok/go-my_pkg", "")
	imp.Add("/go-my_super_pkg/my_file3.go", "my_pkg1", "github.com/iostrovok/go-my_pkg", "alias1")
	imp.Add("/go-my_super_pkg/my_file4.go", "my_pkg1", "github.com/iostrovok/go-my_pkg", "alias2")

	err := imp.UsePath("github.com/iostrovok/go-my_pkg", "")
	c.Assert(err, IsNil)

	err = imp.UsePath("github.com/iostrovok/go-my_pkg", "alias1")
	c.Assert(err, IsNil)

	err = imp.UsePath("github.com/iostrovok/go-my_pkg", "alias2")
	c.Assert(err, IsNil)

	res := imp.List()
	c.Assert(res, DeepEquals, []string{
		`"github.com/iostrovok/go-my_pkg"`,
		`alias1 "github.com/iostrovok/go-my_pkg"`,
		`alias2 "github.com/iostrovok/go-my_pkg"`,
	})
}

func (s *testSuite) TestFunctions2(c *C) {
	//c.Skip("no way")
	imp := New()
	imp.Add("/go-my_pkg/my_file1.go", "my_pkg1", "github.com/iostrovok/go-my_pkg", "")
	imp.Add("/go-my_pkg/my_file2.go", "my_pkg2", "github.com/iostrovok/go-my_pkg", "")
	imp.Add("", "my_pkg1", "github.com/iostrovok/go-my_pkg", "")
	imp.Add("/go-my_pkg/my_file4.go", "my_pkg1", "github.com/iostrovok/go-my_pkg", "alias2")

	err := imp.UseFile("/go-my_pkg/my_file2.go", "")
	c.Assert(err, IsNil)

	err = imp.UsePath("github.com/iostrovok/go-my_pkg", "alias1")
	c.Assert(err, IsNil)

	err = imp.UseFile("/go-my_pkg/my_file4.go", "")
	c.Assert(err, IsNil)

	res := imp.List()
	c.Assert(res, DeepEquals, []string{
		`"github.com/iostrovok/go-my_pkg"`,
		`alias1 "github.com/iostrovok/go-my_pkg"`,
	})
}
