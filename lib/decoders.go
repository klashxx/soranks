package lib

import (
	"encoding/json"
	"io"
)

func JSONDecoder(r io.Reader, data interface{}) error {
	return json.NewDecoder(r).Decode(data)
}
