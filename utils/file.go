package utils

import (
	"bufio"
	"errors"
	"github.com/spf13/cast"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileStruct - File 结构体
type FileStruct struct {
	request  *FileRequest
	response *FileResponse
}

// FileRequest - File 请求
type FileRequest struct {
	// 文件名
	Name string
	// 文件路径（包含文件名）
	Path string
	// 目录路径（不包含文件名）
	Dir string
	// 文件后缀
	Ext string
	// 限制行数
	Limit int
	// 读取偏移量
	Page int
	// 返回结果格式
	Format string
	// 是否包含子目录
	Sub bool
	// 域名 - 用于拼接文件路径
	Domain string
	// 前缀 - 用于过滤前缀
	Prefix string
}

// FileResponse - File 响应
type FileResponse struct {
	Error  error
	Result any
	Text   string
	Byte   []byte
	Slice  []any
}

// File - 文件系统
func File(request ...FileRequest) *FileStruct {

	if len(request) == 0 {
		request = append(request, FileRequest{})
	}

	if IsEmpty(request[0].Limit) {
		request[0].Limit = 10
	}

	if IsEmpty(request[0].Page) {
		request[0].Page = 1
	}

	if IsEmpty(request[0].Format) {
		request[0].Format = "network"
	}

	if IsEmpty(request[0].Sub) {
		request[0].Sub = true
	}

	if IsEmpty(request[0].Ext) {
		request[0].Ext = "*"
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

// Domain 设置域名(用于拼接文件路径)
func (this *FileStruct) Domain(domain any) *FileStruct {
	this.request.Domain = cast.ToString(domain)
	return this
}

// Prefix 设置前缀(用于过滤前缀)
func (this *FileStruct) Prefix(prefix any) *FileStruct {
	this.request.Prefix = cast.ToString(prefix)
	return this
}

// Limit 设置限制行数
func (this *FileStruct) Limit(limit any) *FileStruct {
	this.request.Limit = cast.ToInt(limit)
	return this
}

// Page 设置读取偏移量
func (this *FileStruct) Page(page any) *FileStruct {
	this.request.Page = cast.ToInt(page)
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

	dir := filepath.Dir(this.request.Path)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        // 目录不存在，需要创建
        if err := os.MkdirAll(dir, 0755); err != nil {
            this.response.Error = err
			return this.response
        }
    }

	// 创建文件
    file, err := os.Create(this.request.Path)
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

    return this.response
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

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		this.response.Error = err
		return this.response
	}
	size := info.Size()

	// 小于50MB，整体读取
	if size < 50 * 1024 * 1024 {
		bytes := make([]byte, size)
		_, err = file.Read(bytes)
		if err != nil {
			this.response.Error = err
			return this.response
		}
		this.response.Byte   = bytes
		this.response.Text   = string(bytes)
		this.response.Result = bytes
		return this.response
	}

	// 大于等于50MB，分块读取
	var bytes []byte
	buffer := make([]byte, 1024 * 1024)
	for {
		index, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			this.response.Error = err
			break
		}
		bytes = append(bytes, buffer[:index]...)
		if err == io.EOF {
			break
		}
	}

	this.response.Byte   = bytes
	this.response.Text   = string(bytes)
	this.response.Result = bytes

	return this.response
}

// List 获取指定目录下的所有文件
func (this *FileStruct) List(path ...any) (result *FileResponse) {

	if len(path) != 0 {
		this.request.Path = cast.ToString(path[0])
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
		if !this.request.Sub && filepath.Dir(path) != path {
			return nil
		}
		// []string 转 []any
		var exts []any
		// this.request.Ext 逗号分隔的字符串 转 []string
		for _, val := range strings.Split(this.request.Ext, ",") {
			// 忽略空字符串
			if IsEmpty(val) {
				continue
			}
			// 去除空格
			exts = append(exts, strings.TrimSpace(val))
		}
		// 忽略指定后缀
		if !InArray("*", exts) && !InArray[any](filepath.Ext(path), exts) {
			return nil
		}
		slice = append(slice, path)
		return nil
	})

	// 转码为网络路径
	if this.request.Format == "network" {
		for key, val := range slice {
			slice[key] = filepath.ToSlash(val)
			if !IsEmpty(this.request.Domain) {
				slice[key] = this.request.Domain + slice[key][len(this.request.Prefix):]
			}
		}
	}

	for _, val := range slice {
		this.response.Slice = append(this.response.Slice, val)
	}
	this.response.Result = slice
	this.response.Text   = strings.Join(slice, ",")
	this.response.Byte   = []byte(this.response.Text)

	return this.response
}

// IsExist 判断文件是否存在
func (this *FileStruct) IsExist(path ...any) (ok bool) {

	if len(path) != 0 {
		this.request.Path = cast.ToString(path[0])
	}

	if IsEmpty(this.request.Path) {
		return false
	}

	// 判断文件是否存在
	if _, err := os.Stat(this.request.Path); os.IsNotExist(err) {
		return false
	}

	return true
}

// Line 按行读取文件
func (this *FileStruct) Line(path ...any) (result *FileResponse) {

	if len(path) != 0 {
		this.request.Path = cast.ToString(path[0])
	}

	if IsEmpty(this.request.Path) {
		this.response.Error = errors.New("文件路径不能为空")
		return this.response
	}

	// 读取块
	readBlock := func(file *os.File, start int, end int) ([]string, error) {

		lines   := make([]string, 0)
		scanner := bufio.NewScanner(file)

		// 移动扫描器到指定的起始行
		for i := 1; i < start && scanner.Scan(); i++ {}

		// 开始读取需要的行
		for i := start; i <= end && scanner.Scan(); i++ {
			// 只把需要的行保存到切片中
			if i >= start && i <= end {
				lines = append(lines, scanner.Text())
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		return lines, nil
	}

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

	end   := this.request.Page * this.request.Limit
	start := end - this.request.Limit + 1

	lines     := make([]string, 0)
	// 计算每个块的大小（以10MB为一个块）
	blockSize := 10 * 1024 * 1024
	numBlocks := (end - start + 1) / blockSize

	if (end - start + 1) % blockSize != 0 {
		numBlocks++
	}

	// 并发读取每个块
	var wg sync.WaitGroup
	wg.Add(numBlocks)
	for i := 0; i < numBlocks; i++ {
		go func(i int) {

			startLine := start + i * blockSize
			endLine   := startLine + blockSize - 1
			if endLine > end {
				endLine = end
			}

			blockLines, err := readBlock(file, startLine, endLine)
			if err != nil {
				this.response.Error = err
				return
			} else {
				lines = append(lines, blockLines...)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	this.response.Result = lines
	this.response.Text   = JsonEncode(this.response.Result)
	this.response.Byte   = []byte(this.response.Text)

	for _, v := range lines {
		this.response.Slice = append(this.response.Slice, v)
	}

	return this.response
}