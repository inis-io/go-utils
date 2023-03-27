### Curl

```go
package main

import (
	"fmt"
	"github.com/unti-io/go-utils/utils"
)

func main() {
	// 方式一：
	curl := utils.Curl(utils.CurlRequest{
		Method: "GET",
		Url:    "https://v1.hitokoto.cn/",
		Query: map[string]string{
			"encode": "json",
		},
	}).Send()

	fmt.Println(curl.Json)
	
	// 方式二：
	item := utils.Curl().Get("https://v1.hitokoto.cn/").Query("encode": "json").Send()
	fmt.Println(item.Json)
}
```