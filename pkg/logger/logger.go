package logger

import (
	"archiveFiles/internal/httpcontroller/middleware"
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log глобальный логгер
var Log *zap.Logger

func New(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	// определить log level
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	defer logger.Sync()
	Log = logger
	return logger, nil
}

// FromContext для логирования с XRequestID
func FromContext(ctx context.Context) *zap.Logger {

	if ctx == nil {
		return Log
	}

	if requestID, ok := ctx.Value(middleware.XRequestID).(string); ok {
		return Log.With(zap.String(
			middleware.XRequestID, requestID,
		))
	}

	return Log
}
