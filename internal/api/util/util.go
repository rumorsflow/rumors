package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func IsDebug(log *zap.Logger) bool {
	return zapcore.LevelOf(log.Core()) == zapcore.DebugLevel
}
