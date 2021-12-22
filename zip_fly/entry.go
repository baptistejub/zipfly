package zipfly

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Entry struct {
	Url               string
	ZipPath           string
	CompressionMethod uint16
	ContentReader     io.ReadCloser
}

func NewEntry(urlString string, zipPath string, compress bool) (*Entry, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(url.Scheme, "http") {
		return nil, errors.New("invalid file url")
	}

	zipPath = path.Clean(zipPath)
	zipPath = strings.TrimPrefix(zipPath, "/")

	if strings.HasPrefix(zipPath, "../") {
		return nil, errors.New("invalid zip filename: " + zipPath)
	}

	if filename := path.Base(zipPath); len(filename) == 0 || filename == "." {
		return nil, errors.New("invalid zip filename")
	}

	compressionMethod := zip.Store
	if compress {
		compressionMethod = zip.Deflate
	}

	return &Entry{Url: urlString, ZipPath: zipPath, CompressionMethod: compressionMethod}, nil
}

func (e *Entry) Size() uint64 {
	res, err := http.Head(e.Url)
	if err != nil {
		return 0
	}

	return uint64(res.ContentLength)
}

func (e *Entry) Content() (io.ReadCloser, error) {
	if e.ContentReader != nil {
		return e.ContentReader, nil
	}

	return e.fetchContent()
}

func (e *Entry) fetchContent() (io.ReadCloser, error) {
	resp, err := http.Get(e.Url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("couldn't fetch from URL")
	}

	e.ContentReader = resp.Body

	return e.ContentReader, nil
}
