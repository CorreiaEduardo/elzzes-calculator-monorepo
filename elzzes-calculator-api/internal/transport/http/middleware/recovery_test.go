package middleware

import (
	"net/http"
	"testing"
)

func TestRecover(t *testing.T) {
	log, calls := newLog()
	responded := false
	handler := Recover(log, func(w http.ResponseWriter) {
		responded = true
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))

	rec := serve(handler, http.MethodGet, "/panic")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if !responded || rec.Body.String() != "error" {
		t.Errorf("responded = %v, body = %q, want true/error", responded, rec.Body.String())
	}

	got := calls()
	if len(got) != 1 || got[0].msg != "panic recovered" {
		t.Fatalf("logs = %+v, want one panic entry", got)
	}
	if got[0].args["panic"] != "boom" {
		t.Errorf("panic = %v, want boom", got[0].args["panic"])
	}
	if stack, ok := got[0].args["stack"].(string); !ok || stack == "" {
		t.Errorf("stack = %v, want non-empty string", got[0].args["stack"])
	}
}

func TestRecoverAfterWrite(t *testing.T) {
	log, calls := newLog()
	responded := false
	handler := Recover(log, func(http.ResponseWriter) {
		responded = true
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("partial"))
		panic("late")
	}))

	rec := serve(handler, http.MethodGet, "/late")
	if rec.Code != http.StatusOK || rec.Body.String() != "partial" {
		t.Errorf("response = %d %q, want 200 partial", rec.Code, rec.Body.String())
	}
	if responded {
		t.Error("respond called after the handler wrote a response")
	}
	if len(calls()) != 1 {
		t.Errorf("logs = %+v, want one panic entry", calls())
	}
}

func TestRecoverPassThrough(t *testing.T) {
	log, calls := newLog()
	responded := false
	handler := Recover(log, func(http.ResponseWriter) {
		responded = true
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	rec := serve(handler, http.MethodGet, "/ok")
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Errorf("response = %d %q, want 200 ok", rec.Code, rec.Body.String())
	}
	if responded || len(calls()) != 0 {
		t.Errorf("responded = %v, logs = %+v, want false/none", responded, calls())
	}
}

func TestRecoverLogsRequestStatus(t *testing.T) {
	log, calls := newLog()
	handler := RequestLogger(log)(Recover(log, func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
	})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})))

	serve(handler, http.MethodGet, "/panic")

	got := calls()
	if len(got) != 2 {
		t.Fatalf("logs = %+v, want panic and request entries", got)
	}
	if got[1].msg != "request" || got[1].args["status"] != http.StatusInternalServerError {
		t.Errorf("request entry = %+v, want status 500", got[1])
	}
}
