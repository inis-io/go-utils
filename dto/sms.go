package dto

// SmsConfig - 消息推送服务配置（由调用方传入）
type SmsConfig struct {
	// Engine - 驱动
	Engine  SmsEngine        `json:"engine"`
	// Email - 邮件服务配置
	Email   SmsEmailConfig   `json:"email"`
	// AliYun 阿里云短信服务配置
	AliYun  SmsAliYunConfig  `json:"aliyun"`
	// Tencent 腾讯云短信服务配置
	Tencent SmsTencentConfig `json:"tencent"`
	// Smsbao 短信宝短信服务配置
	Smsbao  SmsBaoConfig     `json:"smsbao"`
	// Hash - 计算配置是否发生变更
	Hash    string           `json:"hash"`
}

// SmsEngine - 短信引擎配置
type SmsEngine struct {
	// Email - 邮件
	Email   string `json:"email" default:"email"`
	// SMS   - 短信
	SMS     string `json:"sms"   default:"aliyun"`
}

// SmsEmailConfig - 邮件服务配置
type SmsEmailConfig struct {
	// Host      - 邮件服务器地址
	Host      string `json:"host" comment:"服务地址" validate:"required,host" default:"smtp.qq.com"`
	// Port      - 邮件服务端口
	Port      int    `json:"port" comment:"端口号" validate:"numeric" default:"465"`
	// Account   - 邮件账号
	Account   string `json:"account"  comment:"账号"  validate:"required"`
	// Password  - 服务密码 - 不是邮箱密码
	Password  string `json:"password" comment:"密码"  validate:"required"`
	// Nickname  - 邮件昵称
	Nickname  string `json:"nickname" comment:"昵称"  validate:"max=64" default:"邮件昵称"`
	// Subject   - 邮件主题
	Subject  string  `json:"subject"  comment:"主题"  validate:"max=64" default:"邮件主题"`
}

// SmsAliYunConfig - 阿里云短信服务配置
type SmsAliYunConfig struct {
	// AccessKeyId     - 阿里云AccessKey ID
	AccessKeyId     string `json:"access_key_id" comment:"AccessKey ID" validate:"required,alphaNum"`
	// AccessKeySecret - 阿里云AccessKey Secret
	AccessKeySecret string `json:"access_key_secret" comment:"AccessKey Secret" validate:"required,alphaNum"`
	// Endpoint        - 阿里云短信服务endpoint
	Endpoint        string `json:"endpoint"    comment:"endpoint" validate:"required,host" default:"dysmsapi.aliyuncs.com"`
	// SignName        - 短信签名
	SignName        string `json:"sign_name"   comment:"短信签名" validate:"required"`
	// VerifyCode      - 验证码模板
	VerifyCode      string `json:"verify_code" comment:"验证码模板" validate:"required,alphaDash"`
}

// SmsTencentConfig - 腾讯云短信服务配置
type SmsTencentConfig struct {
	// SecretId        - 腾讯云SecretId
	SecretId        string `json:"secret_id"   comment:"Secret ID"  validate:"required,alphaNum"`
	// SecretKey       - 腾讯云SecretKey
	SecretKey       string `json:"secret_key"  comment:"Secret Key" validate:"required,alphaNum"`
	// Endpoint        - 腾讯云短信服务endpoint
	Endpoint        string `json:"endpoint"    comment:"endpoint"   validate:"required,host" default:"sms.tencentcloudapi.com"`
	// SmsSdkAppId     - 腾讯云短信服务appid
	SmsSdkAppId     string `json:"sms_sdk_app_id" comment:"短信服务AppID" validate:"required,numeric"`
	// SignName        - 短信签名
	SignName        string `json:"sign_name"   comment:"短信签名" validate:"required"`
	// VerifyCode      - 验证码模板id
	VerifyCode      string `json:"verify_code" comment:"验证码模板" validate:"required,numeric"`
	// Region          - 区域
	Region          string `json:"region" comment:"区域" validate:"required,alphaDash" default:"ap-guangzhou"`
}

// SmsBaoConfig - 短信宝短信服务配置
type SmsBaoConfig struct {
	// Account   - 短信宝账号
	Account   string `json:"account"   comment:"短信宝账号" validate:"required,alphaNum"`
	// ApiKey    - API密钥
	ApiKey    string `json:"api_key"   comment:"API 密钥"  validate:"required,alphaNum"`
	// SignName  - 短信签名
	SignName  string `json:"sign_name" comment:"短信签名" validate:"required"`
	// BaseUrl   - 接口地址
	BaseUrl   string `json:"base_url"  comment:"接口地址" validate:"url" default:"https://api.smsbao.com"`
}

// SmsBody - 短信请求参数
type SmsBody struct {
	// Target - 目标手机号或邮箱
	Target   string
	// Code - 自定义验证码
	Code     string
	// Length - 验证码长度
	Length   int
	// Template - 发送模板
	Template string
	// 主题（标题）
	Subject  string
	// 昵称（发件人昵称）
	Nickname string
	// 用户名（收件人昵称）
	Username string
	// 过期时间（分钟）
	Expired  int64
	// 通信地址
	Address  string
	// 标题
	Title    string
}

// SmsResp - 短信响应
type SmsResp struct {
	// 结果
	Result     any
	// 文本
	Text       string
	// 验证码
	VerifyCode string
}

// TempEmailCode - 临时邮箱验证码脚本模板
const TempEmailCode = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style type="text/css">
        body {
            font-family: 'Helvetica Neue', Arial, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 0;
            color: #333333;
            line-height: 1.6;
        }
        
        .container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            border-radius: 8px;
            font-size: 15px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
        
        .header {
            background-color: #4a6bff;
            padding: 30px 20px;
            text-align: center;
            color: white;
        }
        
        .logo {
            font-size: 24px;
            font-weight: bold;
            margin-bottom: 10px;
        }
        
        .content {
            padding: 20px 15px;
        }
        .verification-code {
            background-color: #f8f9fa;
            border-left: 4px solid #4a6bff;
            padding: 15px;
            margin: 25px 0;
            font-size: 24px;
            font-weight: bold;
            color: #4a6bff;
            text-align: center;
            letter-spacing: 5px;
        }
        
        .footer {
            background-color: #f5f5f5;
            padding: 20px;
            text-align: center;
            font-size: 12px;
            color: #999999;
        }
        
        .button {
            display: inline-block;
            background-color: #4a6bff;
            color: white;
            text-decoration: none;
            padding: 12px 25px;
            border-radius: 4px;
            margin: 15px 0;
            font-weight: bold;
        }
        
        .note {
            font-size: 12px;
            color: #666666;
            margin-top: 30px;
        }
        
        .divider {
            border-top: 1px solid #eeeeee;
            margin: 20px 0;
        }
        
        @media only screen and (max-width: 600px) {
            .container {
                margin: 0;
                border-radius: 0;
            }
            
            .content {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">${title}</div>
            <div>${subject}</div>
        </div>
        
        <div class="content">
            <p>尊敬的 <strong>${username}</strong>，您好！</p>
            
            <p>感谢您使用我们的服务，您正在进行邮箱验证，请使用以下验证码完成验证：</p>
            
            <div class="verification-code">
                ${code}
            </div>
            
            <p>该验证码将在 <strong>${expired}分钟</strong> 后失效，请尽快使用。</p>
            
            <p>如果您并未请求此验证码，请忽略此邮件。</p>
            
            <div class="divider"></div>
            
            <div class="note">
                <p>此为系统自动发送的邮件，请勿直接回复。</p>
            </div>
        </div>
        
        <div class="footer">
            <p>联系邮箱：${email}</p>
            <p>通信地址：${address}</p>
            <p>© ${year} ${title}. 保留所有权利</p>
        </div>
    </div>
</body>
</html>`