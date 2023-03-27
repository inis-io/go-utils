package main

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
)

func main() {
	item := utils.Curl(utils.CurlRequest{
		Method: "GET",
		Url:    "https://api.inis.cn/api/links/all",
		Query: map[string]string{
			"limit": "2",
		},
	}).Send()
	data := cast.ToStringMap(item.Json["data"])["data"]
	for _, val := range cast.ToSlice(data) {
		fmt.Println(utils.Json.Encode(val), "\n")
	}
}
