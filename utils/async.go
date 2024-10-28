package utils

import (
	"github.com/spf13/cast"
	"reflect"
	"sync"
)

type AsyncClass[T any] struct {
	// 读写锁
	Mutex sync.RWMutex
	// 等待组
	Wait sync.WaitGroup
	// 数据
	Data T
}

// Async - 异步数据
func Async[T any]() *AsyncClass[T] {

	var data T
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		data = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(data).Elem()), 0, 0).Interface().(T)
	case reflect.Map:
		data = reflect.MakeMap(reflect.MapOf(reflect.TypeOf(data).Key(), reflect.TypeOf(data).Elem())).Interface().(T)
	default:
		data = reflect.Zero(reflect.TypeOf(data)).Interface().(T)
	}

	return &AsyncClass[T]{
		Mutex: sync.RWMutex{},
		Wait:  sync.WaitGroup{},
		Data:  data,
	}
}

// Get - 获取数据
func (this *AsyncClass[T]) Get(key string) any {

	defer this.Mutex.Unlock()
	this.Mutex.Lock()

	if Is.Empty(this.Data) {
		return nil
	}

	typeof := reflect.TypeOf(this.Data)
	if typeof.Kind() == reflect.Map && typeof.Key().Kind() == reflect.String {
		item := cast.ToStringMap(this.Data)
		return item[key]
	} else if typeof.Kind() == reflect.Slice {
		item := cast.ToSlice(this.Data)
		return item[cast.ToInt(key)]
	}

	return this.Data
}

// Set - 设置数据
func (this *AsyncClass[T]) Set(key string, val any) {

	defer this.Mutex.Unlock()
	this.Mutex.Lock()

	typeof := reflect.TypeOf(this.Data)
	if typeof.Kind() == reflect.Map && typeof.Key().Kind() == reflect.String {
		item := cast.ToStringMap(this.Data)
		item[key] = val
	} else if typeof.Kind() == reflect.Slice {
		index := cast.ToInt(key)
		item := cast.ToSlice(this.Data)
		if len(item) > index {
			item[index] = val
		}
	} else {
		this.Data = val.(T)
	}
}

// Has - 判断是否存在
func (this *AsyncClass[T]) Has(key string) (ok bool) {

	defer this.Mutex.Unlock()
	this.Mutex.Lock()

	if Is.Empty(this.Data) {
		return false
	}

	typeof := reflect.TypeOf(this.Data)
	if typeof.Kind() == reflect.Map && typeof.Key().Kind() == reflect.String {
		item := cast.ToStringMap(this.Data)
		_, ok = item[key]
	} else {
		ok = true
	}

	return ok
}

// Result - 获取所有数据
func (this *AsyncClass[T]) Result() T {
	defer this.Mutex.Unlock()
	this.Mutex.Lock()
	return this.Data
}
