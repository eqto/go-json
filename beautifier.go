package json

import "encoding/json"

type Beautifier interface {
	String() string
	Bytes() []byte
}

type beautifier struct {
	js Object
}

func (b *beautifier) String() string {
	return string(b.Bytes())
}
func (b *beautifier) Bytes() []byte {
	data, e := json.MarshalIndent(b.js, ``, `  `)
	if e != nil {
		return nil
	}
	return data
}
