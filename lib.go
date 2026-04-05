package commons

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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

func Parallel[A any](elements []A, workers int, callback func(A)) {
	if workers == 0 {
		workers = runtime.NumCPU()
	}
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
	for range workers {
		go func() {
			defer wg.Done()
			for task := range elementChan {
				callback(task.element)
			}
		}()
	}
	wg.Wait()
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

func ReadCSVFile(reader io.Reader, skipHeader bool, callback func ([]string)) {
	csvReader := csv.NewReader(reader)
	if skipHeader {
		_, _ = csvReader.Read()
	}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		callback(record)
	}
}

func ReadCSV(path string, skipHeader bool, callback func ([]string)) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to read CSV (%s): %v", path, err)
	}
	defer file.Close()
	ReadCSVFile(file, skipHeader, callback)
	return nil
}

func ReadCSVColumnsFile(reader io.Reader, path string, columns []string, callback func ([]string)) error {
	csvReader := csv.NewReader(reader)
	csvColumns, err := csvReader.Read()
	if err == io.EOF {
		return fmt.Errorf("Failed to read CSV columns from %s: %v", path, err)
	}
	indices := []int{}
	for _, column := range columns {
		index := slices.Index(csvColumns, column)
		if index == -1 {
			return fmt.Errorf("Unable to find column \"%s\" in CSV file %s", column, path)
		}
		indices = append(indices, index)
	}
	line := 2
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		indexRecord := []string{}
		for _, index := range indices {
			if index >= len(record) {
				return fmt.Errorf("Not enough records on line %d in CSV file %s", line, path)
			}
			indexRecord = append(indexRecord, record[index])
		}
		callback(indexRecord)
		line++
	}
	return nil
}

func ReadCSVColumns(path string, columns []string, callback func ([]string)) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to read CSV (%s): %v", path, err)
	}
	defer file.Close()
	ReadCSVColumnsFile(file, path, columns, callback)
	return nil
}

func ReadJSON[T any](path string) T {
	data := ReadFile(path)
	var output T
	err := json.Unmarshal(data, &output)
	if err != nil {
		Fatalf("Failed to deserialize JSON: %v", err)
	}
	return output
}

func LoadConfiguration[T any](path string) *T {
	yamlData := ReadFile(path)
	configuration := new(T)
	err := yaml.Unmarshal(yamlData, configuration)
	if err != nil {
		Fatalf("Failed to unmarshal YAML: %v", err)
	}
	return configuration
}

func Fatalf(format string, arguments ...any) {
	formatted := fmt.Sprintf(format, arguments...)
	fmt.Printf("%s\n", formatted)
	panic(formatted)
}