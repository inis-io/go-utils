package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"hash/fnv"
	rand2 "math/rand"
	"time"
	
	"github.com/spf13/cast"
)

type HashClass struct {}

// Hash - 哈希加密
var Hash *HashClass

// Sum32 - 哈希加密
func (this *HashClass) Sum32(text any) (result string) {
	item := fnv.New32()
	_, err := item.Write([]byte(cast.ToString(text)))
	return cast.ToString(Ternary[any](err != nil, nil, item.Sum32()))
}

// Token 生成指定长度的指纹令牌
/**
 * @param value 值
 * @param length 令牌长度，默认长度为 16
 * @example：
 * token := utils.Hash.Token("test", 16)
 */
func (this *HashClass) Token(value any, length int) (result string) {

	// 计算 MD5 哈希值
	hash := md5.Sum([]byte(cast.ToString(value)))
	MD5Hash := hex.EncodeToString(hash[:])

	if length > len(MD5Hash) { length = len(MD5Hash) }
	result = MD5Hash[:length]

	// 确保结果长度满足要求
	for len(result) < length {
		result += this.Token(result, length - len(result))
	}

	return result
}

// Number 生成指定长度的随机数
/**
 * @param length 长度
 * @return result 随机数
 * @example：
 * 1. number := facade.Hash.Number(6)
 */
func (this *HashClass) Number(length any) (result string) {
	
	// 种子
	seed   := Hash.Sum32(fmt.Sprintf("%s-%d-%d", Get.Mac(), Get.Pid(), time.Now().UnixNano()))
	// 使用当前时间戳创建随机数生成器
	source := rand2.New(rand2.NewSource(cast.ToInt64(seed)))

	// 生成一个随机数
	for i := 0; i < cast.ToInt(length); i++ {
		result += fmt.Sprintf("%d", source.Intn(10))
	}

	return result
}

// AESRequest - 请求输入
type AESRequest struct {
	// 16位密钥
	Key string
	// 16位向量
	Iv string
}

// AESResponse - 响应输出
type AESResponse struct {
	// 加密后的字节
	Byte []byte
	// 加密后的字符串
	Text string
	// 错误信息
	Error error
}

// AES - 对称加密
func AES(key, iv any) *AESRequest {
	return &AESRequest{
		Key: cast.ToString(key),
		Iv:  cast.ToString(iv),
	}
}

// Encrypt 加密
func (this *AESRequest) Encrypt(text any) (result *AESResponse) {

	result = &AESResponse{}

	// 拦截异常
	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Errorf("%v", r)
		}
	}()

	block, err := aes.NewCipher([]byte(this.Key))
	if err != nil {
		result.Error = err
	}

	// 每个块的大小
	blockSize := block.BlockSize()
	// 计算需要填充的长度
	padding   := blockSize - len([]byte(cast.ToString(text)))%blockSize

	// 填充
	fill   := append([]byte(cast.ToString(text)), bytes.Repeat([]byte{byte(padding)}, padding)...)
	encode := make([]byte, len(fill))

	item   := cipher.NewCBCEncrypter(block, []byte(this.Iv))
	item.CryptBlocks(encode, fill)

	result.Byte = encode
	result.Text = base64.StdEncoding.EncodeToString(encode)

	return
}

// Decrypt 解密
func (this *AESRequest) Decrypt(text any) (result *AESResponse) {

	result = &AESResponse{}

	// 拦截异常
	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Errorf("%v", r)
		}
	}()

	newText, err := base64.StdEncoding.DecodeString(cast.ToString(text))
	if err != nil {
		result.Error = err
		return
	}

	block, err := aes.NewCipher([]byte(this.Key))
	if err != nil {
		result.Error = err
		return
	}

	// 确保 newText 是 blockSize 的整数倍
	blockSize := block.BlockSize()
	if len(newText)%blockSize != 0 {
		result.Error = errors.New("invalid ciphertext")
		return
	}

	decode := make([]byte, len(newText))
	item := cipher.NewCBCDecrypter(block, []byte(this.Iv))
	item.CryptBlocks(decode, newText)

	// 去除填充
	padding := decode[len(decode)-1]
	result.Byte = decode[:len(decode)-int(padding)]
	result.Text = string(result.Byte)

	return
}

var RSA *RSAClass

type RSAClass struct {}

type RSAResponse struct {
	// 私钥
	PrivateKey string
	// 公钥
	PublicKey string
	// 错误信息
	Error error
	// 文本
	Text string
}

// Generate 生成 RSA 密钥对
/**
 * @name Generate 生成 RSA 密钥对
 * @param bits 位数 1024, 2048, 4096（一般：2048）
 */
func (this *RSAClass) Generate(bits any) (result *RSAResponse) {

	result = &RSAResponse{}

	private, err := rsa.GenerateKey(rand.Reader, cast.ToInt(bits))
	if err != nil {
		result.Error = err
		return result
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(private)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyBytes,
	}

	// 生成私钥
	result.PrivateKey = string(pem.EncodeToMemory(&privateKeyBlock))

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		result.Error = err
		return result
	}
	publicKeyBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyBytes,
	}

	// 生成公钥
	result.PublicKey = string(pem.EncodeToMemory(&publicKeyBlock))

	return result
}

// Encrypt 加密
func (this *RSAClass) Encrypt(publicKey, text string) (result *RSAResponse) {

	result = &RSAResponse{}

	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Errorf("%v", r)
		}
	}()

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		result.Error = errors.New("public key error")
		return result
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		result.Error = err
		return result
	}

	pub := pubInterface.(*rsa.PublicKey)
	encode, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(text))
	if err != nil {
		result.Error = err
		return result
	}

	result.Text = base64.StdEncoding.EncodeToString(encode)

	return result
}

// Decrypt 解密
func (this *RSAClass) Decrypt(privateKey, text string) (result *RSAResponse) {

	result = &RSAResponse{}

	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Errorf("%v", r)
		}
	}()

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		result.Error = errors.New("private key error")
		return result
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		result.Error = err
		return result
	}

	decode, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		result.Error = err
		return result
	}

	encode, err := rsa.DecryptPKCS1v15(rand.Reader, priv, decode)
	if err != nil {
		result.Error = err
		return result
	}

	result.Text = string(encode)

	return result
}

// PublicPem - 输出完整的 PEM 格式公钥证书
func (this *RSAClass) PublicPem(key string) (cert string) {

	// 创建 PEM 格式块
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: []byte(key),
	}

	// 生成完整的 PEM 证书字符串
	var PEM bytes.Buffer
	if err := pem.Encode(&PEM, block); err != nil { return "" }

	return PEM.String()
}

// PrivatePem - 输出完整的 PEM 格式私钥证书
func (this *RSAClass) PrivatePem(key string) (cert string) {

	// 创建 PEM 格式块
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: []byte(key),
	}

	// 生成完整的 PEM 证书字符串
	var PEM bytes.Buffer
	if err := pem.Encode(&PEM, block); err != nil { return "" }

	return PEM.String()
}

type Md5Class struct {}

// Md5 - MD5 加密
var Md5 *Md5Class

// Encrypt 计算字符串的 MD5 哈希值（返回十六进制字符串）
func (this *Md5Class) Encrypt(value string) string {
	// 创建 MD5 哈希对象
	hash := md5.New()
	// 写入数据（可以多次调用 Write 累加数据）
	hash.Write([]byte(value))
	// 计算哈希值，返回 []byte
	hashBytes := hash.Sum(nil)
	// 转为十六进制字符串
	return hex.EncodeToString(hashBytes)
}