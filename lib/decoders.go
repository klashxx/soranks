package lib

import (
	"encoding/json"
	"io"
	"reflect"
)

func Decoder(r io.Reader, data interface{}) error {
	Trace.Println(reflect.TypeOf(data))
	return json.NewDecoder(r).Decode(data)
}
