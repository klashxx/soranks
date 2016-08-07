package lib

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetKey(path string) (string, error) {

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("Can't find key: %s\n", path)
	}

	strkey, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Can't load key: %s\n", err)
	}

	return strings.TrimRight(string(strkey)[:], "\n"), nil
}

func F2Base64(path string) (string, error) {

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Can't load file to encode: %s", err)
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}
