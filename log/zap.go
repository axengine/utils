package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	infoOnly      struct{}
	infoWithDebug struct{}
	aboveWarn     struct{}
)

var Logger *zap.Logger

func (l infoOnly) Enabled(lv zapcore.Level) bool {
	return lv == zapcore.InfoLevel
}
func (l infoWithDebug) Enabled(lv zapcore.Level) bool {
	return lv == zapcore.InfoLevel || lv == zapcore.DebugLevel
}
func (l aboveWarn) Enabled(lv zapcore.Level) bool {
	return lv >= zapcore.WarnLevel
}

func makeInfoFilter(env string) zapcore.LevelEnabler {
	switch env {
	case "production":
		return infoOnly{}
	default:
		return infoWithDebug{}
	}
}

func makeErrorFilter() zapcore.LevelEnabler {
	return aboveWarn{}
}

// init
// 按照固定的模式初始化日志，调用方直接调用无需初始化
// 默认行为如下
// 存储路径:./log/
// 日志级别：DEBUG(含)以上，不同级别存储在不同文件中(DEBUG中包含DEBUG和INFO，INFO中只包含INFO，ERROR中包含WARN和ERROR)
// 滚动：500Mb
// 压缩：开启
func init() {
	var encoder zapcore.Encoder
	encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	wDebug := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/debug.log",
		MaxSize:    100, // megabytes
		MaxBackups: 50,
		MaxAge:     28, // days
		Compress:   true,
	})
	wInfo := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log",
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		Compress:   true,
	})
	wError := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log",
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		Compress:   true,
	})

	coreDebug := zapcore.NewCore(
		encoder,
		wDebug,
		makeInfoFilter("debug"),
	)

	coreInfo := zapcore.NewCore(
		encoder,
		wInfo,
		makeInfoFilter("production"),
	)
	coreError := zapcore.NewCore(
		encoder,
		wError,
		makeErrorFilter(),
	)

	Logger = zap.New(zapcore.NewTee(coreDebug, coreInfo, coreError), zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))
}
