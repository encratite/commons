package commons

import (
	"log"
	"os"
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