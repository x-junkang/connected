package clog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

const (
	LOG_MAX_SIZE = 10
	LOG_MAX_AGE  = 30
)

func InitLogger(file string, level string) {
	writeSyncer := getLogWriter(file)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, getLevel(level))
	logger = zap.New(core, zap.AddCaller())
}

func getLogWriter(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  file,
		MaxSize:   LOG_MAX_SIZE,
		MaxAge:    LOG_MAX_AGE,
		Compress:  false,
		LocalTime: true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.InfoLevel
}
