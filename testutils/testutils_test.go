package testutils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestUtilsTestSuite struct {
	CfyTestSuite
}

func (s *TestUtilsTestSuite) SetupSuite() {
	t := s.T()
	t.Log("Call SetupTestSuite")
}

func (s *TestUtilsTestSuite) TestNightly() {
	t := s.T()

	t.Log("TestNightly should be run on the nightly mark")
	s.ExecuteOnMarks("nightly")
	t.Log("TestNightly is run")
}

func (s *TestUtilsTestSuite) TestIntegration() {
	t := s.T()

	t.Log("TestIntegration should be run if either integration or nightly mark is specified.")
	s.ExecuteOnMarks("nightly", "integration")
	t.Log("TestIntegration is run")
}

func (s *TestUtilsTestSuite) TestCClientRunning() {
	t := s.T()

	// this test is done during build.  So if it is invoked as part of build (i.e., build is in marks)
	// the test should be skipped if the system has no Centrify Client running or there is no active tenant.
	//
	t.Log("TestCClientRunning requires CClient to be running")
	s.RequiresCClientRunning()

	t.Log("TestCClient is running ")
}

func (s *TestUtilsTestSuite) TestRunOnBuildMark() {
	t := s.T()

	t.Log("TestRunOnBuildMark - test should be run for build mark")
	s.ExecuteOnMarks("build")

	t.Log("Build test is run")

}

func (s *TestUtilsTestSuite) TestRequireActiveTenant() {
	t := s.T()

	t.Log("TestRequireActiveTenant requires system to be connected to running tenant")
	s.RequiresActiveTenant()

	t.Log("TestRequireActiveTenant is running since system is connected to tenant")
}

func (s *TestUtilsTestSuite) TestRequirePrivUser() {
	t := s.T()

	t.Log("TestRequirePrivUser requires the runner to be privileged")
	s.RequiresRunAsPrivilegedUser()

	t.Log("TestRequirePrivUser runs as it is run by privileged user")
}

func (s *TestUtilsTestSuite) TestRequireUnprivUser() {
	t := s.T()

	t.Log("TestRequirUnprivUser requires the runner to be unprivileged")
	s.RequiresRunAsUnprivilegedUser()

	t.Log("TestRequireUnprivUser runs as it is run by unprivileged user")
}
func TestTestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(TestUtilsTestSuite))
}
