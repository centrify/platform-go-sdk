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

// SecretListAllOf struct for SecretListAllOf
type SecretListAllOf struct {
	Items []SecretSparse `json:"items"`
	AdditionalProperties map[string]interface{}
}

type _SecretListAllOf SecretListAllOf

// NewSecretListAllOf instantiates a new SecretListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSecretListAllOf(items []SecretSparse) *SecretListAllOf {
	this := SecretListAllOf{}
	this.Items = items
	return &this
}

// NewSecretListAllOfWithDefaults instantiates a new SecretListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSecretListAllOfWithDefaults() *SecretListAllOf {
	this := SecretListAllOf{}
	return &this
}

// GetItems returns the Items field value
func (o *SecretListAllOf) GetItems() []SecretSparse {
	if o == nil {
		var ret []SecretSparse
		return ret
	}

	return o.Items
}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
func (o *SecretListAllOf) GetItemsOk() (*[]SecretSparse, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Items, true
}

// SetItems sets field value
func (o *SecretListAllOf) SetItems(v []SecretSparse) {
	o.Items = v
}

func (o SecretListAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["items"] = o.Items
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *SecretListAllOf) UnmarshalJSON(bytes []byte) (err error) {
	varSecretListAllOf := _SecretListAllOf{}

	if err = json.Unmarshal(bytes, &varSecretListAllOf); err == nil {
		*o = SecretListAllOf(varSecretListAllOf)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "items")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableSecretListAllOf struct {
	value *SecretListAllOf
	isSet bool
}

func (v NullableSecretListAllOf) Get() *SecretListAllOf {
	return v.value
}

func (v *NullableSecretListAllOf) Set(val *SecretListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableSecretListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableSecretListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSecretListAllOf(val *SecretListAllOf) *NullableSecretListAllOf {
	return &NullableSecretListAllOf{value: val, isSet: true}
}

func (v NullableSecretListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSecretListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


