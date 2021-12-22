package testing

import (
	"testing"

	zipfly "github.com/baptistejub/zipfly/zip_fly"
)

var invalidPayload = []byte(`{"files": "dfdf"}`)
var emptyPayload = []byte(`{}`)
var validPayload = []byte(`{"filename":"test.zip","files": [{"url":"https://a.com/1","filename":"file1.jpg","compress":true},{"url":"https://a.com/2","filename":"file2.jpg","compress":false},{"url":"https://a.com/3","filename":"file3.pdf"}]}`)

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

	if len(r.Files) != 0 {
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
	if len(r.Files) != 3 {
		t.Fatalf("invalid link length %v", len(r.Files))
	}

	if r.Files[0].Url != "https://a.com/1" || r.Files[0].Filename != "file1.jpg" || !r.Files[0].Compress {
		t.Fatal("invalid first file value")
	}

	if r.Files[1].Url != "https://a.com/2" || r.Files[1].Filename != "file2.jpg" || r.Files[1].Compress {
		t.Fatal("invalid second file values")
	}

	if r.Files[2].Url != "https://a.com/3" || r.Files[2].Filename != "file3.pdf" || r.Files[2].Compress {
		t.Fatal("invalid last file values")
	}
}
