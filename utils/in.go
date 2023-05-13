package utils

import (
	"reflect"
	"sync"
)

// InArray - 判断某个值是否在数组中
func InArray[T any](value T, array []T) (ok bool) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	found := false
	for _, val := range array {
		if found {
			break
		}
		wg.Add(1)
		go func(item T) {
			defer wg.Done()
			if reflect.DeepEqual(value, item) {
				mutex.Lock()
				found = true
				mutex.Unlock()
			}
		}(val)
	}
	wg.Wait()
	return found
}