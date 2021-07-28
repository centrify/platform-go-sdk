package utils

import (
	"os/user"
)

// getCagentPath returns the path name of cagent executable
func getCagentPath() string {
	return "/opt/centrify/sbin/cagent"
}

// RunByPrivilegedUser returns true if the program is run by privileged user (root in Linux)
func RunByPrivilegedUser() (bool, error) {
	user, err := user.Current()
	return err == nil && user.Uid == "0", err
}

func getCinfoPath() string {
	return "/usr/bin/cinfo"
}
