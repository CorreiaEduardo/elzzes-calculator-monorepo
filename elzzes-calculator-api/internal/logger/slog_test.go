package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestLoggerWritesStructuredJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)

	logger.Info("request", "method", "GET", "status", 200)
	logger.Error("boom", "error", "failed")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d log lines, want 2: %q", len(lines), buf.String())
	}

	info := decodeLine(t, lines[0])
	if info["level"] != "INFO" || info["msg"] != "request" {
		t.Errorf("info entry level/msg = %v/%v, want INFO/request", info["level"], info["msg"])
	}
	if info["method"] != "GET" {
		t.Errorf("info entry method = %v, want GET", info["method"])
	}
	if status, ok := info["status"].(float64); !ok || status != 200 {
		t.Errorf("info entry status = %v, want 200", info["status"])
	}

	errorEntry := decodeLine(t, lines[1])
	if errorEntry["level"] != "ERROR" || errorEntry["msg"] != "boom" {
		t.Errorf("error entry level/msg = %v/%v, want ERROR/boom", errorEntry["level"], errorEntry["msg"])
	}
	if errorEntry["error"] != "failed" {
		t.Errorf("error entry error = %v, want failed", errorEntry["error"])
	}
}

func decodeLine(t *testing.T, line string) map[string]any {
	t.Helper()
	var entry map[string]any
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("decode log line %q: %v", line, err)
	}
	return entry
}
