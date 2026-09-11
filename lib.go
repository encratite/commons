package commons

import (
	"encoding/json"
	"fmt"
	"runtime"
	"slices"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	httpTimeoutSeconds = 10
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:155.0) Gecko/20100101 Firefox/155.0"
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

func ReadJSON[T any](path string) T {
	data := ReadFile(path)
	var output T
	err := json.Unmarshal(data, &output)
	if err != nil {
		Fatalf("Failed to deserialize JSON: %v", err)
	}
	return output
}

func WriteJSON[T any](data T, path string) {
	output, err := json.Marshal(data)
	if err != nil {
		Fatalf("Failed to serialize JOSN: %v", err)
	}
	WriteFile(path, output)
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