package factory

import (
	"context"
	"fmt"
	"github.com/TemaKut/task-manager-api-gateway-svc/internal/app/config"
	"github.com/TemaKut/task-manager-api-gateway-svc/internal/app/handler/ws"
	"github.com/TemaKut/task-manager-api-gateway-svc/internal/app/logger"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"time"
)

var HttpSet = wire.NewSet(
	ProvideHttpServerProvider,
	ProvideHttpProvider,
	ws.NewHandler,
)

type HttpProvider struct{}

func ProvideHttpProvider(_ *HttpServerProvider) *HttpProvider {
	return &HttpProvider{}
}

type HttpServerProvider struct{}

func ProvideHttpServerProvider(
	cfg *config.Config,
	log *logger.Logger,
	handler *ws.Handler,
) (*HttpServerProvider, func(), error) {
	server := echo.New()

	server.GET(cfg.HttpServer.Websocket.Path, func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("panic handle wabsocket. %+v", err)
			}
		}()

		websocket.Handler(handler.Handle).ServeHTTP(c.Response(), c.Request())

		return nil
	})

	errCh := make(chan error, 1)

	go func() {
		log.Infof("http server starts listening on %s", cfg.HttpServer.Address)

		if err := server.Start(cfg.HttpServer.Address); err != nil {
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
		log.Infof("http server shutdown")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Errorf("error shutting down http server. %s", err)
		}
	}, nil
}
