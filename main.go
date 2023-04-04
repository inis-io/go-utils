package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"os"
)

var ItemToml *utils.ViperResponse

func main() {

	fmt.Println("=========== main ===========")

	ItemToml.Viper.WatchConfig()
	ItemToml.Viper.OnConfigChange(func(event fsnotify.Event) {
		// initToml()
		fmt.Println("文件发生变化：", cast.ToStringMap(ItemToml.Viper.AllSettings()))
		fmt.Println(ItemToml.Get("root", "name"))
	})

	ItemToml.Set("app.name", 111)

	select {}
}

func init() {
	initToml()
}

func initToml() {

	item := utils.Viper(utils.ViperModel{
		Path: ".",
		Mode: "toml",
		Name: "inis",
	}).Read()

	if item.Error != nil {
		// 如果错误中包含文件不存在，则创建文件
		if !os.IsNotExist(item.Error) {
			CreateToml()
		} else {
			fmt.Println("发生错误", item.Error)
			return
		}
	}

	ItemToml = &item
}

func CreateToml() {
	// 文件不存在 - 创建文件
	file, err := os.Create("inis.toml")
	if err != nil {
		fmt.Println("创建文件失败", err)
		return
	}
	// 写入文件
	_, err = file.WriteString(`[app]
name    = "test"
version = "1.0.0"
`)
	if err != nil {
		fmt.Println("写入文件失败", err)
		return
	}
	// 关闭文件
	err = file.Close()
}
