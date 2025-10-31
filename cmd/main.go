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

			_, cleanup, err := factory.InitApp()
			if cleanup != nil {
				defer cleanup()
			}

			if err != nil {
				return fmt.Errorf("error init app: %w", err)
			}

			<-ctx.Done()

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(fmt.Errorf("error run app. %w", err))
	}
}
