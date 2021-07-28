package secret

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// HTTPClientFactory is a factory function that creates the http.Client object to use in the secret
// client.
type HTTPClientFactory func() *http.Client

// Secret is the collection of APIs that manage secrets stored in different secret stores.
// A secret storage implmentation must implement the functions defined here.
// Types PASSecretClient implements the methods specified in this interface.
type Secret interface {
	// Create creates a secret in 'path'. 'description' is an optional description
	// of the secret.  If 'value' is a string, it saves the secret as a
	// secret text string.  If 'value' is type map[string]string, the secret
	// is stored as 'keyvalue' secret.
	// Returns the following information:
	//  bool: whether the secret is created or not.
	//  id: a unique ID of the secret
	//  response: the actual HTTP response
	//
	// The following errors may be returned:
	//
	//	 ErrExists: Secret already exists
	//	 ErrNoCreatePermission: No permission to create secret.
	//	 ErrSecretTypeNotSupported:  Cannot create secret for the specified type.
	//	 ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
	//				technical support.
	Create(path string, description string, value interface{}) (bool, string, *http.Response, error)

	// CreateFolder creates a secret folder in 'path', with optional description.
	// Returns the following information:
	//  bool: whether the secret folder is created or not.
	//  id: a unique ID of the secret folder
	//  response: the actual HTTP response
	//
	// The following errors may be returned:
	//	 ErrExists: Secret folder already exists
	//	 ErrNoCreatePermission: No permission to create secret.
	//	 ErrSecretTypeNotSupported:  Cannot create secret for the specified type.
	//	 ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
	//				technical support.
	CreateFolder(path string, description string) (bool, string, *http.Response, error)

	// Delete deletes the folder/secret specified in 'path'
	// Returns the following information:
	//  response: the actual HTTP response
	// The following errors may be returned:
	//	ErrFolderNotEmpty: Folder is not empty
	//	ErrNoDeletePermission: No permission to delete secret/folder
	//	ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact technical support.
	Delete(path string) (*http.Response, error)

	// Get returns the secret content.
	// If the secret is a keyvalue secret, it returns the secret as map[string]string
	// If the secret is a text string, it returns the secret as string.
	// The following errors may be returned:
	//	 ErrSecretNotFound: Secret specified in path cannot be found.  It is possible that
	//					the caller may not have permission to read the secret.
	//	 ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
	//				technical support.
	Get(path string) (interface{}, *http.Response, error)

	// GetMetaData returns the metadata of a secret.
	// Returns the following information:
	//  MetaData: metadata information for the secret
	//  response: the actual HTTP response
	// The following errors may be returned
	//	ErrSecretNotFound: Secret specified in path cannot be found.  It is possible that
	//		the caller may not have permission to read the secret.
	//	ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
	//		technical support.
	GetMetaData(path string) (*MetaData, *http.Response, error)

	// List lists all secrets in a folder specified in 'path'
	// Returns the following:
	//  items: an array of Item. If the folder is empty, nil is returned.
	//  response: the actual HTTP response
	// the following errors may be returned:
	//	ErrFolderNotFound: Secret specified in path cannot be found.  It is possible that
	//		the caller may not have permission to access the folder.
	//	ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
	//		technical support.
	List(path string) ([]Item, *http.Response, error)

	// Modify modifies a secret in 'path'.
	// If 'description' is not an empty string, it replaces the current secret description.
	// If 'value' is a string, it saves the secret as a
	// secret text string.  If 'value' is type map[string]interface{} or map[string]string, the secret
	// is stored as 'keyvalue' secret.
	//
	// Returns the following information:
	//  bool: whether the secret is modified or not.
	//  id: a unique ID of the secret
	//  response: the actual HTTP response
	// The following errors may be returned:
	//	ErrNoModifyPermission: No permission to modify secret
	//	ErrNotSecretObject: specified path is not a secret
	//	ErrSecretNotFound: secret cannot be found
	//	ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
	//				technical support.
	Modify(path string, description string, value interface{}) (bool, string, *http.Response, error)

	// Additional functions for additional HTTP support

	// SetDebug enables/disables debug messages
	SetDebug(onoff bool)

	// AddDefualtHeaders add additional request headers to default HTTP header
	AddDefaultHeaders(hdrs map[string]string)

	// SetUserAgent sets UserAgent in HTTP header
	SetUserAgent(agent string)
}

// Item represents a secret that is returned in a List operation.
type Item struct {
	Name string // name of secret
	Type string // type of secret, can be KeyValue, Text, or Folder
	ID   string // unique ID of secret
}

// MetaData stores all metadata associated with a secret object that is returned in a GetMetaData operation.
type MetaData struct {
	Item
	CRN          string    // string that can be used in URL path
	WhenCreated  time.Time // creation time
	WhenModified time.Time // last modified time
}

// Common errors
var (
	ErrBadPathName              = errors.New("Invalid secret path name")
	ErrBadServerType            = errors.New("Bad server type")
	ErrCannotModifySecretType   = errors.New("Cannot change type of secret")
	ErrCannotModifySecretFolder = errors.New("Cannot modify a secret folder")
	ErrDeletedSecretExists      = errors.New("A mark-for-delete secret already exists in the same path")
	ErrExists                   = errors.New("Secret/folder already exists")
	ErrFolderNotEmpty           = errors.New("Folder is not empty")
	ErrFolderNotFound           = errors.New("Specified folder cannot be found")
	ErrNoCreatePermission       = errors.New("No permission to create secret")
	ErrNoDeletePermission       = errors.New("No permission to delete secret/folder")
	ErrNoGetMetaDataPermission  = errors.New("No permission to get ")
	ErrNoModifyPermission       = errors.New("No permission to modify secret")
	ErrNoRetrievePermission     = errors.New("No permission to retreive secret")
	ErrNotImplementedYet        = errors.New("Not implemented yet")
	ErrNotSecretObject          = errors.New("Specified path is not a secret")
	ErrNotSecretFolder          = errors.New("Specified path is not a secret folder")
	ErrSecretNotFound           = errors.New("Secret cannot be found")
	ErrSecretTypeNotSupported   = errors.New("Cannot created secret for input type")
	ErrUnexpectedResponse       = errors.New("Unexpected response from PAS")
)

// constant definition for server types
const (
	ServerPAS = "pas" // PAS
	ServerDSV = "dsv" // DSV
	ServerTSS = "tss" // TSS
)

// constant definition for secret types
const (
	SecretTypeFolder = "folder"
	SecretTypeText   = "text"
	SecretTypeKV     = "keyvalue"
)

// NewSecretClient creates a secret client to access secrets stored in 'server' of type 'serverType'.
// 'serverType' must be one of the followings:
//   pas - Centrify PAS
// You can specify the Oauth Token to use in 'accessToken'.
//
// If you need to use a different HTTP Client for the REST API call, you can specify a HTTPClientFactory
// function that returns a http.Client object.
//
func NewSecretClient(server string, serverType string, accessToken string, httpFactory HTTPClientFactory) (Secret, error) {

	var cl Secret
	// validate serverType
	sType := strings.TrimSpace(strings.ToLower(serverType))
	switch sType {
	case ServerPAS:
		cl = newPASSecretClient(server, accessToken, httpFactory)
	case ServerTSS:
		return nil, ErrNotImplementedYet
	case ServerDSV:
		return nil, ErrNotImplementedYet
	default:
		// unknown server type
		return nil, ErrBadServerType
	}
	return cl, nil
}
