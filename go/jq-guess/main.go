package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, in io.Reader, out, errOut io.Writer) error {
	userInput, parserOutput, err := os.Pipe()
	if err != nil {
		return err
	}
	defer parserOutput.Close()
	defer userInput.Close()

	const jobs = 2 // parser + jq with user args
	errs := make(chan error, jobs)
	startJQ(ctx, []string{
		"--raw-input",
		"--unbuffered",
		`. as $line | select(. != "") | try fromjson catch {"_json_parse_error":$line}`,
	}, in, parserOutput, errOut, errs)
	args = append([]string{"--unbuffered"}, args...) // in this situation, only --unbuffered produces useful results in interactive sessions
	startJQ(ctx, args, userInput, out, errOut, errs)
	return <-errs
}

func startJQ(ctx context.Context, args []string, in io.Reader, out io.Writer, errOut io.Writer, errs chan<- error) {
	cmd := exec.CommandContext(ctx, "jq", args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = errOut
	go func() {
		errs <- cmd.Run()
	}()
}
