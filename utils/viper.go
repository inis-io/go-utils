package utils

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
)

type ViperModel struct {
	// 配置文件路径
	Path string
	// 配置文件类型
	Mode string
	// 文件名
	Name string
}

type ViperResponse struct {
	// 配置文件内容
	Result map[string]any
	// 错误信息
	Error error
	// viper实例
	Viper *viper.Viper
}

func Viper(model ...ViperModel) *ViperModel {

	var item *ViperModel

	if len(model) > 0 {
		item = &model[0]
	}

	return item
}

func (this *ViperModel) SetPath(path string) *ViperModel {
	this.Path = path
	return this
}

func (this *ViperModel) SetMode(mode string) *ViperModel {
	this.Mode = mode
	return this
}

func (this *ViperModel) SetName(name string) *ViperModel {
	this.Name = name
	return this
}

func (this *ViperModel) Read() (result ViperResponse) {

	item := viper.New()

	if !IsEmpty(this.Path) {
		item.AddConfigPath(this.Path)
	}

	if !IsEmpty(this.Mode) {
		item.SetConfigType(this.Mode)
	}

	if !IsEmpty(this.Name) {
		item.SetConfigName(this.Name)
	}

	result.Viper = item
	result.Error = item.ReadInConfig()
	result.Result = cast.ToStringMap(item.AllSettings())

	return
}

func (this *ViperResponse) Get(key string, def ...any) (result any) {

	if len(def) > 0 {
		result = def[0]
	}

	if this.Error != nil {
		return
	}

	if this.Result == nil {
		return
	}

	if this.Result[key] != nil {
		result = this.Result[key]
	}

	return
}

func (this *ViperResponse) Set(key string, value any) (result ViperResponse) {

	if this.Error != nil {
		return
	}

	if this.Result == nil {
		return
	}

	file, err := os.OpenFile(this.Viper.ConfigFileUsed(), os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		result.Error = err
	}

	this.Result[key] = value
	this.Viper.Set(key, value)
	result.Error = this.Viper.WriteConfigAs(file.Name())

	result = *this

	// 释放资源
	result.Error = file.Close()

	return
}
