# goson
Java-like json parser for Go.

Goson created to ease json parsing and building without hasle using struct.

# Install

```
go get -u github.com/firabliz/goson

```

# Example

### 1. Json Parsing

Format JSON string and will returns the JSON Struct.

**Parameters :**

jsonString - it should be string format that contain JSON (check example string)

**Returns :**

JSON struct

**Sample :**
```json
{
	"string_data": "string_value",
	"numeric_data": 1
}
```
**How to use**
```go
//parsing json
jsObj   := goson.Parse(jsonVal)

fmt.println(jsObj.GetString(`key`))         //print: string_value

fmt.println(jsObj.GetInt(`numeric_data`))   //print: 1

```


### 2. ToString

Returns JSON text for this JSON value.

**Returns :**

JSON text

**How to use**

```go
//parsing json
json  := goson.Parse(jsonVal)

fmt.println(json.ToString())  //print: '{"string_data": "string_value","numeric_data": 1}'
```

### 3. GetString

Returns the string value to which the specified name is mapped.

**Returns :**

string value

**How to use**

```go
//parsing json
json   := goson.Parse(jsonVal)

fmt.println(json.GetString(`key`))         //print: string_value
```


### 4. GetJsonArray

Returns the array value to which the specified name is mapped.

**Returns :**

array value

**How to use**

```go
//parsing json
json   := goson.Parse(jsonVal)

var arrJson []JsonObject
arrJson = json.GetJsonArray(`key`)

```

### 5. GetJsonObject

Returns the object value to which the specified name is mapped.

**Returns :**

object value

**How to use**

```go
//parsing json
json   := goson.Parse(jsonVal)

var obj JsonObject
obj = json.GetJsonObject(`key`)


```



