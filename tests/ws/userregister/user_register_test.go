package userregister

import (
	"context"
	"fmt"
	"github.com/TemaKut/task-manager-api-gateway-svc/tests/ws"
	taskmanager "github.com/TemaKut/task-manager-client-proto/gen/go"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"testing"
)

func TestUserRegister(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	client, err := ws.NewClient(ctx, ws.ApiGatewayServiceAddr)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating client: %w", err))
	}

	defer client.Close()

	req := &taskmanager.Request{
		Id: uuid.New().String(),
		Data: &taskmanager.Request_UserRegister{
			UserRegister: &taskmanager.UserRegisterRequest{Name: "test"},
		},
	}

	resp, err := client.SendRequest(ctx, req)
	if err != nil {
		t.Fatal(fmt.Errorf("error sending request. %w", err))
	}

	fmt.Println(resp)
}
