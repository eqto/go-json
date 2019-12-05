package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

/**
 * Created by tuxer on 9/6/17.
 */

//Object ...
type Object map[string]interface{}

//ToFormattedBytes ...
func (j Object) ToFormattedBytes() []byte {
	data, e := json.MarshalIndent(j, ``, `  `)
	if e != nil {
		return nil
	}
	return data
}

//ToFormattedString ...
func (j Object) ToFormattedString() string {
	return string(j.ToFormattedBytes())
}

//ToBytes ...
func (j Object) ToBytes() []byte {
	if len(j) == 0 {
		return []byte(`{}`)
	}
	if data, e := json.Marshal(j); e == nil {
		return data
	}
	return nil
}

//ToString ...
func (j Object) ToString() string {
	data := j.ToBytes()
	str := `{}`
	if data != nil {
		str = string(data)
	}
	return str
}

//GetInterface ...
func (j Object) GetInterface(path string) interface{} {
	return j.Get(path)
}

//GetArray ...
func (j Object) GetArray(path string) []Object {
	obj := j.Get(path)
	if values, ok := obj.([]interface{}); ok {
		var arr []Object
		for _, value := range values {
			switch value.(type) {
			case Object:
				arr = append(arr, value.(Object))
			default:
				arr = append(arr, value.(map[string]interface{}))
			}
		}
		return arr
	}
	return nil
}

//GetIntArray ...
func (j Object) GetIntArray(path string) []int {
	obj := j.Get(path)
	if values, ok := obj.([]interface{}); ok {
		var arr []int
		for _, value := range values {
			arr = append(arr, int(value.(float64)))
		}
		return arr
	}
	return nil
}

//GetStringArray ...
func (j Object) GetStringArray(path string) []string {
	obj := j.Get(path)
	if values, ok := obj.([]interface{}); ok {
		var arr []string
		for _, value := range values {
			arr = append(arr, value.(string))
		}
		return arr
	}
	return nil
}

//GetJSONObject ...
func (j Object) GetJSONObject(path string) Object {
	obj := j.Get(path)

	if v, ok := obj.(map[string]interface{}); ok {
		return Object(v)
	}
	return nil
}

//GetFloatNull ...
func (j Object) GetFloatNull(path string) *float64 {
	obj := j.Get(path)

	switch obj.(type) {
	case float64:
		float, _ := obj.(float64)
		return &float
	case string:
		str, _ := obj.(string)
		val, e := strconv.ParseFloat(str, 64)
		if e != nil {
			return nil
		}
		return &val
	default:
		return nil
	}
}

//GetFloatOr ...
func (j Object) GetFloatOr(path string, defValue float64) float64 {
	if val := j.GetFloatNull(path); val != nil {
		return *val
	}
	return defValue
}

//GetFloat ...
func (j Object) GetFloat(path string) float64 {
	return j.GetFloatOr(path, 0)
}

//GetIntNull ...
func (j Object) GetIntNull(path string) *int {
	obj := j.Get(path)

	switch i := obj.(type) {
	case float64:
		float, _ := obj.(float64)
		val := int(float)
		return &val
	case string:
		str, _ := obj.(string)
		f, e := strconv.ParseFloat(str, 64)
		if e != nil {
			return nil
		}
		val := int(f)
		return &val
	case int:
		return &i
	default:
		return nil
	}
}

//GetIntOr ...
func (j Object) GetIntOr(path string, defValue int) int {
	if val := j.GetIntNull(path); val != nil {
		return *val
	}
	return defValue
}

//GetInt ...
func (j Object) GetInt(path string) int {
	return j.GetIntOr(path, 0)
}

//GetBooleanNull ...
func (j Object) GetBooleanNull(path string) *bool {
	obj := j.Get(path)
	if b, ok := obj.(bool); ok {
		return &b
	}
	return nil
}

//GetBooleanOr ...
func (j Object) GetBooleanOr(path string, defValue bool) bool {
	if val := j.GetBooleanNull(path); val != nil {
		return *val
	}
	return defValue
}

//GetBoolean ...
func (j Object) GetBoolean(path string) bool {
	return j.GetBooleanOr(path, false)
}

//GetStringNull ...
func (j Object) GetStringNull(path string) *string {
	obj := j.Get(path)

	str := ``
	switch obj := obj.(type) {
	case string:
		str = obj
	case float64:
		str = strconv.FormatFloat(obj, 'f', -1, 64)
	default:
		return nil
	}
	return &str
}

//GetStringOr ...
func (j Object) GetStringOr(path string, defValue string) string {
	if val := j.GetStringNull(path); val != nil {
		return *val
	}
	return defValue
}

//GetString ...
func (j Object) GetString(path string) string {
	return j.GetStringOr(path, ``)
}

//Put ...
func (j Object) Put(path string, value interface{}) Object {
	j.putE(path, value)
	return j
}

func convertValue(value interface{}) interface{} {
	if val := reflect.ValueOf(value); val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		value = val.Elem().Interface()
	}
	if arr, ok := value.([]Object); ok {
		arrayMap := []interface{}{}
		for _, jo := range arr {
			arrayMap = append(arrayMap, convertValue(jo))
		}
		return arrayMap
	} else if m, ok := value.(map[string]interface{}); ok {
		for key, val := range m {
			m[key] = convertValue(val)
		}
		return m
	} else if js, ok := value.(Object); ok {
		return js
	} else if byteData, ok := value.([]byte); ok {
		return string(byteData)
	} else if intData, ok := value.(int); ok {
		return float64(intData)
	}
	return value
}

func (j Object) putE(path string, value interface{}) error {
	value = convertValue(value)

	rootMap := j
	currentMap := rootMap

	splittedPath := strings.Split(path, `.`)
	for index, pathItem := range splittedPath {
		if index < len(splittedPath)-1 {
			if _, ok := currentMap[pathItem]; !ok {
				currentMap[pathItem] = make(map[string]interface{})
			}
			if curr, ok := currentMap[pathItem].(map[string]interface{}); !ok {
				return errors.New(pathItem + `is not a json object`)
			} else {
				currentMap = curr
			}
		} else {
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Map {
				mapVal := make(map[string]interface{})
				for _, key := range v.MapKeys() {
					strct := v.MapIndex(key)
					mapVal[key.Interface().(string)] = strct.Interface()
				}
				currentMap[pathItem] = mapVal
			} else {
				if m, ok := value.(map[string]interface{}); ok {
					currentMap[pathItem] = m
				} else {
					currentMap[pathItem] = value
				}
			}
		}
	}
	j = rootMap
	return nil
}

//Get ...
func (j Object) Get(path string) interface{} {
	splittedPath := strings.Split(path, `.`)

	var jsonMap Object
	jsonMap = j
	var val interface{}
	for _, pathItem := range splittedPath {
		if jsonMap == nil {
			return nil
		}
		val = jsonMap[pathItem]

		switch val := val.(type) {
		case Object:
			jsonMap = val
		case map[string]interface{}:
			jsonMap = Object(val)
		case []interface{}:
			return val
		default:
			jsonMap = nil
		}
	}
	return val
}

//Remove ...
func (j Object) Remove(path string) {
	index := strings.LastIndex(path, `.`)
	if index >= 0 {
		key := path[index+1:]
		path := path[0:index]
		val := j.Get(path)
		if val, ok := val.(map[string]interface{}); ok {
			delete(val, key)
			j.Put(path, val)
		}
	} else {
		delete(j, path)
	}
}

//Marshal ...
func Marshal(obj interface{}) ([]byte, error) {
	if data, e := json.Marshal(obj); e == nil {
		return data, nil
	} else {
		return nil, e
	}
}

//Parse ...
func Parse(data []byte) Object {
	data = bytes.Trim(data, "\r\n\t ")
	jo := Object{}
	if e := json.Unmarshal(data, &jo); e != nil {
		return nil
	}
	return jo
}

//ParseObject ...
func ParseObject(data interface{}) Object {
	if marshalled, e := json.Marshal(data); e == nil {
		return Parse(marshalled)
	}
	return nil
}

//ParseString ...
func ParseString(data string) Object {
	return Parse([]byte(data))
}

//ParseArray ...
func ParseArray(data []byte) []Object {
	data = []byte(`{"data":` + string(data) + `}`)
	jo := Parse(data)
	return jo.GetArray(`data`)
}

//ParseFile ...
func ParseFile(filename string) (Object, error) {
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	return Parse(data), nil
}
