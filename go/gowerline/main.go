package main

import (
	"context"
	"os"
	"os/signal"
)

const appName = "gowerline"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := run(ctx)
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	return status(ctx, os.Stdout)
}