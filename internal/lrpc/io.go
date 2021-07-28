package lrpc

import (
	"errors"
	"io"
	"net"

	"github.com/centrify/platform-go-sdk/internal/logging"
)

func write(conn net.Conn, bytes []byte) error {
	var err error

	if bytes == nil {
		return errors.New("Internal error: write without bytes")
	}

	// NOTE: Make sure message data is NOT logged. Since message data can
	// contain sensitive information (e.g. password).

	_, err = conn.Write(bytes)

	if err != nil {
		// Seems sometimes the caller is already gone (e.g.
		// ControlLog2). So log error using a lower level.
		switch t := err.(type) {
		case *net.OpError:
			if t.Err.Error() == "i/o timeout" || t.Err.Error() == "broken pipe" {
				logging.Debugf("Failed to write to LRPC2 connection [%p]: %v", conn, t.Err.Error())
			}
		default:
			logging.Debugf("Failed to write to LRPC2 connection [%p]: %v", conn, err)
		}

		return err
	}

	return nil
}

func read(conn net.Conn, bytes []byte) (int, error) {
	var err error
	var n int

	if bytes == nil {
		return 0, errors.New("Internal error: read without bytes")
	}

	// NOTE: Make sure message data is NOT logged. Since message data can
	// contain sensitive information (e.g. password).

	n, err = conn.Read(bytes)

	//
	// EOF is normal for read:
	//
	// "EOF is the error returned by Read when no more input is available.
	// Functions should return EOF only to signal a graceful end of input."
	//
	if err == io.EOF {
		logging.Tracef("No more to read from LRPC2 connection [%p]: %v", conn, err)
		return 0, err
	}

	if err != nil {
		// Seems sometimes the caller is already gone (e.g.
		// ControlLog2). So log error using a lower level.
		switch t := err.(type) {
		case *net.OpError:
			if t.Err.Error() == "i/o timeout" || t.Err.Error() == "broken pipe" {
				logging.Debugf("Failed to read from LRPC2 connection [%p]: %v", conn, t.Err.Error())
			}
		default:
			logging.Debugf("Failed to read from LRPC2 connection [%p]: %v", conn, err)
		}

		return 0, err
	}

	return n, nil
}
