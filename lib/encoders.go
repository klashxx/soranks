package lib

import (
	"bytes"
	"encoding/json"
	"io"
)

func JSONEncoder(data interface{}) (buf io.ReadWriter, err error) {

	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
