package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Get - 获取
var Get *GetClass

type GetClass struct{}

// Type - 获取数据类型
func (this *GetClass) Type(value any) (result string) {
	result, _ = typeof(value)
	return result
}

// Ip - 获取客户端IP
func (this *GetClass) Ip(key ...string) (result any) {

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

// Mac - 获取本机MAC
func (this *GetClass) Mac() (result string) {

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

// Pid - 获取进程ID
func (this *GetClass) Pid() (result int) {
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		return 0
	}
	return process.Pid
}

// Pwd - 获取当前目录
func (this *GetClass) Pwd() (result string) {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// Env - 获取环境变量
func (this *GetClass) Env() (config *any) {
	
	// 初始化 Viper
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return
	}
	// 反序列化配置
	if err := viper.Unmarshal(&config); err != nil {
		return
	}
	
	return config
}

// HostProtocol 使用正则表达式检查地址是否为TLS并返回去除协议头后的内容
// 输入: 地址字符串(可能包含http://, https://或不包含协议头)
// 返回: (是否为TLS, 去除协议头后的地址, 错误)
func (this *GetClass) HostProtocol(url string) (isTLS bool, host string) {
	
	// 编译正则表达式，匹配http://或https://开头，并捕获剩余部分
	item := regexp.MustCompile(`^(?i)(https?://)?(.*)$`)
	
	// 查找匹配项
	matches := item.FindStringSubmatch(url)
	if len(matches) < 3 { return false, "" }
	
	// matches[1]是协议部分(如果有)，matches[2]是主机部分
	protocol := strings.ToLower(matches[1])
	
	// 判断是否为TLS
	return protocol == "https://", matches[2]
}