//go:build wireinject
// +build wireinject

package factory

import "github.com/google/wire"

func InitApp() *App {
	panic(
		wire.Build(
			AppSet,
			ProvideWebsocketProvider,
		),
	)
}
