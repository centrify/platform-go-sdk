package cversion

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	DEBUG_YES  = true
	DEBUG_NO   = false
	NOTYET_YES = true
	NOTYET_NO  = false
)

// cversion test DOES NOT need to parse command line argument
// do this here to avoid errors if -config is passed in
// Anyway we need to declare them here so that go test will not complain
var (
	configPtr      = flag.String("config", "", "configuration file")
	configString   = flag.String("config-string", "", "configuration string")
	VaultRootToken = flag.String("vault-root-token", "root", "Vault root token")
)

func TestParse_InvalidCases(t *testing.T) {
	assert := assert.New(t)

	var invalidCases = []struct {
		input string
		msg   string
	}{
		{"", "Empty string"},
		{"foo", "Arbitrary non-version string"},
		{"1", "Only major version"},
		{"1.", "Only major version and first separator"},
		{"x.y.z", "Non-integer version components"},
		{"x.1", "Non-integer major version"},
		{"1.x", "Non-integer minor version"},
		{"1.2.y", "Non-integer build number"},
		{"1x2.3.4", "Not fully integer version component"},
		{"1-2.3", "Wrong major/minor version separator"},
		{"1.2.", "Trailing . build number separator"},
		{"1.2-", "Trailing - build number separator"},
		{"1.2..100", "More than one . build number separator"},
		{"1.2--100", "More than one - build number separator"},
		{"1.2.-100", "Mixed build number separators (. first)"},
		{"1.2-.100", "Mixed build number separators (- first)"},
		{"1000000000000000000000000.2.3", "Huge major version"},
		{"1.2000000000000000000000000.3", "Huge minor version"},
		{"1.2.3000000000000000000000000", "Huge build number"},
	}
	for _, invalidCase := range invalidCases {
		_, err := Parse(invalidCase.input)
		assert.Error(err,
			"Should fail to parse '%s' for test case '%s'", invalidCase.input, invalidCase.msg)
	}
}

func TestParse_ValidCases(t *testing.T) {
	assert := assert.New(t)

	var validCases = []struct {
		input       string
		major       int
		minor       int
		buildNumber int
		debug       bool
		notYet      bool
		msg         string
	}{
		{"1.2", 1, 2, 0, DEBUG_NO, NOTYET_NO, "Only major and minor version"},
		{"1.2-3", 1, 2, 3, DEBUG_NO, NOTYET_NO, "Build number separated using -"},
		{"1.2.3", 1, 2, 3, DEBUG_NO, NOTYET_NO, "All components present"},
		{"1.2.3.4", 1, 2, 3, DEBUG_NO, NOTYET_NO, "Extra component separated using ."},
		{"1.2-3-4", 1, 2, 3, DEBUG_NO, NOTYET_NO, "Extra component separated using -"},
		{"1.2.3/4", 1, 2, 3, DEBUG_NO, NOTYET_NO, "Extra component not separated using . or -"},
		{"1.2-3-Foo", 1, 2, 3, DEBUG_NO, NOTYET_NO, "Unknown tag"},
		{"1.2-3-Debug", 1, 2, 3, DEBUG_YES, NOTYET_NO, "Debug tag"},
		{"1.2-3-DeBuG", 1, 2, 3, DEBUG_YES, NOTYET_NO, "Debug tag parsing should be case insensitive"},
		{"1.2-3-Release", 1, 2, 3, DEBUG_NO, NOTYET_NO, "Release tag"},
		{"1.2-3-ReLeAsE", 1, 2, 3, DEBUG_NO, NOTYET_NO, "Release tag parsing should be case insensitive"},
		{"1.2-3-NotYet", 1, 2, 3, DEBUG_NO, NOTYET_YES, "NotYet tag"},
		{"1.2-3-NoTyEt", 1, 2, 3, DEBUG_NO, NOTYET_YES, "NotYet parsing should be case insensitive"},
		{"17.3.141-Release-NotYet-2017-02-14T04:16:03-286736", 17, 3, 141, DEBUG_NO, NOTYET_YES,
			"Result.Cloud string obtained by calling SysInfo/Version on devdog"},
		{"1.0.0-Debug-NotYet-2018-03-27T23:24:06-20180328", 1, 0, 0, DEBUG_YES, NOTYET_YES,
			"Result.Cloud string obtained by calling SysInfo/Version on developer build"},
		{"17.3.141-286736", 17, 3, 141, DEBUG_NO, NOTYET_NO,
			". should take precedence over - for the build number separator"},
		// We can't test the return result of actually calling
		// centrify.GetVersion() or centrify.GetFullVersion() works because:
		//   1. We don't know the correct version ahead of time.
		//   2. The functions don't work in a "go test" environment.
		{"16.11-100-NotYet", 16, 11, 100, DEBUG_NO, NOTYET_YES, "centrify.GetFullVersion() result"},
		{"16.11.100-NotYet", 16, 11, 100, DEBUG_NO, NOTYET_YES, "Backend version shown by the Web UI"},
	}
	for _, validCase := range validCases {
		version, err := Parse(validCase.input)
		if assert.NoError(err, "Should be able to parse '%s' for test case '%s'", validCase.input, validCase.msg) {
			assert.Equal(validCase.major, version.Major,
				"Major version should match for test case '%s'", validCase.msg)
			assert.Equal(validCase.minor, version.Minor,
				"Minor version should match for test case '%s'", validCase.msg)
			assert.Equal(validCase.buildNumber, version.BuildNumber,
				"Build number should match for test case '%s'", validCase.msg)
			assert.Equal(validCase.notYet, version.NotYet,
				"NotYet status should match for test case '%s'", validCase.msg)
		}
	}
}

func TestVersionCompare(t *testing.T) {
	assert := assert.New(t)
	comparisonVer := Version{Major: 10, Minor: 5, BuildNumber: 100}

	var testCases = []struct {
		major            int
		minor            int
		buildNumber      int
		comparisonResult ComparisonResult
		msg              string
	}{
		{10, 5, 100, Equal, "10.5.100 should be equal to itself"},
		// Major version
		{11, 0, 0, Earlier, "10.5.100 should be older than 11.0.0"},
		{10, 0, 0, Subsequent, "10.5.100 should be newer than 10.0.0"},
		// Minor version
		{10, 6, 100, Earlier, "10.5.100 should be older than 10.6.100"},
		{10, 4, 100, Subsequent, "10.5.100 should be newer than 10.4.100"},
		// Build number
		{10, 5, 111, Earlier, "10.5.100 should be older than 10.5.111"},
		{10, 5, 99, Subsequent, "10.5.100 should be newer than 10.5.99"},
		// Mix
		{11, 4, 101, Earlier, "10.5.100 should be older than 11.4.101"},
		{9, 6, 101, Subsequent, "10.5.100 should be newer than 9.6.101"},
	}
	for _, testCase := range testCases {
		version := Version{
			Major:       testCase.major,
			Minor:       testCase.minor,
			BuildNumber: testCase.buildNumber,
		}
		assert.Equal(comparisonVer.Compare(version), testCase.comparisonResult, testCase.msg)
	}
}

func TestVersion_MajorMinorString(t *testing.T) {
	assert := assert.New(t)
	version := Version{Major: 10, Minor: 5, BuildNumber: 100}

	assert.Equal("10.5", version.MajorMinorString(), "Format should be x.y")
	_, err := Parse(version.MajorMinorString())
	assert.NoError(err, "Result should be parsable")
}

func TestVersion_String(t *testing.T) {
	assert := assert.New(t)
	version := Version{Major: 10, Minor: 5, BuildNumber: 100}

	var testCases = []struct {
		debug          bool
		notYet         bool
		expectedOutput string
		msg            string
	}{
		{DEBUG_NO, NOTYET_NO, "10.5.100",
			"Format should be x.y.z if Debug=false, NotYet=false"},
		{DEBUG_YES, NOTYET_NO, "10.5.100-Debug",
			"Format should be x.y.z-Debug if Debug=true, NotYet=false"},
		{DEBUG_NO, NOTYET_YES, "10.5.100-NotYet",
			"Format should be x.y.z-NotYet if Debug=false, NotYet=true"},
		{DEBUG_YES, NOTYET_YES, "10.5.100-Debug-NotYet",
			"Format should be x.y.z-Debug-NotYet if Debug=true, NotYet=true"},
	}
	for _, testCase := range testCases {
		version.Debug = testCase.debug
		version.NotYet = testCase.notYet

		assert.Equal(testCase.expectedOutput, version.String(),
			"Result should be as expected for test case '%s'", testCase.msg)

		stringResult := version.String()
		parsedVersion, err := Parse(stringResult)
		if assert.NoError(err, "Result '%s' should be parsable for test case '%s'", stringResult, testCase.msg) {
			assert.Equal(Equal, version.Compare(parsedVersion),
				"Stringified then parsed version should be equal to itself for test case '%s'",
				testCase.msg)
		}
	}
}
