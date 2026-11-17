package scanner

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListFiles(t *testing.T) {
	// Create a temporary directory for testing.
	root := t.TempDir()

	// Create a subdirectory.
	subfolder := filepath.Join(root, "subfolder")

	if err := os.Mkdir(subfolder, 0755); err != nil {
		t.Fatal(err)
	}

	// Define two test files.
	first := filepath.Join(root, "file1.txt")
	second := filepath.Join(subfolder, "file2.txt")

	// Create the files.
	for _, path := range []string{first, second} {
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Execute the function we want to test.
	got, err := ListFiles(root)

	if err != nil {
		t.Fatal(err)
	}

	// Define the expected result.
	want := []string{first, second}

	// Compare actual and expected results.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListFiles() = %v, want %v", got, want)
	}
}
