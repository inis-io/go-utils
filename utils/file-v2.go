package utils

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	
	"github.com/fsnotify/fsnotify"
	"github.com/mholt/archives"
	"github.com/spf13/afero"
)

type FileClass struct {
	Fs afero.Fs
}
var FileV2 *FileClass

func init() {
	// 初始化文件实例
	FileV2 = &FileClass{
		Fs: afero.NewOsFs(),
	}
}

// AutoExtract - 自动解压函数，支持本地和网络压缩包
func(this *FileClass) AutoExtract(sourceURL, dest string) error {

	var (
		err  error
		size int64
		name string
		reader io.ReadCloser
	)

	// 确保目标目录存在
	err = os.MkdirAll(dest, 0755)

	// 判断是本地文件还是网络URL（通过正则表达式）
	if regexp.MustCompile(`^https?://`).MatchString(sourceURL) {

		name, reader, size, err = this.Download(sourceURL)
		if err != nil { return fmt.Errorf("下载压缩包失败: %v", err) }

		defer func(reader io.ReadCloser) { err = reader.Close() }(reader)

	} else {

		// 处理本地文件
		name = filepath.Base(sourceURL)
		reader, err = os.Open(sourceURL)
		if err != nil { return fmt.Errorf("打开本地文件失败: %v", err) }

		defer func(reader io.ReadCloser) { err = reader.Close() }(reader)

		// 获取文件大小
		if file, ok := reader.(*os.File); ok {
			if info, err := file.Stat(); err == nil {
				size = info.Size()
			}
		}
	}

	// 创建临时文件用于存储网络下载的内容（如果需要）
	var temp *os.File
	if regexp.MustCompile(`^https?://`).MatchString(sourceURL) {

		// 创建临时文件
		temp, err = os.CreateTemp("", "archive-*" + filepath.Ext(name))
		if err != nil { return fmt.Errorf("创建临时文件失败: %v", err) }

		defer func(name string) { err = os.Remove(name) }(temp.Name())
		// 清理临时文件
		defer func(temp *os.File) { err = temp.Close() }(temp)

		// 显示下载进度
		progress := &ProgressReader{ Reader: reader, Total: size }

		// 将内容复制到临时文件
		if _, err = io.Copy(temp, progress); err != nil {
			return fmt.Errorf("下载压缩包内容失败: %v", err)
		}

		// 重置文件指针到开始位置
		if _, err = temp.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("重置文件指针失败: %v", err)
		}

		// 使用临时文件作为输入
		reader = temp
	}

	// 自动识别文件格式
	format, stream, err := archives.Identify(context.Background(), name, reader)
	if err != nil { return fmt.Errorf("识别文件格式失败: %v", err) }

	// 检查是否为可提取的格式
	extractor, ok := format.(archives.Extractor)
	if !ok { return fmt.Errorf("不支持的压缩格式: %T", format) }

	// 创建目标目录
	if err = os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	return extractor.Extract(context.Background(), stream, func(ctx context.Context, fileInfo archives.FileInfo) error {

		// 构建目标路径
		destPath := filepath.Join(dest, fileInfo.NameInArchive)

		// 处理目录
		if fileInfo.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// 处理文件
		return this.WriteFile(destPath, fileInfo)
	})
}

// Download - 下载网络压缩包
func(this *FileClass) Download(url string) (fileName string, body io.ReadCloser, length int64, err error) {

	// 创建带有超时的HTTP客户端
	client := &http.Client{
		// 设置较长的超时时间
		Timeout: 10 * time.Minute,
	}

	// 创建HTTP请求
	request, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil { return "", nil, 0, err }

	// 发送请求
	resp, err := client.Do(request)
	if err != nil { return "", nil, 0, err }

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		err = resp.Body.Close()
		if err != nil { return "", nil, 0, err }
		return "", nil, 0, fmt.Errorf("HTTP请求失败: %s", resp.Status)
	}

	// 尝试从URL或响应头中提取文件名
	fileName = filepath.Base(url)
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if idx := strings.Index(cd, "filename="); idx != -1 {
			fileName = strings.Trim(cd[idx+9:], "\"' ")
		}
	}

	return fileName, resp.Body, resp.ContentLength, nil
}

// WriteFile - 写入文件内容
func(this *FileClass) WriteFile(dest string, fileInfo archives.FileInfo) error {

	// 创建目标文件
	file, err := os.Create(dest)
	if err != nil { return fmt.Errorf("创建文件失败: %v", err) }

	defer func(file *os.File) { err = file.Close() }(file)

	// 获取文件内容读取器
	reader, err := fileInfo.Open()
	if err != nil { return fmt.Errorf("打开文件内容失败: %v", err) }
	defer func(reader fs.File) { err = reader.Close() }(reader)

	// 写入文件内容
	_, err = io.Copy(file, reader)
	if err != nil { return fmt.Errorf("写入文件内容失败: %v", err) }

	// 设置文件修改时间（可选）
	if !fileInfo.ModTime().IsZero() {
		err = os.Chtimes(dest, fileInfo.ModTime(), fileInfo.ModTime())
	}

	// 设置 755 权限
	err = os.Chmod(dest, 0755)

	return nil
}

// ProgressReader 进度读取器，用于显示下载进度
type ProgressReader struct {
	Reader     io.Reader
	Total      int64
	Reading    int64
	OnProgress func(read int64)
	LastUpdate time.Time
}

func (this *ProgressReader) Read(b []byte) (n int, err error) {

	n, err = this.Reader.Read(b)
	this.Reading += int64(n)

	// 每500毫秒更新一次进度，避免刷屏
	if time.Since(this.LastUpdate) > 500 * time.Millisecond {
		if this.OnProgress != nil {
			this.OnProgress(this.Reading)
		}
		this.LastUpdate = time.Now()
	}

	return
}

// FileInfo - 文件信息结构体
type FileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Mode    fs.FileMode `json:"mode"`
	ModTime time.Time   `json:"modTime"`
	IsDir   bool `json:"isDir"`
	Sys     any  `json:"sys"`
}

// List - 查询目录下的目录和文件
func (this *FileClass) List(dest string) (files []FileInfo, err error) {

	// 打开指定目录
	dir, err := this.Fs.Open(dest)
	if err != nil { return nil, err }
	defer func(dir afero.File) { err = dir.Close() }(dir)

	// 获取目录下的所有文件信息
	infos, err := dir.Readdir(-1)
	if err != nil { return nil, err }

	for _, info := range infos {
		files = append(files, FileInfo{
			Path:    filepath.ToSlash(filepath.Join(dest, info.Name())),
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
			Sys:     info.Sys(),
		})
	}

	return files, nil
}

// Exist 检查文件或目录是否存在
func (this *FileClass) Exist(dest string) (ok bool) {
	ok, _ = afero.Exists(this.Fs, dest)
	return ok
}

// Chmod - 修改文件或目录权限
func (this *FileClass) Chmod(dest string, mode int64) {
	_, _ = afero.Exists(this.Fs, dest)
	_ = this.Fs.Chmod(dest, os.FileMode(mode))
	return
}

// CreateFile - 创建文件并写入内容
func (this *FileClass) Write(dest string, content []byte) (err error) {

	// 以写入模式打开文件（若不存在则创建，存在则覆盖）
	file, err := this.Fs.Create(dest)

	if err != nil { return fmt.Errorf("创建文件失败: %w", err) }
	// 延迟关闭文件
	defer func(file afero.File) { err = file.Close() }(file)

	// 写入内容到文件
	_, err = file.Write(content)
	if err != nil { return fmt.Errorf("写入文件内容失败: %w", err) }

	// 设置文件权限为 755
	err = this.Fs.Chmod(dest, 0755)

	return nil
}

// Read - 读取文件内容
func(this *FileClass) Read(dest string) (content []byte, err error) {

	// 打开指定文件
	file, err := this.Fs.Open(dest)
	if err != nil { return nil, fmt.Errorf("打开文件失败: %w", err) }
	defer func(file afero.File) { err = file.Close() }(file)

	// 获取文件信息
	info, err := file.Stat()
	if err != nil { return nil, fmt.Errorf("获取文件信息失败: %w", err) }

	// 读取文件内容
	content = make([]byte, info.Size())
	_, err = file.Read(content)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	return content, nil
}

// Delete - 删除文件或目录
func(this *FileClass) Delete(dest string) error {
	return this.Fs.RemoveAll(dest)
}

// FileWatcher 文件监听器
type FileWatcher struct {
	watcher *fsnotify.Watcher
	filePath string
	callback func(event fsnotify.Event)
	stopChan chan bool
	// 防抖相关
	mu sync.Mutex
	timer *time.Timer
	lastEvent fsnotify.Event
	debounce time.Duration
}

// Watcher 创建文件监听器
func(this *FileClass) Watcher(filePath string, callback func(event fsnotify.Event)) (*FileWatcher, error) {
 
	// 创建 fsnotify 监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("创建文件监听器失败: %w", err)
	}
	
	// 获取文件的绝对路径和目录
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		if err := watcher.Close(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("获取文件绝对路径失败: %w", err)
	}
	
	// 获取文件所在目录（fsnotify 监听的是目录）
	dir := filepath.Dir(absPath)
	
	// 添加目录到监听列表
	err = watcher.Add(dir)
	if err != nil {
		if err := watcher.Close(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("添加监听目录失败: %w", err)
	}
	
	return &FileWatcher{
		watcher : watcher,
		filePath: absPath,
		callback: callback,
		stopChan: make(chan bool),
		// 默认防抖时间
		debounce: 200 * time.Millisecond,
	}, nil
}

// Start 开始监听文件变更
func (this *FileWatcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <- this.watcher.Events:
				if !ok { return }
				// 只处理目标文件的变更事件
				if filepath.Clean(event.Name) == this.filePath {
					// 过滤掉一些不需要的事件，比如访问事件
					if event.Op.Has(fsnotify.Write) || event.Op.Has(fsnotify.Remove) || event.Op.Has(fsnotify.Rename) || event.Op.Has(fsnotify.Create) {
						// 使用防抖：如果在短时间内收到多个事件，只执行一次回调
						this.mu.Lock()
						this.lastEvent = event
						if this.timer != nil {
							this.timer.Stop()
						}
						// 创建一个新的定时器，定时器触发时调用回调
						this.timer = time.AfterFunc(this.debounce, func() {
							this.mu.Lock()
							ev := this.lastEvent
							this.mu.Unlock()
							this.callback(ev)
						})
						this.mu.Unlock()
					}
				}
			case _, ok := <- this.watcher.Errors: if !ok { return }
			case <- this.stopChan: return
			}
		}
	}()
}

// Stop 停止监听
func (this *FileWatcher) Stop() {
	this.stopChan <- true
	if this.timer != nil {
		this.timer.Stop()
	}
	err := this.watcher.Close()
	if err != nil { return }
}