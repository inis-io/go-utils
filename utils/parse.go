package utils

import (
	"github.com/spf13/cast"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// 读写锁 - 防止并发写入
type wrLock struct {
	Lock *sync.RWMutex
	Data map[string]any
}

// Get - 获取数据
func (this *wrLock) get(key string) any {
	defer this.Lock.Unlock()
	this.Lock.Lock()
	return this.Data[key]
}

// Set - 设置数据
func (this *wrLock) set(key string, val any) {
	defer this.Lock.Unlock()
	this.Lock.Lock()
	this.Data[key] = val
}

// Has - 判断是否存在
func (this *wrLock) has(key string) (ok bool) {
	defer this.Lock.Unlock()
	this.Lock.Lock()
	_, ok = this.Data[key]
	return ok
}

// ParseParamsBefore - 解析参数
// 把 Content-Type = application/x-www-form-urlencoded 的参数解析成 object.deep.age = 10 的格式
func ParseParamsBefore(params url.Values) (result map[string]any) {

	wg := sync.WaitGroup{}
	wr := wrLock{
		Data: make(map[string]any, 3),
		Lock: &sync.RWMutex{},
	}
	worker := func(keys string, value []string, wg *sync.WaitGroup) {

		defer wg.Done()
		// 加锁
		wr.Lock.Lock()
		// 解锁
		defer wr.Lock.Unlock()

		// 判断是否为 [] 结尾 - 普通数据
		if strings.HasSuffix(keys, "[]") {
			// 将 [] 替换为 .
			key := strings.Replace(keys, "[]", ".", -1)
			wr.Data[key] = value
			return
		}
		keys = strings.Replace(keys, "[", ".", -1)
		keys = strings.Replace(keys, "]", "", -1)

		// 正则匹配末尾是否为 .[0-9]+ - 【一维数组】
		reg := regexp.MustCompile(`\.[0-9]+$`)
		if reg.MatchString(keys) {

			// 取到 .[0-9]+ 的内容
			index := reg.FindString(keys)
			// 去除 . 只要数字
			index = strings.Replace(index, ".", "", -1)

			// 将 .[0-9]+ 替换为 .
			key := reg.ReplaceAllString(keys, ".")
			// 判断 key 是否存在，如果不存在，则创建 - 【初始化一维数组】
			if _, ok := wr.Data[key]; !ok {

				item := make([]any, cast.ToInt(index)+1)
				item[cast.ToInt(index)] = value[0]
				wr.Data[key] = item

			} else {

				// 如果存在，则追加，先判断 index 是否超出数组长度 - 【追加一维数组】
				if cast.ToInt(index) >= len(wr.Data[key].([]any)) {

					// 超出长度，需要扩容
					item := make([]any, cast.ToInt(index)+1)
					// 将原来的数据复制到新数组
					copy(item, wr.Data[key].([]any))
					// 将新数据追加到新数组
					item[cast.ToInt(index)] = value[0]
					wr.Data[key] = item

				} else {

					// 未超出长度，直接追加
					wr.Data[key].([]any)[cast.ToInt(index)] = value[0]
				}
			}
			return
		}
		// 如果不是 [] 结尾，也不是 .[0-9]+ 结尾，则直接赋值 - 【普通数据】
		wr.Data[keys] = value[0]
	}

	for key, value := range params {
		wg.Add(1)
		go worker(key, value, &wg)
	}

	wg.Wait()

	return wr.Data
}

// ParseParams - 解析参数
// 把 Content-Type = application/x-www-form-urlencoded 的参数解析成 map[string]any
func ParseParams(params map[string]any) (result map[string]any) {

	wg := sync.WaitGroup{}
	wr := wrLock{
		Data: make(map[string]any, 3),
		Lock: &sync.RWMutex{},
	}

	var worker func(params map[string]any, key string, val any)
	worker = func(params map[string]any, key string, val any) {

		// 加锁
		wr.Lock.Lock()
		defer wg.Done()
		// 解锁
		defer wr.Lock.Unlock()

		// 如果 key 的末尾为 . ，表示该 key 为一维数组
		if strings.HasSuffix(key, ".") {
			key = strings.TrimSuffix(key, ".")
		}

		// 通过 . 分割 key
		keys := strings.Split(key, ".")

		// 判断是否为最后一个元素
		if len(keys) == 1 {
			params[keys[0]] = val
			return
		}

		// 如果不是最后一个元素，则继续递归
		if _, ok := params[keys[0]]; !ok {
			params[keys[0]] = make(map[string]any, 3)
		}
		wg.Add(1)
		go worker(cast.ToStringMap(params[keys[0]]), strings.Join(keys[1:], "."), val)
		return
	}

	for key, val := range params {
		wg.Add(1)
		go worker(wr.Data, key, val)
	}

	wg.Wait()

	return wr.Data
}