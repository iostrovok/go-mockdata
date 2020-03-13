package mockdata

import (
	"fmt"
	"testing"

	. "github.com/iostrovok/check"

	tc "github.com/iostrovok/go-mockdata/test-code"
	"github.com/iostrovok/go-mockdata/test-code/mmock"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestServiceDiscovery(t *testing.T) { TestingT(t) }

const res1 = `
package mockdata

import (
	"github.com/golang/mock/gomock"
	"github.com/iostrovok/go-mockdata/test-code"
	"github.com/iostrovok/go-mockdata/test-code/mmock"
	"testing"
)

// NewMockIOne returns new mocker with data for *mmock.MockIOne
func NewMockIOne (t *testing.T) *mmock.MockIOne {
	ctrl := gomock.NewController(t)
	m := mmock.NewMockIOne(ctrl)
	m = MockIOneFirstFunc(m)
	m = MockIOneSecondFunc(m)
	return m
}

func MockIOneFirstFunc(m *mmock.MockIOne) *mmock.MockIOne {
	m.EXPECT().FirstFunc("To be or not to be").Return("To be or not to be", &errors.errorString{"test error"})
	return m
}

func MockIOneSecondFunc(m *mmock.MockIOne) *mmock.MockIOne {
	m.EXPECT().SecondFunc("To be or not to be").Return([]string{"To be or not to be"}, nil)
	m.EXPECT().SecondFunc("William Shakespeare").Return([]string{"William Shakespeare"}, nil)
	return m
}

`

func (s *testSuite) TestFunctions(c *C) {

	one := tc.New()

	fmt.Printf("one: %+v\n", one)

	// function for mock
	m := New().SetMMock(mmock.NewMockIOne)

	str := "To be or not to be"
	str2 := "William Shakespeare"

	// save one function data
	m.StartFunction(one.FirstFunc)
	out, err := one.FirstFunc(str)
	m.InOut([]interface{}{str}, []interface{}{out, err})

	m.StartFunction(one.SecondFunc)
	for _, s := range []string{str, str2} {
		res, err := one.SecondFunc(s)
		m.InOut([]interface{}{s}, []interface{}{res, err})
	}

	fmt.Printf("\n\n" + m.Code() + "\n\n")

	c.Assert(m.Code(), Equals, res1)
}
