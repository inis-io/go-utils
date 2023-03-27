package utils

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// GetType - 获取数据类型
func GetType(value any) (result string) {
	result, _ = typeof(value)
	return result
}

// GetIp - 获取客户端IP
func GetIp(key ...string) (result any) {

	type lock struct {
		Lock *sync.RWMutex
		Data map[string]any
	}
	wr := lock{
		Data: make(map[string]any, 2),
		Lock: &sync.RWMutex{},
	}
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		conn, _ := net.Dial("udp", "8.8.8.8:80")
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				fmt.Println("intranet err：", err)
			}
		}(conn)
		localAddr := conn.LocalAddr().String()
		idx := strings.LastIndex(localAddr, ":")
		wr.Lock.Lock()
		wr.Data["intranet"] = localAddr[0:idx]
		wr.Lock.Unlock()
	}()
	
	go func() {
		defer wg.Done()
		// 外网IP - 替代品：https://api.ipify.org https://ipinfo.io/ip https://api.ip.sb/ip
		item := Curl().Get("https://myexternalip.com/raw").Send()
		if item.Error != nil {
			return
		}
		wr.Lock.Lock()
		wr.Data["extranet"] = item.Text
		wr.Lock.Unlock()
	}()

	wg.Wait()

	if len(key) > 0 {
		return wr.Data[key[0]]
	}

	return wr.Data
}