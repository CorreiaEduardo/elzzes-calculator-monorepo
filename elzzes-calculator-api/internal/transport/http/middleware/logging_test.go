package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type logCall struct {
	msg  string
	args map[string]any
}

func newLog() (func(msg string, args ...any), func() []logCall) {
	var calls []logCall
	log := func(msg string, args ...any) {
		fields := make(map[string]any, len(args)/2)
		for i := 0; i+1 < len(args); i += 2 {
			key, ok := args[i].(string)
			if !ok {
				continue
			}
			fields[key] = args[i+1]
		}
		calls = append(calls, logCall{msg: msg, args: fields})
	}
	return log, func() []logCall { return calls }
}

func serve(h http.Handler, method, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestRequestLogger(t *testing.T) {
	log, calls := newLog()
	handler := RequestLogger(log)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	rec := serve(handler, http.MethodGet, "/x?q=1")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}

	got := calls()
	if len(got) != 1 || got[0].msg != "request" {
		t.Fatalf("logs = %+v, want one request entry", got)
	}
	if got[0].args["method"] != http.MethodGet || got[0].args["path"] != "/x" {
		t.Errorf("method/path = %v/%v, want GET//x", got[0].args["method"], got[0].args["path"])
	}
	if got[0].args["status"] != http.StatusCreated {
		t.Errorf("status = %v, want 201", got[0].args["status"])
	}
	if _, ok := got[0].args["duration"].(time.Duration); !ok {
		t.Errorf("duration = %T, want time.Duration", got[0].args["duration"])
	}
}

func TestRequestLoggerImplicit200(t *testing.T) {
	log, calls := newLog()
	handler := RequestLogger(log)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	serve(handler, http.MethodGet, "/x")

	got := calls()
	if len(got) != 1 {
		t.Fatalf("logs = %+v, want one entry", got)
	}
	if got[0].args["status"] != http.StatusOK {
		t.Errorf("status = %v, want 200", got[0].args["status"])
	}
}
