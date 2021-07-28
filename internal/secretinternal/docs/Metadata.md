# Metadata

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique ID, generated upon creation | [optional] 
**Crn** | Pointer to **string** | Unique global reference, generated upon creation, useful for APIs that operate on multiple types of objects. | [optional] 
**Created** | Pointer to **time.Time** | A timestamp in ISO 8601 format of the date the address was created. | [optional] [readonly] 
**Modified** | Pointer to **time.Time** | A timestamp in ISO 8601 format of the date the address was last modified. | [optional] [readonly] 

## Methods

### NewMetadata

`func NewMetadata() *Metadata`

NewMetadata instantiates a new Metadata object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMetadataWithDefaults

`func NewMetadataWithDefaults() *Metadata`

NewMetadataWithDefaults instantiates a new Metadata object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Metadata) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Metadata) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Metadata) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Metadata) HasId() bool`

HasId returns a boolean if a field has been set.

### GetCrn

`func (o *Metadata) GetCrn() string`

GetCrn returns the Crn field if non-nil, zero value otherwise.

### GetCrnOk

`func (o *Metadata) GetCrnOk() (*string, bool)`

GetCrnOk returns a tuple with the Crn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCrn

`func (o *Metadata) SetCrn(v string)`

SetCrn sets Crn field to given value.

### HasCrn

`func (o *Metadata) HasCrn() bool`

HasCrn returns a boolean if a field has been set.

### GetCreated

`func (o *Metadata) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Metadata) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Metadata) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Metadata) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetModified

`func (o *Metadata) GetModified() time.Time`

GetModified returns the Modified field if non-nil, zero value otherwise.

### GetModifiedOk

`func (o *Metadata) GetModifiedOk() (*time.Time, bool)`

GetModifiedOk returns a tuple with the Modified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModified

`func (o *Metadata) SetModified(v time.Time)`

SetModified sets Modified field to given value.

### HasModified

`func (o *Metadata) HasModified() bool`

HasModified returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


