package logger

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInfo(t *testing.T) {
	infoMessage := "logging info message"
	tags := zap.Field{
		Key:    "anotherMessage-Key",
		String: "anotherMessage-Value",
		Type:   zapcore.StringType,
	}

	Info(infoMessage, tags)
}

func TestError(t *testing.T) {
	infoMessage := "logging error message"
	err := errors.New("logging an error message")
	tags := zap.Field{
		Key:    "anotherMessage-Key",
		String: "anotherMessage-Value",
		Type:   zapcore.StringType,
	}

	Error(infoMessage, err, tags)
}

func TestGetLevelLogs(t *testing.T) {
	os.Setenv(LOG_LEVEL, "info")
	logLevel := getLevelLogs()
	assert.Equal(t, logLevel, zap.InfoLevel)

	os.Setenv(LOG_LEVEL, "error")
	logLevel = getLevelLogs()
	assert.Equal(t, logLevel, zap.ErrorLevel)

	os.Setenv(LOG_LEVEL, "debug")
	logLevel = getLevelLogs()
	assert.Equal(t, logLevel, zap.DebugLevel)

	os.Setenv(LOG_LEVEL, "")
	logLevel = getLevelLogs()
	assert.Equal(t, logLevel, zap.InfoLevel)
}
