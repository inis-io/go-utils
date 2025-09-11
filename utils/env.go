package utils

import (
	"fmt"
	"strings"
	"sync"
	
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// EnvToml - 读取配置文件
func EnvToml(args ...string) (result any) {

	key, def, path := "app", "", "./config/app.toml"

	if len(args) > 0 {
		for k, v := range args {
			if k == 0 {
				key = v
			} else if k == 1 {
				def = v
			} else {
				path = v
			}
		}
	}

	keys := strings.Split(key, ".")

	// 文件路径
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("读取配置文件失败：", err)
	}

	if len(keys) == 1 {

		result := viper.GetStringMap(key)

		if empty := Is.Empty(result); empty {
			return def
		}

		return result

	} else {

		result = viper.GetStringMap(keys[0])[keys[1]]

		if empty := Is.Empty(result); empty {
			return def
		}

		return
	}
}

type EnvClass struct {
	path string
	mode string
}

func Env() *EnvClass {
	return &EnvClass{
		path: "./config/app.toml",
		mode: "toml",
	}
}

func (this *EnvClass) Path(path string) *EnvClass {
	this.path = path
	return this
}

func (this *EnvClass) Mode(mode string) *EnvClass {
	this.mode = mode
	return this
}

func (this *EnvClass) All() (result map[string]any) {

	if this.mode == "toml" {
		return cast.ToStringMap(this.TomlAll())
	}

	return nil
}

func (this *EnvClass) Get(key any, def ...any) (result any) {

	if this.mode == "toml" {
		result = this.TomlAll()
	}

	item := cast.ToStringMap(result)

	// 默认值处理
	if item[cast.ToString(key)] == nil {
		if len(def) > 0 {
			return def[0]
		}
		return nil
	}

	return item[cast.ToString(key)]
}

func (this *EnvClass) TomlAll() (result any) {

	// 文件路径
	viper.SetConfigFile(this.path)

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("读取配置文件失败: %v", err)
		return nil
	}

	toml := viper.AllSettings()

	// 读写锁 - 防止并发写入
	type wrLock struct {
		Lock *sync.RWMutex
		Data map[string]any
	}

	wr := wrLock{
		Data: make(map[string]any),
		Lock: &sync.RWMutex{},
	}
	wg := sync.WaitGroup{}

	for key, val := range toml {
		wg.Add(1)
		go func(key string, val any) {
			defer wg.Done()
			for k, v := range cast.ToStringMap(val) {
				wr.Lock.Lock()
				wr.Data[key+"."+k] = v
				wr.Lock.Unlock()
			}
		}(key, val)
	}

	wg.Wait()

	return wr.Data
}