# ModelError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** | A reference URL which contains more information about the type of error  | [optional] [default to "about:blank"]
**Title** | Pointer to **string** | A short, summary of the problem type.  | [optional] 
**Status** | [**FailureStatusCode**](FailureStatusCode.md) |  | 
**Detail** | Pointer to **string** | A human readable explanation specific to this occurrence of the problem.  | [optional] 
**Instance** | Pointer to **string** | A URI to the requested resource  | [optional] 

## Methods

### NewModelError

`func NewModelError(status FailureStatusCode, ) *ModelError`

NewModelError instantiates a new ModelError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewModelErrorWithDefaults

`func NewModelErrorWithDefaults() *ModelError`

NewModelErrorWithDefaults instantiates a new ModelError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *ModelError) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ModelError) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ModelError) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ModelError) HasType() bool`

HasType returns a boolean if a field has been set.

### GetTitle

`func (o *ModelError) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *ModelError) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *ModelError) SetTitle(v string)`

SetTitle sets Title field to given value.

### HasTitle

`func (o *ModelError) HasTitle() bool`

HasTitle returns a boolean if a field has been set.

### GetStatus

`func (o *ModelError) GetStatus() FailureStatusCode`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ModelError) GetStatusOk() (*FailureStatusCode, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ModelError) SetStatus(v FailureStatusCode)`

SetStatus sets Status field to given value.


### GetDetail

`func (o *ModelError) GetDetail() string`

GetDetail returns the Detail field if non-nil, zero value otherwise.

### GetDetailOk

`func (o *ModelError) GetDetailOk() (*string, bool)`

GetDetailOk returns a tuple with the Detail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetail

`func (o *ModelError) SetDetail(v string)`

SetDetail sets Detail field to given value.

### HasDetail

`func (o *ModelError) HasDetail() bool`

HasDetail returns a boolean if a field has been set.

### GetInstance

`func (o *ModelError) GetInstance() string`

GetInstance returns the Instance field if non-nil, zero value otherwise.

### GetInstanceOk

`func (o *ModelError) GetInstanceOk() (*string, bool)`

GetInstanceOk returns a tuple with the Instance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstance

`func (o *ModelError) SetInstance(v string)`

SetInstance sets Instance field to given value.

### HasInstance

`func (o *ModelError) HasInstance() bool`

HasInstance returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


