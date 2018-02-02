package json

import (
    "encoding/json"
    "strings"
    "errors"
    "strconv"
    "go/types"
)

/**
 * Created by tuxer on 9/6/17.
 */

//Object ...
type Object struct {
    dataMap map[string]interface{}
}

//ToFormattedBytes ...
func (j *Object) ToFormattedBytes() []byte {
    data, e := json.MarshalIndent(j.dataMap, ``, `  `)
    if e != nil {
        return nil
    }
    return data
}

//ToBytes ...
func (j *Object) ToBytes() []byte {
    if len(j.dataMap) == 0  {
        return []byte(`{}`)
    }
    if data, e := json.Marshal(j.dataMap); e == nil	{
		return data
	}
	return nil
}

//ToString ...
func (j *Object) ToString() string {
    data := j.ToBytes()
    str := `{}`
    if data != nil  {
        str = string(data)
    }
    return str
}

//GetDataMap ...
func (j *Object) GetDataMap() map[string]interface{}   {
    return j.dataMap
}

//SetDataMap ...
func (j *Object) SetDataMap(dataMap map[string]interface{}) {
	j.dataMap = dataMap
}

//GetInterface ...
func (j *Object) GetInterface(path string) interface{}	{
	return j.get(path)
}

//GetArray ...
func (j *Object) GetArray(path string) []Object    {
    obj := j.get(path)

    values, ok := obj.([]interface{})

    if !ok  {
        return nil
    }
    var arr []Object
    for _, value := range values   {
        jo := Object{ dataMap: value.(map[string]interface{}) }
        arr = append(arr, jo)
    }
    return arr
}

//GetJSONObject ...
func (j *Object) GetJSONObject(path string) *Object    {
    obj := j.get(path)

    if v, ok := obj.(map[string]interface{}); ok	{
        jo := Object{ dataMap: v }
        return &jo
	}
    return nil
}

//GetFloat ...
func (j *Object) GetFloat(path string) *float64 {
    obj := j.get(path)

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

//GetFloatD ...
func (j *Object) GetFloatD(path string, defValue float64) float64 {
    if val := j.GetFloat(path); val != nil	{
		return *val
	}
	return defValue
}

//GetInt ...
func (j *Object) GetInt(path string) *int {
    obj := j.get(path)

    switch obj.(type) {
    case float64:
        float, _ := obj.(float64)
        val := int(float)
        return &val
    case string:
        str, _ := obj.(string)
        val, e := strconv.Atoi(str)
        if e != nil {
            return nil
        }
        return &val
    default:
        return nil
    }
}

//GetIntD ...
func (j *Object) GetIntD(path string, defValue int) int {
    if val := j.GetInt(path); val != nil   {
        return *val
    }
	return defValue
}

//GetBoolean ...
func (j *Object) GetBoolean(path string) *bool {
    obj := j.get(path)
    if b, ok := obj.(bool); ok   {
        return &b
    }
	return nil
}

//GetBooleanD ...
func (j *Object) GetBooleanD(path string, defValue bool) bool {
	if val := j.GetBoolean(path); val != nil	{
		return *val
	}
	return defValue
}

//GetString ...
func (j *Object) GetString(path string) *string {
    obj := j.get(path)

    switch obj.(type) {
    case string:
        str, _ := obj.(string)
        return &str
    case float64:
        float, _ := obj.(float64)
        str := strconv.FormatFloat(float, 'f', -1, 64)
        return &str
    default:
        return nil
    }
}

//GetStringD ...
func (j *Object) GetStringD(path string, defValue string) string {
    if val := j.GetString(path); val != nil   {
        return *val
	}
	return defValue
}

//Put ...
func (j *Object) Put(path string, value interface{}) *Object    {
    j.putE(path, value)
    return j
}

func convertValue(value interface{}) interface{}	{
	//if pointer get the value
    if ptr, ok := value.(types.Pointer); ok	{
        value = ptr.Elem()
	}
	if arr, ok := value.([]Object); ok	{
        arrayMap := []interface{}{}
        for _, jo := range arr {
            arrayMap = append(arrayMap, convertValue(jo.dataMap))
		}
		return arrayMap
	} else if m, ok := value.(map[string]interface{}); ok	{
		for key, val := range m	{
			m[key] = convertValue(val)
		}
		return m
	} else if js, ok := value.(Object); ok	{
		return js.dataMap
	} else if byteData, ok := value.([]byte); ok	{
		return string(byteData)
	}
	return value
}

func (j *Object) putE(path string, value interface{}) error   {
	value = convertValue(value)

    if j.dataMap == nil  {
        j.dataMap = make(map[string]interface{})
    }

    rootMap := j.dataMap
    currentMap := rootMap

    splittedPath := strings.Split(path, `.`)
    for index, pathItem := range splittedPath   {
        if index < len(splittedPath) - 1    {
            _, ok := currentMap[pathItem]
            if !ok   {
                currentMap[pathItem] = make(map[string]interface{})
            }
            currentMap, ok = currentMap[pathItem].(map[string]interface{})
            if !ok  {
                return errors.New(pathItem + `is not a json object`)
            }
        } else {
            currentMap[pathItem] = value
        }
    }
    j.dataMap = rootMap
    return nil
}

func (j *Object) get(path string) interface{} {
    splittedPath := strings.Split(path, `.`)

    if j == nil {
        return nil
    }
    var jsonMap interface{}
    jsonMap = j.dataMap
    var val interface{}
    for _, pathItem := range splittedPath   {
        if jsonMap == nil   {
            return nil
        }
        val = jsonMap.(map[string]interface{})[pathItem]

        switch val.(type) {
        case map[string]interface{}:
            jsonMap = val
        case []interface{}:
            return val
        default:
            jsonMap = nil
        }
    }
    return val
}

//Parse ...
func Parse(data []byte) *Object    {
    jo := Object{}
    if e := json.Unmarshal(data, &jo.dataMap); e == nil	{
		return nil
	}
	return &jo
}