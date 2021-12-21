package testing

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"

	zipfly "github.com/baptistejub/zipfly/zip_fly"
)

var content = io.NopCloser(strings.NewReader("Hello, world!"))

func TestNewZipStreamerNoData(t *testing.T) {
	s, err := zipfly.NewZipStreamer(make([]zipfly.Link, 0))

	if err == nil || s != nil {
		t.Fatalf("created streamer from no data")
	}
}

func TestNewZipStreamerInvalidData(t *testing.T) {
	links := []zipfly.Link{
		{Url: "https://test.com", Filename: "test.jpg"},
		{Url: ""},
	}
	s, err := zipfly.NewZipStreamer(links)

	if err == nil || s != nil {
		t.Fatalf("created streamer from no data")
	}
}

func TestStreamFiles(t *testing.T) {
	entries := []*zipfly.Entry{
		{Url: "https://ignored.com", ZipPath: "test.jpg", ContentReader: content},
		{Url: "https://ignored.com", ZipPath: "test2.json", ContentReader: content},
	}

	s := zipfly.ZipStreamer{Entries: entries, CompressionMethod: zip.Store}

	w := new(bytes.Buffer)

	s.StreamFiles(w)

	if w.String() == "" {
		t.Fatalf("empty stream")
	}
}

func TestStreamFilesInvalid(t *testing.T) {
	entries := []*zipfly.Entry{
		{Url: "https://ignored.com", ZipPath: "test.jpg", ContentReader: content},
		{Url: "invalid", ZipPath: "test2.json"},
	}

	s := zipfly.ZipStreamer{Entries: entries, CompressionMethod: zip.Store}

	w := new(bytes.Buffer)

	err := s.StreamFiles(w)

	if err == nil {
		t.Fatalf("streamed invalid zip")
	}
}
