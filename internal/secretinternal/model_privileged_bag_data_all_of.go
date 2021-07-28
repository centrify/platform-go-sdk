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

// PrivilegedBagDataAllOf struct for PrivilegedBagDataAllOf
type PrivilegedBagDataAllOf struct {
	// The privileged data in a secret bag. Required for vault and modify operations.
	Data *map[string]string `json:"data,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PrivilegedBagDataAllOf PrivilegedBagDataAllOf

// NewPrivilegedBagDataAllOf instantiates a new PrivilegedBagDataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivilegedBagDataAllOf() *PrivilegedBagDataAllOf {
	this := PrivilegedBagDataAllOf{}
	return &this
}

// NewPrivilegedBagDataAllOfWithDefaults instantiates a new PrivilegedBagDataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivilegedBagDataAllOfWithDefaults() *PrivilegedBagDataAllOf {
	this := PrivilegedBagDataAllOf{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *PrivilegedBagDataAllOf) GetData() map[string]string {
	if o == nil || o.Data == nil {
		var ret map[string]string
		return ret
	}
	return *o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivilegedBagDataAllOf) GetDataOk() (*map[string]string, bool) {
	if o == nil || o.Data == nil {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *PrivilegedBagDataAllOf) HasData() bool {
	if o != nil && o.Data != nil {
		return true
	}

	return false
}

// SetData gets a reference to the given map[string]string and assigns it to the Data field.
func (o *PrivilegedBagDataAllOf) SetData(v map[string]string) {
	o.Data = &v
}

func (o PrivilegedBagDataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Data != nil {
		toSerialize["data"] = o.Data
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *PrivilegedBagDataAllOf) UnmarshalJSON(bytes []byte) (err error) {
	varPrivilegedBagDataAllOf := _PrivilegedBagDataAllOf{}

	if err = json.Unmarshal(bytes, &varPrivilegedBagDataAllOf); err == nil {
		*o = PrivilegedBagDataAllOf(varPrivilegedBagDataAllOf)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "data")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePrivilegedBagDataAllOf struct {
	value *PrivilegedBagDataAllOf
	isSet bool
}

func (v NullablePrivilegedBagDataAllOf) Get() *PrivilegedBagDataAllOf {
	return v.value
}

func (v *NullablePrivilegedBagDataAllOf) Set(val *PrivilegedBagDataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivilegedBagDataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivilegedBagDataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivilegedBagDataAllOf(val *PrivilegedBagDataAllOf) *NullablePrivilegedBagDataAllOf {
	return &NullablePrivilegedBagDataAllOf{value: val, isSet: true}
}

func (v NullablePrivilegedBagDataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivilegedBagDataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


