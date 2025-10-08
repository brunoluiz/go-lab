package main

import (
	"log/slog"
	"os"

	"github.com/brunoluiz/go-lab/services/hello-world/internal/service/greet"
)

func main() {
	greeter := greet.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if len(os.Args) < 2 {
		logger.Error("language argument is required")
		return
	}

	helloMsg, err := greeter.Hello(os.Args[1])
	if err != nil {
		logger.Error("unable to greet", slog.String("error", err.Error()))
		return
	}
	logger.Info(helloMsg)
}
