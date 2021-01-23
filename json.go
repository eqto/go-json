package json

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

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

//Has ...
func (j Object) Has(key string) bool {
	split := strings.Split(key, `.`)
	return getFromMap(j, split...) != nil
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

//Array ...
func (j Object) Array(path string) []interface{} {
	obj := j.Get(path)
	if values, ok := obj.([]interface{}); ok {
		return values
	}
	return nil
}

//GetArray ...
func (j Object) GetArray(path string) []Object {
	if objs := j.Array(path); objs != nil {
		var arr []Object
		for _, obj := range objs {
			switch obj := obj.(type) {
			case Object:
				arr = append(arr, obj)
			case map[string]interface{}:
				arr = append(arr, obj)
			default:
				return nil
			}
		}
		return arr
	}
	return nil
}

//GetIntArray ...
func (j Object) GetIntArray(path string) []int {
	if ints := j.Array(path); ints != nil {
		var arr []int
		for _, i := range ints {
			if i, ok := i.(int); ok {
				arr = append(arr, i)
			} else {
				return nil
			}
		}
		return arr
	}
	return nil
}

//GetStringArray ...
func (j Object) GetStringArray(path string) []string {
	if strs := j.Array(path); strs != nil {
		var arr []string
		for _, s := range strs {
			if s, ok := s.(string); ok {
				arr = append(arr, s)
			} else {
				return nil
			}
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
	} else if obj, ok := obj.(Object); ok {
		return obj
	}
	return nil
}

//GetFloatNull ...
func (j Object) GetFloatNull(path string) *float64 {
	obj := j.Get(path)

	switch val := obj.(type) {
	case float64:
		return &val
	case int:
		float := float64(val)
		return &float
	case uint:
		float := float64(val)
		return &float
	case string:
		float, e := strconv.ParseFloat(val, 64)
		if e != nil {
			return nil
		}
		return &float
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
	switch val := j.Get(path).(type) {
	case int:
		return &val
	case uint:
		intVal := int(val)
		return &intVal
	case float64:
		intVal := int(val)
		return &intVal
	case string:
		intVal, e := strconv.Atoi(val)
		if e != nil {
			return nil
		}
		return &intVal
	}
	return nil
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

	switch val := obj.(type) {
	case string:
		return &val
	case float64:
		str := strconv.FormatFloat(val, 'f', -1, 64)
		return &str
	case int:
		str := strconv.Itoa(val)
		return &str
	case uint:
		str := strconv.FormatUint(uint64(val), 10)
		return &str
	}
	return nil
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

//sanitizeValue if value is pointer, then get value from pointer, and convert value to recognizable value
func sanitizeValue(value interface{}) interface{} {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Ptr:
		return sanitizeValue(val.Elem().Interface())
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return uint(val.Uint())
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return int(val.Int())
	case reflect.Map:
		if val.Type().Key().Kind() != reflect.String {
			println(`unsupported map with key type not string`)
			return nil
		}
		m := make(map[string]interface{})
		iter := val.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()
			m[key.Interface().(string)] = sanitizeValue(val.Interface())
		}
		return m
	case reflect.Slice:
		if reflect.TypeOf(value).Elem().Kind() == reflect.Uint8 {
			return string(value.([]uint8))
		}
		length := val.Len()
		slice := make([]interface{}, length)
		for i := 0; i < length; i++ {
			slice[i] = sanitizeValue(val.Index(i).Interface())
		}
		return slice
	}
	return value
}

func (j Object) putE(path string, value interface{}) error {
	value = sanitizeValue(value)

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
			if m, ok := value.(map[string]interface{}); ok {
				for key, val := range m {
					m[key] = sanitizeValue(val)
				}
				currentMap[pathItem] = m
			} else {
				currentMap[pathItem] = value
			}
		}
	}
	j = rootMap
	return nil
}

//Get ...
func (j Object) Get(path string) interface{} {
	split := strings.Split(path, `.`)
	return getFromMap(j, split...)
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

//Clone ...
func (j Object) Clone() Object {
	cp := Object{}
	for k, v := range j {
		if v, ok := v.(Object); ok {
			cp[k] = v.Clone()
			continue
		}
		cp[k] = v
	}
	return cp
}
