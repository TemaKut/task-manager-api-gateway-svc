package factory

import (
	"fmt"
	"github.com/TemaKut/task-manager-api-gateway-svc/internal/app/config"
	"github.com/TemaKut/task-manager-api-gateway-svc/internal/app/handler/ws"
	"github.com/TemaKut/task-manager-api-gateway-svc/internal/app/logger"
	"github.com/google/wire"
)

var AppSet = wire.NewSet(
	ProvideApp,
	config.NewConfig,

	ProvideLogger,
	wire.Bind(new(ws.Logger), new(*logger.Logger)),
)

type App struct{}

func ProvideApp(_ *HttpProvider) *App {
	return &App{}
}

func ProvideLogger(cfg *config.Config) (*logger.Logger, error) {
	var lvl logger.Level

	switch cfg.Logger.Level {
	case config.DebugLevel:
		lvl = logger.DebugLevel
	case config.InfoLevel:
		lvl = logger.InfoLevel
	case config.WarnLevel:
		lvl = logger.WarnLevel
	case config.ErrorLevel:
		lvl = logger.ErrorLevel
	default:
		return nil, fmt.Errorf("error invalid logger level")

	}

	return logger.NewLogger(lvl), nil
}
