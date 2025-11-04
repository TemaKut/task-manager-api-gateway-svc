package session

type Logger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any)
}
