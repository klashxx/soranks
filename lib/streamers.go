package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func StreamHTTP(page int, key string, apiurl string, query string) (users *SOUsers, err error) {

	var reader io.ReadCloser

	url := fmt.Sprintf("%s%d&%s%s", apiurl, page, query, key)
	//Trace.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Trace.Println(err)
		return users, err
	}
	//Trace.Println("Sending header.")

	req.Header.Set("Accept-Encoding", "gzip")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		Trace.Println(err)
		return users, err
	}
	//Trace.Println("Response.")

	defer response.Body.Close()

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			Trace.Println(err)
			return users, err
		}
		defer reader.Close()
	default:
		reader = response.Body
	}
	return Decode(reader)
}

func StreamHTTP2(url string) (repo *Repo, err error) {

	//Trace.Println(url)

	var reader io.ReadCloser
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Trace.Println(err)
	}
	//Trace.Println("Sending header.")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		Trace.Println(err)
	}
	//Trace.Println("Response.")

	defer response.Body.Close()

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			Trace.Println(err)
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	return Decode2(reader)
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
	reader, err := os.Open(jsonfile)
	defer reader.Close()
	return Decode(reader)
}
