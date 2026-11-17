package duplicates

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func hashFile(path string) ([sha256.Size]byte, error) {
	var result [sha256.Size]byte

	file, err := os.Open(path)
	if err != nil {
		return result, fmt.Errorf("opening file %q: %w", path, err)
	}

	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return result, err
	}
	copy(result[:], hasher.Sum(nil))
	return result, nil
}

func GetSameHashFiles(sizeMap map[int64][]string) (map[[sha256.Size]byte][]string, error) {
	hashMap := make(map[[sha256.Size]byte][]string)

	for _, files := range sizeMap {
		for _, path := range files {
			hash, err := hashFile(path)
			if err != nil {
				return nil, fmt.Errorf("hashing file %q: %w", path, err)
			}

			hashMap[hash] = append(hashMap[hash], path)
		}
	}

	for hash, files := range hashMap {
		if len(files) < 2 {
			delete(hashMap, hash)
		}
	}

	return hashMap, nil
}
