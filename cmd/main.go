package main

import (
	"fmt"
	"github.com/TemaKut/task-manager-api-gateway-svc/cmd/factory"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := cli.App{
		Action: func(c *cli.Context) error {
			ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
			defer cancel()

			_ = factory.InitApp()

			<-ctx.Done()

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(fmt.Errorf("error run app. %w", err))
	}
}
