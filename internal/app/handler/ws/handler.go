package ws

import "golang.org/x/net/websocket"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ws *websocket.Conn) {
	// TODO observability of open sessions...
}
