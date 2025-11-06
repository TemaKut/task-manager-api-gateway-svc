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
	conn   *websocket.Conn
	logger Logger

	responseCallbacks map[string]func(resp *taskmanager.Response, respErr *taskmanager.ResponseError)
	done              chan struct{}
	isClosed          bool
}

func NewClient(ctx context.Context, addr string, logger Logger) (*Client, error) {
	client := Client{
		done:              make(chan struct{}),
		responseCallbacks: make(map[string]func(*taskmanager.Response, *taskmanager.ResponseError)),
		logger:            logger,
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

			var serverMessageBytes []byte

			err := websocket.Message.Receive(conn, &serverMessageBytes)
			if client.isClosed {
				return
			}

			if err != nil {
				logger.Logf("error receiving message. %s", err)

				return
			}

			var serverMessage taskmanager.ServerMessage

			if err := proto.Unmarshal(serverMessageBytes, &serverMessage); err != nil {
				logger.Logf("error unmarshal message. %s", err)

				return
			}

			respCallback, ok := client.responseCallbacks[serverMessage.GetResponse().GetRequestId()]
			if !ok {
				logger.Logf("error response has no callback. %s", err)

				return
			}

			respCallback(serverMessage.GetResponse(), serverMessage.GetResponseError())
		}
	}()

	return &client, nil
}

func (c *Client) SendRequest(ctx context.Context, req *taskmanager.Request) (*ResponseContainer, error) {
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request. %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	responseCh := make(chan *ResponseContainer, 1)
	defer close(responseCh)

	c.responseCallbacks[req.GetId()] = func(resp *taskmanager.Response, respErr *taskmanager.ResponseError) {
		select {
		case <-ctx.Done():
		case responseCh <- &ResponseContainer{Response: resp, ResponseErr: respErr}:
		}
	}

	defer delete(c.responseCallbacks, req.GetId())

	if err := websocket.Message.Send(c.conn, reqBytes); err != nil {
		return nil, fmt.Errorf("error sending request. %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.done:
		return nil, nil
	case respContainer := <-responseCh:
		return respContainer, nil
	}
}

func (c *Client) Close() {
	c.isClosed = true
	close(c.done)

	if err := c.conn.Close(); err != nil {
		c.logger.Logf("error close connection. %s", err)

		return
	}
}

type ResponseContainer struct {
	Response    *taskmanager.Response
	ResponseErr *taskmanager.ResponseError
}
