package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// XRequestID номер запроса
const XRequestID = "X-Request-ID"

// XError статус архива
const XError = "X-Error"

// XErrorMessage ошибки
const XErrorMessage = "X-Error-Message"

// Logger - middleware для логирования запросов
func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapper := &responseWriter{w, http.StatusOK}

			requestIDVal := r.Header.Get(XRequestID)
			if requestIDVal == "" {
				requestIDVal = uuid.New().String()
			}

			r.Header.Set(XRequestID, requestIDVal)
			ctx := context.WithValue(r.Context(), XRequestID, requestIDVal)
			next.ServeHTTP(wrapper, r.WithContext(ctx))

			logger.Info("HTTP request",
				zap.String(XRequestID, requestIDVal),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.String("ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Int("status", wrapper.status),
				zap.String(XError, wrapper.Header().Get(XError)),
				zap.String(XErrorMessage, wrapper.Header().Get(XErrorMessage)),
				zap.Duration("duration", time.Since(start)))
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
