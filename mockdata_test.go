package mockdata

import (
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

	// Object which provides data for results
	one := tc.New()

	// NewMockIOne is a function return mock
	//
	m := New().SetMMock(mmock.NewMockIOne)
	defer m.Clean()

	str := "To be or not to be"
	str2 := "William Shakespeare"

	// store FirstFunc parameters and result - "func NewMockIOne(ctrl *gomock.Controller) *MockIOne"
	res, err := one.FirstFunc(str)
	m.StartFunction(one.FirstFunc).InOut([]interface{}{str}, []interface{}{res, err})

	// store SecondFunc multi parameters and result
	m.StartFunction(one.SecondFunc)
	for _, s := range []string{str, str2} {
		res, err := one.SecondFunc(s)
		m.InOut([]interface{}{s}, []interface{}{res, err})
	}

	// out GO-code. We may just to save to file with m.Save(fileName)
	code := m.Code()

	// close current cycle of saving for MockIOne
	m.Clean()

	c.Assert(code, Equals, res1)
}
