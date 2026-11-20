package http

import (
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeLogger struct{}

func (fakeLogger) Info(string, ...any)  {}
func (fakeLogger) Error(string, ...any) {}

func TestHandlerHealth(t *testing.T) {
	handler := NewHandler(fakeLogger{})

	rec := serve(handler, nethttp.MethodGet, "/healthz")
	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"status":"ok"}` {
		t.Errorf("body = %q, want health envelope", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("CORS header = %q, want none", got)
	}
}

func TestHandlerHealthHead(t *testing.T) {
	server := httptest.NewServer(NewHandler(fakeLogger{}))
	defer server.Close()

	req, err := nethttp.NewRequest(nethttp.MethodHead, server.URL+"/healthz", nil)
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

func TestHandlerHealthMethodNotAllowed(t *testing.T) {
	handler := NewHandler(fakeLogger{})

	rec := serve(handler, nethttp.MethodPost, "/healthz")
	if rec.Code != nethttp.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	if got := rec.Header().Get("Allow"); got != "GET, HEAD" {
		t.Errorf("Allow = %q, want %q", got, "GET, HEAD")
	}
	envelope := decodeEnvelope(t, rec)
	if envelope.Error.Code != CodeMethodNotAllowed {
		t.Errorf("code = %q, want METHOD_NOT_ALLOWED", envelope.Error.Code)
	}
}

func TestHandlerUnknownAPIRoute(t *testing.T) {
	handler := NewHandler(fakeLogger{})

	rec := serve(handler, nethttp.MethodPost, "/api/v1/calculate")
	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	envelope := decodeEnvelope(t, rec)
	if envelope.Error.Code != CodeNotFound {
		t.Errorf("code = %q, want NOT_FOUND", envelope.Error.Code)
	}
}
