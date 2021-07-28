# SecretList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | [**[]SecretSparse**](SecretSparse.md) |  | 
**Object** | **string** | What kind of resource does this list contain? | 
**NextUrl** | **NullableString** | Url of next page of items in list. | 
**PreviousUrl** | **NullableString** | Url of previous page of items in list. | 

## Methods

### NewSecretList

`func NewSecretList(items []SecretSparse, object string, nextUrl NullableString, previousUrl NullableString, ) *SecretList`

NewSecretList instantiates a new SecretList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSecretListWithDefaults

`func NewSecretListWithDefaults() *SecretList`

NewSecretListWithDefaults instantiates a new SecretList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *SecretList) GetItems() []SecretSparse`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *SecretList) GetItemsOk() (*[]SecretSparse, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *SecretList) SetItems(v []SecretSparse)`

SetItems sets Items field to given value.


### GetObject

`func (o *SecretList) GetObject() string`

GetObject returns the Object field if non-nil, zero value otherwise.

### GetObjectOk

`func (o *SecretList) GetObjectOk() (*string, bool)`

GetObjectOk returns a tuple with the Object field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObject

`func (o *SecretList) SetObject(v string)`

SetObject sets Object field to given value.


### GetNextUrl

`func (o *SecretList) GetNextUrl() string`

GetNextUrl returns the NextUrl field if non-nil, zero value otherwise.

### GetNextUrlOk

`func (o *SecretList) GetNextUrlOk() (*string, bool)`

GetNextUrlOk returns a tuple with the NextUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextUrl

`func (o *SecretList) SetNextUrl(v string)`

SetNextUrl sets NextUrl field to given value.


### SetNextUrlNil

`func (o *SecretList) SetNextUrlNil(b bool)`

 SetNextUrlNil sets the value for NextUrl to be an explicit nil

### UnsetNextUrl
`func (o *SecretList) UnsetNextUrl()`

UnsetNextUrl ensures that no value is present for NextUrl, not even an explicit nil
### GetPreviousUrl

`func (o *SecretList) GetPreviousUrl() string`

GetPreviousUrl returns the PreviousUrl field if non-nil, zero value otherwise.

### GetPreviousUrlOk

`func (o *SecretList) GetPreviousUrlOk() (*string, bool)`

GetPreviousUrlOk returns a tuple with the PreviousUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreviousUrl

`func (o *SecretList) SetPreviousUrl(v string)`

SetPreviousUrl sets PreviousUrl field to given value.


### SetPreviousUrlNil

`func (o *SecretList) SetPreviousUrlNil(b bool)`

 SetPreviousUrlNil sets the value for PreviousUrl to be an explicit nil

### UnsetPreviousUrl
`func (o *SecretList) UnsetPreviousUrl()`

UnsetPreviousUrl ensures that no value is present for PreviousUrl, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


