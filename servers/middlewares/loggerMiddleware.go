package middlewares

import (
	"bytes"
	"errors"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"net/http"
	"time"
)

// responseRecorder is a wrapper for http.ResponseWriter that records the status code written to the response.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// Write captures the data written to the buffer.
func (rr *responseRecorder) Write(data []byte) (int, error) {
	rr.body.Write(data)
	return rr.ResponseWriter.Write(data)
}

// WriteHeader captures the status code for later inspection.
func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// LoggerMiddleware logs incoming HTTP requests, their duration, and the response body in the event of errors.
func LoggerMiddleware(logger common.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			ip := r.RemoteAddr // This might be a proxy IP address

			// Attempt to get the real IP if the request is forwarded from a proxy.
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			next.ServeHTTP(rr, r) // Call the handler

			if rr.statusCode >= http.StatusBadRequest {
				logger.Error("request failed", errors.New(rr.body.String()), "ip", ip, "method", r.Method, "uri", r.RequestURI, "duration", time.Since(start))
			} else {
				logger.Debug("request handled", "ip", ip, "method", r.Method, "uri", r.RequestURI, "duration", time.Since(start))
			}
		})
	}
}
