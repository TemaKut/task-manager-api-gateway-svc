package config

type Config struct {
	Logger struct {
		Level LoggerLevel
	}
	HttpServer struct {
		Address   string
		Websocket struct {
			Path string
		}
	}
}

func NewConfig() *Config {
	var cfg Config

	// Default state
	cfg.HttpServer.Address = ":8000"
	cfg.HttpServer.Websocket.Path = "/ws"

	cfg.Logger.Level = DebugLevel
	// ^^^^^^^^^^^^^^

	// TODO parse from .env

	return &cfg
}
