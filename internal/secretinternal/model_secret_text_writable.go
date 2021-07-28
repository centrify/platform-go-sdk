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
	"reflect"
	"strings"
)

// SecretTextWritable struct for SecretTextWritable
type SecretTextWritable struct {
	SecretWritable
	// The privileged data in a text secret. Required for vault and modify operations.
	Data string `json:"data"`
	AdditionalProperties map[string]interface{}
}

type _SecretTextWritable SecretTextWritable

// NewSecretTextWritable instantiates a new SecretTextWritable object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSecretTextWritable(data string, type_ Secrettypes, name string) *SecretTextWritable {
	this := SecretTextWritable{}
	this.Type = type_
	this.Name = name
	this.Data = data
	return &this
}

// NewSecretTextWritableWithDefaults instantiates a new SecretTextWritable object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSecretTextWritableWithDefaults() *SecretTextWritable {
	this := SecretTextWritable{}
	return &this
}

// GetData returns the Data field value
func (o *SecretTextWritable) GetData() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *SecretTextWritable) GetDataOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *SecretTextWritable) SetData(v string) {
	o.Data = v
}

func (o SecretTextWritable) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	serializedSecretWritable, errSecretWritable := json.Marshal(o.SecretWritable)
	if errSecretWritable != nil {
		return []byte{}, errSecretWritable
	}
	errSecretWritable = json.Unmarshal([]byte(serializedSecretWritable), &toSerialize)
	if errSecretWritable != nil {
		return []byte{}, errSecretWritable
	}
	if true {
		toSerialize["data"] = o.Data
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *SecretTextWritable) UnmarshalJSON(bytes []byte) (err error) {
	type SecretTextWritableWithoutEmbeddedStruct struct {
		// The privileged data in a text secret. Required for vault and modify operations.
		Data string `json:"data"`
	}

	varSecretTextWritableWithoutEmbeddedStruct := SecretTextWritableWithoutEmbeddedStruct{}

	err = json.Unmarshal(bytes, &varSecretTextWritableWithoutEmbeddedStruct)
	if err == nil {
		varSecretTextWritable := _SecretTextWritable{}
		varSecretTextWritable.Data = varSecretTextWritableWithoutEmbeddedStruct.Data
		*o = SecretTextWritable(varSecretTextWritable)
	} else {
		return err
	}

	varSecretTextWritable := _SecretTextWritable{}

	err = json.Unmarshal(bytes, &varSecretTextWritable)
	if err == nil {
		o.SecretWritable = varSecretTextWritable.SecretWritable
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "data")

		// remove fields from embedded structs
		reflectSecretWritable := reflect.ValueOf(o.SecretWritable)
		for i := 0; i < reflectSecretWritable.Type().NumField(); i++ {
			t := reflectSecretWritable.Type().Field(i)

			if jsonTag := t.Tag.Get("json"); jsonTag != "" {
				fieldName := ""
				if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
					fieldName = jsonTag[:commaIdx]
				} else {
					fieldName = jsonTag
				}
				if fieldName != "AdditionalProperties" {
					delete(additionalProperties, fieldName)
				}
			}
		}

		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableSecretTextWritable struct {
	value *SecretTextWritable
	isSet bool
}

func (v NullableSecretTextWritable) Get() *SecretTextWritable {
	return v.value
}

func (v *NullableSecretTextWritable) Set(val *SecretTextWritable) {
	v.value = val
	v.isSet = true
}

func (v NullableSecretTextWritable) IsSet() bool {
	return v.isSet
}

func (v *NullableSecretTextWritable) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSecretTextWritable(val *SecretTextWritable) *NullableSecretTextWritable {
	return &NullableSecretTextWritable{value: val, isSet: true}
}

func (v NullableSecretTextWritable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSecretTextWritable) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


