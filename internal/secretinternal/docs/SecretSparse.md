# SecretSparse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique ID, generated upon creation | 
**Name** | **string** | Names of secrets may contain a alphanumeric characters separated by / to indicate folder membership. Vaulting a bag or text secret within folders will implicitely create the folders. | 

## Methods

### NewSecretSparse

`func NewSecretSparse(id string, name string, ) *SecretSparse`

NewSecretSparse instantiates a new SecretSparse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSecretSparseWithDefaults

`func NewSecretSparseWithDefaults() *SecretSparse`

NewSecretSparseWithDefaults instantiates a new SecretSparse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SecretSparse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SecretSparse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SecretSparse) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *SecretSparse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SecretSparse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SecretSparse) SetName(v string)`

SetName sets Name field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


