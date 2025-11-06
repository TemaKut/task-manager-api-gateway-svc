package session

import (
	"context"
	"errors"
	"fmt"
	taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/proto"
	"io"
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

		err := websocket.Message.Receive(s.conn, &reqBytes)
		switch {
		case errors.Is(err, io.EOF):
			println(1)
			return nil
		default:
			if err != nil {
				return fmt.Errorf("error receiving request. %w", err)
			}
		}

		var req taskmanager.Request

		if err := proto.Unmarshal(reqBytes, &req); err != nil {
			return fmt.Errorf("error unmarshalling request. %w", err)
		}

		resp, err := s.handleRequest(ctx, &req)
		if err != nil {
			s.logger.Errorf("error handle request id <%s>. %s", req.GetId(), err)

			respErrorBytes, err := proto.Marshal(&taskmanager.ServerMessage{
				Data: &taskmanager.ServerMessage_ResponseError{
					ResponseError: encodeError(req.GetId(), err),
				},
			})
			if err != nil {
				return fmt.Errorf("error marshalling response error. %w", err)
			}

			if err := websocket.Message.Send(s.conn, respErrorBytes); err != nil {
				return fmt.Errorf("error sending response error. %w", err)
			}

			continue
		}

		respBytes, err := proto.Marshal(&taskmanager.ServerMessage{
			Data: &taskmanager.ServerMessage_Response{
				Response: resp,
			},
		})
		if err != nil {
			return fmt.Errorf("error marshalling response. %w", err)
		}

		if err := websocket.Message.Send(s.conn, respBytes); err != nil {
			return fmt.Errorf("error sending response. %w", err)
		}
	}
}

// HandleUpdates TODO Вижу реализацию так: брокер кидает сообщения в сервис, затем сервис использует fan-out для всех сессий.
// Сессия подписывается на сервис при помощи коллбэка на различного типа эвенты
func (s *Session) HandleUpdates(ctx context.Context) error {
	return nil
}

func (s *Session) handleRequest(ctx context.Context, req *taskmanager.Request) (*taskmanager.Response, error) {
	var handler handlerFn

	switch {
	case req.GetUserRegister() != nil:
		handler = s.withMiddleware(
			[]middlewareFn{
				s.requestIdMiddleware,
			},
			s.handleUserRegisterRequest,
		)
	default:
		return nil, ErrUnknownRequestType
	}

	if handler == nil {
		return nil, fmt.Errorf("error handler is nil")
	}

	return handler(ctx, req)
}

func (s *Session) withMiddleware(middlewares []middlewareFn, handler handlerFn) handlerFn {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func (s *Session) handleUserRegisterRequest(ctx context.Context, req *taskmanager.Request) (*taskmanager.Response, error) {
	//request := req.GetUserRegister()

	return newResponse(req.GetId()), nil
}
