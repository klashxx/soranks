package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func StreamHTTP(url string, data interface{}, gzipped bool) error {

	var reader io.ReadCloser

	Trace.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Trace.Println(err)
		return err
	}
	Trace.Println("Sending header.")

	if gzipped {
		req.Header.Set("Accept-Encoding", "gzip")
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		Trace.Println(err)
		return err
	}
	Trace.Println("Response.")

	defer response.Body.Close()

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			Trace.Println(err)
			return err
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	return Decoder(reader, data)
}

func DataToGihub(data interface{}) (buf io.ReadWriter, err error) {

	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		return nil, err
	}
	return buf, nil

}

func StreamFile(jsonfile string) (users *SOUsers, err error) {
	reader, _ := os.Open(jsonfile)
	defer reader.Close()

	users = new(SOUsers)
	return users, Decoder(reader, users)
}
