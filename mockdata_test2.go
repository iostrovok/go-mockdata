package mockdata

import (
	"fmt"

	. "github.com/iostrovok/check"

	tc "github.com/iostrovok/go-mockdata/testcode"
	"github.com/iostrovok/go-mockdata/testcode/mmock"
)

type testSuite2 struct{}

var _ = Suite(&testSuite2{})

const checkValue2 = `
package mockdata

import (
	"github.com/golang/mock/gomock"
	tc "github.com/iostrovok/go-mockdata/testcode"
	"github.com/iostrovok/go-mockdata/testcode/mmock"
	"testing"
)

// NewMockIOne returns new mocker with data for *mmock.MockIOne
func NewMockIOne (t *testing.T) *mmock.MockIOne {
	ctrl := gomock.NewController(t)
	m := mmock.NewMockIOne(ctrl)
	m = MockIOneThirdFunc(m)
	m = MockIOneFourthFunc(m)
	return m
}

func MockIOneThirdFunc(m *mmock.MockIOne) *mmock.MockIOne {
	m.EXPECT().ThirdFunc("To be or not to be").Return([]string{"To be or not to be"}, &tc.LocalError(errors.errorString{"one-more-error"}))
	return m
}

func MockIOneFourthFunc(m *mmock.MockIOne) *mmock.MockIOne {
	m.EXPECT().FourthFunc("To be or not to be").Return([]string{"To be or not to be"}, &otherpackage.Error(errors.errorString{"one-more-error"}))
	return m
}

`

func (s *testSuite2) TestFunctions2(c *C) {

	//c.Skip("not now")

	// Object which provides data for results
	one := tc.New()

	// NewMockIOne is a function return mock
	//
	m := New().SetMMock(mmock.NewMockIOne)
	defer m.Clean()

	str := "To be or not to be"

	// store ThirdFunc parameters and result - "func NewMockIOne(ctrl *gomock.Controller) *MockIOne"
	res1, err := one.ThirdFunc(str)
	m.StartFunction(one.ThirdFunc).Add([]interface{}{str}, []interface{}{res1, err})

	// store ThirdFunc parameters and result - "func NewMockIOne(ctrl *gomock.Controller) *MockIOne"
	res2, err := one.FourthFunc(str)
	m.StartFunction(one.FourthFunc).Add([]interface{}{str}, []interface{}{res2, err})

	// store ThirdFunc parameters and result - "func NewMockIOne(ctrl *gomock.Controller) *MockIOne"
	res3, err := one.FifthFunc(str)
	m.StartFunction(one.FifthFunc).Add([]interface{}{str}, []interface{}{res3, err})

	// out GO-code. We may just to save to file with m.Save(fileName)
	code := m.Code()

	fmt.Printf("\n\ncode: \n%s\n\n\n", code)

	// close current cycle of saving for MockIOne
	m.Clean()

	c.Assert(code, Equals, checkValue2)
}
