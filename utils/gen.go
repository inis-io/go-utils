package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

// Gen - 生成
var Gen *GenClass

type GenClass struct {}

// SerialNo 生成指定前缀和长度的序列号
/**
 * prefix: 前缀字符串
 * len: 序列号总长度
 * 格式: 前缀 + 日期(8位) + 时间(6位) + 随机数(动态长度)
 * 当指定长度小于前缀+日期+时间的长度时，会截断到指定长度
 */
func (this *GenClass) SerialNo(prefix any, length int) string {
	
	// 种子
	seed   := Hash.Sum32(fmt.Sprintf("%s-%s-%d-%d", cast.ToString(prefix), Get.Mac(), Get.Pid(), time.Now().UnixNano()))
	// 使用当前时间戳创建随机数生成器
	source := rand.New(rand.NewSource(cast.ToInt64(seed)))
	
	// 获取当前日期时间
	now := time.Now()
	datePart := now.Format("20060102") // 8位日期
	timePart := now.Format("150405")   // 6位时间
	
	// 计算固定部分的总长度
	fixedPart   := cast.ToString(prefix) + datePart + timePart
	fixedLength := len(fixedPart)
	
	var serialNo string
	
	// 如果指定长度小于等于固定部分长度，直接截断到指定长度
	if length <= fixedLength {
		
		serialNo = fixedPart[:length]
		
	} else {
		
		// 计算需要的随机数长度
		randomLength := length - fixedLength
		
		// 生成指定长度的随机数
		// 计算10的 randomLength 次方，作为随机数的上限
		maxLimit := 1
		for i := 0; i < randomLength; i++ { maxLimit *= 10 }
		
		// 生成随机数并格式化到指定长度
		randomPart := fmt.Sprintf("%0" + fmt.Sprintf("%dd", randomLength), source.Intn(maxLimit))
		
		// 组合所有部分
		serialNo = fixedPart + randomPart
	}
	
	return serialNo
}

// SerialDate 生成年月日 + (当天累计毫秒数)的字符串
/*
 * prefix: 前缀字符串
 * length: 结果总长度（包括前缀），如果结果长度超过指定长度，会被截断；如果结果长度不足指定长度，会在累计毫秒数前面补0
 * 格式: 前缀 + 年月日(8位) + 当天累计毫秒数(8位，固定长度，不足部分前面补0)
 * 建议：除去前缀部分，长度至少为16位（8位日期 + 8位累计毫秒数），否则可能会导致结果不唯一
 */
func (this *GenClass) SerialDate(prefix any, length int) string {
	
	// 1. 获取当前本地时间
	now := time.Now()
	
	// 2. 构造今天00:00:00的时间对象
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	// 3. 计算当前毫秒戳和今日零点毫秒戳
	currentMilli := now.UnixNano() / 1e6       // 当前毫秒戳
	todayStartMilli := todayStart.UnixNano() / 1e6 // 今日零点毫秒戳
	
	// 4. 计算当天累计的毫秒数（当前毫秒戳 - 今日零点毫秒戳）
	elapsedMilli    := currentMilli - todayStartMilli
	// elapsedMilli 固定长度为 8 位，超过部分会被截断，不足部分前面补0
	elapsedMilliStr := fmt.Sprintf("%08d", elapsedMilli)
	
	// 5. 格式化年月日（例如20260304）
	dateStr := now.Format("20060102")
	
	// 6. 将累计毫秒数转为字符串，拼接最终结果
	result := dateStr + elapsedMilliStr
	
	// 如果总长度超过指定长度，截断到指定长度
	if length > 0 && len(cast.ToString(prefix) + result) > length {
		result = result[:length - len(cast.ToString(prefix))]
	}
	// 如果总长度不足指定长度，在 elapsedMilliStr 前面补0，直到达到指定长度
	if length > 0 && len(cast.ToString(prefix) + result) < length {
		paddingLength := length - len(cast.ToString(prefix) + result)
		result = dateStr + fmt.Sprintf("%0*d", paddingLength + 8, elapsedMilli)
	}
	// 最终结果 = 前缀 + 年月日 + 累计毫秒数
	return cast.ToString(prefix) + result
}

// SnowFlakeID 生成雪花ID字符串
func (this *GenClass) SnowFlakeID(node int64) string {
	
	if node < 0 || node > 1023 {
		fmt.Printf("节点ID必须在0-1023范围内，当前值：%d\n", node)
		return ""
	}
	
	if node == 0 {
		// 生成随机节点ID
		source := rand.New(rand.NewSource(time.Now().UnixNano()))
		// 0-1023范围内的随机数
		node   = int64(source.Intn(1024))
	}
	
	// 初始化节点（分布式场景中，每个节点ID唯一，范围0-1023）
	item, err := snowflake.NewNode(node)
	if err != nil {
		// 如果雪花算法初始化失败，回退到普通序列号生成
		return this.SerialNo(nil, 24)
	}
	
	return item.Generate().String()
}

// BatchNo - 批次号
func (this *GenClass) BatchNo(prefix string, value any) string {
	date := time.Now().Format("20060102")
	if prefix != "" {
		return fmt.Sprintf("%s-%s-%04d", prefix, date, value)
	}
	return fmt.Sprintf("%s-%04d", date, value)
}
// IP 生成随机公网IP地址
// 排除内网、保留地址等非公网IP范围
func (this *GenClass) IP() string {
	
	// 初始化随机数生成器
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	var first, second, third, fourth int
	
	// 生成第一个字节，排除内网和特殊地址范围
	for {
		first = source.Intn(256)
		// 排除内网地址段: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
		// 排除特殊地址段: 0.0.0.0/8, 127.0.0.0/8, 169.254.0.0/16, 192.0.0.0/24等
		if !(first == 0 || first == 10 || first == 127 ||
			(first >= 169 && first <= 171) || first == 172 ||
			(first >= 192 && first <= 193)) {
			break
		}
	}
	
	// 处理特殊情况的第一个字节
	switch {
	case first == 172:
		// 172.16.0.0/12 是内网，所以第二个字节不能在16-31范围内
		for {
			second = source.Intn(256)
			if second < 16 || second > 31 {
				break
			}
		}
	case first == 192:
		// 192.168.0.0/16 是内网，所以第二个字节不能是168
		for {
			second = source.Intn(256)
			if second != 168 {
				break
			}
		}
	default:
		second = source.Intn(256)
	}
	
	third  = source.Intn(256)
	fourth = source.Intn(256)
	
	return fmt.Sprintf("%d.%d.%d.%d", first, second, third, fourth)
}

// UA 生成随机用户代理字符串
func (this *GenClass) UA() string {
	
	// 初始化随机数生成器
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// 常见浏览器类型
	browsers := []struct {
		name       string
		versions []string
	}{
		{"Chrome", []string{"91.0.4472.124", "92.0.4515.107", "93.0.4577.63", "94.0.4606.61", "95.0.4638.54", "96.0.4664.45", "97.0.4692.71"}},
		{"Firefox", []string{"89.0", "90.0", "91.0", "92.0", "93.0", "94.0", "95.0"}},
		{"Safari", []string{"14.1.2", "15.0", "15.1", "15.2", "15.3", "15.4"}},
		{"Edge", []string{"91.0.864.59", "92.0.902.73", "93.0.961.38", "94.0.992.50", "95.0.1020.44"}},
	}
	
	// 常见操作系统
	osList := []struct {
		name    string
		version string
	}{
		{"Windows NT 11.0", "Win64; x64"},
		{"Windows NT 10.0", "Win64; x64"},
		{"Windows NT 6.1", "WOW64"},
		{"Windows NT 6.3", "Win64; x64"},
		{"Macintosh", "Intel Mac OS X 10_15_7"},
		{"Macintosh", "Intel Mac OS X 12_0_1"},
		{"X11", "Linux x86_64"},
		{"iPhone", "CPU iPhone OS 15_0 like Mac OS X"},
		{"iPad", "CPU OS 15_0 like Mac OS X"},
	}
	
	// 随机选择浏览器
	browser := browsers[source.Intn(len(browsers))]
	browserVersion := browser.versions[source.Intn(len(browser.versions))]
	
	// 随机选择操作系统
	osInfo := osList[source.Intn(len(osList))]
	
	// 构建UA字符串
	var ua string
	switch browser.name {
	case "Chrome":
		webkitVersion := fmt.Sprintf("537.%d.%d", source.Intn(10), source.Intn(10))
		ua = fmt.Sprintf("Mozilla/5.0 (%s; %s) AppleWebKit/%s (KHTML, like Gecko) Chrome/%s Safari/%s", osInfo.name, osInfo.version, webkitVersion, browserVersion, webkitVersion)
	case "Firefox":
		geckoVersion  := fmt.Sprintf("20100101 Firefox/%s", browserVersion)
		ua = fmt.Sprintf("Mozilla/5.0 (%s; %s; rv:%s) Gecko/%s", osInfo.name, osInfo.version, browserVersion, geckoVersion)
	case "Safari":
		webkitVersion := fmt.Sprintf("605.%d.%d", source.Intn(10), source.Intn(10))
		safariVersion := browserVersion
		ua = fmt.Sprintf("Mozilla/5.0 (%s; %s) AppleWebKit/%s (KHTML, like Gecko) Version/%s Safari/%s", osInfo.name, osInfo.version, webkitVersion, safariVersion, webkitVersion)
	case "Edge":
		webkitVersion := fmt.Sprintf("537.%d.%d", source.Intn(10), source.Intn(10))
		edgeVersion   := browserVersion
		ua = fmt.Sprintf("Mozilla/5.0 (%s; %s) AppleWebKit/%s (KHTML, like Gecko) Chrome/%s Safari/%s Edg/%s", osInfo.name, osInfo.version, webkitVersion, browserVersion, webkitVersion, edgeVersion)
	}
	
	return ua
}

// Domain 生成随机域名
func (this *GenClass) Domain() string {
	
	// 初始化随机数生成器
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// 常见域名前缀
	prefixes := []string{
		"apple", "banana", "cherry", "date", "elderberry", "fig", "grape",
		"house", "car", "tree", "mountain", "river", "sun", "moon",
		"happy", "quick", "slow", "smart", "brave", "calm", "daring",
		"tech", "soft", "web", "net", "data", "info", "cloud",
	}
	
	// 常见顶级域名
	tlds := []string{
		"com", "org", "net", "info", "biz", "co.uk", "de", "fr", "jp",
		"us", "ca", "au", "io", "dev", "app", "shop", "site", "online",
	}
	
	// 随机选择前缀数量（1-3个）
	prefixCount := source.Intn(3) + 1
	var parts []string
	
	for i := 0; i < prefixCount; i++ {
		parts = append(parts, prefixes[source.Intn(len(prefixes))])
	}
	
	// 随机选择是否添加数字
	if source.Float64() < 0.3 { // 30%的概率添加数字
		parts[len(parts)-1] += fmt.Sprintf("%d", source.Intn(100))
	}
	
	// 组合前缀和顶级域名
	domain := strings.Join(parts, "") + "." + tlds[source.Intn(len(tlds))]
	
	return domain
}

// UUID 生成UUID字符串
func (this *GenClass) UUID() string {
	return uuid.New().String()
}