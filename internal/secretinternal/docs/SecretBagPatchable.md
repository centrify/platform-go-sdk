# SecretBagPatchable

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **map[string]string** | The privileged data in a secret bag. Required for vault and modify operations. | 

## Methods

### NewSecretBagPatchable

`func NewSecretBagPatchable(data map[string]string, ) *SecretBagPatchable`

NewSecretBagPatchable instantiates a new SecretBagPatchable object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSecretBagPatchableWithDefaults

`func NewSecretBagPatchableWithDefaults() *SecretBagPatchable`

NewSecretBagPatchableWithDefaults instantiates a new SecretBagPatchable object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *SecretBagPatchable) GetData() map[string]string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *SecretBagPatchable) GetDataOk() (*map[string]string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *SecretBagPatchable) SetData(v map[string]string)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


