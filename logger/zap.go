package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

// var logger *zap.Logger

// InitLogger 初始化全局 logger 并将其设置为 zap 的默认 logger
// 使用 logger 的地方只需要 zap.L().XXX() 就可以了
func InitZapLogger(env, logLevel, format string, opts ...zap.Option) error {
	var conf zap.Config
	if env == "production" {
		conf = zap.NewProductionConfig()
	} else {
		conf = zap.NewDevelopmentConfig()
	}
	conf.Encoding = format
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	// change stderr all to stdout
	conf.OutputPaths = []string{"stdout"}
	conf.ErrorOutputPaths = []string{"stdout"}
	if logLevel == "debug" {
		level.SetLevel(zapcore.DebugLevel)
	} else if logLevel == "warn" {
		level.SetLevel(zapcore.WarnLevel)
	} else if logLevel == "fatal" {
		level.SetLevel(zapcore.FatalLevel)
	} else {
		level.SetLevel(zapcore.InfoLevel)
	}

	conf.Level = level
	var l, err = conf.Build(opts...)
	if err != nil {
		return err
	}
	_, err = zap.RedirectStdLogAt(l, zap.DebugLevel)
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(l)
	//logger = l
	return nil
}

// 设置全局 Logger 对应的级别
// 这个函数是并发安全的
func SetLevel(lv zapcore.Level) {
	level.SetLevel(lv)
}

// 获取全局 Logger 的级别
func GetLevel() zapcore.Level {
	return level.Level()
}

func NoStacktrace(l *zap.Logger) *zap.Logger {
	return l.WithOptions(zap.AddStacktrace(zap.FatalLevel))
}
