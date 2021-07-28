package testutils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestUtilsMoreTestSuite struct {
	CfyTestSuite
	Config *TestConfig // configuration
}

func (s *TestUtilsMoreTestSuite) SetupSuite() {
	t := s.T()
	t.Log("Call SetupTestSuite for TestUtilsMoreTestSuite")
}

func (s *TestUtilsMoreTestSuite) TestStress() {
	t := s.T()

	t.Log("TestStressMark is run when stress is specified in mark")
	s.ExecuteOnMarks("stress")
	t.Log("TestStress is run")
}

func (s *TestUtilsMoreTestSuite) TestUnit() {
	t := s.T()

	t.Log("TestUnit is always run")
}

func (s *TestUtilsMoreTestSuite) TestRequireEnrollment() {
	t := s.T()

	t.Log("This test requires system to be enrolled")
	s.RequiresEnrolledTenant()

	t.Log("TestRequireEnrollment runs as system is enrolled")
}

func (s *TestUtilsMoreTestSuite) TestRequireInstallClient() {
	t := s.T()

	t.Log("This test requires Centrify Client to be installed")
	s.RequiresCClientInstalled()

	t.Log("TestRequireInstallClient runs as Centrify Client is installed")
}
func TestTestUtilsMoreTestSuite(t *testing.T) {
	suite.Run(t, new(TestUtilsMoreTestSuite))
}
