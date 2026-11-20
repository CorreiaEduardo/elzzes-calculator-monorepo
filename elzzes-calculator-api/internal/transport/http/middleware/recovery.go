package middleware

import (
	"net/http"
	"runtime/debug"
)

// Recover logs the panic and stack, then calls respond when the handler has not
// written a response yet.
func Recover(log func(msg string, args ...any), respond func(http.ResponseWriter)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			rec := &recorder{ResponseWriter: w}
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}
				log("panic recovered",
					"method", req.Method,
					"path", req.URL.Path,
					"panic", recovered,
					"stack", string(debug.Stack()),
				)
				if !rec.wroteHeader {
					respond(rec)
				}
			}()
			next.ServeHTTP(rec, req)
		})
	}
}
