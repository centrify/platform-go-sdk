# Unit test framework and utilities

Package testutils contains all support functions that may be used in Centrify "go tests". The
test framework is based on the packages provided in testing and github.com/stretchr/testify.

## Test suite organization

Developers should try to organize tests into test suites.  Each test suite should be named something like <suite_feature>Suite and there
should be a function name Test<suite_feature>Suite in the test program.   This allows anyone to select individual
test suites to run by specifying "go test -run XYZ" which will run all the tests in all test suites that includes XYZ as part of the test suite name.

## Requirements for test files

1. The test file must import github.com/centrify/testutils
2. Each test suite must has testutils.CfyTestSuite as an anonymous member.

## Test Configuration file

Each test suite should read a test configuration file  (in JSON) in its own implementation of SetupTestSuite() method.   The path of the configuration file
is specified by "-config" parameter.

Current supported parameters in the configuration parameters are:

 - **TenantURL**:  specify that the suite requires the system be enrooled to this tenant URL (e.g., devdog.centrify.com).
			 If this is blank or not specified, tests will use the tenant that the current system is enrolled.
			 If the system is enrolled to a tenant that is different from this URL setting, all tests that calls
			 RequiresActiveTenant() or RequiresEnrolledTenant() will fail.
 - **Marks**:      A json string array that specifies the matching test marks.  A test can check if a mark exists in
             the configuration file to determine whether the test should be skipped or not.

Supported test marks:
- **scalability**:   indicates running scalability tests.
- **integration**:   indicates running integration tests.  Requires connectivity to tenant.
- **build**:         indicates the test is running as part of build.  Should not require connectivity to tenant, or a running Centrify Client.

If no configuration file is specified:
- "build" is used as the only test mark
- if any test requires the system to be enrolled to a tenant and/or the tenant is active, there is no restriction on what that tenant is.

## Individual test function

Each test function must do the followings:

  1. All test function names must start with Test.  Otherwise the test WILL NOT be picked up by "go test".
  2. If you want the test to execute only when any of the marks is specified in the configuration file, the test needs to call ExecuteOnMarks().
     For example,  ExecuteOnMarks("nightly", "integration") means that the test is run if any of "nightly"/"integration" is specified in "Marks"
     in the test configuration file.
  3. If your test requires the test user to be privileged user (root in Linux), call RequiresRunAsPrivilegedUser().  The test will be
     marked as skipped if it is not run by root in Linux.  Similarly, if your test requires the test user to be unprivileged (i.e., not 
	 root in Linux), call RequiresRunAsUnprivilegedUser().
  4. Call one of the following RequiresXXX() function if you test requires certain environment setup:
	 - RequiresActiveTenant(): The test system must be connected to an active tenant.
	 - RequiresCClientRunning(): Centrify Client must be running in the system.  It may or may not be connected to the tenant.
	 - RequiresEnrolledTenant(): The test system is enrolled to the tenant.  Centrify Client may or may not be running.
	      Also, the system may/may not be connected to the tenant.
	 - RequiresCClientInstalled(): Centrify Client must be installed in the system.  It may or not may be enrolled.
  5. If RequiresActiveTenant()/RequiresCClientRunning()/RequiresEnrolledTenant() is called, AND TenantURL is specified in the configuration
     file, AND the system is enrolled to a different tenant, the test is not run and marked as FAILED.

You can review the test files testutils_test.go and test2_test.go to review how these tests are used.

## How to run tests

You can use the two test files in this directory as examples for tests.

 go test ./... -args -config=/tmp/config.json<br>
	This will run all tests using the configuration file /tmp/config.json

 go test ./... -run=MoreTest -args -config=/tmp/config.json<br>
	This will run all the tests in the test suite TestUtilsMoreTestSuite but not TestUtilsTestSuite

 go test ./... -run=UtilsTest -args -config=/tmp/config.json<br>
	This will run all the tests in the test suite TestUtilsTestSuite but not TestUtilsMoreTestSuite

If you change the system configuration (e.g., starts Centrify Client, enroll the system etc), you can force "go test" to not
used previous test result by running the command "go clean -testcache"