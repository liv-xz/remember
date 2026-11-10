package applog

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Config 日志配置结构体
// 用于配置日志系统的各项参数
type Config struct {
	Level      string // 日志输出级别: debug, info, warn, error, fatal
	Formatter  string // 日志格式: json 或 text
	OutputPath string // 日志文件输出路径，为空则输出到控制台
	MaxSize    int    // 单个日志文件最大大小（MB）
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 旧日志文件保留天数
	Compress   bool   // 是否压缩旧日志文件
}

// Logger 全局日志实例
var Logger *logrus.Logger

// LogConfig 日志默认配置
// 可在程序启动前修改此配置来定制日志行为
var LogConfig = &Config{
	Level:      "info",                      // 默认日志级别为 info
	Formatter:  "json",                      // 默认使用 json 格式
	OutputPath: "./initialize_logs/app.log", // 默认日志文件路径
	MaxSize:    100,                         // 单个文件最大 100MB
	MaxBackups: 30,                          // 保留 30 个旧文件
	MaxAge:     30,                          // 保留 30 天
	Compress:   true,                        // 压缩旧日志
}

// Init 初始化日志系统
// 在程序启动时调用一次，设置日志级别、格式和输出位置
func Init() {
	// 创建新的 logrus 实例
	Logger = logrus.New()

	// 解析日志级别配置，默认为 info
	level, err := logrus.ParseLevel(LogConfig.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 根据配置选择日志格式
	// json 格式适合写入文件便于分析，text 格式适合开发时查看
	switch LogConfig.Formatter {
	case "json":
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05", // 统一时间戳格式
		})
	default:
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,                  // 显示完整时间戳
			TimestampFormat:  "2006-01-02 15:04:05", // 时间戳格式
			DisableQuote:     true,                  // 不使用引号包裹
			DisableTimestamp: false,                 // 显示时间戳
		})
	}

	// 设置日志输出位置
	// 如果配置了 OutputPath 则写入文件，否则输出到控制台
	if LogConfig.OutputPath != "" {
		file, err := os.OpenFile(LogConfig.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err == nil {
			Logger.SetOutput(file)
		} else {
			// 文件创建失败时降级到控制台输出
			Logger.Warnf("日志文件创建失败: %v，使用标准输出", err)
			Logger.SetOutput(os.Stdout)
		}
	} else {
		Logger.SetOutput(os.Stdout)
	}

	Logger.Info("日志系统初始化成功")
}

// Debug 输出调试级别日志
// 仅在开发环境或需要详细排查问题时使用
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Info 输出信息级别日志
// 记录程序正常运行时的关键信息
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Warn 输出警告级别日志
// 记录可能存在问题但不影响运行的情况
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Error 输出错误级别日志
// 记录发生错误但程序仍可继续运行的情况
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Fatal 输出致命错误日志
// 记录严重错误并导致程序退出
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Debugf 格式化输出调试级别日志
// 支持占位符格式化，如: applog.Debugf("用户 %s 登录成功", username)
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Infof 格式化输出信息级别日志
// 支持占位符格式化，如: applog.Infof("服务器启动于 %s", addr)
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warnf 格式化输出警告级别日志
// 支持占位符格式化
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Errorf 格式化输出错误级别日志
// 支持占位符格式化，如: applog.Errorf("数据库连接失败: %v", err)
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatalf 格式化输出致命错误日志
// 支持占位符格式化，输出后程序将退出
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// WithFields 创建带有自定义字段的日志条目
// 用于结构化日志，方便日志分析和检索
// 用法: applog.WithFields(logrus.Fields{"user": "张三", "action": "登录"}).Info("用户操作")
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

// GetRequestLog 记录 HTTP 请求日志
// 自动根据响应状态码选择日志级别
// start: 请求开始时间
// method: HTTP 方法 (GET, POST 等)
// path: 请求路径
// status: 响应状态码
// duration: 请求处理耗时
func GetRequestLog(start time.Time, method, path string, status int, duration time.Duration) {
	fields := logrus.Fields{
		"method":   method,            // HTTP 请求方法
		"path":     path,              // 请求路径
		"status":   status,            // 响应状态码
		"duration": duration.String(), // 请求耗时
	}

	// 根据状态码判断日志级别
	// 5xx: 服务器内部错误
	// 4xx: 客户端请求错误
	// 其他: 正常请求
	switch {
	case status >= 500:
		Logger.WithFields(fields).Error("请求处理失败")
	case status >= 400:
		Logger.WithFields(fields).Warn("请求处理异常")
	default:
		Logger.WithFields(fields).Info("请求处理完成")
	}
}
