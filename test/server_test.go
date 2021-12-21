package testing

import (
	"testing"

	zipfly "github.com/baptistejub/zipfly/zip_fly"
)

var invalidPayload = []byte(`{"links": "dfdf"}`)
var emptyPayload = []byte(`{}`)
var validPayload = []byte(`{"filename":"test.zip","links": [{"url":"https://a.com/1","filename":"file1.jpg"},{"url":"https://a.com/2","filename":"file2.jpg"}]}`)

func TestUnmarshalBodyInvalid(t *testing.T) {
	p, err := zipfly.UnmarshalPayload(invalidPayload)

	if err == nil || p != nil {
		t.Fatalf("allowed invalid payload: %v", p)
	}
}

func TestUnmarshalBodyEmpty(t *testing.T) {
	r, err := zipfly.UnmarshalPayload(emptyPayload)

	if err != nil {
		t.Fatal("error on empty payload")
	}

	if len(r.Links) != 0 {
		t.Fatal("empty payload produces data")
	}
}

func TestUnmarshalBodyValid(t *testing.T) {
	r, err := zipfly.UnmarshalPayload(validPayload)

	if err != nil {
		t.Fatalf("unparsable valid payload: %v", err)
	}
	if r.Filename == "" {
		t.Fatalf("filename not parsed")
	}
	if len(r.Links) != 2 {
		t.Fatalf("invalid link length %v", len(r.Links))
	}
	if r.Links[0].Url != "https://a.com/1" || r.Links[1].Filename != "file2.jpg" {
		t.Fatal("invalid link values")
	}
}
