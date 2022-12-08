package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	const maxBackupDepth = 20
	r, oldBackupPath, newBackupPath, err := generateCompare(ctx, maxBackupDepth)
	if err != nil {
		return err
	}
	result, err := decode(r)
	if err != nil {
		return err
	}
	report := createReport(result, oldBackupPath, newBackupPath)
	fmt.Println(report)
	return nil
}
