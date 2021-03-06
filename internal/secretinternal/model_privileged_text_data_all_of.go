/*
 * Centrify Vault REST API
 *
 * Vault REST API specification 
 *
 * API version: 1.0
 * Contact: support@centrify.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package secretinternal

import (
	"encoding/json"
)

// PrivilegedTextDataAllOf struct for PrivilegedTextDataAllOf
type PrivilegedTextDataAllOf struct {
	// The privileged data in a text secret. Required for vault and modify operations.
	Data *string `json:"data,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PrivilegedTextDataAllOf PrivilegedTextDataAllOf

// NewPrivilegedTextDataAllOf instantiates a new PrivilegedTextDataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivilegedTextDataAllOf() *PrivilegedTextDataAllOf {
	this := PrivilegedTextDataAllOf{}
	return &this
}

// NewPrivilegedTextDataAllOfWithDefaults instantiates a new PrivilegedTextDataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivilegedTextDataAllOfWithDefaults() *PrivilegedTextDataAllOf {
	this := PrivilegedTextDataAllOf{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *PrivilegedTextDataAllOf) GetData() string {
	if o == nil || o.Data == nil {
		var ret string
		return ret
	}
	return *o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivilegedTextDataAllOf) GetDataOk() (*string, bool) {
	if o == nil || o.Data == nil {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *PrivilegedTextDataAllOf) HasData() bool {
	if o != nil && o.Data != nil {
		return true
	}

	return false
}

// SetData gets a reference to the given string and assigns it to the Data field.
func (o *PrivilegedTextDataAllOf) SetData(v string) {
	o.Data = &v
}

func (o PrivilegedTextDataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Data != nil {
		toSerialize["data"] = o.Data
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *PrivilegedTextDataAllOf) UnmarshalJSON(bytes []byte) (err error) {
	varPrivilegedTextDataAllOf := _PrivilegedTextDataAllOf{}

	if err = json.Unmarshal(bytes, &varPrivilegedTextDataAllOf); err == nil {
		*o = PrivilegedTextDataAllOf(varPrivilegedTextDataAllOf)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "data")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePrivilegedTextDataAllOf struct {
	value *PrivilegedTextDataAllOf
	isSet bool
}

func (v NullablePrivilegedTextDataAllOf) Get() *PrivilegedTextDataAllOf {
	return v.value
}

func (v *NullablePrivilegedTextDataAllOf) Set(val *PrivilegedTextDataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivilegedTextDataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivilegedTextDataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivilegedTextDataAllOf(val *PrivilegedTextDataAllOf) *NullablePrivilegedTextDataAllOf {
	return &NullablePrivilegedTextDataAllOf{value: val, isSet: true}
}

func (v NullablePrivilegedTextDataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivilegedTextDataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


