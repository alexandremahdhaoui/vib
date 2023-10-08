/*
Copyright 2023 Alexandre Mahdhaoui

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), logLevel)
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
	if zap.S().Level() == zap.DebugLevel {
		SugaredLoggerWithCallerSkip(1).Fatalf(err.Error())
		return
	}
	SugaredLoggerWithCallerSkip(1).Errorf(err.Error())
}

func Error(err error) {
	SugaredLoggerWithCallerSkip(1).Errorf(err.Error())
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
