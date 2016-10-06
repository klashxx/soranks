package lib

import (
	"bytes"
	"encoding/json"
	"io"
)

func JSONDecoder(r io.Reader, data interface{}) error {
	return json.NewDecoder(r).Decode(data)
}

func JSONEncoder(data interface{}) (buf io.ReadWriter, err error) {

	buf = new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}
	return buf, nil
}
