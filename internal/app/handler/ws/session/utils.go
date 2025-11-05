package session

import taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"

func newResponse(requestId string) *taskmanager.Response {
	return &taskmanager.Response{
		RequestId: requestId,
		Data:      nil,
	}
}
