// secretcli is a sample program that demonstrates how to use the secret package
// to manage secrets stored in Centrify PAS, Thycotic Secret Server and Thycotic Devops Secret Vault.
package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/sys/unix"

	"github.com/centrify/platform-go-sdk/secret"
)

// convertErrToExitStatus converts error to Unix exit status
// There are the exit status code and the corresponding errors:
// EPERM (1): ErrSecretTypeNotSupported, ErrCannotModifySecretType, ErrCannotModifySecretFolder
// ENOENT (2):	ErrFolderNotFound, ErrSecretNotFound
// EACCES (13):	ErrNoCreatePermission, ErrNoDeletePermission, ErrNoModifyPermission, ErrNoGetMetaDataPermission, ErrNoRetrievePermission
// EEXIST (17): ErrExists, ErrDeletedSecretExists
// ENOTDIR (20): ErrNotSecretFolder
// EISDIR (21): ErrNotSecretObject
// EINVAL (22): ErrBadPathName, ErrBadServerType
// ENOSYS (38): ErrNotImplementedYet
// ENOTEMPTY (39): ErrFolderNotEmpty
// EPROTO(72): ErrUnexpectedResponse
func convertErrToExitStatus(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, secret.ErrSecretTypeNotSupported) || errors.Is(err, secret.ErrCannotModifySecretType) ||
		errors.Is(err, secret.ErrCannotModifySecretFolder) {
		return int(unix.EPERM)
	} else if errors.Is(err, secret.ErrFolderNotFound) || errors.Is(err, secret.ErrSecretNotFound) {
		return int(unix.ENOENT)
	} else if errors.Is(err, secret.ErrNoCreatePermission) || errors.Is(err, secret.ErrNoDeletePermission) ||
		errors.Is(err, secret.ErrNoModifyPermission) || errors.Is(err, secret.ErrNoRetrievePermission) ||
		errors.Is(err, secret.ErrNoGetMetaDataPermission) {
		return int(unix.EACCES)
	} else if errors.Is(err, secret.ErrExists) || errors.Is(err, secret.ErrDeletedSecretExists) {
		return int(unix.EEXIST)
	} else if errors.Is(err, secret.ErrNotSecretFolder) {
		return int(unix.ENOTDIR)
	} else if errors.Is(err, secret.ErrNotSecretObject) {
		return int(unix.EISDIR)
	} else if errors.Is(err, secret.ErrBadPathName) || errors.Is(err, secret.ErrBadServerType) {
		return int(unix.EINVAL)
	} else if errors.Is(err, secret.ErrNotImplementedYet) {
		return int(unix.ENOSYS)
	} else if errors.Is(err, secret.ErrFolderNotEmpty) {
		return int(unix.ENOTEMPTY)
	} else if errors.Is(err, secret.ErrUnexpectedResponse) {
		return int(unix.EPROTO)
	}
	return -2 // unknown error
}
func main() {
	params, err := getConfiguration()
	if err != nil {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	var accessToken string
	var clientFactory secret.HTTPClientFactory
	var cl secret.Secret

	// use the newLogClient client factory when -log is specified
	if params.Log {
		clientFactory = newLogClient
	}

	// try to get access token
	if params.ServerType == secret.ServerPAS {
		accessToken, err = getAccessToken(params)
		if err != nil {
			fmt.Printf("Error in getting access token: %v\n", err)
			os.Exit(-3)
		}

	}
	// create a client handle to access secrets backend
	cl, err = secret.NewSecretClient(params.ServerPath, params.ServerType, accessToken, clientFactory)
	if err != nil {
		fmt.Printf("Error in setting up client: %v\n", err)
		os.Exit(-3)
	}

	if params.Debug {
		// turn of debug
		cl.SetDebug(true)
	}

	if params.UserAgent != "" {
		cl.SetUserAgent(params.UserAgent)
	}

	if params.ExtraHeadersMap != nil {
		cl.AddDefaultHeaders(params.ExtraHeadersMap)
	}

	switch params.Operation {
	case create:
		err = doCreate(cl, params)

	case createFolder:
		err = doCreateFolder(cl, params)

	case delete:
		err = doDelete(cl, params)

	case get:
		err = doGet(cl, params)

	case getMetaData:
		err = doGetMetaData(cl, params)

	case list:
		err = doList(cl, params)

	case modify:
		err = doModify(cl, params)
	}
	os.Exit(convertErrToExitStatus(err))
}

func doCreate(cl secret.Secret, params *Parameters) error {
	fmt.Printf("Creating secret of type %s in path [%s]\n", params.SecretType, params.SecretPath)
	var success bool
	var id string
	var r *http.Response
	var err error
	if params.SecretType == secret.SecretTypeKV {
		success, id, r, err = cl.Create(params.SecretPath, params.Description, params.KVSecret)
	} else {
		success, id, r, err = cl.Create(params.SecretPath, params.Description, params.TextValue)
	}
	if success {
		fmt.Printf("Secret created. ID: %s\n", id)
		return nil
	}

	if err != nil {
		fmt.Printf("Error in creating secret: %v\n", err)
	}
	if r != nil {
		fmt.Printf("HTTP response: %v\n", *r)
	}

	return err
}

func doCreateFolder(cl secret.Secret, params *Parameters) error {
	fmt.Printf("Creating secret folder in path [%s]\n", params.SecretPath)
	var success bool
	var id string
	var r *http.Response
	var err error

	success, id, r, err = cl.CreateFolder(params.SecretPath, params.Description)
	if success {
		fmt.Printf("Secret folder created. ID: %s\n", id)
		return nil
	}
	if err != nil {
		fmt.Printf("Error in creating secret folder: %v\n", err)
	}
	if r != nil {
		fmt.Printf("HTTP response: %v\n", *r)
	}

	return err
}
func doDelete(cl secret.Secret, params *Parameters) error {

	fmt.Printf("Deleting secret in path [%s]\n", params.SecretPath)
	r, reqError := cl.Delete(params.SecretPath)

	if reqError == nil {
		fmt.Println("Secret deleted")
		return nil
	}

	fmt.Printf("Error received: %v\n", reqError)
	if r != nil {
		fmt.Printf("HTTP response: %v\n", *r)
	}

	return reqError
}
func doGet(cl secret.Secret, params *Parameters) error {
	fmt.Printf("Getting secret from path [%s]\n", params.SecretPath)
	value, r, err := cl.Get(params.SecretPath)
	if err == nil {
		switch value.(type) {
		case string:
			fmt.Printf("Secret is a text string. Value: [%s]\n", value)
			return nil
		case map[string]string:
			fmt.Println("Secret is key value pair collection:")
			for k, v := range value.(map[string]string) {
				fmt.Printf("Key: %s\tValue:%v\n", k, v)
			}
			return nil
		default:
			fmt.Printf("Unknown secret type: %T\n", value)
			return secret.ErrUnexpectedResponse
		}
	}
	if r != nil {
		fmt.Printf("HTTP response: %v\n", r)
	}
	return err
}

func doGetMetaData(cl secret.Secret, params *Parameters) error {
	fmt.Printf("Getting secret metadata from path [%s]\n", params.SecretPath)
	value, r, err := cl.GetMetaData(params.SecretPath)
	if err == nil {
		fmt.Printf("Metadata returned:\nType: %s\n", value.Type)
		fmt.Printf("ID: %s\n", value.ID)
		fmt.Printf("CRN: %s\n", value.CRN)
		fmt.Printf("Created: %v\n", value.WhenCreated)
		if value.WhenModified.IsZero() {
			fmt.Println("Not modified")
		} else {
			fmt.Printf("Last modified time: %v\n", value.WhenModified)
		}
		return nil
	}
	fmt.Printf("Error in getting secret metadata: %v\n", err)
	if r != nil {
		fmt.Printf("HTTP response: %v\n", *r)
	}

	return err
}

func doList(cl secret.Secret, params *Parameters) error {
	fmt.Printf("Listing contents of [%s]\n", params.SecretPath)
	items, r, err := cl.List(params.SecretPath)
	if err != nil {
		fmt.Printf("Error in listing secrets: %v\n", err)
		if r != nil {
			fmt.Printf("Full HTTP response: %v\n", r)
		}
		return err
	}
	fmt.Printf("Number of items in folder: %d\n", len(items))
	for _, item := range items {
		fmt.Printf("ID: %s\tType: %s\tName: %s\n", item.ID, item.Type, item.Name)
	}
	return nil
}
func doModify(cl secret.Secret, params *Parameters) error {
	fmt.Printf("Modifying secret of type %s in path [%s]\n", params.SecretType, params.SecretPath)
	var success bool
	var id string
	var r *http.Response
	var err error
	if params.SecretType == secret.SecretTypeKV {
		success, id, r, err = cl.Modify(params.SecretPath, params.Description, params.KVSecret)
	} else {
		success, id, r, err = cl.Modify(params.SecretPath, params.Description, params.TextValue)
	}
	if success {
		fmt.Printf("Secret modified. ID: %s\n", id)
		return nil
	}
	if err != nil {
		fmt.Printf("Error in modifying secret: %v\n", err)
	}
	if r != nil {
		fmt.Printf("HTTP response: %v\n", *r)
	}
	return err
}
