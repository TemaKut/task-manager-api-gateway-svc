package session

import "errors"

var (
	ErrRequestHasNoId     = errors.New("request has no id")
	ErrUnknownRequestType = errors.New("unknown request type")
)
