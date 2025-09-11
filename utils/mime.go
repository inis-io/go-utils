package utils

import (
	"strings"
	
	"github.com/spf13/cast"
)

var Mime *MimeClass

type MimeClass struct{}

var MimeMap = map[string]string{
	"js":   "application/javascript",
	"json": "application/json",
	"xml":  "application/xml",
	"css":  "text/css",
	"html": "text/html",
	"txt":  "text/plain",
	"gif":  "image/gif",
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"svg":  "image/svg+xml",
	"ico":  "image/x-icon",
	"pdf":  "application/pdf",
	"zip":  "application/zip",
	"rar":  "application/x-rar-compressed",
	"gz":   "application/x-gzip",
	"tar":  "application/x-tar",
	"7z":   "application/x-7z-compressed",
	"mp3":  "audio/mpeg",
	"mp4":  "video/mp4",
	"avi":  "video/x-msvideo",
	"doc":  "application/msword",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"xls":  "application/vnd.ms-excel",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"ppt":  "application/vnd.ms-powerpoint",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"csv":  "text/csv",
	"md":   "text/markdown",
}

// Type 获取后缀对应的 mime
func (this *MimeClass) Type(suffix any) (mime string) {
	// 获取后缀
	suffix = strings.ToLower(cast.ToString(suffix))
	// 取出 . 后面的内容
	suffix = strings.TrimPrefix(cast.ToString(suffix), ".")
	// 获取后缀对应的 mime
	mime = MimeMap[cast.ToString(suffix)]
	return
}
