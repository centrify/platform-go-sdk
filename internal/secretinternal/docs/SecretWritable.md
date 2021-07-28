# SecretWritable

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | [**Secrettypes**](Secrettypes.md) |  | 
**Name** | **string** | Names of secrets may contain a alphanumeric characters separated by / to indicate folder membership. Vaulting a bag or text secret within folders will implicitely create the folders. | 

## Methods

### NewSecretWritable

`func NewSecretWritable(type_ Secrettypes, name string, ) *SecretWritable`

NewSecretWritable instantiates a new SecretWritable object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSecretWritableWithDefaults

`func NewSecretWritableWithDefaults() *SecretWritable`

NewSecretWritableWithDefaults instantiates a new SecretWritable object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *SecretWritable) GetType() Secrettypes`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *SecretWritable) GetTypeOk() (*Secrettypes, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *SecretWritable) SetType(v Secrettypes)`

SetType sets Type field to given value.


### GetName

`func (o *SecretWritable) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SecretWritable) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SecretWritable) SetName(v string)`

SetName sets Name field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


