package commons

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"runtime"
	"slices"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	httpTimeoutSeconds = 10
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:142.0) Gecko/20100101 Firefox/142.0"
)

type taskTuple[T any] struct {
	index int
	element T
}

func Contains[T comparable](slice []T, element T) bool {
	for _, x := range slice {
		if x == element {
			return true
		}
	}
	return false
}

func ContainsFunc[T any](slice []T, match func (T) bool) bool {
	for _, x := range slice {
		if match(x) {
			return true
		}
	}
	return false
}

func Find[T any](slice []T, match func (T) bool) (T, bool) {
	index := slices.IndexFunc(slice, func (element T) bool {
		return match(element)
	})
	if index >= 0 {
		return slice[index], true
	} else {
		var zeroValue T
		return zeroValue, false
	}
}

func FindPointer[T any](slice []T, match func (T) bool) (*T, bool) {
	index := slices.IndexFunc(slice, func (element T) bool {
		return match(element)
	})
	if index >= 0 {
		return &slice[index], true
	} else {
		return nil, false
	}
}

func ParallelMap[A, B any](elements []A, callback func(A) B) []B {
	workers := runtime.NumCPU()
	elementChan := make(chan taskTuple[A], len(elements))
	for i, x := range elements {
		elementChan <- taskTuple[A]{
			index: i,
			element: x,
		}
	}
	close(elementChan)
	var wg sync.WaitGroup
	wg.Add(workers)
	output := make([]B, len(elements))
	for range workers {
		go func() {
			defer wg.Done()
			for task := range elementChan {
				output[task.index] = callback(task.element)
			}
		}()
	}
	wg.Wait()
	return output
}

func ReadCSV(path string, callback func ([]string)) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to read CSV (%s): %v", path, err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	_, _ = reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		callback(record)
	}
}

func ReadCSVColumns(path string, columns []string, callback func ([]string)) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to read CSV file from %s: %v", path, err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	csvColumns, err := reader.Read()
	if err == io.EOF {
		log.Fatalf("Failed to read CSV columns from %s: %v", path, err)
	}
	indices := []int{}
	for _, column := range columns {
		index := slices.Index(csvColumns, column)
		if index == -1 {
			log.Fatalf("Unable to find column \"%s\" in CSV file %s", column, path)
		}
		indices = append(indices, index)
	}
	line := 2
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		indexRecord := []string{}
		for _, index := range indices {
			if index >= len(record) {
				log.Fatalf("Not enough records on line %d in CSV file %s", line, path)
			}
			indexRecord = append(indexRecord, record[index])
		}
		callback(indexRecord)
		line++
	}
}

func ReadJSON[T any](path string) T {
	data := ReadFile(path)
	var output T
	err := json.Unmarshal(data, &output)
	if err != nil {
		log.Fatalf("Failed to deserialize JSON: %v", err)
	}
	return output
}

func LoadConfiguration[T any](path string, configuration *T) *T {
	if configuration != nil {
		panic("Configuration had already been loaded")
	}
	yamlData := ReadFile(path)
	configuration = new(T)
	err := yaml.Unmarshal(yamlData, configuration)
	if err != nil {
		log.Fatal("Failed to unmarshal YAML:", err)
	}
	return configuration
}