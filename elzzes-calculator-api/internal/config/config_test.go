package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    Config
		wantErr string
	}{
		{name: "unset uses default", raw: "", want: Config{Port: 8080}},
		{name: "valid port", raw: "8080", want: Config{Port: 8080}},
		{name: "minimum port", raw: "1", want: Config{Port: 1}},
		{name: "maximum port", raw: "65535", want: Config{Port: 65535}},
		{name: "zero rejected", raw: "0", wantErr: "invalid PORT"},
		{name: "negative rejected", raw: "-1", wantErr: "invalid PORT"},
		{name: "too large rejected", raw: "65536", wantErr: "invalid PORT"},
		{name: "letters rejected", raw: "abc", wantErr: "invalid PORT"},
		{name: "trailing junk rejected", raw: "8080x", wantErr: "invalid PORT"},
		{name: "leading space rejected", raw: " 8080", wantErr: "invalid PORT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string {
				if key == "PORT" {
					return tt.raw
				}
				return ""
			}

			got, err := Load(getenv)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Load() error = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestConfigAddr(t *testing.T) {
	tests := []struct {
		port int
		want string
	}{
		{port: 8080, want: ":8080"},
		{port: 1, want: ":1"},
	}

	for _, tt := range tests {
		got := (Config{Port: tt.port}).Addr()
		if got != tt.want {
			t.Errorf("Config{Port: %d}.Addr() = %q, want %q", tt.port, got, tt.want)
		}
	}
}
