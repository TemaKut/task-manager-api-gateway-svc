package ws

import (
	"context"
	"fmt"
	taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/proto"
	"time"
)

type Client struct {
	conn *websocket.Conn

	responses map[string]func(*taskmanager.Response)
	done      chan struct{}
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	client := Client{
		done:      make(chan struct{}),
		responses: make(map[string]func(*taskmanager.Response)),
	}

	cfg, err := websocket.NewConfig(addr, "http://localhost")
	if err != nil {
		return nil, fmt.Errorf("error make config. %w", err)
	}

	conn, err := cfg.DialContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error dial context. %w", err)
	}

	client.conn = conn

	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, fmt.Errorf("error setting read deadline. %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
			case <-client.done:
			default:
			}

			var respBytes []byte

			if err := websocket.Message.Receive(conn, &respBytes); err != nil {
				return
			}

			var resp taskmanager.Response

			if err := proto.Unmarshal(respBytes, &resp); err != nil {
				return
			}

			client.responses[resp.GetId()](&resp)
		}
	}()

	return &client, nil
}

// TODO учесть что может респонс может прилететь не для текущего запроса + апдейты.
// TODO учесть кейс когда сервак не отправил ничего (Например сеть мигнула)
func (c *Client) SendRequest(ctx context.Context, req *taskmanager.Request) (*taskmanager.Response, error) {
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request. %w", err)
	}

	if err := websocket.Message.Send(c.conn, reqBytes); err != nil {
		return nil, fmt.Errorf("error sending request. %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	responseCh := make(chan *taskmanager.Response, 1)
	defer close(responseCh)

	c.responses[req.GetId()] = func(resp *taskmanager.Response) {
		select {
		case <-ctx.Done():
		case responseCh <- resp:
		default:
		}
	}

	defer delete(c.responses, req.GetId())

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.done:
		return nil, nil
	case resp := <-responseCh:
		return resp, nil
	}
}

func (c *Client) Close() error {
	close(c.done)

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("error close connection. %w", err)
	}

	return nil
}
