package lib

import (
	"fmt"
	"os"
	"testing"
)

type SOError struct {
	ErrorID      int    `json:"error_id"`
	ErrorMessage string `json:"error_message"`
	ErrorName    string `json:"error_name"`
}

func TestJSONDecoder(t *testing.T) {
	pwd, _ := os.Getwd()

	testCases := []struct {
		file     string
		hasError bool
	}{
		{"../samples/coders_ok.json", false},
		{"../samples/coders_ko.json", true},
	}

	soTestError := new(SOError)

	for _, test := range testCases {
		reader, err := os.Open(fmt.Sprintf("%s/%s", pwd, test.file))
		if err != nil {
			t.Fatalf("Cant't open json: %s", err)
		}
		defer reader.Close()

		err = JSONDecoder(reader, soTestError)
		if err != nil && !test.hasError {
			t.Fatalf("JSONDecoder should not return an error: %s", err)
		}
		if err == nil && test.hasError {
			t.Fatalf("JSONDecoder should return an error: %s", err)
		}
	}
}
