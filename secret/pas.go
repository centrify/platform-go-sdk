package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/centrify/platform-go-sdk/internal/secretinternal"
)

// PASSecretClient implements the Secrets interface where the secret is stored in PAS
type PASSecretClient struct {
	apiClient   *secretinternal.APIClient
	httpClient  *http.Client // http client
	accessToken string       // access token
	tenantURL   string       // tenant URL
	debug       bool         // whether debug is on/off
}

// newPASSecretClient creates a new client handle for calling other functions in the secret package to
// access objects in 'tenantURL' with the OAuth 'accessToken'.
func newPASSecretClient(tenantURL string, accessToken string, httpFactory HTTPClientFactory) *PASSecretClient {
	// get clean tenant URL information

	// set up configuration
	cfg := secretinternal.NewConfiguration()
	if httpFactory != nil {
		cfg.HTTPClient = httpFactory()
	} else {
		cfg.HTTPClient = http.DefaultClient
	}

	// update tenantHost information in configuration
	host := cfg.Servers[0].Variables["tenantHost"]
	cleanHost := strings.TrimPrefix(tenantURL, "https://")
	cleanHost = strings.TrimPrefix(cleanHost, "http://")
	cleanHost = strings.TrimSuffix(cleanHost, "/")
	host.DefaultValue = cleanHost
	cfg.Servers[0].Variables["tenantHost"] = host

	// add default header for X-CENTRIFY-NATIVE-CLIENT
	cfg.AddDefaultHeader("X-CENTRIFY-NATIVE-CLIENT", "Yes")

	// add Oauth token
	cfg.AddDefaultHeader("Authorization", "Bearer "+accessToken)

	apiClient := secretinternal.NewAPIClient(cfg)
	cl := &PASSecretClient{
		apiClient:   apiClient,
		accessToken: accessToken,
		tenantURL:   tenantURL,
		httpClient:  cfg.HTTPClient,
	}
	return cl
}

// Get returns the secret content.
// If the secret is a keyvalue secret, it returns the secret as map[string]string
// If the secret is a text string, it returns the secret as string.
// The following errors may be returned:
//	 ErrSecretNotFound: Secret specified in path cannot be found.  It is possible that
//					the caller may not have permission to read the secret.
//	 ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
//				technical support.
func (c *PASSecretClient) Get(path string) (interface{}, *http.Response, error) {
	data, r, err := c.apiClient.SecretsApi.RetrieveExecute(c.apiClient.SecretsApi.Retrieve(context.Background(), path))
	if err != nil {
		if r != nil {
			// handle common error cases
			switch r.StatusCode {
			case 401: // 401 - unauthorized
				return nil, r, ErrNoRetrievePermission
			case 404: // 404 - not found
				return nil, r, ErrSecretNotFound
			default:
				return nil, r, ErrUnexpectedResponse
			}
		}
		return nil, r, err
	}

	switch data.Type {
	case secretinternal.TEXT:
		return data.AdditionalProperties["data"], r, nil
	case secretinternal.KEYVALUE:
		// convert each item into a string map
		results, ok := data.AdditionalProperties["data"].(map[string]interface{})
		if ok {
			res := make(map[string]string, len(results))
			for k, v := range results {
				str, ok := v.(string)
				if ok {
					res[k] = str
				}
			}
			return res, r, nil
		}
	}
	return nil, r, fmt.Errorf("unsupported data type: %v", data.Type)
}

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
//   ErrBadPathName: Invalid secret path name
//	 ErrExists: Secret already exists
//	 ErrNoCreatePermission: No permission to create secret.
//	 ErrSecretTypeNotSupported:  Cannot create secret for the specified type.
//	 ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
//				technical support.
func (c *PASSecretClient) Create(path string, description string, value interface{}) (bool, string, *http.Response, error) {

	var secretType secretinternal.Secrettypes

	switch value.(type) {
	case string:
		secretType = secretinternal.TEXT
	case map[string]string:
		secretType = secretinternal.KEYVALUE
	default:
		return false, "", nil, ErrSecretTypeNotSupported
	}

	req := c.apiClient.SecretsApi.SecretsCreate(context.Background())
	if secretType == secretinternal.TEXT {
		textSecret := secretinternal.NewSecretTextWritable(value.(string), secretType, path)
		req = req.SecretWritable(textSecret)
	} else {
		bagSecret := secretinternal.NewSecretBagWritable(value.(map[string]string), secretType, path)
		req = req.SecretWritable(bagSecret)
	}
	resp, r, err := req.Execute()
	if err == nil {
		if resp.Meta.Id != nil {
			return true, *resp.Meta.Id, r, nil
		}
		hasID, id := c.getIDFromObject(&resp.SecretWritable)
		if hasID {
			return true, id, r, nil
		}
		return false, "", r, ErrUnexpectedResponse
	}
	if r != nil {
		switch r.StatusCode {
		case 400: // bad request
			return false, "", r, ErrBadPathName
		case 401: // unauthorized access
			return false, "", r, ErrNoCreatePermission
		case 409: // conflict - object already exists
			return false, "", r, ErrExists
		case 500:
			// some request error is returned as 500 errors by backend
			// check for the ones that we know
			isAPIErr, summary := c.handleOpenAPIError(err)
			if isAPIErr {
				if summary == "A set must have a name" {
					return false, "", r, ErrBadPathName
				}
			}
			return false, "", r, ErrUnexpectedResponse
		default:
			return false, "", r, ErrUnexpectedResponse
		}
	}

	return false, "", r, err
}

// CreateFolder creates a secret folder in 'path', with optional description.
// Returns the following information:
//  bool: whether the secret folder is created or not.
//  id: a unique ID of the secret folder
//  response: the actual HTTP response
//
// The following errors may be returned:
//   ErrBadPathName: Invalid secret path name
//	 ErrExists: Secret folder already exists
//	 ErrNoCreatePermission: No permission to create secret.
//	 ErrSecretTypeNotSupported:  Cannot create secret for the specified type.
//	 ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
//				technical support.
func (c *PASSecretClient) CreateFolder(path string, description string) (bool, string, *http.Response, error) {

	secretType := secretinternal.FOLDER
	req := c.apiClient.SecretsApi.SecretsCreate(context.Background())
	writable := secretinternal.NewSecretFolderWritable(secretType, path)
	req = req.SecretWritable(writable)

	resp, r, err := req.Execute()
	if err != nil {
		if r != nil {
			// convert HTTP status code into specific error
			switch r.StatusCode {
			case 400: // bad request
				return false, "", r, ErrBadPathName
			case 409: // conflict
				return false, "", r, ErrExists
			case 401:
				return false, "", r, ErrNoCreatePermission
			case 500:
				// some request error is returned as 500 errors by backend
				// check for the ones that we know
				isAPIErr, summary := c.handleOpenAPIError(err)
				if isAPIErr {
					if summary == "A set must have a name" {
						return false, "", r, ErrBadPathName
					}
				}
				return false, "", r, ErrUnexpectedResponse

			default:
				return false, "", r, ErrUnexpectedResponse
			}
		}

		return false, "", r, err
	}

	if resp.Meta.Id != nil {
		return true, *resp.Meta.Id, r, err
	}
	// extract the ID from resposne
	hasID, id := c.getIDFromObject(&resp.SecretWritable)
	if hasID {
		return true, id, r, err
	}
	return false, "", r, ErrUnexpectedResponse
}

// List lists all secrets in a folder specified in 'path'
// Returns the following:
//  items: an array of Item. If the folder is empty, nil is returned.
//  response: the actual HTTP response
// the following errors may be returned:
//	ErrFolderNotFound: Secret specified in path cannot be found.  It is possible that
//		the caller may not have permission to access the folder.
//	ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
//		technical support.
func (c *PASSecretClient) List(path string) ([]Item, *http.Response, error) {
	req := c.apiClient.SecretsApi.Get(context.Background(), path)

	resp, r, err := c.apiClient.SecretsApi.GetExecute(req)
	if err != nil {
		if r != nil {
			// map HTTP status into specific error
			switch r.StatusCode {
			case 404: // not found
				return nil, r, ErrFolderNotFound
			default:
				return nil, r, ErrUnexpectedResponse
			}
		}

		return nil, r, err
	}
	res := resp.SecretWritable
	if res.Type != secretinternal.FOLDER {
		// object is not a folder
		return nil, r, ErrNotSecretFolder
	}
	contents, ok := res.AdditionalProperties["items"]
	if !ok {
		// expect backend to return at least 0 items
		return nil, r, fmt.Errorf("No items returned in list response: %w", ErrUnexpectedResponse)
	}

	items, ok := contents.([]interface{})
	if !ok {
		// unexpected type
		return nil, r, fmt.Errorf("Returned items is not type []interface{} as expected. Got type %T.  %w", contents, ErrUnexpectedResponse)
	}
	if c.debug {
		log.Printf("Number of items returned: %d\n", len(items))
	}
	retItems := make([]Item, len(items))

	for i, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			// unexpected type
			return nil, r, fmt.Errorf("List item [%v] has type %T which is not map[string]interface{}: %w", item, item, ErrUnexpectedResponse)
		}
		// verify existence of various fields
		if obj["name"] != nil {
			retItems[i].Name = obj["name"].(string)
		} else {
			return nil, r, fmt.Errorf("Name is specified in returned item [%v]: %w", obj, ErrUnexpectedResponse)
		}
		if obj["type"] != nil {
			retItems[i].Type = obj["type"].(string)
		} else {
			return nil, r, fmt.Errorf("Type is not specified in returned item [%v]: %w", obj, ErrUnexpectedResponse)
		}
		if obj["meta"] != nil {
			// check type
			meta, ok := obj["meta"].(map[string]interface{})
			if !ok {
				return nil, r, fmt.Errorf("Metadata has type %T which is not map[string]interface{} as expected: %w", obj["meta"], ErrUnexpectedResponse)
			}
			if meta["id"] != nil {
				retItems[i].ID = meta["id"].(string)
			} else {
				return nil, r, fmt.Errorf("ID is not specified in metadata [%v]: %w", meta, ErrUnexpectedResponse)
			}
		} else {
			return nil, r, fmt.Errorf("Metadata is not specified in returned item [%v]: %w", obj, ErrUnexpectedResponse)
		}
	}

	return retItems, r, nil
}

// Delete deletes the folder/secret specified in 'path'
// Returns the following information:
//  response: the actual HTTP response
// The following errors may be returned:
//	ErrFolderNotEmpty: Folder is not empty
//	ErrNoDeletePermission: No permission to delete secret/folder
//	ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact technical support.
func (c *PASSecretClient) Delete(path string) (*http.Response, error) {
	req := c.apiClient.SecretsApi.Delete(context.Background(), path)
	resp, err := c.apiClient.SecretsApi.DeleteExecute(req)
	if err == nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 401: // unauthorized
		return resp, ErrNoDeletePermission
	case 404: // Not found
		return resp, ErrSecretNotFound
	case 409: // conflict
		return resp, ErrFolderNotEmpty
	default:
		return resp, ErrUnexpectedResponse
	}
}

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
//	ErrCannotModifySecretFolder:  Modification of a secret folder is not supported.
//	ErrCannotModifySecretType:  Modification of keyvalue secret to text or vice versa is not supported.
//	ErrNoModifyPermission: No permission to modify secret
//	ErrNotSecretObject: specified path is not a secret
//	ErrSecretNotFound: secret cannot be found
//	ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
//				technical support.
func (c *PASSecretClient) Modify(path string, description string, value interface{}) (bool, string, *http.Response, error) {

	var secretType secretinternal.Secrettypes

	switch value.(type) {
	case string:
		secretType = secretinternal.TEXT
	case map[string]string:
		secretType = secretinternal.KEYVALUE
	default:
		return false, "", nil, ErrSecretTypeNotSupported
	}

	req := c.apiClient.SecretsApi.Modify(context.Background(), path)
	if secretType == secretinternal.TEXT {
		textSecret := secretinternal.NewSecretTextPatchable(value.(string), secretType)
		req = req.SecretPatchable(textSecret)
	} else {
		bagSecret := secretinternal.NewSecretBagPatchable(value.(map[string]string), secretType)
		req = req.SecretPatchable(bagSecret)
	}
	resp, r, err := req.Execute()
	if err == nil {
		return true, *resp.Meta.Id, r, nil
	}
	// map error status code to error
	if r != nil {

		switch r.StatusCode {
		case 401: // unauthorized
			return false, "", r, ErrNoModifyPermission
		case 404: // not found, or it may be that the path specifies a folder
			metadata, _, err2 := c.GetMetaData(path)
			if err2 != nil {
				// cannot get metadata for the object
				return false, "", r, ErrSecretNotFound
			}
			if metadata.Type == SecretTypeFolder {
				// it's a folder which we cannot modify
				return false, "", r, ErrCannotModifySecretFolder
			}
			// in this case, there is no object of the new content type
			// which means that we cannot change the type of secret
			return false, "", r, ErrCannotModifySecretType
		case 409: // conflict...e.g., cannot convert secret type
			return false, "", r, ErrCannotModifySecretType
		}
	}
	return false, "", r, err
}

// GetMetaData returns the metadata of a secret.
// Returns the following information:
//  MetaData: metadata information for the secret
//  response: the actual HTTP response
// The following errors may be returned:
//	ErrNoGetMetaDataPermission:  The caller has no permission to get metadata information
//	ErrSecretNotFound: Secret specified in path cannot be found.  It is possible that
//		the caller may not have permission to view the secret.
//	ErrUnexpectedResponse:  The response for the REST API is not expected.  Please contact
//		technical support.
func (c *PASSecretClient) GetMetaData(path string) (*MetaData, *http.Response, error) {
	data, r, err := c.apiClient.SecretsApi.GetExecute(c.apiClient.SecretsApi.Get(context.Background(), path))
	if err != nil {
		// error
		if r != nil {
			// convert HTTP status code to error
			switch r.StatusCode {
			case 401: // unauthorized
				return nil, r, ErrNoGetMetaDataPermission
			case 404: // not found
				return nil, r, ErrSecretNotFound
			}
		}
		return nil, r, err
	}

	result := &MetaData{}
	result.Name = data.Name
	result.Type = string(data.Type)
	if data.Meta.Id == nil {
		// error...should never be nil for object found
		return nil, r, fmt.Errorf("Id should never be empty: %w", ErrUnexpectedResponse)
	}
	result.ID = *data.Meta.Id
	if data.Meta.Crn == nil {
		// error... should never be nil for object found
		return nil, r, fmt.Errorf("CRN should never be empty:%w", ErrUnexpectedResponse)
	}
	result.CRN = *data.Meta.Crn
	if data.Meta.Created != nil {
		result.WhenCreated = *data.Meta.Created
	}
	if data.Meta.Modified != nil {
		result.WhenModified = *data.Meta.Modified
	}

	// extract information about secret metadata into result
	return result, r, nil
}

// SetDebug enables/disables debug messages.  For PASSecretClient, it dumps the
// HTTP request and response to the standard logger.
// DO NOT enable debugging in production environment as the full HTTP request
// and response that may contain secret information are logged.
func (c *PASSecretClient) SetDebug(onoff bool) {
	c.apiClient.SetDebug(onoff)
	c.debug = onoff
}

// SetUserAgent sets UserAgent in HTTP header
func (c *PASSecretClient) SetUserAgent(agent string) {
	c.apiClient.SetUserAgent(agent)
}

// AddDefaultHeaders add extra headers to default HTTP request header
func (c *PASSecretClient) AddDefaultHeaders(hdrs map[string]string) {
	c.apiClient.AddDefaultHeaders(hdrs)
}

// getIDFromObject returns the ID that is returned in a secretinternal.SecretWritable object
// in response
func (c *PASSecretClient) getIDFromObject(obj *secretinternal.SecretWritable) (bool, string) {
	if obj == nil {
		return false, ""
	}
	meta := obj.AdditionalProperties["meta"]
	if meta == nil {
		// no meta data
		return false, ""
	}
	attribs, ok := meta.(map[string]interface{})
	if !ok {
		// wrong type
		return false, ""
	}
	id := attribs["id"]
	if id == nil {
		return false, ""
	}
	return true, id.(string)
}

// handleOpenAPIError checks if the error is a GenericOpenAPIError.
// It returns true if it is and the associated "title" field in the error response.
// Otherwise it returns false.
// This is required as PAS returns "500 Server Error" for some use cases where
// it should be an error with the request itself.
func (c *PASSecretClient) handleOpenAPIError(err error) (bool, string) {
	// check if it is openAPI error
	openAPIErr, ok := err.(secretinternal.GenericOpenAPIError)
	if !ok {
		return false, "" // not OpenAPI error..cannot return error information
	}

	var errDetails map[string]interface{}

	marshalErr := json.Unmarshal(openAPIErr.Body(), &errDetails)
	if marshalErr != nil {
		return false, "" // cannot unmarshal body...not really OpenAPI error
	}
	if title, ok := errDetails["title"]; ok {
		if tstr, ok := title.(string); ok {
			return true, tstr
		}
	}

	// cannot get error summary
	return false, ""
}
