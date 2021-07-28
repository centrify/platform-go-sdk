// Package oauthhelper implements helper functions for Centrify Vault software to get OAuth access tokens
package oauthhelper

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"fmt"

	"github.com/centrify/platform-go-sdk/internal/lrpc"
	"github.com/centrify/platform-go-sdk/internal/securemessage"
	"github.com/centrify/platform-go-sdk/utils"
)

// common status return from Centrify Client
const (
	ErrExpiredPublicKey       = 2
	ErrGetResourceOwnerDenied = 3
)

// TODO: wrap error details in returned errors

// setupLrpcClient tries to setup and connect to the LRPC server in Centrify Client
func setupLrpcClient() (lrpc.MessageClient, error) {

	// check if Centrify Client is installed

	installed, err := utils.IsCClientInstalled()
	if err != nil {
		return nil, fmt.Errorf("Cannot get Centrify Client installation status: %w", utils.ErrClientNotInstalled)
	}
	if !installed {
		return nil, utils.ErrClientNotInstalled
	}

	// create a lrpc client and connect to it
	cl := lrpc.NewLrpc2ClientSession(utils.GetDMCEndPoint())
	if cl == nil {
		return nil, utils.ErrCannotSetupConnection
	}
	err = cl.Connect()
	if err != nil {
		return nil, utils.ErrCannotSetupConnection
	}
	return cl, nil
}

func getPublicKey(cl lrpc.MessageClient) (*rsa.PublicKey, uint32, error) {
	// password is sensitive information, get the public key from Centrify Client so that we can encrypt it

	results, err := lrpc.DoRequest(cl, lrpc.Lrpc2MsgIDGetPublicKey, nil)
	if err != nil {
		return nil, 0, utils.ErrCommunicationError
	}

	// check results of LRPC call
	// results[0]: status
	// results[1]: error message (if status != 0)
	// results[1]: keyID (uint32) (if status == 0)
	// results[2]: []byte - public key encoded in gob (only present if status == 0)
	if len(results) < 2 {
		return nil, 0, utils.ErrCommunicationError
	}

	status, ok := results[0].(int32)
	if !ok {
		return nil, 0, utils.ErrCommunicationError
	}
	if status != 0 {
		return nil, 0, utils.ErrCommunicationError
	}
	if len(results) != 3 {
		// always expect 3 return values for success
		return nil, 0, utils.ErrCommunicationError
	}
	keyID, ok := results[1].(uint32)
	if !ok {
		return nil, 0, utils.ErrCommunicationError
	}
	blob, ok := results[2].([]byte)
	if !ok {
		return nil, 0, utils.ErrGettingPublicKey
	}

	// now decode the bytestream into a public key
	buf := bytes.NewBuffer(blob)
	var key rsa.PublicKey
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&key)
	if err != nil {
		return nil, 0, utils.ErrGettingPublicKey
	}
	return &key, keyID, nil
}

// GetResourceOwnerToken sends a request to Centrify Client to get an OAuth token using "resource owner" grant request
func GetResourceOwnerToken(appID string, scope string, user string, passwd string) (accessToken string, tokenType string,
	expiresIn uint32, refreshToken string, err error) {

	var cl lrpc.MessageClient
	var pubKey *rsa.PublicKey
	var keyID uint32

	cl, err = setupLrpcClient()
	if err != nil {
		return
	}
	defer cl.Close()

	pubKey, keyID, err = getPublicKey(cl)
	if err != nil {
		return
	}

	// LRPC request parameters:
	// args[0]: applicationID
	// args[1]: scope
	// args[2]: username
	// args[3]: keyID
	// args[4]: encrypted password (in []string)

	args := make([]interface{}, 5)

	// encrypt the password, use the username as label
	args[4], err = securemessage.EncryptString(passwd, user, pubKey)
	if err != nil {
		return
	}
	args[0] = appID
	args[1] = scope
	args[2] = user
	args[3] = keyID

	// send LRPC request
	var results []interface{}

	results, err = lrpc.DoRequest(cl, lrpc.Lrpc2MsgIDGetResourceOwnerCred, args)

	if err != nil {
		return
	}

	// results should have:
	// ret[0]: status (0 for success, 1 for error)
	// For success:
	// ret[1]: access Token (string)
	// ret[2]: token type (string)
	// ret[3]: expiresIn (uint32)
	// ret[4]: refresh token
	// For error:
	// ret[1]: error message

	if len(results) < 2 {
		// not expected
		err = utils.ErrCommunicationError
		return
	}
	status, ok := results[0].(int32)
	if !ok {
		err = utils.ErrCommunicationError
		return
	}

	// handle error conditions
	switch status {
	case ErrExpiredPublicKey:
		err = utils.ErrExpiredPublicKey
		return

	case ErrGetResourceOwnerDenied:
		err = utils.ErrInvalidCredential
		return
	}

	if status != 0 {
		err = utils.ErrGettingResourceOwner
		return
	}

	// now get good status...try to unmarshal response
	if len(results) != 5 {
		// unexpected result
		err = utils.ErrCommunicationError
		return
	}

	var ok1, ok2, ok3, ok4 bool
	accessToken, ok1 = results[1].(string)
	tokenType, ok2 = results[2].(string)
	expiresIn, ok3 = results[3].(uint32)
	refreshToken, ok4 = results[4].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		err = utils.ErrCommunicationError
	}
	// return result....
	return
}
