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

// PrivilegedData struct for PrivilegedData
type PrivilegedData struct {
	Type Secrettypes `json:"type"`
	AdditionalProperties map[string]interface{}
}

type _PrivilegedData PrivilegedData

// NewPrivilegedData instantiates a new PrivilegedData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivilegedData(type_ Secrettypes) *PrivilegedData {
	this := PrivilegedData{}
	this.Type = type_
	return &this
}

// NewPrivilegedDataWithDefaults instantiates a new PrivilegedData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivilegedDataWithDefaults() *PrivilegedData {
	this := PrivilegedData{}
	return &this
}

// GetType returns the Type field value
func (o *PrivilegedData) GetType() Secrettypes {
	if o == nil {
		var ret Secrettypes
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *PrivilegedData) GetTypeOk() (*Secrettypes, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *PrivilegedData) SetType(v Secrettypes) {
	o.Type = v
}

func (o PrivilegedData) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["type"] = o.Type
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *PrivilegedData) UnmarshalJSON(bytes []byte) (err error) {
	varPrivilegedData := _PrivilegedData{}

	if err = json.Unmarshal(bytes, &varPrivilegedData); err == nil {
		*o = PrivilegedData(varPrivilegedData)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "type")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePrivilegedData struct {
	value *PrivilegedData
	isSet bool
}

func (v NullablePrivilegedData) Get() *PrivilegedData {
	return v.value
}

func (v *NullablePrivilegedData) Set(val *PrivilegedData) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivilegedData) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivilegedData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivilegedData(val *PrivilegedData) *NullablePrivilegedData {
	return &NullablePrivilegedData{value: val, isSet: true}
}

func (v NullablePrivilegedData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivilegedData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


