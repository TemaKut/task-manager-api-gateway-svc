package session

import (
	"context"
	"fmt"
	taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/proto"
)

type Session struct {
	conn *websocket.Conn

	id string

	logger Logger
}

func NewSession(conn *websocket.Conn, logger Logger) *Session {
	return &Session{
		id:     uuid.NewString(),
		conn:   conn,
		logger: logger,
	}
}

func (s *Session) HandleRequests(ctx context.Context) error {
	for {
		var reqBytes []byte

		if err := websocket.Message.Receive(s.conn, &reqBytes); err != nil {
			return fmt.Errorf("error receiving request. %w", err)
		}

		var req taskmanager.Request

		if err := proto.Unmarshal(reqBytes, &req); err != nil {
			return fmt.Errorf("error unmarshalling request. %w", err)
		}

		_, err := s.handleRequest(ctx, &req)
		if err != nil {
			return fmt.Errorf("error handling request. %w", err)
		}
		// TODO send response

	}
}

// HandleUpdates TODO Вижу реализацию так: брокер кидает сообщения в сервис, затем сервис использует fan-out для всех сессий.
// Сессия подписывается на сервис при помощи коллбэка на различного типа эвенты
func (s *Session) HandleUpdates(ctx context.Context) error {
	return nil
}

func (s *Session) handleRequest(ctx context.Context, req *taskmanager.Request) (*taskmanager.Response, error) {
	if err := uuid.Validate(req.GetId()); err != nil {
		return nil, fmt.Errorf("error validating request id <%s>. %w", req.GetId(), err)
	}

	switch {
	case req.GetUserRegister() != nil:
		return someMiddleware(ctx, func(ctx context.Context) (*taskmanager.Response, error) {
			resp, err := s.handleUserRegisterRequest(ctx, req.GetUserRegister())
			if err != nil {
				return nil, fmt.Errorf("error handle user register request. %w", err)
			}

			return &taskmanager.Response{Data: &taskmanager.Response_UserRegister{UserRegister: resp}}, nil
		})
	default:
		return nil, fmt.Errorf("unknown request type")
	}
}

type handlerFn func(ctx context.Context) (*taskmanager.Response, error)

func someMiddleware(ctx context.Context, handler handlerFn, middlewares ...func(next handlerFn) handlerFn) (*taskmanager.Response, error) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler(ctx)
}

func (s *Session) handleUserRegisterRequest(ctx context.Context, req *taskmanager.UserRegisterRequest) (*taskmanager.UserRegisterResponse, error) {
	return nil, nil
}
