package duplicates

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGroupBySize(t *testing.T) {
	root := t.TempDir()

	// Define test file paths.
	file1 := filepath.Join(root, "a.txt")
	file2 := filepath.Join(root, "b.txt")
	file3 := filepath.Join(root, "c.txt")
	file4 := filepath.Join(root, "d.txt")
	file5 := filepath.Join(root, "e.txt")

	// Create files with different sizes.
	testFiles := map[string]string{
		file1: "abc",     // 3 bytes
		file2: "xyz",     // 3 bytes
		file3: "hello",   // 5 bytes
		file4: "world",   // 5 bytes
		file5: "1234567", // 7 bytes
	}

	for path, content := range testFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	paths := []string{file1, file2, file3, file4, file5}

	// Execute the function.
	groups, repeatedSizes, err := GetRepeatedSizes(paths)
	if err != nil {
		t.Fatal(err)
	}

	// Expected groups: only sizes with multiple files.
	wantGroups := map[int64][]string{
		3: {file1, file2},
		5: {file3, file4},
	}

	// Check the groups.
	if !reflect.DeepEqual(groups, wantGroups) {
		t.Errorf("groups = %v, want %v", groups, wantGroups)
	}

	// Check the number of repeated sizes.
	if len(repeatedSizes) != len(wantGroups) {
		t.Fatalf("expected %d repeated sizes, got %d",
			len(wantGroups), len(repeatedSizes))
	}

	// Check repeated sizes without depending on their order.
	seen := make(map[int64]bool)

	for _, size := range repeatedSizes {
		if _, exists := wantGroups[size]; !exists {
			t.Errorf("unexpected repeated size: %d", size)
		}

		if seen[size] {
			t.Errorf("repeated size appears twice: %d", size)
		}

		seen[size] = true
	}
}
