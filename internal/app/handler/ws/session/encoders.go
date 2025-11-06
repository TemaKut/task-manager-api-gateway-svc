package session

import (
	"errors"
	taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func encodeError(requestId string, err error) *taskmanager.ResponseError {
	if err == nil {
		return nil
	}

	responseError := &taskmanager.ResponseError{
		RequestId:  requestId,
		ServerTime: timestamppb.Now(),
	}

	switch {
	case errors.Is(err, ErrRequestHasNoId):
		responseError.ErrorType = taskmanager.ResponseErrorType_RESPONSE_ERROR_TYPE_REQUEST_HAS_NO_ID
		responseError.Description = "request has no id"
	default:
		responseError.ErrorType = taskmanager.ResponseErrorType_RESPONSE_ERROR_TYPE_UNSPECIFIED
		responseError.Description = "unknown error"
	}

	return responseError
}
