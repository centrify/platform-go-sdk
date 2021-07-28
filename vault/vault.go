package vault

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/centrify/platform-go-sdk/internal/lrpc"
	"github.com/centrify/platform-go-sdk/internal/securemessage"
	"github.com/centrify/platform-go-sdk/utils"
)

// GetHashiVaultToken returns an HashiCorp Vault token for the requested scope that has the
// identity of the current machine account.
//
// Possible error returns:
//
//  ErrCannotGetToken  - other errors in getting the token
//  ErrCannotSetupConnection - Cannot setup connection to Centrify Client
//  ErrClientNotInstalled - Centrify Client is not installed in system
//  ErrCommunicationError - Communication error with Centrify Client
func GetHashiVaultToken(scope string, vaultURL string) (string, error) {

	installed, err := utils.IsCClientInstalled()
	if err != nil {
		return "", fmt.Errorf("Cannot get Centrify Client installation status: %w", utils.ErrClientNotInstalled)
	}
	if !installed {
		return "", utils.ErrClientNotInstalled
	}

	// create a lrpc2 client and connect to it
	cl := lrpc.NewLrpc2ClientSession(utils.GetDMCEndPoint())
	if cl == nil {
		return "", utils.ErrCannotSetupConnection
	}
	err = cl.Connect()
	if err != nil {
		return "", utils.ErrCannotSetupConnection
	}
	defer cl.Close()

	// generate a public key so that Centrify Client can encrypt the vault
	// token in the reply
	pubKey, _, err := securemessage.GetPublicKey()
	if err != nil {
		return "", err
	}

	// success, convert the public key to a bytestream
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(*pubKey)
	if err != nil {
		return "", err
	}

	// send LRPC message to Centrify Client
	var args []interface{}
	args = append(args, scope, vaultURL, buffer.Bytes())
	results, err := lrpc.DoRequest(cl, lrpc.Lrpc2MsgGetHashicorpVaultToken, args)

	if err != nil {
		return "", utils.ErrCommunicationError
	}

	// Check results of LRPC call
	// return message should have
	// results[0] status
	// results[1] error message
	// results[2] access token (encrypted)
	if len(results) != 3 {
		return "", utils.ErrCommunicationError
	}

	status, ok := results[0].(int32)

	if ok {
		if status == 0 {
			// good status, token is in results[2]

			cipher, ok1 := results[2].([]string)
			if !ok1 {
				return "", utils.ErrCommunicationError
			}

			// now decrypt the token
			token, err := securemessage.DecryptString(cipher, "")
			if err != nil {
				return "", utils.ErrCannotDecryptToken
			}
			return token, nil
		}
		// error message return in second value
		errmsg, ok := results[1].(string)
		if ok {
			return "", fmt.Errorf("%w: %s", utils.ErrCannotGetToken, errmsg)
		}

	}
	return "", utils.ErrCommunicationError
}
