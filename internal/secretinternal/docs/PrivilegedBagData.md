# PrivilegedBagData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **map[string]string** | The privileged data in a secret bag. Required for vault and modify operations. | 

## Methods

### NewPrivilegedBagData

`func NewPrivilegedBagData(data map[string]string, ) *PrivilegedBagData`

NewPrivilegedBagData instantiates a new PrivilegedBagData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPrivilegedBagDataWithDefaults

`func NewPrivilegedBagDataWithDefaults() *PrivilegedBagData`

NewPrivilegedBagDataWithDefaults instantiates a new PrivilegedBagData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *PrivilegedBagData) GetData() map[string]string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *PrivilegedBagData) GetDataOk() (*map[string]string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *PrivilegedBagData) SetData(v map[string]string)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


