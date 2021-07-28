// Package dmc provides application with APIs for obtaining Delegated Machine Credentials (DMC)
//
// Run go tests
//
// You need to do the followings to have a successful run of go unit tests:
//  1. Install Centrify Client (version 21.5 or later) on the system.
//  2. Enroll Centrify Client to a PAS tenant:
//    - enable DMC feature by specifying "-F all" or "-F dmc" in cenroll command line.
//    - specify the DMC scope "testsdk" by specifying "-d testsdk:security/whoami" in cenroll command line.
//  3. Run the unit test as root
//
//
// Sample Program
//
// A sample program can be found in https://github.com/centrify/platform-go-sdk/examples/dmc
package dmc

import (
	"fmt"
	"strings"

	"github.com/centrify/platform-go-sdk/internal/lrpc"
	"github.com/centrify/platform-go-sdk/utils"
)

func getDMCEndPoint() string {
	return utils.GetDMCEndPoint()
}

// GetDMCToken returns an oauth token for the requested scope that has the
// identity of the current machine account.
//
// Possible error returns:
//
//  ErrCannotGetToken  - other errors in getting the token
//  ErrCannotSetupConnection - Cannot setup connection to Centrify Client
//  ErrClientNotInstalled - Centrify Client is not installed in system
//  ErrCommunicationError - Communication error with Centrify Client
func GetDMCToken(scope string) (string, error) {

	installed, err := utils.IsCClientInstalled()
	if err != nil {
		return "", fmt.Errorf("Cannot get Centrify Client installation status: %w", utils.ErrClientNotInstalled)
	}
	if !installed {
		return "", utils.ErrClientNotInstalled
	}

	// create a lrpc2 client and connect to it
	cl := lrpc.NewLrpc2ClientSession(getDMCEndPoint())
	if cl == nil {
		return "", utils.ErrCannotSetupConnection
	}
	err = cl.Connect()
	if err != nil {
		return "", utils.ErrCannotSetupConnection
	}
	defer cl.Close()

	// send LRPC message to Centrify Client

	var args []interface{}
	args = append(args, scope)
	results, err := lrpc.DoRequest(cl, lrpc.Lrpc2MsgIDAdminClientGetToken, args)

	if err != nil {
		return "", utils.ErrCommunicationError
	}

	// Check results of LRPC call
	// return message should have
	// results[0] status
	// results[1] error message
	// results[2] access token
	if len(results) != 3 {
		// error...
		return "", utils.ErrCommunicationError
	}

	status, ok := results[0].(int32)

	if ok {
		if status == 0 {
			// good status, token is in results[2]

			token, ok1 := results[2].(string)
			if !ok1 {
				return "", utils.ErrCommunicationError
			}
			return token, nil
		} else {
			// error message return in second value
			errmsg, ok := results[1].(string)
			if ok {
				return "", fmt.Errorf("%w: %s", utils.ErrCannotGetToken, errmsg)
			}
		}
	}
	return "", utils.ErrCommunicationError
}

// GetEnrollmentInfo returns information about Centrify Client enrollment information
func GetEnrollmentInfo() (string, string, error) {

	installed, err := utils.IsCClientInstalled()
	if err != nil {
		return "", "", fmt.Errorf("Cannot get Centrify Client installation status: %w", utils.ErrClientNotInstalled)
	}
	if !installed {
		return "", "", utils.ErrClientNotInstalled
	}

	// create a lrpc2 client and connect to it
	cl := lrpc.NewLrpc2ClientSession(getDMCEndPoint())
	if cl == nil {
		return "", "", utils.ErrCannotSetupConnection
	}
	err = cl.Connect()
	if err != nil {
		return "", "", utils.ErrCannotSetupConnection
	}
	defer cl.Close()

	// send LRPC message to Centrify Client

	results, err := lrpc.DoRequest(cl, lrpc.LrpcMsgIDClientInfo, nil)

	if err != nil {
		return "", "", utils.ErrCommunicationError
	}

	// Check results of LRPC call
	if len(results) < 1 {
		// error...
		return "", "", utils.ErrCommunicationError
	}

	info, ok := results[0].(map[string]string)

	if ok {
		tenantURL, ok1 := info["ServiceURI"]
		oAuthClientID, ok2 := info["OauthClientID"]

		if !ok1 || !ok2 {
			// error...
			return "", "", utils.ErrCommunicationError
		}

		// strip off https:// or http:// prefix
		tenantURL = strings.TrimPrefix(tenantURL, "https://")
		tenantURL = strings.TrimPrefix(tenantURL, "http://")
		tenantURL = strings.TrimSuffix(tenantURL, "/")
		return tenantURL, oAuthClientID, nil
	}
	return "", "", utils.ErrCommunicationError
}
