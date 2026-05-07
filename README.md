# JSON Package Documentation

This package provides utility functions for working with JSON data in Go. It offers a convenient wrapper around Go's standard `encoding/json` package with additional helper methods for type-safe JSON manipulation.

## Overview

The package defines two main types:
- `Object`: A map-based representation of JSON objects with helper methods
- `ParseResult`: A result wrapper for JSON parsing operations

The `Object` type provides methods for:
- Basic operations: `Set`, `Get`, `Has`
- Type-safe getters: `GetInt`, `GetFloat`, `GetString`
- Array operations: `GetArray`, `GetObjectArray`, `GetIntArray`
- Nested object access: `GetObject`
- Serialization: `String`, `Bytes`

## Types

### Object

`Object` is defined as `map[string]any` and represents a JSON object. It provides convenient methods for accessing and manipulating JSON data.

### ParseResult

`ParseResult` contains the result of JSON parsing operations:
```go
type ParseResult struct {
    isArray bool  // indicates if the parsed data is a JSON array
    result  any   // the parsed JSON data (either []Object or Object)
    err     error // any error that occurred during parsing
}
```

## API Reference

### Object Methods

#### Basic Operations

#### `Set(key string, val any)`
Adds or updates a key-value pair in the JSON object.
- **Parameters:**
  - `key`: the string key to set in the object
  - `val`: the value to associate with the key (can be any type)

#### `Get(key string) any`
Retrieves a value from the JSON object by key.
- **Parameters:**
  - `key`: the string key to look up in the object
- **Returns:** the value associated with the key, or `nil` if the key doesn't exist

#### `GetObject(key string) Object`
Retrieves a nested JSON object from the current object by key. It handles both `Object` type and `map[string]any` type conversions.
- **Parameters:**
  - `key`: the string key to look up in the object
- **Returns:** the nested Object if found and convertible, or an empty Object if not found or not convertible

#### `Has(key string) bool`
Checks if a key exists in the JSON object.
- **Parameters:**
  - `key`: the string key to check for existence
- **Returns:** `true` if the key exists in the object, `false` otherwise

#### Type-Safe Getters

#### `GetInt(key string) int`
Retrieves an integer value from the JSON object by key. It attempts to convert various types (`float64`, `string`, `bool`) to `int`.
- **Type Conversion Rules:**
  - `float64`: Direct cast to `int`
  - `string`: Uses `strconv.Atoi()` for conversion
  - `bool`: `true` returns `1`, `false` returns `0`
- **Parameters:**
  - `key`: the string key to look up in the object
- **Returns:** the integer value if convertible, or `0` if not found or not convertible

#### `GetFloat(key string) float64`
Retrieves a float64 value from the JSON object by key. It attempts to convert various types (`float64`, `string`) to `float64`.
- **Type Conversion Rules:**
  - `float64`: Returns as-is
  - `string`: Uses `strconv.ParseFloat()` for conversion
- **Parameters:**
  - `key`: the string key to look up in the object
- **Returns:** the float64 value if convertible, or `0` if not found or not convertible

#### `GetString(key string) string`
Retrieves a string value from the JSON object by key. It converts various types (`string`, `float64`, `bool`) to string representation.
- **Type Conversion Rules:**
  - `string`: Returns as-is
  - `float64`: Formats without decimal places using `strconv.FormatFloat(v, 'f', 0, 64)`
  - `bool`: Returns `"true"` or `"false"`
- **Parameters:**
  - `key`: the string key to look up in the object
- **Returns:** the string value if convertible, or empty string if not found or not convertible

#### `String() string`
Converts the JSON object to its string representation using `json.Marshal`.
- **Returns:** the JSON string representation of the object, or `"{}"` if marshaling fails

#### Array Operations

#### `GetObjectArray(key string) []Object`
Retrieves an array of JSON objects from the JSON object by key. It handles JSON arrays where each element should be a JSON object (`map[string]any`).
- **Parameters:**
  - `key`: the string key to look up in the object
- **Returns:** a slice of Object if the value is an array of objects, or empty slice if not found or not convertible

#### `GetArray(key string) []any`
Retrieves an array of any type from the JSON object by key. It handles JSON arrays and returns them as `[]any`.
- **Parameters:**
  - `key`: the string key to look up in the object
- **Returns:** a slice of any if the value is an array, or empty slice if not found or not convertible

#### `GetIntArray(key string) []int`
Retrieves an array of integers from the JSON object by key. It handles JSON arrays where numbers are typically parsed as `float64` and converts them to `int`.
- **Parameters:**
  - `key`: the string key to look up in the object
- **Returns:** a slice of integers if the value is a convertible array, or empty slice if not found or not convertible

#### `Bytes() []byte`
Converts the JSON object to its byte slice representation using `json.Marshal`.
- **Returns:** the JSON byte slice representation of the object, or empty byte slice if marshaling fails

### ParseResult Methods

#### `IsArray() bool`
Checks if the parsed result contains a JSON array.
- **Returns:** `true` if the parsed data is an array, `false` otherwise

#### `IsValid() bool`
Checks if the parsing operation was successful (no errors occurred).
- **Returns:** `true` if parsing was successful, `false` if an error occurred

#### `Error() error`
Returns the error that occurred during parsing, if any.
- **Returns:** the parsing error, or `nil` if parsing was successful

#### `Array() []Object`
Retrieves the parsed result as a slice of JSON objects. This should only be called after checking `IsArray()` returns `true`.
- **Returns:** the array of Objects if the result is an array, or empty slice otherwise

#### `Object() Object`
Retrieves the parsed result as a JSON object. This should only be called after checking `IsArray()` returns `false`.
- **Returns:** the Object if the result is an object, or empty Object otherwise

### Package Functions

#### `Parse(data []byte) *ParseResult`
Attempts to parse JSON data from a byte slice. It automatically detects whether the JSON is an array or object based on the first character (`[` for array, `{` for object). The function trims whitespace and validates the JSON structure.

- **Parameters:**
  - `data`: the byte slice containing JSON data to parse
- **Returns:** a `ParseResult` containing the parsed data, type information, and any parsing errors
- **Error Cases:**
  - Returns error for invalid JSON (doesn't start with `[` or `{`)
  - Returns error for empty JSON data
  - Returns unmarshaling errors from `encoding/json` package

## Usage Examples

### Creating and Manipulating JSON Objects
```go
// Create a new JSON object
obj := json.Object{
    "name": "John Doe",
    "age": 30,
    "active": true,
}

// Set values
obj.Set("email", "john@example.com")
obj.Set("score", 95.5)

// Get values
name := obj.GetString("name")    // "John Doe"
age := obj.GetInt("age")         // 30
active := obj.GetString("active") // "true"
email := obj.GetString("email") // "john@example.com"
score := obj.GetInt("score")     // 95 (converted from float64)

// Check existence
if obj.Has("name") {
    fmt.Println("Name exists")
}

// Convert to string and bytes
jsonStr := obj.String() // {"name":"John Doe","age":30,"active":true,"email":"john@example.com","score":95.5}
jsonBytes := obj.Bytes() // []byte with the same JSON content
```

### Working with Arrays
```go
// Create JSON object with arrays
obj := json.Object{
    "tags": ["go", "json", "api"],
    "scores": [95.5, 88.0, 76.5],
    "users": [
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"}
    ],
}

// Get generic array
tags := obj.GetArray("tags") // []any{"go", "json", "api"}

// Get integer array (converts float64 to int)
scores := obj.GetIntArray("scores") // []int{95, 88, 76}

// Get object array
users := obj.GetObjectArray("users") // []json.Object
for _, user := range users {
    id := user.GetInt("id")
    name := user.GetString("name")
    fmt.Printf("User %d: %s\n", id, name)
}
```

### Parsing JSON Data
```go
// Parse JSON object
jsonData := []byte(`{"name":"Alice","age":25}`)
result := json.Parse(jsonData)

if result.IsValid() {
    obj := result.Object()
    name := obj.GetString("name") // "Alice"
    age := obj.GetInt("age")      // 25
} else {
    fmt.Println("Parse error:", result.Error())
}

// Parse JSON array
jsonArray := []byte(`[{"id":1,"name":"Item 1"},{"id":2,"name":"Item 2"}]`)
result := json.Parse(jsonArray)

if result.IsArray() && result.IsValid() {
    items := result.Array()
    for _, item := range items {
        id := item.GetInt("id")
        name := item.GetString("name")
        fmt.Printf("ID: %d, Name: %s\n", id, name)
    }
}
```

### Working with Nested Objects
```go
jsonData := []byte(`{
    "user": {
        "profile": {
            "name": "Bob",
            "settings": {
                "theme": "dark",
                "notifications": true
            }
        }
    }
}`)

result := json.Parse(jsonData)
if result.IsValid() {
    obj := result.Object()
    user := obj.GetObject("user")
    profile := user.GetObject("profile")
    settings := profile.GetObject("settings")
    
    theme := settings.GetString("theme")         // "dark"
    notifications := settings.GetString("notifications") // "true"
}
```

## Error Handling

Always check for parsing errors and validate results:
```go
result := json.Parse(jsonData)

if !result.IsValid() {
    log.Printf("JSON parsing failed: %v", result.Error())
    return
}

if result.IsArray() {
    items := result.Array()
    // Process array
} else {
    obj := result.Object()
    // Process object
}
```

## Type Safety Notes

- **Default Values**: Methods like `GetInt()` and `GetString()` return zero values (`0`, `""`) when keys don't exist or types can't be converted
- **Type Conversions**: The package performs automatic type conversions where reasonable, but may lose precision (e.g., `float64` to `int`)
- **Nested Objects**: Always check if nested objects exist before accessing them to avoid empty object returns