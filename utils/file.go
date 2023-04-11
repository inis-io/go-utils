package utils

import (
	"errors"
	"github.com/spf13/cast"
	"io"
	"os"
	"path/filepath"
)

type FileRequest struct {
	// 文件名
	Name string
	// 文件路径
	Path string
	// 目录路径
	Dir string
	// 文件后缀
	Ext string
}

// FileStruct - File 结构体
type FileStruct struct {
	request  *FileRequest
	response *FileResponse
}

type FileResponse struct {
	Error  error
	Result string
	Byte   []byte
	Slice  []any
}

func File(request ...FileRequest) *FileStruct {

	if len(request) == 0 {
		request = append(request, FileRequest{})
	}

	return &FileStruct{
		request : &request[0],
		response: &FileResponse{},
	}
}

// Path 设置文件路径(包含文件名，如：/tmp/test.txt)
func (this *FileStruct) Path(path any) *FileStruct {
	this.request.Path = cast.ToString(path)
	return this
}

// Dir 设置目录路径(不包含文件名，如：/tmp)
func (this *FileStruct) Dir(dir any) *FileStruct {
	this.request.Dir = cast.ToString(dir)
	return this
}

// Name 设置文件名(不包含路径，如：test.txt)
func (this *FileStruct) Name(name any) *FileStruct {
	this.request.Name = cast.ToString(name)
	return this
}

// Ext 设置文件后缀(如：.txt)
func (this *FileStruct) Ext(ext any) *FileStruct {
	this.request.Ext = cast.ToString(ext)
	return this
}

// Save 保存文件
func (this *FileStruct) Save(reader io.Reader, path ...string) (result *FileResponse) {

	if len(path) != 0 {
		this.request.Path = path[0]
	}

	if IsEmpty(this.request.Path) {
		this.response.Error = errors.New("文件路径不能为空")
		return this.response
	}

	filePath := cast.ToString(path)

	dir := filepath.Dir(filePath)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        // 目录不存在，需要创建
        if err := os.MkdirAll(dir, 0755); err != nil {
            this.response.Error = err
			return this.response
        }
    }

	// 创建文件
    file, err := os.Create(filePath)
    if err != nil {
		this.response.Error = err
		return this.response
    }

	// 关闭文件
    defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			this.response.Error = err
			return
		}
	}(file)

	// 将文件通过流的方式写入磁盘
	_, err = io.Copy(file, reader)
	if err != nil {
		this.response.Error = err
		return this.response
	}

    return nil
}

// Byte 获取文件字节
func (this *FileStruct) Byte(path ...any) (result *FileResponse) {

	if len(path) != 0 {
		this.request.Path = cast.ToString(path[0])
	}

	if IsEmpty(this.request.Path) {
		this.response.Error = errors.New("文件路径不能为空")
		return this.response
	}

	// 读取文件
	file, err := os.Open(this.request.Path)
	if err != nil {
		this.response.Error = err
		return this.response
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	// // 获取文件大小
	// info, _ := file.Stat()
	// size := info.Size()
	// // 读取文件
	// data := make([]byte, size)
	// file.Read(data)
	// return data

	var bytes []byte

	// 分批次读取
	buf := make([]byte, 1024)
	for {
		line, err := file.Read(buf)
		if err != nil {
			break
		}
		bytes = append(bytes, buf[:line]...)
	}

	this.response.Byte   = bytes
	this.response.Result = string(bytes)

	return this.response
}

// FileList 获取指定目录下的所有文件
func (this *FileStruct) List(opt ...map[string]any) (result *FileResponse) {

	// 默认参数
	defOpt := map[string]any{
		// 获取指定后缀的文件
		"ext": []string{"*"},
		// 包含子目录
		"sub": true,
		// 返回路径格式
		"format": "network",
		// 域名
		"domain": "",
		// 过滤前缀
		"prefix": "",
		// 目录路径
		"dir": "",
	}

	if len(opt) != 0 {
		// 合并参数
		for key, val := range defOpt {
			if opt[0][key] == nil {
				opt[0][key] = val
			}
		}
	} else {
		// 默认参数
		opt = append(opt, defOpt)
	}

	conf := opt[0]

	if !IsEmpty(conf["dir"]) {
		this.request.Dir = cast.ToString(conf["dir"])
	}

	if IsEmpty(this.request.Dir) {
		this.response.Error = errors.New("目录路径不能为空")
		return this.response
	}
	var slice []string
	this.response.Error = filepath.Walk(this.request.Dir, func(path string, info os.FileInfo, err error) error {
		// 忽略当前目录
		if info.IsDir() {
			return nil
		}
		// 忽略子目录
		if !conf["sub"].(bool) && filepath.Dir(path) != path {
			return nil
		}
		// []string 转 []any
		var exts []any
		for _, v := range conf["ext"].([]string) {
			exts = append(exts, v)
		}
		// 忽略指定后缀
		if !InArray("*", exts) && !InArray(filepath.Ext(path), exts) {
			return nil
		}
		slice = append(slice, path)
		return nil
	})

	// 转码为网络路径
	if conf["format"] == "network" {
		for key, val := range slice {
			slice[key] = filepath.ToSlash(val)
			if !IsEmpty(conf["domain"]) {
				slice[key] = cast.ToString(conf["domain"]) + slice[key][len(cast.ToString(conf["prefix"])):]
			}
		}
	}

	this.response.Slice = cast.ToSlice(slice)
	return this.response
}