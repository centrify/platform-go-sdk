package lrpc

import (
	"errors"
	"fmt"
	"net"

	ps "github.com/mitchellh/go-ps"
	"golang.org/x/sys/unix"
)

// type session represents the LRPC session
// It needs to implement the SessionCtxt interface
type session struct {
	SessionCtxtBase
	// conn - network connection information
	conn              net.Conn
	cred              *unix.Ucred
	clientProcessName string
}

// NewSessionCtxt creates a session context object
func NewSessionCtxt(conn net.Conn) (SessionCtxt, error) {
	p := new(session)
	p.conn = conn
	cred, err := getPeerConn(conn)
	if err != nil {
		return nil, err
	}
	p.cred = cred
	p.Set("_uid", cred.Uid)
	return p, err
}

// IsPrivileged returns true if the calling is privileged.
//
// This function is OS dependent.
//
// In Linux, it means the effective UID of the calling process is 0.
func (p *session) IsPrivileged() (bool, error) {
	if p.cred == nil {
		return false, errors.New("No peer credential")
	}
	return p.cred.Uid == 0, nil
}

// GetProcesId returns the PID of the process that sends the LRPC request
func (p *session) GetProcessID() (int32, error) {
	if p.cred == nil {
		return 0, errors.New("No peer credential")
	}
	return p.cred.Pid, nil
}

// GetCallerUserID returns the UID (as a string) of the process that sends the LRPC request
func (p *session) GetCallerUserID() (string, error) {
	if p.cred == nil {
		return "", errors.New("No peer credential")
	}
	return fmt.Sprintf("%d", p.cred.Uid), nil
}

func getPeerConn(conn net.Conn) (*unix.Ucred, error) {

	// Convert to net.UnixConn from generic net.Conn type. See
	// newUnixConn() in golang src/pkg/net/unixsock_posix.go
	unixConn := conn.(*net.UnixConn)

	//
	// Note that from golang source code, this call will also put the
	// original file descriptor into blocking mode (was non-blocking when
	// accepted):
	//
	// File() will set fd into blocking mode:
	// https://golang.org/src/net/fd_unix.go?h=SetNonblock#L504
	//
	// Accept() will mark the return fd as non-blocking and close-on-exec:
	// https://golang.org/src/net/sock_cloexec.go?h=SetNonblock#L82
	// https://golang.org/src/net/sys_cloexec.go?h=SetNonblock#L52
	//
	// Why File() uses dup() and why set fd into blocking mode:
	// https://github.com/golang/go/issues/5052
	//
	dupFile, err := unixConn.File()
	if err != nil {
		return nil, fmt.Errorf("Cannot clone file to get peer information: %v", err)
	}
	defer dupFile.Close()

	// Get the duplicate Unix file descriptor
	dupFd := int(dupFile.Fd())

	// File() sets the duplicate fd into blocking mode. Unfortunately this
	// property is shared among all fds. So the original fd in net.Conn
	// will also be affected. It is suggested to work around this issue by
	// setting non-blocking again. And we follow here.
	unix.SetNonblock(dupFd, true)

	// Get user credential here. We can revisit later if the performance is
	// a problem, e.g. check if the LRPC2 message type requires process
	// auth before getting user credential through socket.

	cred, err := unix.GetsockoptUcred(dupFd, unix.SOL_SOCKET, unix.SO_PEERCRED)
	if err != nil {
		return nil, fmt.Errorf("Cannot get socket information: %v", err)
	}

	return cred, nil
}

func (p *session) GetProgram() (string, error) {
	if p.clientProcessName != "" {
		return p.clientProcessName, nil
	}

	pid, err1 := p.GetProcessID()
	if err1 != nil {
		return "", err1
	}

	proc, err2 := ps.FindProcess(int(pid))
	if err2 != nil {
		return "", err2
	}

	if proc == nil {
		// process exits already
		return "", errors.New("Process is terminated already")
	}

	p.clientProcessName = proc.Executable()

	return p.clientProcessName, nil
}
