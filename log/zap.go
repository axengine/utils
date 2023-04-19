package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	aboveDebug struct{}
	aboveWarn  struct{}
)

var Logger *zap.Logger

func (l aboveDebug) Enabled(lv zapcore.Level) bool {
	return lv >= zapcore.DebugLevel
}

func (l aboveWarn) Enabled(lv zapcore.Level) bool {
	return lv >= zapcore.WarnLevel
}

// init 初始化默认Logger
func init() {
	var encoder zapcore.Encoder
	encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	wDebug := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/debug.log",
		MaxSize:    100,
		MaxBackups: 50,
		MaxAge:     365,
		LocalTime:  true,
		Compress:   true,
	})

	wError := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     365,
		LocalTime:  true,
		Compress:   true,
	})

	coreDebug := zapcore.NewCore(
		encoder,
		wDebug,
		aboveDebug{},
	)

	coreError := zapcore.NewCore(
		encoder,
		wError,
		aboveWarn{},
	)

	Logger = zap.New(zapcore.NewTee(coreDebug, coreError), zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))
}
