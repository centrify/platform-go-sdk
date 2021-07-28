# SecretTextPatchable

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **string** | The privileged data in a text secret. Required for vault and modify operations. | 

## Methods

### NewSecretTextPatchable

`func NewSecretTextPatchable(data string, ) *SecretTextPatchable`

NewSecretTextPatchable instantiates a new SecretTextPatchable object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSecretTextPatchableWithDefaults

`func NewSecretTextPatchableWithDefaults() *SecretTextPatchable`

NewSecretTextPatchableWithDefaults instantiates a new SecretTextPatchable object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *SecretTextPatchable) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *SecretTextPatchable) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *SecretTextPatchable) SetData(v string)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


