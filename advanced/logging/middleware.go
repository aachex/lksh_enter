package logging

import (
	"log/slog"
	"net/http"
)

type logReponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *logReponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func Middleware(next http.HandlerFunc, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(
			"New request",
			"endpoint", r.URL.String(),
			"method", r.Method,
			"user-agent", r.UserAgent(),
		)

		lrw := &logReponseWriter{w, http.StatusOK}
		next(lrw, r)

		logger.Debug(
			"Request processed",
			"statusCode", lrw.statusCode,
		)
	}
}
