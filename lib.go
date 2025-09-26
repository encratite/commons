package commons

import (
	"runtime"
	"slices"
	"sync"
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