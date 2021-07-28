package utils

import (
	"strings"
	"testing"

	"github.com/centrify/platform-go-sdk/testutils"
	"github.com/centrify/platform-go-sdk/utils"
	"github.com/stretchr/testify/suite"
)

type TestUtilsSuite struct {
	testutils.CfyTestSuite
}

func (s *TestUtilsSuite) TestCClientVersionReqCurrentVersion() {
	t := s.T()

	t.Log("Test Centrify Client version requirements when requirement is for current version")

	s.RequiresCClientInstalled()
	expected, err := utils.GetCClientVersion()
	s.Assert().NoError(err, "Should not get error in getting current version requirement")

	// strip off build number
	parts := strings.Split(expected, "-")
	res, err := utils.VerifyCClientVersionReq(parts[0])
	s.Assert().NoError(err, "Should not get error in checking version requirement")
	s.Assert().Truef(res, "version requirement checking failed. Req version: %s. Expect: true, got %v", parts[0], res)
}

func (s *TestUtilsSuite) TestCClientVersionReqirements() {
	t := s.T()

	t.Log("Test Centrify Client version requirements")

	s.RequiresCClientInstalled()

	type testcase struct {
		reqVersion     string
		expectedResult bool
	}
	var testcases = []testcase{
		{"1.1", true},
		{"99.99", false}, // large version number
		{"21.1", true},
		{"21.4", true},
	}

	for _, tc := range testcases {
		res, err := utils.VerifyCClientVersionReq(tc.reqVersion)
		s.Assert().NoError(err, "Should not get error in checking version requirement")
		s.Assert().Equalf(tc.expectedResult, res, "version requirement checking failed. Req version: %s. Expect: %v, got %v", tc.reqVersion, tc.expectedResult, res)
	}
}

func (s *TestUtilsSuite) TestCClientVersionReqNotInstalled() {
	t := s.T()

	t.Log("TestGetCClient version when it is not installed")

	installed, err := utils.IsCClientInstalled()
	s.Assert().NoError(err, "Should not get error when checking if Centrify Client is installed")
	if installed {
		t.Skip("Skip test as Centrify Client is installed")
		return
	}

	_, err = utils.VerifyCClientVersionReq("21.1")
	s.Assert().ErrorIs(err, utils.ErrClientNotInstalled, "Should get error in checking version requirement")

}

func (s *TestUtilsSuite) TestGetCClientVersionWhenNotInstalled() {
	t := s.T()

	t.Log("TestGetCClient version when it is not installed")

	installed, err := utils.IsCClientInstalled()
	s.Assert().NoError(err, "Should not get error when checking if Centrify Client is installed")
	if installed {
		t.Skip("Skip test as Centrify Client is installed")
		return
	}

	_, err = utils.GetCClientVersion()
	s.Assert().ErrorIs(err, utils.ErrClientNotInstalled, "Should get error ErrClientNotInstalled in getting version")
}

func (s *TestUtilsSuite) TestGetCClientVersion() {
	t := s.T()

	t.Log("TestGetCClient version")
	s.RequiresCClientInstalled()

	ver, err := utils.GetCClientVersion()
	s.Assert().NoError(err, "Should not get error in getting version")
	s.Assert().NotEmpty(ver, "Version string should not be empty.")
	t.Logf("Version: %s\n", ver)
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(TestUtilsSuite))
}
