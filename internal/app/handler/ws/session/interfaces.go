package session

import (
	"context"
	taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"
)

type Logger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any)
}

type handlerFn func(ctx context.Context, req *taskmanager.Request) (*taskmanager.Response, error)

type middlewareFn func(next handlerFn) handlerFn
