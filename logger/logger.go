package logger

import (
	"fmt"
	"videohub/config"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// FileLogger 文件日志
var (
	FileLogger *logrus.Logger
)

// logMessage 格式化日志信息
func logMessage(entry *logrus.Entry, isColor bool) string {
	caller := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	color := ""
	reset := ""
	if isColor {
		color = logColor(entry)
		reset = "\033[0m"
	}
	switch entry.Level {
	case logrus.InfoLevel, logrus.WarnLevel:
		return fmt.Sprintf("%s[%s] %s \"%s\"%s",
			color,
			entry.Level.String(),
			entry.Time.Format("2006/01/02 - 15:04:05"),
			entry.Message,
			reset)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return fmt.Sprintf("%s[%s] %s %s \"%s\"%s",
			color,
			entry.Level.String(),
			entry.Time.Format("2006/01/02 - 15:04:05"),
			caller,
			entry.Message,
			reset)
	default:
		return fmt.Sprintf("%s[%s] %v %s", color, entry.Level.String(), entry.Message, reset)
	}
}

// logColor 根据日志级别返回颜色
func logColor(entry *logrus.Entry) string {
	switch entry.Level {
	case logrus.InfoLevel:
		return "\033[32m" // 绿色
	case logrus.WarnLevel:
		return "\033[33m" // 黄色
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return "\033[31m" // 红色
	default:
		return "\033[37m" // 白色
	}
}

// fileHook 文件钩子
type fileHook struct {
	currentDate string
}

// Levels 返回支持的日志级别
func (hook *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 输出日志
func (hook *fileHook) Fire(entry *logrus.Entry) error {
	if hook.currentDate != entry.Time.Format("2006-01-02") {
		hook.currentDate = entry.Time.Format("2006-01-02")
		myLogger := &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s.log", config.AppConfig.Log.Path, hook.currentDate),
			MaxSize:    config.AppConfig.Log.MaxSize,
			MaxBackups: config.AppConfig.Log.MaxBackups,
			MaxAge:     config.AppConfig.Log.MaxAge,
			Compress:   config.AppConfig.Log.Compress,
		}
		hook.currentDate = entry.Time.Format("2006-01-02")
		FileLogger.SetOutput(myLogger)
	}
	return nil
}

// myHook 自定义钩子
type myHook struct {
}

// Levels 返回支持的日志级别
func (hook *myHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 输出日志
func (hook *myHook) Fire(entry *logrus.Entry) error {
	msg := ""
	if entry.Message != "" {
		msg = logMessage(entry, false)
	}
	msg += "\n"
	switch entry.Level {
	case logrus.InfoLevel:
		FileLogger.Info(msg)
	case logrus.WarnLevel:
		FileLogger.Warn(msg)
	case logrus.ErrorLevel:
		FileLogger.Error(msg)
	case logrus.FatalLevel:
		FileLogger.Fatal(msg)
	case logrus.PanicLevel:
		FileLogger.Panic(msg)
	}
	return nil
}

// fileFormatter 文件格式化
type fileFormatter struct{}

// Format 格式化日志信息
func (f *fileFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

// myFormatter 自定义格式化
type myFormatter struct{}

// Format 格式化日志信息
func (g *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := ""
	if entry.Message != "" {
		msg = logMessage(entry, true)
	}
	return []byte(msg + "\n"), nil
}

// InitLogger 初始化日志
func InitLogger(debug bool) {
	hook := &fileHook{}
	FileLogger = logrus.New()
	FileLogger.AddHook(hook)
	FileLogger.SetFormatter(&fileFormatter{})
	FileLogger.SetLevel(logrus.InfoLevel)

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&myFormatter{})
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	hook2 := &myHook{}
	logrus.AddHook(hook2)
	logrus.Info("Logger initialized successfully")
}
