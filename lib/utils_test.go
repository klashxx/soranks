package lib

import "testing"

func TestGetKey(t *testing.T) {
	testCases := []struct {
		file     string
		expected string
		hasError bool
	}{
		{"../samples/test.key", "ThisIsSecret", false},
		{"../samples/test", "", true},
	}

	for _, test := range testCases {
		key, err := GetKey(test.file)
		if err == nil && test.hasError {
			t.Fatalf("GetKey (%s) should return an error.", test.file)

		}
		if err != nil && !test.hasError {
			t.Fatalf("GetKey (%s) should not return an error.", test.file)
		}
		if key != test.expected {
			t.Fatalf("Got: %s, want %s", key, test.expected)
		}
	}
}

func TestF2Base64(t *testing.T) {
	testCases := []struct {
		file     string
		expected string
		hasError bool
	}{
		{"../samples/lorem.dat", `TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdCwgc2VkIGRvIGVpdXNtb2QgdGVtcG9yIGluY2lkaWR1bnQgdXQgbGFib3JlIGV0IGRvbG9yZSBtYWduYSBhbGlxdWEu`, false},
		{"../samples/test", "", true},
	}
	for _, test := range testCases {
		base64, err := F2Base64(test.file)
		if err == nil && test.hasError {
			t.Fatalf("F2Base64 (%s) should return an error.", test.file)

		}
		if err != nil && !test.hasError {
			t.Fatalf("F2Base64 (%s) should not return an error.", test.file)
		}

		if base64 != test.expected {
			t.Fatalf("Got: %s, want %s", base64, test.expected)

		}

	}

}
