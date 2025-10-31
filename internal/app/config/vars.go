package config

type LoggerLevel int

const (
	DebugLevel LoggerLevel = 1
	InfoLevel  LoggerLevel = 2
	WarnLevel  LoggerLevel = 3
	ErrorLevel LoggerLevel = 4
)
