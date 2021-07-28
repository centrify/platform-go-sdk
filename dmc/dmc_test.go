package dmc

import (
	"testing"

	"github.com/centrify/platform-go-sdk/testutils"
	"github.com/centrify/platform-go-sdk/utils"
	"github.com/stretchr/testify/suite"
)

type DMCTestSuite struct {
	testutils.CfyTestSuite
}

func TestDMCTestSuite(t *testing.T) {
	suite.Run(t, new(DMCTestSuite))
}

func (s *DMCTestSuite) SetupSuite() {

}

func (s *DMCTestSuite) TestRunByUnprivilegedUser() {
	t := s.T()

	t.Log("Test request DMC by unprivileged user")
	s.RequiresRunAsUnprivilegedUser()
	s.RequiresActiveTenant()

	_, err := GetDMCToken("scope")
	s.Assert().Error(err, "Expect to have error when run by unprivileged user")
	s.Assert().ErrorIs(err, utils.ErrCannotGetToken, "Expect ErrCannotGetToken error")
}
func (s *DMCTestSuite) TestCClientNotInstalled() {
	t := s.T()

	t.Log("Test Centrify Client not installed")
	installed, err := utils.IsCClientInstalled()
	s.Assert().NoError(err, "Should not have error when checking client is installed")

	if installed {
		t.Skip("Test skipped as Centrify Client is installed")
	}
	_, err = GetDMCToken("scope")
	if !installed {
		s.Assert().ErrorIs(err, utils.ErrClientNotInstalled, "Expect ErrClientNotInstalled error.  Got error %v", err)
	}
}

func (s *DMCTestSuite) TestNotEnrolled() {
	t := s.T()

	t.Log("Test Centrify Client not enrolled")
	s.RequiresRunAsPrivilegedUser()
	s.RequiresCClientInstalled()

	_, err := utils.GetEnrolledTenant()
	if err == nil {
		t.Skip("Test skipped as system is enrolled")
	}
	s.Assert().ErrorIs(err, utils.ErrNotEnrolled, "Unexpected error in getting enrolled tenant information")

	// now do the test
	_, err = GetDMCToken("scope")
	s.Assert().Error(err, "Expect error when not rolled")
	s.Assert().ErrorIs(err, utils.ErrCannotSetupConnection, "Expect ErrCannotSetupConnection error when not enrolled. Got error %v", err)
}

func (s *DMCTestSuite) TestClientStop() {
	t := s.T()

	t.Log("Test Centrify Client stopped")
	s.RequiresRunAsPrivilegedUser()
	s.RequiresCClientInstalled()

	status, err := utils.GetCClientStatus()
	s.Assert().NoError(err, "Should not get error when getting Centrify Client status")

	if status != utils.AgentStatusStopped {
		// unexpected client status
		t.Skipf("Skip test as client status [%v] is not \"stopped\"", status)
	}

	// now do the test
	_, err = GetDMCToken("scope")
	s.Assert().Error(err, "Expect error when Centrify Client is stopped")
	s.Assert().ErrorIs(err, utils.ErrCannotSetupConnection, "Expect ErrCannotSetupConnection error when not enrolled. Got error %v", err)
}

func (s *DMCTestSuite) TestDisconnect() {
	t := s.T()

	t.Log("Test Centrify Client disconnected")
	s.RequiresRunAsPrivilegedUser()
	s.RequiresCClientInstalled()

	status, err := utils.GetCClientStatus()
	s.Assert().NoError(err, "Should not get error when getting Centrify Client status")

	if status != utils.AgentStatusDisconnected {
		// unexpected client status
		t.Skipf("Skip test as client status [%v] is not \"disconnected\"", status)
	}

	// now do the test
	_, err = GetDMCToken("scope")
	s.Assert().Error(err, "Expect error when Centrify Client is stopped")
	s.Assert().ErrorIs(err, utils.ErrCannotGetToken, "Expect ErrCannotGetToken error disconnected. Got error %v", err)
}

func (s *DMCTestSuite) TestScopeNotDefined() {

	t := s.T()

	t.Log("Test scope not defined")
	s.RequiresRunAsPrivilegedUser()
	s.RequiresActiveTenant()

	// now do the test
	_, err := GetDMCToken("undefined_scope")
	s.Assert().Error(err, "Expect error when scope is not defined")
	s.Assert().ErrorIs(err, utils.ErrCannotGetToken, "Expect ErrCannotGetToken error disconnected. Got error %v", err)
}

func (s *DMCTestSuite) TestGetGoodToken() {
	t := s.T()

	t.Log("Test good case of getting the token")
	s.RequiresRunAsPrivilegedUser()
	s.RequiresActiveTenant()

	// now do the test
	token, err := GetDMCToken("testsdk")
	if err != nil {
		s.Assert().ErrorIs(err, utils.ErrCannotGetToken, "Should get ErrCannotGetToken for invalid scope")
		reason := err.Error()
		s.Assert().Contains(reason, "Unknown scope", "reason should be unknown scope")
		s.Assert().FailNow("You need to enroll the system that enables DMC and specifies testsdk as a scope")
	}

	// verify if token has valid accessToken
	s.Assert().NotEmpty(token, "Access token should not be empty")
}

func (s *DMCTestSuite) TestGetReservedScope() {

	s.RequiresRunAsPrivilegedUser()
	s.RequiresActiveTenant()

	reserved := []string{
		"cagent",
		"zso",
		"mobile",
		"__centrify_vault",
	}

	for _, scope := range reserved {
		_, err := GetDMCToken(scope)
		s.Assert().ErrorIs(err, utils.ErrCannotGetToken, "Expect error when getting scope [%s]", scope)
		reason := err.Error()
		s.Assert().Contains(reason, "reserved", "reason should be about reserved scopes")
	}
}
func (s *DMCTestSuite) TestGetCEnrollmentInfo() {
	t := s.T()

	t.Log("Test getting Centrify Client enrollment information")
	s.RequiresCClientRunning()

	tenantURL, clientID, err := GetEnrollmentInfo()
	s.Assert().NoError(err, "Should not get error in getting Centrify Client information")
	s.Assert().NotEmpty(tenantURL, "Tenant URL should not be empty")
	s.Assert().NotEmpty(clientID, "Client ID should not be empty")
	t.Logf("tenant: [%s]   Client ID: [%s]", tenantURL, clientID)
}
