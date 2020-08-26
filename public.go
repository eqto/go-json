package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

//Marshal ...
func Marshal(obj interface{}) ([]byte, error) {
	if data, e := json.Marshal(obj); e == nil {
		return data, nil
	} else {
		return nil, e
	}
}

//Parse ...
func Parse(data []byte) (Object, error) {
	data = bytes.Trim(data, "\r\n\t ")
	jo := Object{}
	if e := json.Unmarshal(data, &jo); e != nil {
		return nil, e
	}
	return jo, nil
}

//ParseObject ...
func ParseObject(data interface{}) (Object, error) {
	if marshalled, e := json.Marshal(data); e == nil {
		return Parse(marshalled)
	} else {
		return nil, e
	}
}

//ParseString ...
func ParseString(data string) (Object, error) {
	return Parse([]byte(data))
}

//ParseArray ...
func ParseArray(data []byte) ([]Object, error) {
	data = []byte(`{"data":` + string(data) + `}`)
	jo, e := Parse(data)
	if e != nil {
		return nil, e
	}
	return jo.GetArray(`data`), nil
}

//ParseFile ...
func ParseFile(filename string) (Object, error) {
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	return Parse(data)
}
