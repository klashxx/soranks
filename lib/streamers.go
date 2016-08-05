package lib

import (
	"compress/gzip"
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

func StreamFile(jsonfile string, data interface{}) error {
	reader, _ := os.Open(jsonfile)
	defer reader.Close()

	return Decoder(reader, data)
}
