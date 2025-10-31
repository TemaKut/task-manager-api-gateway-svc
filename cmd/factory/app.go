package factory

import (
	"github.com/google/wire"
)

var AppSet = wire.NewSet(
	ProvideApp,
)

type App struct{}

func ProvideApp(_ *WebsocketProvider) *App {
	return &App{}
}
