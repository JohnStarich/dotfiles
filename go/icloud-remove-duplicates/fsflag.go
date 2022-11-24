package main

import (
	"path/filepath"

	osfs "github.com/hack-pad/hackpadfs/os"
)

type fsPathFlag struct {
	fromOSPath func(string) (string, error)
	value      *string
}

func fsPathFlagVar(osFS *osfs.FS, value *string, defaultValue string) *fsPathFlag {
	flag := &fsPathFlag{
		fromOSPath: osFS.FromOSPath,
		value:      value,
	}
	err := flag.Set(defaultValue)
	if err != nil {
		panic(err)
	}
	return flag
}

func (f *fsPathFlag) String() string {
	return *f.value
}

func (f *fsPathFlag) Set(s string) error {
	s, err := filepath.Abs(s)
	if err != nil {
		return err
	}
	*f.value, err = f.fromOSPath(s)
	return err
}
