### 简介
> 这是一个 [GoLang](https://golang.org/) 的工具包，包含了许多常用的函数，用于简化开发过程中的重复性工作。

### 安装
```bash
go get github.com/inis-io/go-utils
```

### 使用
> 详细的使用方法请参考 [文档](./document/README.md)

### Storage 快速使用

```go
package main

import (
	"os"

	"github.com/inis-io/aide/dto"
	"github.com/inis-io/aide/facade"
)

func main() {
	// 1) 初始化全局存储（推荐在应用启动时执行一次）
	facade.StorageInst.Init(dto.StorageConfig{
		Engine: "local",
		Local: dto.LocalStorageConfig{Domain: "http://localhost:2000"},
	})

	// 2) 使用全局实例
	file, _ := os.Open("./avatar.png")
	defer file.Close()
	resp := facade.Storage.Dir("avatar").Ext("png").Upload(file)
	_ = resp

	// 3) 按配置创建独立实例（适合多租户或临时切换引擎）
	custom := facade.Storage.NewStorage(dto.StorageConfig{Engine: "local"})
	_ = custom
}
```
