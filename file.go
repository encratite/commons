package commons

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func ReadFile(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file (%s): %v", path, err)
	}
	return content
}

func WriteFile(path, data string) {
	bytes := []byte(data)
	err := os.WriteFile(path, bytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write file (%s): %v", path, err)
	}
}

func CreateDirectory(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalf("Failed to create directory: %v", err)
	}
}

func GetFiles(directory string, extension string) []string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Unable to read files from directory (%s): %v", directory, err)
	}
	files := []string{}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && filepath.Ext(name) == extension {
			path := filepath.Join(directory, name)
			files = append(files, path)
		}
	}
	return files
}

func GetDirectories(directory string) []string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Unable to read directory (%s): %v", directory, err)
	}
	directories := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(directory, entry.Name())
			directories = append(directories, path)
		}
	}
	return directories
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}