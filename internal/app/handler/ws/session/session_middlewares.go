package session

import (
	"context"
	taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"
	"github.com/google/uuid"
)

func (s *Session) requestIdMiddleware(next handlerFn) handlerFn {
	return func(ctx context.Context, req *taskmanager.Request) (*taskmanager.Response, error) {
		if req.GetId() == "" || uuid.Validate(req.GetId()) != nil {
			return nil, ErrRequestHasNoId
		}

		return next(ctx, req)
	}
}
