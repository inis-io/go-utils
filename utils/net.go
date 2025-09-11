package utils

import (
	"net"
	"sync"
	"time"
	
	"github.com/spf13/cast"
)

// Net - 网络
var Net *NetClass

type NetClass struct {}

func (this *NetClass) Tcping(host any, opts ...map[string]any) (ok bool, detail []map[string]any) {

	if len(opts) == 0 {
		opts = append(opts, map[string]any{
			"count":   4,
			"timeout": 5,
		})
	}

	opt := opts[0]
	count := cast.ToInt(opt["count"])
	timeout := cast.ToDuration(opt["timeout"]) * time.Second

	// 读写锁 - 防止并发写入
	type wrLock struct {
		Lock *sync.RWMutex
		Data []map[string]any
	}

	wr := wrLock{
		Data: make([]map[string]any, 0),
		Lock: &sync.RWMutex{},
	}
	wg := sync.WaitGroup{}

	wg.Add(count)

	for i := 0; i < count; i++ {

		go func() {

			defer wg.Done()
			// 加锁
			wr.Lock.Lock()
			// 解锁
			defer wr.Lock.Unlock()

			start := time.Now()
			_, err := net.DialTimeout("tcp", cast.ToString(host), timeout)

			if err != nil {
				wr.Data = append(wr.Data, map[string]any{
					"host":    host,
					"status":  false,
					"waist":   time.Now().Sub(start).Milliseconds(),
					"message": "Site unreachable, error: " + err.Error(),
				})
				return
			}

			wr.Data = append(wr.Data, map[string]any{
				"host":    host,
				"status":  true,
				"waist":   time.Now().Sub(start).Milliseconds(),
				"message": "tcp server is ok",
			})
		}()
	}

	wg.Wait()

	// 只要有一个 ping 成功就返回 true
	for _, val := range wr.Data {
		if cast.ToBool(val["status"]) {
			return true, wr.Data
		}
	}

	return false, wr.Data
}
