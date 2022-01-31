package log

import (
	"isbnbook/app/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger interface {
	Info(message string)
	Warn(message string)
	Error(message string, err error)
}

type zapLogger struct {
}

func initLogger() *zap.SugaredLogger {

	// change work dir as app dir for logging
	utils.ChangeWorkDir()

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)
	
	myConfig := zap.Config{
		Level: level,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "./app.log"},
		ErrorOutputPaths: []string{"stderr", "./error.log"},
	}

	logger, _ := myConfig.Build()
	defer logger.Sync()
	sugar := logger.Sugar()
	return sugar
}

var sugar = initLogger()

func (z *zapLogger) Info(message string) {
	sugar.Info(message)
}

func (z *zapLogger) Warn(message string) {
	sugar.Warn(message)
}

func (z *zapLogger) Error(message string, err error) {
	if err != nil {
		sugar.Error(message, err)
	} else {
		sugar.Error(message)
	}
}

var myLogger = zapLogger{}

func GetLogger() AppLogger {
	return &myLogger
}
