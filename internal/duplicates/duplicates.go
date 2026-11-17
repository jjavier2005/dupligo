package duplicates

import (
	"fmt"
	"os"
)

func GetRepeatedSizes(paths []string) (map[int64][]string, []int64, error) {
	sizeMap := make(map[int64][]string)

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"getting file info for %q: %w", path, err,
			)
		}

		if !info.Mode().IsRegular() {
			return nil, nil, fmt.Errorf("%q is not a regular file", path)
		}

		size := info.Size()
		sizeMap[size] = append(sizeMap[size], path)
	}

	// Remove groups with only one file.
	// Store the sizes of the remaining groups.
	var repeatedSizes []int64

	for size, files := range sizeMap {
		if len(files) < 2 {
			delete(sizeMap, size)
		} else {
			repeatedSizes = append(repeatedSizes, size)
		}
	}

	return sizeMap, repeatedSizes, nil
}
