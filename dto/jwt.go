package dto

// JwtResp - JWT响应
type JwtResp struct {
	No       string   `json:"no"`
	Value    string   `json:"value"`
	Data     any      `json:"data"`
	Expired  int64    `json:"expired"`
}

type JwtBody struct {
	// 过期时间（秒）
	Expired  int64  `json:"expired"`
	// 颁发者签名
	Issuer   string `json:"issuer"`
	// 主题
	Subject  string `json:"subject"`
	// 密钥
	Key 	 string `json:"key"`
}