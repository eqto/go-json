# go-json
go-json created to ease json parsing and building without hasle using struct.

# Install

```
go get -u github.com/eqto/go-json

```

# Example

### 1. Json Parsing

Format JSON string and will returns the JSON Struct.

**Parameters :**

jsonBytes - it should be []byte format that contain JSON

**Returns :**

JSON struct

**How to use**
```go
/*
data:
{
  "title": "Learning go",
  "num_chapters": 10,
  "author": {
    "first_name": "John",
    "last_name": "Doe"
  }
}
*/

obj := json.Parse(data)

fmt.println(obj.GetString(`title`))         //print: Learning go

fmt.println(obj.GetInt(`num_chapters`))   //print: 10

fmt.println(obj.GetString(`author.first_name`))   //print: John

```

### 2. GetJsonArray

Returns the array value to which the specified name is mapped.

**Returns :**

array value

**How to use**

```go
/*
data:
{
  "books": [
    {
      "title": "Learning go",
      "num_chapters": 10
    },
    {
      "title": "Basic of go",
      "num_chapters": 5
    }
  ]
}
*/

obj := goson.Parse(data)

arr := obj.GetArray(`books`)

```

