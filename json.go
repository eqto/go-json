package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

// Object represents a JSON object as a map with string keys and any type values.
// It provides convenient methods for accessing and manipulating JSON data.
type Object map[string]interface{}

// Set adds or updates a key-value pair in the JSON object.
// Parameters:
//   - key: the string key to set in the object
//   - val: the value to associate with the key (can be any type)
func (o Object) Set(key string, val interface{}) {
	o[key] = val
}

// Get retrieves a value from the JSON object by key.
// Parameters:
//   - key: the string key to look up in the object
//
// Returns:
//   - the value associated with the key, or nil if the key doesn't exist
func (o Object) Get(key string) interface{} {
	if v, ok := o[key]; ok {
		return v
	}
	return nil
}

// GetObject retrieves a nested JSON object from the current object by key.
// It handles both Object type and map[string]any type conversions.
// Parameters:
//   - key: the string key to look up in the object
//
// Returns:
//   - the nested Object if found and convertible, or an empty Object if not found or not convertible
func (o Object) GetObject(key string) Object {
	if v, ok := o[key]; ok {
		switch v := v.(type) {
		case Object:
			return v
		case map[string]interface{}:
			return Object(v)
		}
	}
	return Object{}
}

// Has checks if a key exists in the JSON object.
// Parameters:
//   - key: the string key to check for existence
//
// Returns:
//   - true if the key exists in the object, false otherwise
func (o Object) Has(key string) bool {
	_, ok := o[key]
	return ok
}

// GetInt retrieves an integer value from the JSON object by key.
// It attempts to convert various types (float64, string, bool) to int.
// For strings, it uses strconv.Atoi for conversion.
// For booleans, true returns 1, false returns 0.
// Parameters:
//   - key: the string key to look up in the object
//
// Returns:
//   - the integer value if convertible, or 0 if not found or not convertible
func (o Object) GetInt(key string) int {
	if v, ok := o[key]; ok {
		switch v := v.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, e := strconv.Atoi(v); e == nil {
				return i
			}
		case bool:
			if v {
				return 1
			}
			return 0
		}
	}
	return 0
}

// GetFloat retrieves a float64 value from the JSON object by key.
// It attempts to convert various types (float64, string) to float64.
// For strings, it uses strconv.ParseFloat for conversion.
// Parameters:
//   - key: the string key to look up in the object
//
// Returns:
//   - the float64 value if convertible, or 0 if not found or not convertible
func (o Object) GetFloat(key string) float64 {
	if v, ok := o[key]; ok {
		switch v := v.(type) {
		case float64:
			return v
		case string:
			if f, e := strconv.ParseFloat(v, 64); e == nil {
				return f
			}
		}
	}
	return 0
}

// GetString retrieves a string value from the JSON object by key.
// It converts various types (string, float64, bool) to string representation.
// For float64, it formats without decimal places.
// For bool, it returns "true" or "false".
// Parameters:
//   - key: the string key to look up in the object
//
// Returns:
//   - the string value if convertible, or empty string if not found or not convertible
func (o Object) GetString(key string) string {
	if v, ok := o[key]; ok {
		switch v := v.(type) {
		case string:
			return v
		case float64:
			return strconv.FormatFloat(v, 'f', 0, 64)
		case bool:
			if v {
				return `true`
			}
			return `false`
		}
	}
	return ``
}

// GetArray retrieves an array of any type from the JSON object by key.
// It handles JSON arrays and returns them as []any.
// Parameters:
//   - key: the string key to look up in the object
//
// Returns:
//   - a slice of any if the value is an array, or empty slice if not found or not convertible
func (o Object) GetArray(key string) []interface{} {
	if v, ok := o[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			return arr
		}
	}
	return []interface{}{}
}

// GetObjectArray retrieves an array of JSON objects from the JSON object by key.
// It handles JSON arrays where each element should be a JSON object (map[string]any).
// Parameters:
//   - key: the string key to look up in the object
//
// Returns:
//   - a slice of Object if the value is an array of objects, or empty slice if not found or not convertible
func (o Object) GetObjectArray(key string) []Object {
	if v, ok := o[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			result := make([]Object, 0, len(arr))
			for _, item := range arr {
				if obj, ok := item.(map[string]interface{}); ok {
					result = append(result, Object(obj))
				}
			}
			return result
		}
	}
	return []Object{}
}

// GetIntArray retrieves an array of integers from the JSON object by key.
// It handles JSON arrays where numbers are typically parsed as float64.
// Parameters:
//   - key: the string key to look up in the object
//
// Returns:
//   - a slice of integers if the value is a convertible array, or empty slice if not found or not convertible
func (o Object) GetIntArray(key string) []int {
	if v, ok := o[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			result := make([]int, 0, len(arr))
			for _, item := range arr {
				switch val := item.(type) {
				case float64:
					result = append(result, int(val))
				case int:
					result = append(result, val)
				case int64:
					result = append(result, int(val))
				}
			}
			return result
		}
	}
	return []int{}
}

// String converts the JSON object to its string representation.
// It uses json.Marshal for serialization and returns a formatted JSON string.
// Returns:
//   - the JSON string representation of the object, or "{}" if marshaling fails
func (o Object) String() string {
	return string(o.Bytes())
}

// Bytes converts the JSON object to its byte slice representation.
// It uses json.Marshal for serialization and returns the raw bytes.
// Returns:
//   - the JSON byte slice representation of the object, or empty byte slice if marshaling fails
func (o Object) Bytes() []byte {
	if b, e := json.Marshal(o); e == nil {
		return b
	}
	return []byte(`{}`)
}

// ParseResult represents the result of parsing JSON data.
// It contains information about whether the parsed data is an array or object,
// the parsed result, and any error that occurred during parsing.
type ParseResult struct {
	isArray bool        // indicates if the parsed data is a JSON array
	result  interface{} // the parsed JSON data (either []Object or Object)
	err     error       // any error that occurred during parsing
}

// IsArray checks if the parsed result contains a JSON array.
// Returns:
//   - true if the parsed data is an array, false otherwise
func (p *ParseResult) IsArray() bool {
	return p.isArray
}

// IsValid checks if the parsing operation was successful (no errors occurred).
// Returns:
//   - true if parsing was successful, false if an error occurred
func (p *ParseResult) IsValid() bool {
	return p.err == nil
}

// Error returns the error that occurred during parsing, if any.
// Returns:
//   - the parsing error, or nil if parsing was successful
func (p *ParseResult) Error() error {
	return p.err
}

// Array retrieves the parsed result as a slice of JSON objects.
// This should only be called after checking IsArray() returns true.
// Returns:
//   - the array of Objects if the result is an array, or empty slice otherwise
func (p *ParseResult) Array() []Object {
	if arr, ok := p.result.([]Object); ok {
		return arr
	}
	return []Object{}
}

// Object retrieves the parsed result as a JSON object.
// This should only be called after checking IsArray() returns false.
// Returns:
//   - the Object if the result is an object, or empty Object otherwise
func (p *ParseResult) Object() Object {
	if obj, ok := p.result.(Object); ok {
		return obj
	}
	return Object{}
}

// Parse attempts to parse JSON data from a byte slice.
// It automatically detects whether the JSON is an array or object based on the first character.
// The function trims whitespace and validates the JSON structure.
// Parameters:
//   - data: the byte slice containing JSON data to parse
//
// Returns:
//   - a ParseResult containing the parsed data, type information, and any parsing errors
func Parse(data []byte) *ParseResult {
	p := &ParseResult{}
	data = bytes.TrimSpace(data)
	if len(data) > 0 {
		switch data[0] {
		case '[': //array
			p.isArray = true
			obj := []Object{}
			if err := json.Unmarshal(data, &obj); err != nil {
				p.err = err
			}
			p.result = obj
		case '{': //object
			obj := make(Object)
			if err := json.Unmarshal(data, &obj); err != nil {
				p.err = err
			}
			p.result = obj
		default:
			p.err = errors.New(`invalid json`)
		}
	} else {
		p.err = errors.New(`empty json`)
	}
	return p
}

func ParseString(data string) *ParseResult {
	return Parse([]byte(data))
}
