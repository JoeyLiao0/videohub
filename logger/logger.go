package logger

import (
	"fmt"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var (
	FileLogger *logrus.Logger
)

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

type fileHook struct {
	currentDate string
}

func (hook *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *fileHook) Fire(entry *logrus.Entry) error {
	if hook.currentDate != entry.Time.Format("2006-01-02") {
		hook.currentDate = entry.Time.Format("2006-01-02")
		// file, err := os.OpenFile("log/"+hook.currentDate+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		myLogger := &lumberjack.Logger{
			Filename:   "log/" + hook.currentDate + ".log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		}
		hook.currentDate = entry.Time.Format("2006-01-02")
		FileLogger.SetOutput(myLogger)
	}
	return nil
}

type myHook struct {
	currentDate string
}

func (hook *myHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *myHook) Fire(entry *logrus.Entry) error {
	msg := logMessage(entry, false)
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

type fileFormatter struct{}

func (f *fileFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

type myFormatter struct{}

func (g *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := logMessage(entry, true)
	return []byte(msg + "\n"), nil
}

func InitLogger(debug bool) {
	hook := &fileHook{}
	FileLogger = logrus.New()
	FileLogger.AddHook(hook)
	FileLogger.SetFormatter(&fileFormatter{})

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
