package scanner


import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func ListFiles(root string) ([]string,err error){
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry,walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type().IsRegular() {
			files = append(files,path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error scanning %q: %w",root,err)
	}
	
	return files, nil
}