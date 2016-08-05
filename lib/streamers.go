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

	url := fmt.Sprintf("%s/%s%s", apiurl, fmt.Sprintf(query, page), key)
	Trace.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Trace.Println(err)
		return users, err
	}
	Trace.Println("Sending header.")

	req.Header.Set("Accept-Encoding", "gzip")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		Trace.Println(err)
		return users, err
	}
	Trace.Println("Response.")

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
	users = new(SOUsers)
	return users, Decoder(reader, users)
}

func StreamHTTP2(url string) (repo *Repo, err error) {

	Trace.Println(url)

	var reader io.ReadCloser
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Trace.Println(err)
	}
	Trace.Println("Sending header.")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		Trace.Println(err)
	}
	Trace.Println("Response.")

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

	repo = new(Repo)
	return repo, Decoder(reader, repo)
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
