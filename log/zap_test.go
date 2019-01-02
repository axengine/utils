package log

import (
	"go.uber.org/zap"
	"testing"
)

// type

func Test_logLevel(t *testing.T) {
	Logger.Debug("debug", zap.String("k", "v"))
	Logger.Info("info", zap.String("k", "v"))
	Logger.Warn("warn", zap.String("k", "v"))
	Logger.Error("error", zap.String("k", "v"))
}

func Benchmark_logs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Logger.Debug("debug", zap.Int("i", i))
		if i%10 == 0 {
			Logger.Info("info", zap.Int("i", i))
		}
		if i%100 == 0 {
			Logger.Warn("warn", zap.Int("i", i))
		}
		if i%1000 == 0 {
			Logger.Error("error", zap.Int("i", i))
		}
	}
}
