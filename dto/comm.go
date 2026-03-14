package dto

// AesResp - 响应输出
type AesResp struct {
	// 加密后的字节
	Byte []byte
	// 加密后的字符串
	Text string
	// 错误信息
	Error error
}