package session

import (
	"context"
	taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"
)

type Logger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any)
}

type requestHandler func(ctx context.Context, req *taskmanager.Request) error

type requestMiddleware func(next requestHandler) requestHandler
