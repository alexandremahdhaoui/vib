package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New(debug bool) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.ConsoleSeparator = " "
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logLevel := zapcore.InfoLevel
	var options []zap.Option

	if debug {
		logLevel = zapcore.DebugLevel
		config.EncodeCaller = zapcore.ShortCallerEncoder
		config.CallerKey = "callerKey"

		options = append(
			options,
			zap.WithCaller(true),
			zap.AddStacktrace(zap.ErrorLevel),
		)
	}

	encoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel)
	logger := zap.New(core, options...)
	_ = logger.Sync()

	zap.ReplaceGlobals(logger)
}

// ---------------------------------------------------------------------------------------------------------------------
// Logging
// ---------------------------------------------------------------------------------------------------------------------

func SugaredLoggerWithCallerSkip(n int) *zap.SugaredLogger {
	return zap.S().WithOptions(zap.AddCallerSkip(n))
}

func Fatal(err error) {
	if err != nil {
		SugaredLoggerWithCallerSkip(1).Fatalf(err.Error())
	}
}

func Error(err error) {
	if err != nil {
		SugaredLoggerWithCallerSkip(1).Errorf(err.Error())
	}
}

func Warn(template string, a ...interface{}) {
	SugaredLoggerWithCallerSkip(1).Warnf(template, a...)
}

func Info(template string, a ...interface{}) {
	SugaredLoggerWithCallerSkip(1).Infof(template, a...)
}

func Debug(template string, a ...interface{}) {
	SugaredLoggerWithCallerSkip(1).Debugf(template, a...)
}