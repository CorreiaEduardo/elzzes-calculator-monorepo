// Package http adapts net/http to the service: routing, middleware, JSON errors,
// and health. It has no calculator knowledge.
package http

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"strings"
)

// Router wraps http.ServeMux so unmatched paths and wrong methods produce the
// contract's JSON envelopes. All Handle calls must complete before serving.
type Router struct {
	mux     *nethttp.ServeMux
	allowed map[string][]string
}

func NewRouter() *Router {
	r := &Router{
		mux:     nethttp.NewServeMux(),
		allowed: make(map[string][]string),
	}
	r.mux.HandleFunc("/", func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		writeError(w, nethttp.StatusNotFound, CodeNotFound, "not found")
	})
	return r
}

// Handle registers handler for method and path, plus a path fallback that returns
// 405 with the accumulated allowed methods for any other method.
func (r *Router) Handle(method, path string, handler nethttp.Handler) {
	methods, registered := r.allowed[path]
	if !registered {
		r.mux.HandleFunc(path, func(w nethttp.ResponseWriter, req *nethttp.Request) {
			w.Header().Set("Allow", strings.Join(r.allowed[path], ", "))
			writeError(w, nethttp.StatusMethodNotAllowed, CodeMethodNotAllowed,
				fmt.Sprintf("method %s is not allowed", req.Method))
		})
	}

	if method == nethttp.MethodGet {
		methods = append(methods, nethttp.MethodGet, nethttp.MethodHead)
	} else {
		methods = append(methods, method)
	}
	r.allowed[path] = methods
	r.mux.Handle(method+" "+path, handler)
}

func (r *Router) ServeHTTP(w nethttp.ResponseWriter, req *nethttp.Request) {
	r.mux.ServeHTTP(w, req)
}

func writeJSON(w nethttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
