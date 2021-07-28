// Package cversion provides utilities for working with Centrify version strings.
package cversion

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Major       int
	Minor       int
	BuildNumber int
	Debug       bool
	NotYet      bool
}

type ComparisonResult int

const (
	Earlier ComparisonResult = iota - 1
	Equal
	Subsequent
)

func (version *Version) Compare(otherVersion Version) ComparisonResult {
	if version.Major < otherVersion.Major {
		return Earlier
	}
	if version.Major > otherVersion.Major {
		return Subsequent
	}

	if version.Minor < otherVersion.Minor {
		return Earlier
	}
	if version.Minor > otherVersion.Minor {
		return Subsequent
	}

	if version.BuildNumber < otherVersion.BuildNumber {
		return Earlier
	}
	if version.BuildNumber > otherVersion.BuildNumber {
		return Subsequent
	}

	return Equal
}

func (version *Version) MajorMinorString() string {
	return fmt.Sprintf("%d.%d", version.Major, version.Minor)
}

func (version *Version) String() string {
	versionStr := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.BuildNumber)
	if version.Debug {
		versionStr += "-Debug"
	}
	if version.NotYet {
		versionStr += "-NotYet"
	}
	return versionStr
}

var versionRegex = regexp.MustCompile("(?i)^(\\d+)\\.(\\d+)([\\.-])?(\\d+)?(-Debug|-Release)?(-NotYet)?")

// Parses version strings into a Version struct.
// Strings like the following are accepted:
//   "17.2", "17.2.100", "17.2-100", "17.2.100-NotYet", "17.2-100-NotYet", etc.
// Strings like the following are (currently) not accepted:
//   "", "foo", "17.2.", "17.2-NotYet", etc.
// If the build number is missing from the input, it is treated as if it were 0.
func Parse(versionStr string) (Version, error) {
	var version Version

	submatches := versionRegex.FindStringSubmatch(versionStr)
	if submatches == nil {
		return version, errors.New("Input does not match expected format")
	}

	// Golang seems to use the empty string to denote unmatched optional groups.
	// Hence, the number of submatches on success will always be:
	//   1 for the entire match.
	//   2 for the major and minor version numbers.
	//   1 for the build number separator.
	//   1 for the build number.
	//   1 for "-Debug" or "-Release".
	//   1 for "-NotYet".
	// ... even if neither the build number separator nor the build number is
	//     present.
	// For example, submatchCount will be 7 for "17.2", "17.2.", and "17.2.100".
	expectedSubmatchCount := 7
	submatchCount := len(submatches)
	if submatchCount != expectedSubmatchCount {
		err := fmt.Errorf("Expected %d sub-matches, got %d",
			expectedSubmatchCount, submatchCount)
		return version, err
	}

	// This block of code catches invalid versions like "17.2.", "17.2-", and
	// "17.2..100".
	// This block also catches versions that have the format "17.2-Foo".
	// Nothing currently returns versions in this format, so we don't support
	// the format to avoid unnecessary complexity.
	buildNumberPresent := submatches[4] != ""
	var expectedBuildNumSeparatorCount int
	if buildNumberPresent {
		expectedBuildNumSeparatorCount = 1
	}
	buildNumSeparatorCount := len(submatches[3])
	if buildNumSeparatorCount != expectedBuildNumSeparatorCount {
		err := fmt.Errorf("Expected %d build number separators, got %d",
			expectedBuildNumSeparatorCount, buildNumSeparatorCount)
		return version, err
	}

	var err error
	version.Major, err = strconv.Atoi(submatches[1])
	if err != nil {
		return version, err
	}

	version.Minor, err = strconv.Atoi(submatches[2])
	if err != nil {
		return version, err
	}

	version.BuildNumber = 0
	if buildNumberPresent {
		version.BuildNumber, err = strconv.Atoi(submatches[4])
		if err != nil {
			return version, err
		}
	}

	version.Debug = strings.EqualFold(submatches[5], "-Debug")

	version.NotYet = submatches[6] != ""

	return version, nil
}
