package zipfly

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"
)

type ZipStreamer struct {
	Entries []*Entry
}

func NewZipStreamer(files []File) (*ZipStreamer, error) {
	if len(files) == 0 {
		return nil, errors.New("no file to zip")
	}

	entries := make([]*Entry, 0)
	for _, file := range files {
		entry, err := NewEntry(file.Url, file.Filename, file.Compress)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	z := ZipStreamer{Entries: entries}

	return &z, nil
}

func (z *ZipStreamer) StreamFiles(w io.Writer) error {
	zipWriter := zip.NewWriter(w)

	for _, entry := range z.Entries {
		err := z.writeEntry(zipWriter, entry)
		if err != nil {
			fmt.Println("Error while writing file to stream", entry.ZipPath, ":", err.Error())
			return err
		}
	}

	return zipWriter.Close()
}

func (z *ZipStreamer) writeEntry(zipWriter *zip.Writer, entry *Entry) error {
	content, err := entry.Content()
	if err != nil {
		return err
	}

	defer content.Close()

	header := &zip.FileHeader{
		Name:     entry.ZipPath,
		Method:   entry.CompressionMethod,
		Modified: time.Now(),
	}
	entryWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(entryWriter, bufio.NewReader(content))
	if err != nil {
		return err
	}

	return nil
}
