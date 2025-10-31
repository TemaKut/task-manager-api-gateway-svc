package factory

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

var WebsocketSet = wire.NewSet()

type WebsocketProvider struct{}

func ProvideWebsocketProvider() *WebsocketProvider {
	return &WebsocketProvider{}
}

type HttpServerProvider struct{}

func ProvideHttpServerProvider() (*HttpServerProvider, func(), error) {
	server := echo.New()

	server.GET("/ws", func(c echo.Context) error {
		websocket2.Handler(handler.Handle).ServeHTTP(c.Response(), c.Request())

		return nil
	})

	errCh := make(chan error, 1)

	go func() {
		if err := server.Start(cfg.Server.Http.Addr); err != nil {
			errCh <- fmt.Errorf("error starting http server. %w", err)
		}
	}()

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	select {
	case err := <-errCh:
		return nil, nil, fmt.Errorf("error from errCh. %w", err)
	case <-ticker.C:
	}

	return &HttpServerProvider{}, func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Printf("error shutting down http server. %s", err) // TODO
		}
	}, nil
}
