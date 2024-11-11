package archive_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Darkness4/toshokan/archive"
)

func TestFindFile(t *testing.T) {
	tt := []struct {
		archivePath string
		found       bool
		fileName    string
	}{
		{"fixtures/test.cbz", true, "info.yaml"},
		{"fixtures/test.cbz", false, "subdir/koushoku.yaml"},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%s/%s", tc.archivePath, tc.fileName), func(t *testing.T) {
			_, found, err := archive.FindFile(tc.archivePath, tc.fileName)
			if err != nil {
				t.Fatal(err)
			}
			if found != tc.found {
				t.Fatalf("expected %v, got %v", tc.found, found)
			}
		})
	}
}

func TestExtractFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "toshokan")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	tt := []struct {
		archivePath string
		fileName    string
		error       bool
	}{
		{"fixtures/test.cbz", "info.yaml", false},
		{"fixtures/test.cbz", "subdir/koushoku.yaml", true},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%s/%s", tc.archivePath, tc.fileName), func(t *testing.T) {
			err := archive.ExtractFile(tc.archivePath, tc.fileName, filepath.Join(tmpDir, tc.fileName))
			if err != nil && !tc.error {
				t.Fatal(err)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "toshokan")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	tt := []struct {
		archivePath string
	}{
		{"fixtures/test.cbz"},
	}

	for _, tc := range tt {
		t.Run(tc.archivePath, func(t *testing.T) {
			err := archive.Extract(tc.archivePath, tmpDir)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
