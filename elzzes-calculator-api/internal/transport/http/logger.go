package http

// Logger is the logging capability the HTTP transport needs. The port lives here,
// owned by its consumer; internal/logger provides the adapter.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
