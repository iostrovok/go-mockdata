package onefunction

import (
	"fmt"
	"testing"

	. "github.com/iostrovok/check"

	tc "github.com/iostrovok/go-mockdata/test-code"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestServiceDiscovery(t *testing.T) { TestingT(t) }

const res1 = `
func MockIVerificationCountByStatus(m *mmock.MockIVerification) *mmock.MockIVerification {
	m.EXPECT().CountByStatus().Return([]*testdata.CountGroups{&testdata.CountGroups{9, 0}, &testdata.CountGroups{10, 1}})
	return m
}
`

const res2 = `
func MockIVerificationCountByStatus(m *mmock.MockIVerification) *mmock.MockIVerification {
	m.EXPECT().CountByStatus(1).Return([]*testdata.CountGroups{&testdata.CountGroups{9, 0}, &testdata.CountGroups{10, 1}})
	m.EXPECT().CountByStatus(2).Return([]*testdata.CountGroups{&testdata.CountGroups{9, 0}, &testdata.CountGroups{10, 1}}, nil)
	m.EXPECT().CountByStatus(3).Return(nil, []*testdata.CountGroups{&testdata.CountGroups{9, 0}, &testdata.CountGroups{10, 1}})
	m.EXPECT().CountByStatus(4).Return(nil, nil)
	return m
}
`

func (s *testSuite) TestToEscape(c *C) {

	in := `"Fran & Freddie's Diner ☺
	☺"`
	expected := `"\"Fran & Freddie's Diner ☺\n\t☺\""`

	c.Assert(Escape(in), DeepEquals, expected)
}

func (s *testSuite) TestToString1(c *C) {

	out := []*tc.CountGroups{
		{
			Count:   9,
			GroupBy: 0,
		},
		{
			Count:   10,
			GroupBy: 1,
		},
	}
	check := `[]*testdata.CountGroups{&testdata.CountGroups{9, 0}, &testdata.CountGroups{10, 1}}`

	str := ToString(out)

	c.Assert(str, Equals, check)
}

func (s *testSuite) TestToString2(c *C) {

	out := map[string]*tc.CountGroups{
		"one": {
			Count:   9,
			GroupBy: 0,
		},
		"two": {
			Count:   10,
			GroupBy: 1,
		},
	}
	check := `map[string]*testdata.CountGroups{"two":&testdata.CountGroups{10, 1}, "one":&testdata.CountGroups{9, 0}}`
	check2 := `map[string]*testdata.CountGroups{"one":&testdata.CountGroups{9, 0}, "two":&testdata.CountGroups{10, 1}}`

	str := ToString(out)

	fmt.Printf("str: %s\n", str)

	c.Assert(str == check || str == check2, Equals, true)
}

func (s *testSuite) TestToString3(c *C) {

	out := map[int][]bool{
		0: {true, false},
		2: {true, false, false},
	}
	check := `map[int][]bool{0:[]bool{true, false}, 2:[]bool{true, false, false}}`
	check2 := `map[int][]bool{2:[]bool{true, false, false}, 0:[]bool{true, false}}`

	str := ToString(out)

	c.Assert(str == check || str == check2, Equals, true)
}

func (s *testSuite) TestFullFunctionName(c *C) {
	w := New().FunctionName("CountByStatus").MockType("*mmock.MockIVerification")
	c.Assert(w.FullFunctionName(), Equals, "MockIVerificationCountByStatus")
}

func (s *testSuite) TestFunctions_1(c *C) {

	out := []*tc.CountGroups{
		{
			Count:   9,
			GroupBy: 0,
		},
		{
			Count:   10,
			GroupBy: 1,
		},
	}

	w := New().FunctionName("CountByStatus").MockType("*mmock.MockIVerification").
		Add([]interface{}{}, []interface{}{out})
	c.Assert(w.Code(), Equals, res1)
}

func (s *testSuite) TestFunctions_2(c *C) {

	out := []*tc.CountGroups{
		{
			Count:   9,
			GroupBy: 0,
		},
		{
			Count:   10,
			GroupBy: 1,
		},
	}

	w := New().FunctionName("CountByStatus").MockType("*mmock.MockIVerification")

	w.Add([]interface{}{1}, []interface{}{out})
	w.Add([]interface{}{2}, []interface{}{out, nil})
	w.Add([]interface{}{3}, []interface{}{nil, out})
	w.Add([]interface{}{4}, []interface{}{nil, nil})

	c.Assert(w.Code(), Equals, res2)
}
