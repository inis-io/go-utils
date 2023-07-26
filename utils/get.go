package utils

import (
	"fmt"
	"github.com/spf13/cast"
	"net"
	"os"
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

// GetMac - 获取本机MAC
func GetMac() (result string) {

	interfaces, err := net.Interfaces()

	if err != nil {
		return ""
	}

	for _, item := range interfaces {
		// 过滤掉非物理接口类型
		if item.Flags&net.FlagUp != 0 && item.Flags&net.FlagLoopback == 0 && item.Flags&net.FlagPointToPoint == 0 {
			if value, err := item.Addrs(); err == nil {
				for _, val := range value {
					if IPNet, ok := val.(*net.IPNet); ok && !IPNet.IP.IsLoopback() {
						if mac := item.HardwareAddr; len(mac) > 0 {
							return cast.ToString(mac)
						}
					}
				}
			}
		}
	}

	return ""
}

// GetPid - 获取进程ID
func GetPid() (result int) {
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		return 0
	}
	return process.Pid
}

// GetPwd - 获取当前目录
func GetPwd() (result string) {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}
