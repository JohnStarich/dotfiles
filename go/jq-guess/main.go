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
	const jobs = 3 // parser + jq with user args + stdout copier
	errs := make(chan error, jobs)
	parserOutput, err := startJQ(ctx, []string{
		"--raw-input",
		"--unbuffered",
		`. as $line | select(. != "") | try fromjson catch {"_json_parse_error":$line}`,
	}, in, errOut, errs)
	if err != nil {
		return err
	}
	args = append([]string{"--unbuffered"}, args...) // in this situation, only --unbuffered produces useful results in interactive sessions
	output, err := startJQ(ctx, args, parserOutput, errOut, errs)
	if err != nil {
		return err
	}
	go func() {
		defer parserOutput.Close()
		defer output.Close()
		_, err := io.Copy(out, output)
		if err != nil {
			errs <- err
		}
	}()
	return <-errs
}

func startJQ(ctx context.Context, args []string, in io.Reader, errOut io.Writer, errs chan<- error) (io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, "jq", args...)
	cmd.Stdin = in
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = errOut
	go func() {
		errs <- cmd.Run()
	}()
	return out, nil
}
