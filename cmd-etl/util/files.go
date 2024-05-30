package util

import (
	"os"
)

func ReadDir(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, file := range files {
		name := file.Name()
		if !file.IsDir() {
			results = append(results, name)
		}
	}

	return results, nil
}
