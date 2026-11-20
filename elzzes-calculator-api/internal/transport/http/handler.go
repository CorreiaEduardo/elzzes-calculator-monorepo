package http

import (
	nethttp "net/http"

	"elzzes-calculator-api/internal/transport/http/middleware"
)

// NewHandler builds the HTTP adapter: calculator routes arrive with specs 001-002.
func NewHandler(logger Logger) nethttp.Handler {
	router := NewRouter()
	router.Handle(nethttp.MethodGet, "/healthz", healthHandler())

	var handler nethttp.Handler = router
	handler = middleware.Recover(logger.Error, func(w nethttp.ResponseWriter) {
		writeError(w, nethttp.StatusInternalServerError, CodeInternalError, "internal error")
	})(handler)
	handler = middleware.RequestLogger(logger.Info)(handler)
	return handler
}

type healthResponse struct {
	Status string `json:"status"`
}

func healthHandler() nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		writeJSON(w, nethttp.StatusOK, healthResponse{Status: "ok"})
	})
}
