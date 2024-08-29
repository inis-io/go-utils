package utils

// Resp - 返回数据结构
type Resp struct {
	// 消息
	Msg   string `json:"msg"`
	// 状态码
	Code  int    `json:"code"`
	// 数据
	Data  any    `json:"data"`
}