package utils

import (
	"time"
	
	JWT "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/inis-io/aide/dto"
)

type JwtClass struct {
	// 参数
	Body    dto.JwtBody `json:"body"`
	// 目标
	Dest    any         `json:"dest"`
	// 记录配置 Hash 值，用于检测配置文件是否有变化
	Hash    string      `json:"hash"`
}

var Jwt = &JwtClass{
	Body: dto.JwtBody{
		Key: "inis-jwt-key",
		Issuer : "inis-issuer",
		Subject: "inis-subject",
		Expired: 7 * 24 * 60 * 60,
	},
}

// SetBody -
func (this *JwtClass) SetBody(body dto.JwtBody) *JwtClass {
	this.Body = body
	return this
}

// GetHash - 获取Hash值
func (this *JwtClass) GetHash() string {
	return this.Hash
}

// SetHash - 设置Hash值
func (this *JwtClass) SetHash(hash string) *JwtClass {
	this.Hash = hash
	return this
}

// Encode - 创建JWT
func (this *JwtClass) Encode(data any) (resp dto.JwtResp, err error) {
	
	type JwtClaims struct {
		Data  any `json:"data"`
		JWT.RegisteredClaims
	}
	
	IssuedAt   := JWT.NewNumericDate(time.Now())
	ExpiresAt  := JWT.NewNumericDate(this.GetExpired())
	
	value, err := JWT.NewWithClaims(JWT.SigningMethodHS256, JwtClaims{
		Data: data,
		RegisteredClaims: JWT.RegisteredClaims{
			IssuedAt:  IssuedAt,           // 签发时间戳
			ExpiresAt: ExpiresAt,          // 过期时间戳
			Issuer:    this.Body.Issuer,   // 颁发者签名
			Subject:   this.Body.Subject,  // 签名主题
		},
	}).SignedString([]byte(this.Body.Key))
	
	if err != nil { return resp, err }
	
	return dto.JwtResp{
		Data : data,
		Value: value,
		No: uuid.NewString(),
		Expired: this.GetExpired().Unix(),
	}, nil
}

// Decode - 解析JWT
func (this *JwtClass) Decode(token string) (resp dto.JwtResp, err error) {
	
	type JwtClaims struct {
		Data  any `json:"data"`
		JWT.RegisteredClaims
	}
	
	item, err := JWT.ParseWithClaims(token, &JwtClaims{}, func(token *JWT.Token) (any, error) {
		return []byte(this.Body.Key), nil
	})
	
	if err != nil { return resp, err }
	
	if row, _ := item.Claims.(*JwtClaims); item.Valid {
		
		resp.Data    = row.Data
		resp.Expired = row.RegisteredClaims.ExpiresAt.Time.Unix()
		
		// 解析目标
		if this.Dest != nil { _ = Json.Unmarshal([]byte(Json.Encode(row.Data)), this.Dest) }
	}
	
	return resp, nil
}

// Unmarshal - 解析JWT
func (this *JwtClass) Unmarshal(dest any) *JwtClass {
	this.Dest = dest
	return this
}

// GetExpired - 获取JWT过期时间
func (this *JwtClass) GetExpired() time.Time {
	return time.Now().Add(time.Second * time.Duration(this.Body.Expired))
}