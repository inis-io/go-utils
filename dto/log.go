package dto

// LogConfig - 日志配置
type LogConfig struct {
	// Enable      - 是否启用日志
	Enable  bool   `json:"enable"   default:"true"`
	// Size    - 单个日志文件大小（MB）
	Size    int    `json:"size" default:"2"`
	// Age     - 日志文件保存天数
	Age     int    `json:"age"  default:"7"`
	// Backups - 日志文件最大保存数量
	Backups int    `json:"backups" default:"20"`
	// Hash - 计算配置是否发生变更
	Hash    string `json:"hash"`
}