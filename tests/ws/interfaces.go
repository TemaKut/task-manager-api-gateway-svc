package ws

type Logger interface {
	Logf(format string, args ...any)
}
