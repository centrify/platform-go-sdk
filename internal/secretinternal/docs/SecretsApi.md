# \SecretsApi

All URIs are relative to *https://api.my.centrify.net/api/v1.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Delete**](SecretsApi.md#Delete) | **Delete** /secrets/{nameOrId} | Delete a secret
[**Get**](SecretsApi.md#Get) | **Get** /secrets/{nameOrId} | Get secrets
[**Modify**](SecretsApi.md#Modify) | **Patch** /secrets/{nameOrId} | Modify a secret
[**Retrieve**](SecretsApi.md#Retrieve) | **Get** /privilegeddata/secrets/{nameOrId} | Retrieve privileged data
[**SecretsCreate**](SecretsApi.md#SecretsCreate) | **Post** /secrets | Create a secret
[**SecretsList**](SecretsApi.md#SecretsList) | **Get** /secrets | List secrets



## Delete

> Delete(ctx, nameOrId).Execute()

Delete a secret



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    nameOrId := "passwords/mine" // string | The name or id of a secret. Note a name can be a path, and contained / characters should not be url encoded.

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.SecretsApi.Delete(context.Background(), nameOrId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SecretsApi.Delete``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nameOrId** | **string** | The name or id of a secret. Note a name can be a path, and contained / characters should not be url encoded. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Get

> SecretDense Get(ctx, nameOrId).Execute()

Get secrets



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    nameOrId := "passwords/mine" // string | The name or id of a secret. Note a name can be a path, and contained / characters should not be url encoded.

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.SecretsApi.Get(context.Background(), nameOrId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SecretsApi.Get``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Get`: SecretDense
    fmt.Fprintf(os.Stdout, "Response from `SecretsApi.Get`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nameOrId** | **string** | The name or id of a secret. Note a name can be a path, and contained / characters should not be url encoded. | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SecretDense**](SecretDense.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Modify

> SecretDense Modify(ctx, nameOrId).SecretPatchable(secretPatchable).Execute()

Modify a secret



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    nameOrId := "passwords/mine" // string | The name or id of a secret. Note a name can be a path, and contained / characters should not be url encoded.
    secretPatchable := *openapiclient.NewSecretPatchable(openapiclient.secrettypes("text")) // SecretPatchable | A modify operation will update only the properties included in the request body. The request body for a text and bag secret can differ. See examples.

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.SecretsApi.Modify(context.Background(), nameOrId).SecretPatchable(secretPatchable).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SecretsApi.Modify``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Modify`: SecretDense
    fmt.Fprintf(os.Stdout, "Response from `SecretsApi.Modify`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nameOrId** | **string** | The name or id of a secret. Note a name can be a path, and contained / characters should not be url encoded. | 

### Other Parameters

Other parameters are passed through a pointer to a apiModifyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **secretPatchable** | [**SecretPatchable**](SecretPatchable.md) | A modify operation will update only the properties included in the request body. The request body for a text and bag secret can differ. See examples. | 

### Return type

[**SecretDense**](SecretDense.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Retrieve

> PrivilegedData Retrieve(ctx, nameOrId).Execute()

Retrieve privileged data



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    nameOrId := "passwords/mine" // string | The name or id of a secret. Note a name can be a path, and contained / characters should not be url encoded.

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.SecretsApi.Retrieve(context.Background(), nameOrId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SecretsApi.Retrieve``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Retrieve`: PrivilegedData
    fmt.Fprintf(os.Stdout, "Response from `SecretsApi.Retrieve`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nameOrId** | **string** | The name or id of a secret. Note a name can be a path, and contained / characters should not be url encoded. | 

### Other Parameters

Other parameters are passed through a pointer to a apiRetrieveRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**PrivilegedData**](PrivilegedData.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SecretsCreate

> SecretDense SecretsCreate(ctx).SecretWritable(secretWritable).Execute()

Create a secret



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    secretWritable := *openapiclient.NewSecretWritable(openapiclient.secrettypes("text"), "passwords/my_password") // SecretWritable | 

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.SecretsApi.SecretsCreate(context.Background()).SecretWritable(secretWritable).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SecretsApi.SecretsCreate``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SecretsCreate`: SecretDense
    fmt.Fprintf(os.Stdout, "Response from `SecretsApi.SecretsCreate`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSecretsCreateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **secretWritable** | [**SecretWritable**](SecretWritable.md) |  | 

### Return type

[**SecretDense**](SecretDense.md)

### Authorization

[OAuth2](../README.md#OAuth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SecretsList

> SecretList SecretsList(ctx).Limit(limit).OrderBy(orderBy).Search(search).Filter(filter).Execute()

List secrets



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    limit := int32(5) // int32 | How many results to return. (optional) (default to 10)
    orderBy := []string{"name desc"} // []string | A comma separated list of properties to sort by.  (optional)
    search := "john" // string | Provide search text to use default search capabilities.  For more advanced filtering capabilities, use the filter parameter.  (optional)
    filter := "activeCheckouts gt 0" // string | Conditional filtering of a list (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.SecretsApi.SecretsList(context.Background()).Limit(limit).OrderBy(orderBy).Search(search).Filter(filter).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SecretsApi.SecretsList``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SecretsList`: SecretList
    fmt.Fprintf(os.Stdout, "Response from `SecretsApi.SecretsList`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSecretsListRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | How many results to return. | [default to 10]
 **orderBy** | **[]string** | A comma separated list of properties to sort by.  | 
 **search** | **string** | Provide search text to use default search capabilities.  For more advanced filtering capabilities, use the filter parameter.  | 
 **filter** | **string** | Conditional filtering of a list | 

### Return type

[**SecretList**](SecretList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

