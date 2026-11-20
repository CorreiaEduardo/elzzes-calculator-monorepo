// Package middleware holds transport middleware that depends only on the small
// capabilities it uses.
package middleware

import (
	"net/http"
	"time"
)

type recorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *recorder) WriteHeader(status int) {
	if r.wroteHeader {
		return
	}
	r.status = status
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *recorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

func (r *recorder) statusCode() int {
	if r.status == 0 {
		return http.StatusOK
	}
	return r.status
}

// RequestLogger writes one structured line per request. It must wrap Recover so a
// recovered panic is reported with status 500.
func RequestLogger(log func(msg string, args ...any)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			rec := &recorder{ResponseWriter: w}
			defer func() {
				log("request",
					"method", req.Method,
					"path", req.URL.Path,
					"status", rec.statusCode(),
					"duration", time.Since(start),
				)
			}()
			next.ServeHTTP(rec, req)
		})
	}
}
