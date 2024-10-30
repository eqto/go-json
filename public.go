package json

import (
	"encoding/json"
	"os"
)

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func Parse(data []byte) (Object, error) {
	jo := Object{}
	if e := json.Unmarshal(data, &jo); e != nil {
		return nil, e
	}
	return jo, nil
}

func ParseObject(data interface{}) (Object, error) {
	data = sanitizeValue(data)
	if marshalled, e := json.Marshal(data); e == nil {
		return Parse(marshalled)
	} else {
		return nil, e
	}
}

func ParseString(data string) (Object, error) {
	return Parse([]byte(data))
}

func ParseArray(data []byte) ([]Object, error) {
	data = []byte(`{"data":` + string(data) + `}`)
	jo, e := Parse(data)
	if e != nil {
		return nil, e
	}
	return jo.GetArray(`data`), nil
}

func ParseFile(filename string) (Object, error) {
	data, e := os.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	return Parse(data)
}
