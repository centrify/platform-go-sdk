package lrpc

import (
	"fmt"
	"net"
	"os"
)

/*
StartListener is OS-dependent function to establish an internal listener.
 Input parameter:
 	endpoint (string):          The name of the endpoint
	acl (AccessControlLevel):   The ACL for starting listener, normally is null for linux

 Return values:
 	net.Listener:	network listener object established.
 	error

Note:

This function should remove the endpoint if it already exists.
*/
func StartListener(endpoint string, acl AccessControlLevel) (net.Listener, error) {

	// remove existing one first
	_, err := os.Stat(endpoint)
	if err == nil || os.IsExist(err) {
		// already exist
		err = os.Remove(endpoint)
		if err != nil {
			return nil, fmt.Errorf("Cannot remove existing endpoint: %v", err)
		}
	}

	unixAddr := &net.UnixAddr{
		Name: endpoint,
		Net:  "unix",
	}

	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		return nil, fmt.Errorf("Cannot create listen socket %s: %v", endpoint, err)
	}
	err = os.Chmod(endpoint, 0777)
	return listener, err
}

/*
ConnectToServer is OS-dependent code for connecting to an internal endpoint.  It is used by a RPC client.

 Input parameters:
 	endpoint (string):   The endpoint to connect to.

 Return values:
 	net.Conn:   The established network connection.
	error:		err
*/
func ConnectToServer(endpoint string) (net.Conn, error) {
	conn, err := net.Dial("unix", endpoint)
	if err != nil {
		return nil, fmt.Errorf("Error in connecting to server: %v", err)
	}
	return conn, nil

}
