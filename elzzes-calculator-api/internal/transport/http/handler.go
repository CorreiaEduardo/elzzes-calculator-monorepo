package http

import (
	nethttp "net/http"

	"elzzes-calculator-api/internal/transport/http/middleware"
)

// NewHandler builds the HTTP adapter: health and the calculator endpoint.
func NewHandler(logger Logger, calc Calculator) nethttp.Handler {
	router := NewRouter()
	router.Handle(nethttp.MethodGet, "/healthz", healthHandler())
	router.Handle(nethttp.MethodPost, "/api/v1/calculate", calculateHandler(logger, calc))

	var handler nethttp.Handler = router
	handler = middleware.Recover(logger.Error, func(w nethttp.ResponseWriter) {
		writeError(w, nethttp.StatusInternalServerError, CodeInternalError, msgInternalError)
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
