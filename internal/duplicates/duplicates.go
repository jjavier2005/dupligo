package duplicates

import (
	"fmt"
	"os"
)

func GetRepeatedSizes(paths []string) (map[int64][]string, error) {
	sizeMap := make(map[int64][]string)

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf(
				"getting file info for %q: %w", path, err,
			)
		}

		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("%q is not a regular file", path)
		}

		size := info.Size()
		sizeMap[size] = append(sizeMap[size], path)
	}

	return sizeMap, nil
}
