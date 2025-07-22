// internal/logger/logger.go
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 日志配置
type Config struct {
	Level      string `yaml:"level" json:"level"`           // 日志级别: debug, info, warn, error
	Format     string `yaml:"format" json:"format"`         // 日志格式: json, console
	Output     string `yaml:"output" json:"output"`         // 输出方式: file, stdout, both
	LogDir     string `yaml:"logDir" json:"logDir"`         // 日志目录
	TimeFormat string `yaml:"timeFormat" json:"timeFormat"` // 时间格式
	MaxSize    int    `yaml:"maxSize" json:"maxSize"`       // 单个文件最大MB
	MaxBackups int    `yaml:"maxBackups" json:"maxBackups"` // 最多保留旧文件数量
	MaxAge     int    `yaml:"maxAge" json:"maxAge"`         // 保留天数
	Compress   bool   `yaml:"compress" json:"compress"`     // 是否压缩
	Caller     bool   `yaml:"caller" json:"caller"`         // 是否显示调用者信息
}

// DefaultConfig 默认配置 - 7天文件流转，按日期命名
func DefaultConfig() *Config {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "." // 如果获取失败，使用当前目录
	}
	
	// 如果当前在cmd目录，需要回到项目根目录
	if filepath.Base(workDir) == "cmd" {
		workDir = filepath.Dir(workDir)
	}
	
	return &Config{
		Level:      "info",
		Format:     "json",
		Output:     "both",
		LogDir:     filepath.Join(workDir, "logs"), // 使用绝对路径，确保在项目根目录
		TimeFormat: "20060102 15:04:05", // 时间格式
		MaxSize:    100,              // 100MB
		MaxBackups: 7,                // 保留7个备份文件
		MaxAge:     7,                // 保留7天
		Compress:   true,             // 压缩旧文件
		Caller:     false,            // 不显示调用者信息
	}
}

var (
	logInstance *zap.Logger
	config      *Config
	once        sync.Once
)

// getLogFilename 获取当天的日志文件名
func getLogFilename() string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("app-%s.log", today)
}

// Init 初始化日志器
func Init(cfg *Config) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	
	// 处理日志目录路径
	if !filepath.IsAbs(cfg.LogDir) {
		// 如果是相对路径，转换为绝对路径
		workDir, err := os.Getwd()
		if err != nil {
			workDir = "."
		}
		
		// 如果当前在cmd目录，需要回到项目根目录
		if filepath.Base(workDir) == "cmd" {
			workDir = filepath.Dir(workDir)
		}
		
		cfg.LogDir = filepath.Join(workDir, cfg.LogDir)
	}
	
	config = cfg
	
	// 确保日志目录存在
	if cfg.Output == "file" || cfg.Output == "both" {
		if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			panic(fmt.Sprintf("创建日志目录失败: %v", err))
		}
	}
	
	// 初始化日志器
	Log()
}

// Log 初始化或返回全局 logger
func Log() *zap.Logger {
	once.Do(func() {
		if config == nil {
			config = DefaultConfig()
		}
		
		// 解析日志级别
		level := parseLevel(config.Level)
		
		// 创建编码器配置
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			timeFormat := config.TimeFormat
			if timeFormat == "" {
				timeFormat = "20060102 15:04:05"
			}
			enc.AppendString(t.Format(timeFormat))
		}
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		
		// 选择编码器
		var encoder zapcore.Encoder
		if config.Format == "console" {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		}
		
		// 创建写入器
		var writeSyncer []zapcore.WriteSyncer
		
		if config.Output == "stdout" || config.Output == "both" {
			writeSyncer = append(writeSyncer, zapcore.AddSync(os.Stdout))
		}
		
		if config.Output == "file" || config.Output == "both" {
			// 使用当天日期作为文件名
			filename := filepath.Join(config.LogDir, getLogFilename())
			fileSyncer := zapcore.AddSync(&lumberjack.Logger{
				Filename:   filename,
				MaxSize:    config.MaxSize,
				MaxBackups: config.MaxBackups,
				MaxAge:     config.MaxAge,
				Compress:   config.Compress,
			})
			writeSyncer = append(writeSyncer, fileSyncer)
		}
		
		// 创建核心
		core := zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(writeSyncer...),
			level,
		)
		
		// 创建选项
		opts := []zap.Option{}
		if config.Caller {
			opts = append(opts, zap.AddCaller())
		}
		
		logInstance = zap.New(core, opts...)
	})
	return logInstance
}

// parseLevel 解析日志级别
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// ==================== 业务日志API ====================

// Info 输出 Info 级别日志
func Info(msg string, fields ...zap.Field) {
	Log().Info(msg, fields...)
}

// InfoObject 输出对象信息
func InfoObject(msg string, obj interface{}) {
	Log().Info(msg, zap.Any("data", obj))
}

// InfoMap 输出Map信息
func InfoMap(msg string, data map[string]interface{}) {
	fields := make([]zap.Field, 0, len(data))
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}
	Log().Info(msg, fields...)
}

// InfoStruct 输出结构体信息
func InfoStruct(msg string, obj interface{}) {
	Log().Info(msg, zap.Any("data", obj))
}

// Warn 输出 Warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	Log().Warn(msg, fields...)
}

// Debug 输出 Debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	Log().Debug(msg, fields...)
}

// DebugObject 输出对象信息
func DebugObject(msg string, obj interface{}) {
	Log().Debug(msg, zap.Any("data", obj))
}

// Error 输出 Error 级别日志
func Error(msg string, fields ...zap.Field) {
	Log().Error(msg, fields...)
}

// ErrorObject 输出对象信息
func ErrorWithObject(msg string, obj interface{},err error,fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Any("data", obj),zap.Error(err)}, fields...)
	Log().Error(msg, allFields...)
}

// ErrorWithErr 输出带错误信息的Error级别日志
func ErrorWithErr(msg string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Error(err)}, fields...)
	Log().Error(msg, allFields...)
}

// ErrorObject 输出带对象的Error级别日志
func ErrorObject(msg string, obj interface{}, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Any("object", obj)}, fields...)
	Log().Error(msg, allFields...)
}

// ErrorMap 输出带Map的Error级别日志
func ErrorMap(msg string, data map[string]interface{}, fields ...zap.Field) {
	allFields := make([]zap.Field, 0, len(data)+len(fields))
	for k, v := range data {
		allFields = append(allFields, zap.Any(k, v))
	}
	allFields = append(allFields, fields...)
	Log().Error(msg, allFields...)
}

// ==================== 字段创建函数 ====================

// Field 创建任意类型字段
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// String 创建字符串字段
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int 创建整数字段
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 创建64位整数字段
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Float64 创建浮点数字段
func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

// Bool 创建布尔字段
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Time 创建时间字段
func Time(key string, value time.Time) zap.Field {
	return zap.Time(key, value)
}

// Duration 创建持续时间字段
func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

// ErrorField 创建错误字段
func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

// ==================== 上下文日志器 ====================

// WithFields 创建带字段的日志器
func WithFields(fields ...zap.Field) *zap.Logger {
	return Log().With(fields...)
}

// WithContext 创建带上下文的日志器
func WithContext(ctx map[string]interface{}) *zap.Logger {
	fields := make([]zap.Field, 0, len(ctx))
	for k, v := range ctx {
		fields = append(fields, Field(k, v))
	}
	return Log().With(fields...)
}

// WithCaller 创建带调用者信息的日志器
func WithCaller() *zap.Logger {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return Log().With(String("caller", fmt.Sprintf("%s:%d", filepath.Base(file), line)))
	}
	return Log()
}

// WithRequestID 创建带请求ID的日志器
func WithRequestID(requestID string) *zap.Logger {
	return Log().With(String("request_id", requestID))
}

// WithUser 创建带用户信息的日志器
func WithUser(userID string, username string) *zap.Logger {
	return Log().With(
		String("user_id", userID),
		String("username", username),
	)
}

// WithOperation 创建带操作信息的日志器
func WithOperation(operation string) *zap.Logger {
	return Log().With(String("operation", operation))
}

// WithResource 创建带资源信息的日志器
func WithResource(resourceType, resourceID string) *zap.Logger {
	return Log().With(
		String("resource_type", resourceType),
		String("resource_id", resourceID),
	)
}

// WithPerformance 创建带性能信息的日志器
func WithPerformance(duration time.Duration) *zap.Logger {
	return Log().With(Duration("duration", duration))
}

// WithHTTPRequest 创建带HTTP请求信息的日志器
func WithHTTPRequest(method, path, remoteAddr string) *zap.Logger {
	return Log().With(
		String("http_method", method),
		String("http_path", path),
		String("remote_addr", remoteAddr),
	)
}

// WithHTTPResponse 创建带HTTP响应信息的日志器
func WithHTTPResponse(statusCode int, responseSize int64) *zap.Logger {
	return Log().With(
		Int("http_status", statusCode),
		Int64("response_size", responseSize),
	)
}

// ==================== 工具函数 ====================

// Sync 同步日志缓冲区
func Sync() error {
	return Log().Sync()
}

// GetLevel 获取当前日志级别
func GetLevel() string {
	if config == nil {
		return "info"
	}
	return config.Level
}

// IsDebug 检查是否为调试级别
func IsDebug() bool {
	return parseLevel(GetLevel()) <= zapcore.DebugLevel
}

// IsInfo 检查是否为信息级别
func IsInfo() bool {
	return parseLevel(GetLevel()) <= zapcore.InfoLevel
}

// IsWarn 检查是否为警告级别
func IsWarn() bool {
	return parseLevel(GetLevel()) <= zapcore.WarnLevel
}

// IsError 检查是否为错误级别
func IsError() bool {
	return parseLevel(GetLevel()) <= zapcore.ErrorLevel
}

// GetCurrentLogFile 获取当前日志文件路径
func GetCurrentLogFile() string {
	if config == nil {
		return ""
	}
	return filepath.Join(config.LogDir, getLogFilename())
}
