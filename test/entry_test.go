package testing

import (
	"testing"

	zipfly "github.com/baptistejub/zipfly/zip_fly"
)

func TestNewEntryValid(t *testing.T) {
	p, err := zipfly.NewEntry("https://my_file.com", "my_file.epub")

	if err != nil || p == nil {
		t.Fatalf("invalid entry")
	}

	if p.Url != "https://my_file.com" {
		t.Fatalf("invalid entry url: %s", p.Url)
	}

	if p.ZipPath != "my_file.epub" {
		t.Fatalf("invalid entry path: %s", p.ZipPath)
	}
}

func TestNewEntryValidWithAbosolutePath(t *testing.T) {
	p, err := zipfly.NewEntry("https://my_file.com", "/my_file.epub")

	if err != nil || p == nil {
		t.Fatalf("invalid entry")
	}

	if p.ZipPath != "my_file.epub" {
		t.Fatalf("invalid entry path: %s", p.ZipPath)
	}
}

func TestNewEntryValidWithSubDir(t *testing.T) {
	p, err := zipfly.NewEntry("https://my_file.com", "test-dir/my_file.epub")

	if err != nil || p == nil {
		t.Fatalf("invalid entry")
	}

	if p.ZipPath != "test-dir/my_file.epub" {
		t.Fatalf("invalid entry path: %s", p.ZipPath)
	}
}

func TestNewEntryValidWithWeirdPath(t *testing.T) {
	p, err := zipfly.NewEntry("https://my_file.com", "/test-dir/../.my_file.epub")

	if err != nil || p == nil {
		t.Fatalf("invalid entry")
	}

	if p.ZipPath != ".my_file.epub" {
		t.Fatalf("invalid entry path: %s", p.ZipPath)
	}
}

func TestNewEntryInvalidUrl(t *testing.T) {
	p, err := zipfly.NewEntry("gs://my_file.com", "my_file.epub")

	if err == nil || p != nil {
		t.Fatalf("invalid entry accepted")
	}
}

func TestNewEntryMissingPath(t *testing.T) {
	p, err := zipfly.NewEntry("http://my_file.com", "")

	if err == nil || p != nil {
		t.Fatalf("invalid entry accepted")
	}
}

func TestNewEntryInvalidPath(t *testing.T) {
	p, err := zipfly.NewEntry("http://my_file.com", "../test.jpg")

	if err == nil || p != nil {
		t.Fatalf("invalid entry accepted %s", p.ZipPath)
	}
}
