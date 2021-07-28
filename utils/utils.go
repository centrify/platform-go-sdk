// Package utils implements various functions about the system environment.
package utils

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/centrify/platform-go-sdk/internal/cversion"
)

// Centrify Client Status
const (
	AgentStatusUnknown      = "unknown"
	AgentStatusConnected    = "connected"
	AgentStatusDisconnected = "disconnected"
	AgentStatusStopped      = "stopped"
	AgentStatusStarting     = "starting"
	AgentStatusDisabled     = "disabled"
)

// exit status for cinfo

// ExitCodeNotEnrolled means system is not enrolled
const ExitCodeNotEnrolled = 10

// Errors
var (
	ErrCannotGetToken        = errors.New("Cannot obtain token")
	ErrCannotSetupConnection = errors.New("Cannot setup connection to Centrify Client")
	ErrCannotDecryptToken    = errors.New("Cannot decrypt token")
	ErrClientNotInstalled    = errors.New("Centrify client Not installed")
	ErrCommunicationError    = errors.New("Communication error with Centrify Client")
	ErrExpiredPublicKey      = errors.New("Expired public key")
	ErrGettingPublicKey      = errors.New("Error in getting Centrify Client's public key")
	ErrGettingResourceOwner  = errors.New("Cannot get resource owner credential")
	ErrInvalidCredential     = errors.New("Invalid username or password")
	ErrNotEnrolled           = errors.New("Not enrolled")
)

// try to parse out version string like 21.1-101
var versionRegEx = regexp.MustCompile(`\d+\.\d+(-\d+)?(-Debug|-Release)?(-NotYet)?`)

// GetCClientVersion returns the version of Centrify Client
func GetCClientVersion() (string, error) {
	installed, err := IsCClientInstalled()
	if err != nil {
		return "", err
	}
	if !installed {
		return "", ErrClientNotInstalled
	}
	cmd := getCinfoPath()
	out, err := exec.Command(cmd, "--version").Output()
	if err != nil {
		return "", err
	}

	res := string(out)

	// res is something like "Centrify cagent (Centrify Client 21.4-106)"
	// we want to extract the version string 21.4-106
	verStr := versionRegEx.FindString(res)
	return verStr, nil
}

// VerifyCClientVersionReq verifies if the version of Centrify Client meets the version requirement
func VerifyCClientVersionReq(req string) (bool, error) {

	ver, err := GetCClientVersion()
	if err != nil {
		// error in getting version
		return false, err
	}

	curVersion, err := cversion.Parse(ver)
	if err != nil {
		// error in parsing version
		return false, err
	}
	reqVersion, err := cversion.Parse(req)
	if err != nil {
		// error in parsing version
		return false, err
	}

	res := curVersion.Compare(reqVersion)
	return res != cversion.Earlier, nil
}

// IsCClientInstalled returns a bool value that shows whether the Centrify Client is installed on the system
// by checking for the existence of cagent.
func IsCClientInstalled() (bool, error) {
	_, err := os.Stat(getCagentPath())
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetCClientStatus returns the status of Centrify Client
func GetCClientStatus() (string, error) {
	cmd := getCinfoPath()
	out, err := exec.Command(cmd, "--agent-status").Output()

	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			// Something bad happend.
			return "", err
		}
		// Otherwise, it is an ExitError, which means that command
		// exited with non-zero exit code, which in turn means that
		// cagent is not connected.
	}

	res := string(out)
	res = strings.TrimSuffix(res, "\n")
	return res, nil
}

// GetEnrolledTenant returns the tenant that the system is enrolled to.  If the system is not
// enrolled to any tenant, ErrNotEnrolled is returned as error.
func GetEnrolledTenant() (string, error) {
	cmd := getCinfoPath()
	out, err := exec.Command(cmd, "--tenant").Output()

	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			if exitError.ProcessState.ExitCode() == ExitCodeNotEnrolled {
				// return specific error for this case
				return "", ErrNotEnrolled
			}
		}
		return "", err
	}

	// get tenant URL, strip off https:// and trailing "/" and newline
	res := string(out)
	res = strings.TrimSuffix(res, "/\n")
	return strings.TrimPrefix(res, "https://"), nil
}

// GetDMCEndPoint returns the endpoint used by applications to request machine credential
func GetDMCEndPoint() string {
	return endpointDMC
}
