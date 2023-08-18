package utils

import "net/url"

func init() {
	Is.Ip = IsIp
	Is.Url = IsUrl
	Is.Email = IsEmail
	Is.Phone = IsPhone
	Is.Mobile = IsMobile
	Is.Empty = IsEmpty
	Is.True = IsTrue
	Is.False = IsFalse
	Is.Number = IsNumber
	Is.Float = IsFloat
	Is.Bool = IsBool
	Is.Accepted = IsAccepted
	Is.Date = IsDate
	Is.Alpha = IsAlpha
	Is.AlphaNum = IsAlphaNum
	Is.AlphaDash = IsAlphaDash
	Is.Chs = IsChs
	Is.ChsAlpha = IsChsAlpha
	Is.ChsAlphaNum = IsChsAlphaNum
	Is.ChsDash = IsChsDash
	Is.Cntrl = IsCntrl
	Is.Graph = IsGraph
	Is.Lower = IsLower
	Is.Upper = IsUpper
	Is.Space = IsSpace
	Is.Xdigit = IsXdigit
	Is.ActiveUrl = IsActiveUrl
	Is.Domain = IsDomain
	Is.IdCard = IsIdCard
	Is.MacAddr = IsMacAddr
	Is.Zip = IsZip
	Is.String = IsString
	Is.Slice = IsSlice
	Is.Array = IsArray
	Is.JsonString = IsJsonString
	Is.Map = IsMap
	Is.SliceSlice = IsSliceSlice
	Is.MapAny = IsMapAny
	Get.Type = GetType
	Get.Ip = GetIp
	Get.Mac = GetMac
	Get.Pid = GetPid
	Get.Pwd = GetPwd
	In.Array = InArray[any]
	Array.Filter = ArrayFilter
	Array.Remove = ArrayRemove
	Array.Unique = ArrayUnique[any]
	Array.Empty = ArrayEmpty[any]
	Array.MapWithField = ArrayMapWithField
	Password.Create = PasswordCreate
	Password.Verify = PasswordVerify
	Rand.Int = RandInt
	Rand.String = RandString
	Rand.Slice = RandSlice
	Struct.Set = StructSet
	Struct.Get = StructGet
	Struct.Del = StructDel
	Struct.Has = StructHas
	Struct.Keys = StructKeys
	Struct.Values = StructValues
	Struct.Len = StructLen
	Struct.Map = StructMap
	Struct.Slice = StructSlice
	Json.Encode = JsonEncode
	Json.Decode = JsonDecode
	Json.Get = JsonGet
	Json.String = JsonString
	Format.Query = FormatQuery
	Parse.ParamsBefore = ParseParamsBefore
	Parse.Params = ParseParams
	Parse.Domain = ParseDomain
	Net.Tcping = NetTcping
	Mime.Type = MimeType
	Map.WithField = MapWithField[map[string]any]
	Map.WithoutField = MapWithoutField[map[string]any]
	Map.ToURL = MapToURL
	Map.Keys = MapKeys[map[string]any]
	Map.Values = MapValues[map[string]any]
	Version.Go = VersionGo
	Version.Compare = VersionCompare
	Unity.Ids = UnityIds
	Unity.Keys = UnityKeys
}

var Is struct {
	Ip          func(ip any) (ok bool)
	Url         func(url any) (ok bool)
	Email       func(email any) (ok bool)
	Phone       func(phone any) (ok bool)
	Mobile      func(mobile any) (ok bool)
	Empty       func(value any) (ok bool)
	True        func(value any) (ok bool)
	False       func(value any) (ok bool)
	Number      func(value any) (ok bool)
	Float       func(value any) (ok bool)
	Bool        func(value any) (ok bool)
	Accepted    func(value any) (ok bool)
	Date        func(date any) (ok bool)
	Alpha       func(value any) (ok bool)
	AlphaNum    func(value any) (ok bool)
	AlphaDash   func(value any) (ok bool)
	Chs         func(value any) (ok bool)
	ChsAlpha    func(value any) (ok bool)
	ChsAlphaNum func(value any) (ok bool)
	ChsDash     func(value any) (ok bool)
	Cntrl       func(value any) (ok bool)
	Graph       func(value any) (ok bool)
	Lower       func(value any) (ok bool)
	Upper       func(value any) (ok bool)
	Space       func(value any) (ok bool)
	Xdigit      func(value any) (ok bool)
	ActiveUrl   func(value any) (ok bool)
	Domain      func(domain any) (ok bool)
	IdCard      func(value any) (ok bool)
	MacAddr     func(value any) (ok bool)
	Zip         func(value any) (ok bool)
	String      func(value any) (ok bool)
	Slice       func(value any) (ok bool)
	Array       func(value any) (ok bool)
	JsonString  func(value any) (ok bool)
	Map         func(value any) (ok bool)
	SliceSlice  func(value any) (ok bool)
	MapAny      func(value any) (ok bool)
}

var Get struct {
	Type func(value any) (result string)
	Ip   func(key ...string) (result any)
	Mac  func() (result string)
	Pid  func() (result int)
	Pwd  func() (result string)
}

var In struct {
	Array func(value any, array []any) (ok bool)
}

var Array struct {
	Filter       func(array []string) (slice []string)
	Remove       func(array []string, args ...string) (slice []string)
	Unique       func(array []any) (slice []any)
	Empty        func(array []any) (slice []any)
	MapWithField func(array []map[string]any, field any) (slice []any)
}

var Password struct {
	Create func(password any) (result string)
	Verify func(encode any, password any) (ok bool)
}

var Rand struct {
	Int    func(max int, min ...int) (result int)
	String func(length int, chars ...string) (result string)
	Slice  func(slice []any, limit any) (result []any)
}

var Struct struct {
	Set    func(obj any, key string, val any)
	Get    func(obj any, key string) (result any)
	Del    func(obj any, key string)
	Has    func(obj any, key string) (ok bool)
	Keys   func(obj any) (slice []string)
	Values func(obj any) (slice []any)
	Len    func(obj any) (length int)
	Map    func(obj any) (result map[string]any)
	Slice  func(obj any) (slice []any)
}

var Json struct {
	Encode func(value any) (result string)
	Decode func(value any) (result any)
	Get    func(value any, key any) (result any, err error)
	String func(value any) (result string)
}

var Format struct {
	Query func(data any) (result string)
}

var Parse struct {
	ParamsBefore func(params url.Values) (result map[string]any)
	Params       func(params map[string]any) (result map[string]any)
	Domain       func(value any) (domain string)
}

var Net struct {
	Tcping func(host any, opts ...map[string]any) (ok bool, detail []map[string]any)
}

var Mime struct {
	Type func(suffix any) (mime string)
}

var Map struct {
	WithField    func(data map[string]any, field []string) (result map[string]any)
	WithoutField func(data map[string]any, field []string) (result map[string]any)
	ToURL        func(data map[string]any) (result string)
	Keys         func(data map[string]any) (result []string)
	Values       func(data map[string]any) (result []any)
}

var Version struct {
	Go      func() (version string)
	Compare func(v1 any, v2 any) (result int)
}

var Unity struct {
	Ids  func(param ...any) (result []any)
	Keys func(param any, reg ...any) (result []any)
}
