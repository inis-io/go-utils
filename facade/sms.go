package facade

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	AliYunOpenApi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	AliYunSmsApi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	AliYunUtilV2 "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	AliYunCredential "github.com/aliyun/credentials-go/credentials"
	"github.com/inis-io/aide/dto"
	"github.com/inis-io/aide/utils"
	"github.com/spf13/cast"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	TencentCloud "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"gopkg.in/gomail.v2"
)

var SmsInst = &SmsClass{}

type SmsClass struct {
	// 记录配置 Hash 值，用于检测配置文件是否有变化
	Hash string `json:"hash"`
	// 当前短信配置（由调用方注入）
	Config dto.SmsConfig `json:"config"`
	// 是否已经注入过配置
	HasConfig bool `json:"hasConfig"`
}

func init() { SmsInst.Init() }

// normalizeSmsConfig - 统一配置默认值，避免不同项目接入时行为不一致
func normalizeSmsConfig(config dto.SmsConfig) dto.SmsConfig {
	
	config.Engine.Email = strings.ToLower(strings.TrimSpace(config.Engine.Email))
	if utils.Is.Empty(config.Engine.Email) || config.Engine.Email != "email" {
		config.Engine.Email = "email"
	}

	config.Engine.SMS = strings.ToLower(strings.TrimSpace(config.Engine.SMS))
	switch config.Engine.SMS {
	case "aliyun", "tencent", "smsbao":
	default:
		config.Engine.SMS = "aliyun"
	}

	if utils.Is.Empty(config.Email.Host) {
		config.Email.Host = "smtp.qq.com"
	}
	if config.Email.Port <= 0 {
		config.Email.Port = 465
	}
	if utils.Is.Empty(config.Email.Nickname) {
		config.Email.Nickname = "邮件昵称"
	}
	if utils.Is.Empty(config.Email.Subject) {
		config.Email.Subject = "邮件主题"
	}

	if utils.Is.Empty(config.AliYun.Endpoint) {
		config.AliYun.Endpoint = "dysmsapi.aliyuncs.com"
	}
	if utils.Is.Empty(config.Tencent.Endpoint) {
		config.Tencent.Endpoint = "sms.tencentcloudapi.com"
	}
	if utils.Is.Empty(config.Tencent.Region) {
		config.Tencent.Region = "ap-guangzhou"
	}
	if utils.Is.Empty(config.Smsbao.BaseUrl) {
		config.Smsbao.BaseUrl = "https://api.smsbao.com"
	}

	if utils.Is.Empty(config.Hash) {
		config.Hash = utils.Hash.Sum32(utils.Json.Encode(config))
	}

	return config
}

// defaultSmsConfig - 获取默认短信配置
func defaultSmsConfig() dto.SmsConfig {
	return normalizeSmsConfig(dto.SmsConfig{})
}

// useDefaultSms - 使用默认配置激活短信服务
func useDefaultSms() {
	conf := defaultSmsConfig()
	setActiveSms(conf)
}

// setActiveSms - 按配置切换当前活动短信实现
func setActiveSms(config dto.SmsConfig) {
	conf := normalizeSmsConfig(config)

	GoMail = &GoMailClass{}
	
	if strings.ToLower(conf.Engine.Email) == "email" { GoMail.Init() }

	SmsAliYun  = nil
	SmsTencent = nil
	SmsBao     = nil
	SMS        = GoMail

	switch strings.ToLower(conf.Engine.SMS) {
	case "tencent":
		SmsTencent = &SmsTencentClass{}
		SmsTencent.Init()
	case "smsbao":
		SmsBao = &SmsBaoClass{}
		SmsBao.Init()
	default:
		SmsAliYun = &SmsAliYunClass{}
		SmsAliYun.Init()
	}
}

// NewSMS - 创建SMS实例
/**
 * @param mode 驱动模式
 * @return SmsAPI
 * @example：
 * 1. sms := facade.NewSMS("email")
 * 2. sms := facade.NewSMS(facade.SMSModeEmail)
 */
func NewSMS(mode any) SmsAPI {
	switch strings.ToLower(cast.ToString(mode)) {
	case "email":
		SMS = GoMail
	case "aliyun":
		if SmsAliYun != nil {
			SMS = SmsAliYun
		} else {
			SMS = GoMail
		}
	case "tencent":
		if SmsTencent != nil {
			SMS = SmsTencent
		} else {
			SMS = GoMail
		}
	case "smsbao":
		if SmsBao != nil {
			SMS = SmsBao
		} else {
			SMS = GoMail
		}
	default:
		SMS = GoMail
	}
	return SMS
}

// newSMSWithConfig - 使用传入配置重建短信服务并返回指定驱动
func newSMSWithConfig(config dto.SmsConfig, mode string) SmsAPI {
	SmsInst.Init(config)
	return NewSMS(mode)
}

// setConfig - 注入短信配置
func (this *SmsClass) setConfig(config dto.SmsConfig) *SmsClass {
	this.Config = normalizeSmsConfig(config)
	this.HasConfig = true
	return this
}

// ReloadIfChanged - 当配置发生变化时重新加载短信服务
func (this *SmsClass) ReloadIfChanged(config ...dto.SmsConfig) {

	if len(config) > 0 { this.setConfig(config[0]) }

	if !this.HasConfig { return }

	// hash 变化，说明配置有更新
	if this.Hash != this.Config.Hash { this.Init() }
}

// Init 初始化 SMS
func (this *SmsClass) Init(config ...dto.SmsConfig) {
	
	if len(config) > 0 {
		this.setConfig(config[0])
	}

	if !this.HasConfig {
		useDefaultSms()
		return
	}

	this.Config = normalizeSmsConfig(this.Config)
	this.Hash = this.Config.Hash
	setActiveSms(this.Config)
}

// SMS - SMS实例
/**
 * @return SmsAPI
 * @example：
 * sms := facade.SMS.VerifyCode("手机号", "验证码")
 */
var SMS SmsAPI

// GoMail 	   - GoMail邮件服务
var GoMail = &GoMailClass{}

// SmsAliYun   - 阿里云短信
var SmsAliYun  *SmsAliYunClass

// SmsTencent  - 腾讯云短信
var SmsTencent *SmsTencentClass

// SmsBao 	   - 短信宝
var SmsBao *SmsBaoClass

// SmsAPI - 短信接口
type SmsAPI interface {
	// Target - 目标手机号
	Target(target string) SmsAPI
	// Code - 自定义验证码
	Code(code string) SmsAPI
	// Len - 验证码长度
	Len(length int) SmsAPI
	// Send - 发送验证码
	Send(target ...any) *dto.SmsResp
	// Subject - 主题（标题）
	Subject(subject string) SmsAPI
	// SetBody - 设置参数体
	SetBody(body dto.SmsBody) SmsAPI
	// NewCache - 新建缓存
	NewCache(config dto.SmsConfig) SmsAPI
}

// ================================== GoMail邮件服务 - 开始 ==================================

// GoMailClass - GoMail邮件服务
type GoMailClass struct {
	// 邮件客户端
	Client *gomail.Dialer
	// 参数
	Body   dto.SmsBody
}

// Init 初始化 邮件服务
func (this *GoMailClass) Init() {
	this.Client = gomail.NewDialer(SmsInst.Config.Email.Host, SmsInst.Config.Email.Port, SmsInst.Config.Email.Account, SmsInst.Config.Email.Password)
	// 确保在STARTTLS/implicit-TLS流程中，TLS握手使用正确的服务器名称。
	this.Client.TLSConfig = &tls.Config{ServerName: SmsInst.Config.Email.Host}
	this.SetBody(dto.SmsBody{Length: 6, Expired: 5})
	this.Body.Template = dto.TempEmailCode
}

// Code - 自定义验证码
func (this *GoMailClass) Code(code string) SmsAPI {
	this.Body.Code = code
	return this
}

// Len - 验证码长度
func (this *GoMailClass) Len(length int) SmsAPI {
	this.Body.Length = length
	return this
}

// Target - 目标邮箱
func (this *GoMailClass) Target(target string) SmsAPI {
	this.Body.Target = target
	return this
}

// Subject - 主题（标题）
func (this *GoMailClass) Subject(subject string) SmsAPI {
	this.Body.Subject = subject
	return this
}

// Nickname - 昵称（发件人）
func (this *GoMailClass) Nickname(nickname string) SmsAPI {
	this.Body.Nickname = nickname
	return this
}

// Send - 发送验证码
func (this *GoMailClass) Send(target ...any) (response *dto.SmsResp) {

	response = &dto.SmsResp{}

	// 这里的 target 是邮箱地址 - 优先级最高
	if len(target) > 0 {
		this.Body.Target = cast.ToString(target[0])
	}

	socialType, err := utils.Identify.EmailOrPhone(this.Body.Target)
	// 如果不是邮箱或手机号
	if err != nil {
		response.Error = err
		return
	}

	if socialType == "phone" {
		return NewSMS(SmsInst.Config.Engine.SMS).Send(this.Body.Target)
	}

	// 如果自定义验证码为空，则生成一个验证码
	if utils.Is.Empty(this.Body.Code) {
		this.Body.Code = utils.Rand.Code(this.Body.Length)
	}

	subject := utils.Default(this.Body.Subject, SmsInst.Config.Email.Subject)
	nickname := utils.Default(this.Body.Nickname, SmsInst.Config.Email.Nickname)

	item := gomail.NewMessage()
	// 设置邮件内容类型
	item.SetHeader("Content-Type", "text/html; charset=UTF-8")
	// 设置发件人
	item.SetAddressHeader("From", SmsInst.Config.Email.Account, utils.Default(this.Body.Nickname, SmsInst.Config.Email.Nickname))
	// 发送给多个用户
	item.SetHeader("To", this.Body.Target)
	// 设置邮件主题
	item.SetHeader("Subject", subject)
	// 替换验证码
	temp := utils.Replace(this.Body.Template, map[string]any{
		"${title}":    this.Body.Title,
		"${code}":     this.Body.Code,
		"${subject}":  subject,
		"${nickname}": nickname,
		"${username}": this.Body.Username,
		"${expired}":  this.Body.Expired,
		"${email}":    SmsInst.Config.Email.Account,
		"${address}":  this.Body.Address,
		"${year}":     time.Now().Format("2006"),
	})
	// 设置邮件正文
	item.SetBody("text/html", temp)

	// 发送邮件
	// 使用明确的 “拨号 + 发送” 操作，以确保在执行AUTH和MAIL命令之前完成SMTP握手（EHLO/STARTTLS）。
	sender, err := this.Client.Dial()
	if err != nil {
		response.Error = err
		return response
	}
	defer func() { _ = sender.Close() }()

	if err := gomail.Send(sender, item); err != nil {
		response.Error = err
		return response
	}
	response.VerifyCode = this.Body.Code

	// 重置参数
	this.reset()
	return response
}

// SetBody - 设置参数体
func (this *GoMailClass) SetBody(body dto.SmsBody) SmsAPI {
	this.Body.Code = utils.Default(body.Code, this.Body.Code)
	this.Body.Length = utils.Default(body.Length, this.Body.Length)
	this.Body.Target = utils.Default(body.Target, this.Body.Target)
	this.Body.Title = utils.Default(body.Title, this.Body.Title)
	this.Body.Subject = utils.Default(body.Subject, this.Body.Subject)
	this.Body.Nickname = utils.Default(body.Nickname, this.Body.Nickname)
	this.Body.Username = utils.Default(body.Username, this.Body.Username)
	if body.Expired != 0 {
		this.Body.Expired = body.Expired
	}
	this.Body.Address = utils.Default(body.Address, this.Body.Address)
	return this
}

// NewCache - 使用传入配置创建短信实例
func (this *GoMailClass) NewCache(config dto.SmsConfig) SmsAPI {
	return newSMSWithConfig(config, "email")
}

// 重置参数
func (this *GoMailClass) reset() {
	this.Body.Code = ""
	this.Body.Length = 6
	this.Body.Target = ""
	this.Body.Subject = ""
	this.Body.Nickname = ""
	this.Body.Username = ""
	this.Body.Expired = 5
	this.Body.Address = ""
}

// ================================== 阿里云短信 - 开始 ==================================

// SmsAliYunClass - 阿里云短信
type SmsAliYunClass struct {
	Client *AliYunSmsApi.Client
	// 短信请求
	Body dto.SmsBody
}

// Init 初始化 阿里云短信
func (this *SmsAliYunClass) Init() {

	// 创建访问凭证
	credential, err := AliYunCredential.NewCredential(nil)
	// 凭证创建失败
	if err != nil {
		return
	}

	// 创建客户端
	client, err := AliYunSmsApi.NewClient(&AliYunOpenApi.Config{
		Credential: credential,
		// 访问的域名 dysmsapi.aliyuncs.com
		Endpoint: tea.String(utils.Default(SmsInst.Config.AliYun.Endpoint, "dysmsapi.aliyuncs.com")),
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(SmsInst.Config.AliYun.AccessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(SmsInst.Config.AliYun.AccessKeySecret),
	})
	// 客户端创建失败
	if err != nil {
		return
	}

	this.Client = client
	this.SetBody(dto.SmsBody{Length: 6, Expired: 5})
}

// ApiInfo - 接口信息
func (this *SmsAliYunClass) ApiInfo() (result *AliYunOpenApi.Params) {
	return &AliYunOpenApi.Params{
		// 接口名称
		Action: tea.String("SendSms"),
		// 接口版本
		Version: tea.String("2017-05-25"),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// 接口 PATH
		Pathname: tea.String("/"),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}
}

// Code - 自定义验证码
func (this *SmsAliYunClass) Code(code string) SmsAPI {
	this.Body.Code = code
	return this
}

// Len - 验证码长度
func (this *SmsAliYunClass) Len(length int) SmsAPI {
	this.Body.Length = length
	return this
}

// Target - 目标手机号
func (this *SmsAliYunClass) Target(target string) SmsAPI {
	this.Body.Target = target
	return this
}

// Subject - 主题（标题）
func (this *SmsAliYunClass) Subject(subject string) SmsAPI {
	this.Body.Subject = subject
	return this
}

// Nickname - 昵称（发件人）
func (this *SmsAliYunClass) Nickname(nickname string) SmsAPI {
	this.Body.Nickname = nickname
	return this
}

// Send - 发送验证码
func (this *SmsAliYunClass) Send(target ...any) (response *dto.SmsResp) {

	response = &dto.SmsResp{}

	// 这里的 target 是手机号 - 优先级最高
	if len(target) > 0 {
		this.Body.Target = cast.ToString(target[0])
	}

	// 如果不是邮箱或手机号
	if attr, err := utils.Identify.EmailOrPhone(this.Body.Target); err != nil {
		response.Error = err
		return
	} else if attr == "email" {
		return NewSMS(SmsInst.Config.Engine.Email).Send(this.Body.Target)
	}

	// 如果自定义验证码为空，则生成一个验证码
	if utils.Is.Empty(this.Body.Code) {
		this.Body.Code = utils.Rand.Code(this.Body.Length)
	}

	params := &AliYunSmsApi.SendSmsRequest{
		PhoneNumbers: tea.String(this.Body.Target),
		SignName:     tea.String(SmsInst.Config.AliYun.SignName),
		TemplateCode: tea.String(SmsInst.Config.AliYun.VerifyCode),
		TemplateParam: tea.String(utils.Json.Encode(map[string]any{
			"code": this.Body.Code,
			"time": this.Body.Expired,
		})),
	}

	resp, err := this.Client.SendSmsWithOptions(params, &AliYunUtilV2.RuntimeOptions{})
	if err != nil {
		response.Error = err
		return response
	}

	if strings.ToLower(*resp.Body.Code) != "ok" {
		response.Error = errors.New(cast.ToString(*resp.Body.Message))
		return response
	}

	response.Result = cast.ToStringMap(*resp.Body)
	response.Text = utils.Json.Encode(*resp.Body)
	response.VerifyCode = this.Body.Code

	// 重置参数
	this.reset()
	return response
}

// SetBody - 设置参数体
func (this *SmsAliYunClass) SetBody(body dto.SmsBody) SmsAPI {
	this.Body.Code = utils.Default(body.Code, this.Body.Code)
	this.Body.Length = utils.Default(body.Length, this.Body.Length)
	this.Body.Target = utils.Default(body.Target, this.Body.Target)
	this.Body.Title = utils.Default(body.Title, this.Body.Title)
	this.Body.Subject = utils.Default(body.Subject, this.Body.Subject)
	this.Body.Nickname = utils.Default(body.Nickname, this.Body.Nickname)
	this.Body.Username = utils.Default(body.Username, this.Body.Username)
	if body.Expired != 0 {
		this.Body.Expired = body.Expired
	}
	this.Body.Address = utils.Default(body.Address, this.Body.Address)
	return this
}

// NewCache - 使用传入配置创建短信实例
func (this *SmsAliYunClass) NewCache(config dto.SmsConfig) SmsAPI {
	return newSMSWithConfig(config, "aliyun")
}

// 重置参数
func (this *SmsAliYunClass) reset() {
	this.Body.Code = ""
	this.Body.Length = 6
	this.Body.Target = ""
	this.Body.Subject = ""
	this.Body.Nickname = ""
	this.Body.Username = ""
	this.Body.Expired = 5
	this.Body.Address = ""
}

// ================================== 腾讯云短信 - 开始 ==================================

// SmsTencentClass - 腾讯云短信
type SmsTencentClass struct {
	Client *TencentCloud.Client
	// 短信请求
	Body dto.SmsBody
}

// Init 初始化 腾讯云短信
func (this *SmsTencentClass) Init() {

	credential := common.NewCredential(SmsInst.Config.Tencent.SecretId, SmsInst.Config.Tencent.SecretKey)
	clientProfile := profile.NewClientProfile()
	// sms.tencentcloudapi.com
	clientProfile.HttpProfile.Endpoint = SmsInst.Config.Tencent.Endpoint
	// ap-guangzhou
	client, err := TencentCloud.NewClient(credential, SmsInst.Config.Tencent.Region, clientProfile)
	
	if err != nil { return }
	
	this.Client = client
	this.SetBody(dto.SmsBody{Length: 6, Expired: 5})
}

// Code - 自定义验证码
func (this *SmsTencentClass) Code(code string) SmsAPI {
	this.Body.Code = code
	return this
}

// Len - 验证码长度
func (this *SmsTencentClass) Len(length int) SmsAPI {
	this.Body.Length = length
	return this
}

// Target - 目标手机号
func (this *SmsTencentClass) Target(target string) SmsAPI {
	this.Body.Target = target
	return this
}

// Subject - 主题（标题）
func (this *SmsTencentClass) Subject(subject string) SmsAPI {
	this.Body.Subject = subject
	return this
}

// Nickname - 昵称（发件人）
func (this *SmsTencentClass) Nickname(nickname string) SmsAPI {
	this.Body.Nickname = nickname
	return this
}

// Send - 发送验证码
func (this *SmsTencentClass) Send(target ...any) (response *dto.SmsResp) {

	response = &dto.SmsResp{}

	// 这里的 target 是手机号 - 优先级最高
	if len(target) > 0 {
		this.Body.Target = cast.ToString(target[0])
	}

	socialType, err := utils.Identify.EmailOrPhone(this.Body.Target)
	// 如果不是邮箱或手机号
	if err != nil {
		response.Error = err
		return
	}

	if socialType == "email" {
		return NewSMS(SmsInst.Config.Engine.Email).Send(this.Body.Target)
	}

	// 如果自定义验证码为空，则生成一个验证码
	if utils.Is.Empty(this.Body.Code) {
		this.Body.Code = utils.Rand.Code(this.Body.Length)
	}

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := TencentCloud.NewSendSmsRequest()

	request.PhoneNumberSet = common.StringPtrs([]string{this.Body.Target})
	request.SmsSdkAppId = common.StringPtr(SmsInst.Config.Tencent.SmsSdkAppId)
	request.SignName = common.StringPtr(SmsInst.Config.Tencent.SignName)
	request.TemplateId = common.StringPtr(SmsInst.Config.Tencent.VerifyCode)
	request.TemplateParamSet = common.StringPtrs([]string{this.Body.Code})

	item, err := this.Client.SendSms(request)

	if err != nil {
		response.Error = err
		return response
	}

	if item.Response == nil {
		response.Error = errors.New("response is nil")
		return response
	}

	if len(item.Response.SendStatusSet) == 0 {
		response.Error = errors.New("response send status set is nil")
		return response
	}

	if *item.Response.SendStatusSet[0].Code != "Ok" {
		response.Error = errors.New(cast.ToString(item.Response.SendStatusSet[0].Message))
		return response
	}

	response.VerifyCode = this.Body.Code
	response.Text = item.ToJsonString()
	response.Result = utils.Json.Decode(item.ToJsonString())

	// 重置参数
	this.reset()
	return response
}

// SetBody - 设置参数体
func (this *SmsTencentClass) SetBody(body dto.SmsBody) SmsAPI {
	this.Body.Code = utils.Default(body.Code, this.Body.Code)
	this.Body.Length = utils.Default(body.Length, this.Body.Length)
	this.Body.Target = utils.Default(body.Target, this.Body.Target)
	this.Body.Title = utils.Default(body.Title, this.Body.Title)
	this.Body.Subject = utils.Default(body.Subject, this.Body.Subject)
	this.Body.Nickname = utils.Default(body.Nickname, this.Body.Nickname)
	this.Body.Username = utils.Default(body.Username, this.Body.Username)
	if body.Expired != 0 {
		this.Body.Expired = body.Expired
	}
	this.Body.Address = utils.Default(body.Address, this.Body.Address)
	return this
}

// NewCache - 使用传入配置创建短信实例
func (this *SmsTencentClass) NewCache(config dto.SmsConfig) SmsAPI {
	return newSMSWithConfig(config, "tencent")
}

// 重置参数
func (this *SmsTencentClass) reset() {
	this.Body.Code = ""
	this.Body.Length = 6
	this.Body.Target = ""
	this.Body.Subject = ""
	this.Body.Nickname = ""
	this.Body.Username = ""
	this.Body.Expired = 5
	this.Body.Address = ""
}

// ================================== 短信宝 - 开始 ==================================

// SmsBaoClass - 短信宝
type SmsBaoClass struct {
	// 短信宝账号
	Account string
	// 短信宝API密钥
	ApiKey string
	// 短信签名
	SignName string
	// 接口地址
	BaseUrl string
	// 短信请求
	Body dto.SmsBody
}

// Init 初始化 短信宝
func (this *SmsBaoClass) Init() {
	this.Account = SmsInst.Config.Smsbao.Account
	this.ApiKey = SmsInst.Config.Smsbao.ApiKey
	this.SignName = SmsInst.Config.Smsbao.SignName
	this.BaseUrl = SmsInst.Config.Smsbao.BaseUrl
	this.SetBody(dto.SmsBody{Length: 6, Expired: 5})
	this.Body.Template = fmt.Sprintf("【%s】您的验证码是：${code}，有效期5分钟。（打死也不要把验证码告诉别人）", this.SignName)
}

// Code - 自定义验证码
func (this *SmsBaoClass) Code(code string) SmsAPI {
	this.Body.Code = code
	return this
}

// Len - 验证码长度
func (this *SmsBaoClass) Len(length int) SmsAPI {
	this.Body.Length = length
	return this
}

// Target - 目标手机号
func (this *SmsBaoClass) Target(target string) SmsAPI {
	this.Body.Target = target
	return this
}

// Subject - 主题（标题）
func (this *SmsBaoClass) Subject(subject string) SmsAPI {
	this.Body.Subject = subject
	return this
}

// Nickname - 昵称（发件人）
func (this *SmsBaoClass) Nickname(nickname string) SmsAPI {
	this.Body.Nickname = nickname
	return this
}

// Send - 发送验证码
func (this *SmsBaoClass) Send(target ...any) (response *dto.SmsResp) {

	response = &dto.SmsResp{}

	// 这里的 target 是手机号 - 优先级最高
	if len(target) > 0 {
		this.Body.Target = cast.ToString(target[0])
	}

	socialType, err := utils.Identify.EmailOrPhone(this.Body.Target)
	// 如果不是邮箱或手机号
	if err != nil {
		response.Error = err
		return
	}

	if socialType == "email" {
		return NewSMS(SmsInst.Config.Engine.Email).Send(this.Body.Target)
	}

	// 如果自定义验证码为空，则生成一个验证码
	if utils.Is.Empty(this.Body.Code) {
		this.Body.Code = utils.Rand.Code(this.Body.Length)
	}

	if utils.Is.Empty(this.ApiKey) {
		response.Error = errors.New("API密钥不能为空")
		return
	}

	if utils.Is.Empty(this.Account) {
		response.Error = errors.New("账号不能为空")
		return
	}

	item := utils.Curl(utils.CurlRequest{
		Method: "GET",
		Url:    fmt.Sprintf("%s/sms", this.BaseUrl),
		Query: map[string]any{
			"u": this.Account,
			"p": this.ApiKey,
			"m": this.Body.Target,
			"c": utils.Replace(this.Body.Template, map[string]any{
				"${code}": this.Body.Code,
			}),
		},
	}).Send()

	if item.Error != nil {
		response.Error = item.Error
		return
	}

	if cast.ToInt(item.Text) != 0 {
		response.Error = errors.New("发送失败")
		return
	}

	response.VerifyCode = this.Body.Code
	response.Text = item.Text

	// 重置参数
	this.reset()
	return response
}

// SetBody - 设置参数体
func (this *SmsBaoClass) SetBody(body dto.SmsBody) SmsAPI {
	this.Body.Code = utils.Default(body.Code, this.Body.Code)
	this.Body.Length = utils.Default(body.Length, this.Body.Length)
	this.Body.Target = utils.Default(body.Target, this.Body.Target)
	this.Body.Title = utils.Default(body.Title, this.Body.Title)
	this.Body.Subject = utils.Default(body.Subject, this.Body.Subject)
	this.Body.Nickname = utils.Default(body.Nickname, this.Body.Nickname)
	this.Body.Username = utils.Default(body.Username, this.Body.Username)
	if body.Expired != 0 {
		this.Body.Expired = body.Expired
	}
	this.Body.Address = utils.Default(body.Address, this.Body.Address)
	return this
}

// NewCache - 使用传入配置创建短信实例
func (this *SmsBaoClass) NewCache(config dto.SmsConfig) SmsAPI {
	return newSMSWithConfig(config, "smsbao")
}

// 重置参数
func (this *SmsBaoClass) reset() {
	this.Body.Code = ""
	this.Body.Length = 6
	this.Body.Target   = ""
	this.Body.Subject = ""
	this.Body.Nickname = ""
	this.Body.Username = ""
	this.Body.Expired = 5
	this.Body.Address = ""
}
