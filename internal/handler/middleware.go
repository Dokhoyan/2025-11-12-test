package handler

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (h *Handler) LoggingMiddleware(logger *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(wrapped, r)

		duration := time.Since(start)
		logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", wrapped.statusCode),
			zap.Duration("duration", duration),
		)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (h *Handler) MethodMiddleware(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			newErrorResponse("Method not allowed", w, http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
