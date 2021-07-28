/*
Package testutils contains all support functions that may be used in Centrify "go tests". The
test framework is based on the packages provided in testing and github.com/stretchr/testify.

Test suite organization

Developers should try to organize tests into test suites.  Each test suite should be named something like <suite_feature>Suite and there
should be a function name Test<suite_feature>Suite in the test program.   This allows anyone to select individual
test suites to run by specifying "go test -run XYZ" which will run all the tests in all test suites that includes XYZ as part of the test suite name.

Requirements for test files

 1. The test file must import github.com/centrify/testutils
 2. Each test suite must has testutils.CfyTestSuite as an anonymous member.

Test Configuration file

Each test suite should read a test configuration file  (in JSON) in its own implementation of SetupTestSuite() method.   The path of the configuration file
is specified by "-config" parameter.

Current supported parameters in the configuration parameters are:

 TenantURL:  specify that the suite requires the system be enrooled to this tenant URL (e.g., devdog.centrify.com).
			 If this is blank or not specified, tests will use the tenant that the current system is enrolled.
			 If the system is enrolled to a tenant that is different from this URL setting, all tests that calls
			 RequiresActiveTenant() or RequiresEnrolledTenant() will fail.
 Marks:      A json string array that specifies the matching test marks.  A test can check if a mark exists in
             the configuration file to determine whether the test should be skipped or not.

Supported test marks:
	scalability:   indicates running scalability tests
	integration:   indicates running integration tests.  Requires connectivity to tenant
	build:         indicates the test is running as part of build.  Should not require connectivity to tenant, or a running Centrify Client

If no configuration file is specified:
	- "build" is used as the only test mark
	- if any test requires the system to be enrolled to a tenant and/or the tenant is active, there is no restriction on what that tenant is.

Individual test function

Each test function must do the followings:

  1. All test function names must start with Test.  Otherwise the test WILL NOT be picked up by "go test".
  2. If you want the test to execute only when any of the marks is specified in the configuration file, the test needs to call ExecuteOnMarks().
     For example,  ExecuteOnMarks("nightly", "integration") means that the test is run if any of "nightly"/"integration" is specified in "Marks"
     in the test configuration file.
  3. If your test requires the test user to be privileged user (root in Linux), call RequiresRunAsPrivilegedUser().  The test will be
     marked as skipped if it is not run by root in Linux.
  4. Call one of the following RequiresXXX() function if you test requires certain environment setup:
	 - RequiresActiveTenant(): The test system must be connected to an active tenant.
	 - RequiresCClientRunning(): Centrify Client must be running in the system.  It may or may not be connected to the tenant.
	 - RequiresEnrolledTenant(): The test system is enrolled to the tenant.  Centrify Client may or may not be running.
	      Also, the system may/may not be connected to the tenant.
	 - RequiresCClientInstalled(): Centrify Client must be installed in the system.  It may or not may be enrolled.
  5. If RequiresActiveTenant()/RequiresCClientRunning()/RequiresEnrolledTenant() is called, AND TenantURL is specified in the configuration
     file, AND the system is enrolled to a different tenant, the test is not run and marked as FAILED.

You can review the test files testutils_test.go and test2_test.go to review how these tests are used.

How to run tests

You can use the two test files in this directory as example.

 go test ./... -args -config=/tmp/config.json
	This will run all tests using the configuration file /tmp/config.json

 go test ./... -run=MoreTest -args -config=/tmp/config.json
	This will run all the tests in the test suite TestUtilsMoreTestSuite but not TestUtilsTestSuite

 go test ./... -run=UtilsTest -args -config=/tmp/config.json
	This will run all the tests in the test suite TestUtilsTestSuite but not TestUtilsMoreTestSuite

If you change the system configuration (e.g., starts Centrify Client, enroll the system etc), you can force "go test" to not
used previous test result by running the command "go clean -testcache"
*/
package testutils

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os/exec"
	"testing"

	"github.com/centrify/platform-go-sdk/utils"
	"github.com/stretchr/testify/suite"
)

type UserCredentials struct {
	Username string
	Password string
}

// TestConfig stores the test configuration information
type TestConfig struct {
	TenantURL        string   // tenant URL
	Marks            []string // test marks
	ClientID         string
	ClientSecret     string
	AppID            string
	Scope            string
	PASuser          UserCredentials
	PolicyChangeUser UserCredentials // For the test where we are changing policies for the role
	HTTPProxyURL     string
}

func defConfig() *TestConfig {
	return &TestConfig{
		TenantURL: "",
		Marks:     []string{"build"},
	}
}

// loadConfiguration reads the file specified in 'path' and returns the configuration specified in the file.
// An empty object is returned for any issue (e.g., file not found, invalid format etc)
func loadConfiguration(t *testing.T, path string, configString string) *TestConfig {

	c := defConfig()
	if path == "" && configString == "" {
		return c
	}

	if path != "" && configString != "" {
		t.Fatalf("Config file path and config string provided in parameters")
	}

	if path != "" {
		filebuf, err := ioutil.ReadFile(path)
		if err != nil {
			return c // error in read
		}

		err = json.Unmarshal(filebuf, c)
		if err != nil {
			return defConfig()
		}
		return c
	} else if configString != "" {
		err := json.Unmarshal([]byte(configString), c)
		if err != nil {
			return defConfig()
		}
		return c
	}
	return c
}

// CfyTestSuite is the base struct for all Centrify go tests
type CfyTestSuite struct {
	suite.Suite             // testify suite
	Config      *TestConfig // configuration

	// internal information
	enrolledTenant  string // enrolled tenant information
	agentInstalled  bool   // flag to show if Centrify Client is installed
	agentEnrolled   bool   // flag to show if Centrify Client is enrolled to tenant
	vaultRunning    bool   // flag to show if HashiCorp Vault is running
	agentStatus     string // agent status
	platformVersion string // get platform version (requires run as privileged user)
	setupDone       bool   // flag to show if configuration file has been read
}

// Note:  there is no need to call flag.Parse().  "go test" already calls it
// but this needs to be declared here first.
var (
	configPtr      = flag.String("config", "", "configuration file")
	configString   = flag.String("config-string", "", "configuration string")
	VaultRootToken = flag.String("vault-root-token", "root", "Vault root token")
)

// LoadConfig is called by all public methods to ensure that the configuration parameters
// are loaded
func (s *CfyTestSuite) LoadConfig() {
	if !s.setupDone {
		s.T().Logf("Loading test configuration from %s\n", *configPtr)
		s.Config = loadConfiguration(s.T(), *configPtr, *configString)
		s.reloadAgentStatus()
		s.checkVault()
		s.setupDone = true
	}
}

// reloadAgentStatus updates agent status information
func (s *CfyTestSuite) reloadAgentStatus() {
	isInstalled, err := utils.IsCClientInstalled()
	if err != nil {
		s.T().Fatalf("Cannot get Centrify Client Status: %v\n", err)
	}
	s.agentInstalled = isInstalled

	// reset all other status
	s.platformVersion = ""
	s.agentStatus = utils.AgentStatusUnknown

	if !isInstalled {
		s.enrolledTenant = ""
		s.agentEnrolled = false
		return
	}

	var tenant string
	status, err := utils.GetCClientStatus()
	s.T().Logf("Agent status: %s\n", status)
	if err != nil {
		s.T().Fatalf("Cannot get agent status: %v\n", err)
	} else {
		s.agentStatus = status
		switch status {
		case utils.AgentStatusConnected,
			utils.AgentStatusDisabled,
			utils.AgentStatusDisconnected:
			tenant, err = utils.GetEnrolledTenant()
			if err != nil {
				s.T().Fatalf("Cannot get tenant URL: %v\n", err)
			}
			s.agentEnrolled = true
			s.enrolledTenant = tenant

		case utils.AgentStatusStarting:
			// just fails test and tell user to try again later
			s.T().Fatalf("Centrify client starting, try again later")

		case utils.AgentStatusStopped:
			tenant, err = utils.GetEnrolledTenant()
			if err == nil {
				s.enrolledTenant = tenant
				s.agentEnrolled = true
			} else if err == utils.ErrNotEnrolled {
				s.agentEnrolled = false
				s.enrolledTenant = ""
			} else {
				s.T().Fatalf("Cannot get tenant URL: %v\n", err)
			}

		case utils.AgentStatusUnknown:
			s.T().Fatalf("Cannot get Centrify Client status")
		}
	}

	if s.agentEnrolled {
		s.T().Logf("Enrolled to [%s]\n", s.enrolledTenant)
	}
}

// GetTenantURL returns the tenant URL stored in the configuration file
func (s *CfyTestSuite) GetTenantURL() string {
	s.LoadConfig()
	return s.Config.TenantURL
}

// HasMark checks if a mark exists in the test configuration
func (s *CfyTestSuite) HasMark(mark string) bool {
	s.LoadConfig()
	if s.Config.Marks == nil {
		return false // no marks
	}
	for _, v := range s.Config.Marks {
		if v == mark {
			return true
		}
	}
	return false
}

// ExecuteOnMarks checks if any of the marks is specified in the configuration.
// and skip the test if none of the marks is specified
func (s *CfyTestSuite) ExecuteOnMarks(marks ...string) {
	s.LoadConfig()

	okToRun := false

	for _, mark := range marks {
		if s.HasMark(mark) {
			// mark is specified in config
			okToRun = true
			break
		}
	}
	if !okToRun {
		s.T().Skipf("Test skipped as %v is not specified in test configuration", marks)
	}
}

// RequiresRunAsPrivilegedUser checks and skips the test if it is not run by privileged
// user
func (s *CfyTestSuite) RequiresRunAsPrivilegedUser() {
	priv, err := utils.RunByPrivilegedUser()
	if err != nil {
		s.T().Errorf("Cannot get information about current user: %v", err)
	}
	if !priv {
		// not run by privileged user, skip
		s.T().Skip("Test skipped as it is not run by privileged user")
	}
}

// RequiresRunAsUnprivilegedUser checks and skips the test if it is not run by unprivileged
// user
func (s *CfyTestSuite) RequiresRunAsUnprivilegedUser() {
	priv, err := utils.RunByPrivilegedUser()
	if err != nil {
		s.T().Errorf("Cannot get information about current user: %v", err)
	}
	if priv {
		// run by privileged user, skip
		s.T().Skip("Test skipped as it is run by privileged user")
	}
}

// RequiresCClientInstalled checks and skips the test if Centrify Client is not
// installed.
func (s *CfyTestSuite) RequiresCClientInstalled() {
	s.LoadConfig()

	if !s.agentInstalled {
		s.T().Skip("Test skipped as Centrify Client is not installed")
	}
}

// RequiresCClientRunning checks and skips the test if CClient is not running
func (s *CfyTestSuite) RequiresCClientRunning() {
	s.LoadConfig()

	if !s.agentInstalled {
		s.T().Skip("Test skipped as Centrify Client is not installed")
	}
	if s.agentStatus == utils.AgentStatusUnknown || s.agentStatus == utils.AgentStatusStopped {
		s.T().Skipf("Test skipped as Centrify Client is not running. Agent status: [%s]\n", s.agentStatus)
	}
	// check if tenant is same as requested
	if s.Config.TenantURL != "" && s.Config.TenantURL != s.enrolledTenant {
		s.T().Fatalf("Configuration requests to use tenant %v but system is connected to %v\n",
			s.Config.TenantURL, s.enrolledTenant)
	}
}

// RequiresCClientNotRunning checks and skips the test if CClient is running
func (s *CfyTestSuite) RequiresCClientNotRunning() {
	s.LoadConfig()

	if !s.agentInstalled {
		// If it is not installed then it is definitely not running :).
		return
	}

	if s.agentStatus == utils.AgentStatusConnected {
		s.T().Skipf("Test skipped as Centrify Client is running. Agent status: [%s]", s.agentStatus)
	}
	s.T().Logf("Centrify Client is not running. Agent status: [%s]", s.agentStatus)
}

// RequiresActiveTenant checks and skips the test if system is not enrolled and
// connected to the tenant. Also, FAILS test if the current tenant is different
// from the one specified in the config file.
func (s *CfyTestSuite) RequiresActiveTenant() {
	s.LoadConfig()

	if s.agentStatus != utils.AgentStatusConnected {
		s.T().Skipf("Test skipped as system is not connected to tenant. Agent status: [%s]\n", s.agentStatus)
	}
	// check if tenant is same as requested
	if s.Config.TenantURL != "" && s.Config.TenantURL != s.enrolledTenant {
		s.T().Fatalf("Configuration requests to use tenant %v but system is connected to %v\n",
			s.Config.TenantURL, s.enrolledTenant)
	}
}

// RequiresEnrolledTenant checks and skips the test if system is not
// enrolled.  Also, FAILS test if the current tenant is different from
// the one specified in the config file
func (s *CfyTestSuite) RequiresEnrolledTenant() {
	s.LoadConfig()

	if !s.agentEnrolled {
		s.T().Skip("Test skipped as system is not enrolled")
	}
	// check if tenant is same as requested
	if s.Config.TenantURL != "" && s.Config.TenantURL != s.enrolledTenant {
		s.T().Fatalf("Configuration requests to use tenant %v but system is connected to %v\n",
			s.Config.TenantURL, s.enrolledTenant)
	}
}

// RequiresVault checks and skips test if HashiCorp Vault is not installed and running.
func (s *CfyTestSuite) RequiresVault() {
	s.LoadConfig()

	if !s.vaultRunning {
		s.T().Skip("Test skipped as HashiCorp Vault is not running")
	}
}

// checkVault checks if HashiCorp Vault is not installed and running.
func (s *CfyTestSuite) checkVault() {
	_, err := exec.LookPath("vault")
	if err != nil {
		s.T().Log(err)
		return
	}

	resp, err := http.Get("http://localhost:8200")
	if err != nil || resp.StatusCode != http.StatusOK {
		s.T().Log(err)
		return
	}

	s.vaultRunning = true
}
