package lib

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetKey(path string) (key string) {

	_, err := os.Stat(path)
	if err != nil {
		Error.Printf("Can't find key: %s\n", path)
		return ""
	}

	strkey, err := ioutil.ReadFile(path)
	if err != nil {
		Error.Printf("Can't load key: %s\n", err)
		return ""
	}

	return strings.TrimRight(string(strkey)[:], "\n")
}

func Markdown2Base64(path string) (b64 string, err error) {

	mdraw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Can't load markdown: %s", err)
	}
	return base64.StdEncoding.EncodeToString(mdraw), nil
}
