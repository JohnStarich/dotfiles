//go:build unix

package main

import (
	"syscall"

	"github.com/pkg/errors"
)

const maxOpenFiles = 1 << 14 // 2^14 appears to be the maximum value for macOS and 2^20 on linux

func increaseNoFile() error {
	// increase soft open file limit to maximum
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		return errors.Wrap(err, "Failed to get open file limit")
	}
	if rLimit.Cur < rLimit.Max && rLimit.Cur < maxOpenFiles {
		rLimit.Cur = min(rLimit.Max, maxOpenFiles)
		err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			return errors.Wrap(err, "Failed to increase soft open file limit")
		}
	}
	return nil
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
