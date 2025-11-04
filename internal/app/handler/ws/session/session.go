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

		resp, err := s.handleRequest(ctx, &req)
		if err != nil {
			s.logger.Errorf("error handle request id <%s>. %s", req.GetId(), err)
			continue
		}

		respBytes, err := proto.Marshal(resp)
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
	switch {
	case req.GetUserRegister() != nil:

	default:
		return nil, fmt.Errorf("unknown request type")
	}
}

func (s *Session) handleUserRegisterRequest(ctx context.Context, req *taskmanager.Request) (*taskmanager.Response, error) {
	//request := req.GetUserRegister()

	return newResponse(req.GetId()), nil
}
