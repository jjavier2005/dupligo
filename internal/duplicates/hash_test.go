package duplicates

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Creates a temporary file with the given content.
func createTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	return path
}

// Tests that files with the same hash are grouped together.
func TestGetSameHashFiles(t *testing.T) {
	root := t.TempDir()

	// All five files have the same size (3 bytes).
	file1 := createTestFile(t, root, "a.txt", "abc")
	file2 := createTestFile(t, root, "b.txt", "abc")
	file3 := createTestFile(t, root, "c.txt", "xyz")
	file4 := createTestFile(t, root, "d.txt", "xyz")
	file5 := createTestFile(t, root, "e.txt", "123")

	sizeMap := map[int64][]string{
		3: {file1, file2, file3, file4, file5},
	}

	got, err := GetSameHashFiles(sizeMap)
	if err != nil {
		t.Fatal(err)
	}

	// Expected groups based on file content.
	want := map[[sha256.Size]byte][]string{
		sha256.Sum256([]byte("abc")): {file1, file2},
		sha256.Sum256([]byte("xyz")): {file3, file4},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetSameHashFiles() = %v, want %v", got, want)
	}
}

// Tests that files with different hashes are discarded.
func TestGetSameHashFiles_NoMatches(t *testing.T) {
	root := t.TempDir()

	file1 := createTestFile(t, root, "a.txt", "abc")
	file2 := createTestFile(t, root, "b.txt", "xyz")

	sizeMap := map[int64][]string{
		3: {file1, file2},
	}

	got, err := GetSameHashFiles(sizeMap)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 0 {
		t.Errorf("expected an empty map, got %v", got)
	}
}

// Tests error handling when a file does not exist.
func TestGetSameHashFiles_MissingFile(t *testing.T) {
	root := t.TempDir()

	missingPath := filepath.Join(root, "missing.txt")

	sizeMap := map[int64][]string{
		3: {missingPath},
	}

	got, err := GetSameHashFiles(sizeMap)

	if err == nil {
		t.Error("expected an error for a missing file")
	}

	if got != nil {
		t.Errorf("expected nil map, got %v", got)
	}
}
