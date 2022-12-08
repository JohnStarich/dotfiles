package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func generateCompare(ctx context.Context, maxDepth int) (r io.Reader, oldBackupPath, newBackupPath string, err error) {
	const tmutil = "tmutil"
	cmd := exec.CommandContext(ctx, tmutil, "listbackups", "-m")
	cmd.Stderr = os.Stderr
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err = cmd.Run()
	if err != nil {
		return nil, "", "", err
	}
	lines := strings.Split(buf.String(), "\n")
	const compareNBackups = 2
	var backups []string
	for i := len(lines) - 1; len(backups) < compareNBackups && i >= 0; i-- {
		backup := strings.TrimSpace(lines[i])
		if backup != "" {
			backups = append(backups, backup)
		}
	}
	if len(backups) != compareNBackups {
		return nil, "", "", errors.New("failed to find latest 2 backups")
	}

	fmt.Println("Comparing latest 2 backups:")
	for _, b := range backups {
		fmt.Println(b)
	}

	args := append([]string{"compare", "-n", "-s", "-X", "-D", strconv.Itoa(maxDepth)}, backups...)
	cmd = exec.CommandContext(ctx, "tmutil", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, "", "", err
	}
	cmd.Stderr = os.Stderr
	fmt.Println("Running compare...")
	err = cmd.Start()
	if err != nil {
		return nil, "", "", err
	}
	return stdout, backups[0], backups[1], nil
}
