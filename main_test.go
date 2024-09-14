package hoard

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(suiteTest))
}

type suiteTest struct {
	suite.Suite
}

// compile-time check whether the suiteTest implements SetupTest
var _ suite.SetupSubTest = (*suiteTest)(nil)

// compile-time check whether the suiteTest implements SetupTest
var _ suite.SetupTestSuite = (*suiteTest)(nil)

func (s *suiteTest) SetupTest() {
}

func (s *suiteTest) SetupSubTest() {
}
