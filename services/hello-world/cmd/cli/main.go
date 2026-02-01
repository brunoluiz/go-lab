package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/brunoluiz/go-lab/lib/app"
	"github.com/brunoluiz/go-lab/services/hello-world/internal/service/greet"
)

type CLI struct {
	Language string `kong:"arg,required,name=language,help='Language code for the greeting (e.g., en, es, fr)'"`
}

func (cli *CLI) Run(ctx context.Context, logger *slog.Logger) error {
	greeter := greet.New()

	helloMsg, err := greeter.Hello(cli.Language)
	if err != nil {
		return fmt.Errorf("unable to greet: %w", err)
	}
	logger.InfoContext(ctx, helloMsg)
	return nil
}

func main() {
	app.CLI(&CLI{})
}
