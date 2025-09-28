package commons

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"
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

func CreateDirectory(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalf("Failed to create directory: %v", err)
	}
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Download(url string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout: httpTimeoutSeconds * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err := fmt.Errorf("Failed to create HTTP request (%s): %v", url, err)
		return "", err
	}
	request.Header.Set("User-Agent", userAgent)
	response, err := client.Do(request)
	if err != nil {
		err := fmt.Errorf("Failed to GET data (%s): %v", url, err)
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		err := fmt.Errorf("Failed to read response (%s): %v", url, err)
		return "", err
	}
	return string(body), nil
}

func DownloadFile(url string, path string) error {
	data, err := Download(url)
	if err != nil {
		return err
	}
	WriteFile(path, data)
	return nil
}