package main

import (
	"context"
	goos "os"
	"os/signal"
	"path"
	"path/filepath"

	"github.com/hack-pad/hackpadfs/os"
)

const appName = "gowerline"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), goos.Interrupt)
	defer cancel()
	err := run(ctx)
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	fs := os.NewFS()
	cacheDir, err := goos.UserCacheDir()
	if err != nil {
		return err
	}
	absoluteCacheDir, err := filepath.Abs(cacheDir)
	if err != nil {
		return err
	}
	cacheSubPath, err := fs.FromOSPath(absoluteCacheDir)
	if err != nil {
		return err
	}
	appCacheSubPath := path.Join(cacheSubPath, appName)
	cacheFS, err := fs.Sub(appCacheSubPath)
	if err != nil {
		return err
	}
	return status(ctx, goos.Stdout, cacheFS)
}
