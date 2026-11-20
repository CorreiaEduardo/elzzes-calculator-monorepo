package http

import (
	"encoding/json"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

func serve(h nethttp.Handler, method, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func decodeEnvelope(t *testing.T, rec *httptest.ResponseRecorder) errorEnvelope {
	t.Helper()
	var envelope errorEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode error envelope %q: %v", rec.Body.String(), err)
	}
	return envelope
}

func TestRouterNotFound(t *testing.T) {
	router := NewRouter()

	tests := []struct {
		method string
		path   string
	}{
		{method: nethttp.MethodGet, path: "/api/v1/unknown"},
		{method: nethttp.MethodPost, path: "/api/v1/calculate"},
		{method: nethttp.MethodGet, path: "/api/"},
		{method: nethttp.MethodGet, path: "/api"},
		{method: nethttp.MethodGet, path: "/"},
		{method: nethttp.MethodGet, path: "/nope"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			rec := serve(router, tt.method, tt.path)
			if rec.Code != nethttp.StatusNotFound {
				t.Fatalf("status = %d, want 404", rec.Code)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			if got := rec.Header().Get("Allow"); got != "" {
				t.Errorf("Allow = %q, want empty", got)
			}
			envelope := decodeEnvelope(t, rec)
			if envelope.Error.Code != CodeNotFound || envelope.Error.Message != "not found" {
				t.Errorf("envelope = %+v, want NOT_FOUND/not found", envelope.Error)
			}
		})
	}
}

func TestRouterMethodNotAllowed(t *testing.T) {
	router := NewRouter()
	router.Handle(nethttp.MethodPost, "/api/v1/calculate", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		w.WriteHeader(nethttp.StatusNoContent)
	}))
	router.Handle(nethttp.MethodGet, "/healthz", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		w.WriteHeader(nethttp.StatusNoContent)
	}))

	tests := []struct {
		method    string
		path      string
		wantAllow string
	}{
		{method: nethttp.MethodGet, path: "/api/v1/calculate", wantAllow: "POST"},
		{method: nethttp.MethodPut, path: "/api/v1/calculate", wantAllow: "POST"},
		{method: nethttp.MethodOptions, path: "/api/v1/calculate", wantAllow: "POST"},
		{method: nethttp.MethodPost, path: "/healthz", wantAllow: "GET, HEAD"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			rec := serve(router, tt.method, tt.path)
			if rec.Code != nethttp.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want 405", rec.Code)
			}
			if got := rec.Header().Get("Allow"); got != tt.wantAllow {
				t.Errorf("Allow = %q, want %q", got, tt.wantAllow)
			}
			envelope := decodeEnvelope(t, rec)
			wantMessage := "method " + tt.method + " is not allowed"
			if envelope.Error.Code != CodeMethodNotAllowed || envelope.Error.Message != wantMessage {
				t.Errorf("envelope = %+v, want METHOD_NOT_ALLOWED/%q", envelope.Error, wantMessage)
			}
		})
	}
}

func TestRouterMultipleMethods(t *testing.T) {
	router := NewRouter()
	noContent := nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		w.WriteHeader(nethttp.StatusNoContent)
	})
	router.Handle(nethttp.MethodGet, "/x", noContent)
	router.Handle(nethttp.MethodPost, "/x", noContent)

	rec := serve(router, nethttp.MethodDelete, "/x")
	if rec.Code != nethttp.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	if got := rec.Header().Get("Allow"); got != "GET, HEAD, POST" {
		t.Errorf("Allow = %q, want %q", got, "GET, HEAD, POST")
	}
}

func TestRouterHeadMatchesGet(t *testing.T) {
	router := NewRouter()
	router.Handle(nethttp.MethodGet, "/greet", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		_, _ = w.Write([]byte("hello"))
	}))

	server := httptest.NewServer(router)
	defer server.Close()

	req, err := nethttp.NewRequest(nethttp.MethodHead, server.URL+"/greet", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HEAD request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if resp.StatusCode != nethttp.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if len(body) != 0 {
		t.Errorf("HEAD body = %q, want empty", body)
	}
}
