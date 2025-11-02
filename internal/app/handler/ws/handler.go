package ws

import (
	"fmt"
	"github.com/TemaKut/task-manager-api-gateway-svc/internal/app/handler/ws/session"
	"golang.org/x/net/websocket"
	"golang.org/x/sync/errgroup"
)

type Handler struct {
	logger Logger
}

func NewHandler(logger Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Handle(conn *websocket.Conn) {
	clientAddr := conn.Request().RemoteAddr // TODO Real addr!

	h.logger.Debugf("handle new websocket connection from %s", clientAddr)

	defer func() {
		h.logger.Debugf("close websocket connection from %s", clientAddr)

		if err := conn.Close(); err != nil {
			h.logger.Errorf("error close websocket connection from %s", clientAddr)
		}
	}()

	sess := session.NewSession(conn, h.logger)

	eg, egCtx := errgroup.WithContext(conn.Request().Context())

	eg.Go(func() error {
		if err := sess.HandleRequests(egCtx); err != nil {
			return fmt.Errorf("error handle requests. %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		if err := sess.HandleUpdates(egCtx); err != nil {
			return fmt.Errorf("error handle updates. %w", err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		h.logger.Errorf("error wait session %s wait group. %s", clientAddr, err)

		return
	}
}
