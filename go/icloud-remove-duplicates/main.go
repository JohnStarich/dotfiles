package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"unicode"

	"github.com/hack-pad/hackpadfs"
	osfs "github.com/hack-pad/hackpadfs/os"
	"github.com/pkg/errors"
)

func main() {
	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v\n", err)
		os.Exit(1)
	}
}

type Args struct {
	DryRun bool
	FS     hackpadfs.FS
	Root   string
}

func run(osArgs []string) error {
	fs := osfs.NewFS()
	args := Args{
		FS: fs,
	}
	set := flag.NewFlagSet("icloud-remove-duplicates", flag.ContinueOnError)
	set.BoolVar(&args.DryRun, "dry-run", true, "Set to --dry-run=false to actually clean up the files printed from the dry run.")
	set.Var(fsPathFlagVar(fs, &args.Root, "."), "root", "The directory to scan and remove duplicates from.")
	err := set.Parse(osArgs)
	if err != nil {
		return err
	}
	return runArgs(args)
}

func runArgs(args Args) error {
	duplicateCandidates := make(map[string]bool)
	err := hackpadfs.WalkDir(args.FS, args.Root, func(filePath string, dirEntry hackpadfs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := dirEntry.Name()
		nameNoExt := strings.TrimSuffix(name, path.Ext(name))
		if endInSpaceAndPositiveInteger(nameNoExt) {
			duplicateCandidates[filePath] = true
		}
		return nil
	})
	if err != nil {
		return err
	}

	var duplicateFiles []string
	for filePath := range duplicateCandidates {
		isDuplicate, err := isDuplicate(filePath, args.FS)
		if err != nil {
			return err
		}
		if isDuplicate {
			duplicateFiles = append(duplicateFiles, filePath)
		}
	}

	for _, filePath := range duplicateFiles {
		if args.DryRun {
			fmt.Println(filePath)
		} else {
			err := hackpadfs.Remove(args.FS, filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func endInSpaceAndPositiveInteger(s string) bool {
	lastSpace := strings.LastIndexByte(s, ' ')
	if lastSpace == -1 {
		return false
	}
	lastWord := s[lastSpace+1:]
	lastWordInt, err := strconv.Atoi(lastWord)
	return err == nil && lastWordInt > 0
}

func isDuplicate(candidateFilePath string, fs hackpadfs.FS) (_ bool, err error) {
	defer func() { err = errors.WithStack(err) }()

	name := path.Base(candidateFilePath)
	ext := path.Ext(name)
	nameNoExt := strings.TrimSuffix(name, ext)
	rootName := strings.TrimRightFunc(nameNoExt, func(r rune) bool {
		return unicode.IsDigit(r) || r == ' '
	}) + ext
	rootFilePath := path.Join(path.Dir(candidateFilePath), rootName)

	rootInfo, err := hackpadfs.LstatOrStat(fs, rootFilePath)
	if err != nil && !errors.Is(err, hackpadfs.ErrNotExist) {
		return false, err
	}
	if errors.Is(err, hackpadfs.ErrNotExist) {
		// Failed to find root file, so attempt brctl download.
		err = brctlDownload(fs, rootFilePath)
		if err != nil {
			return false, err
		}
		rootInfo, err = hackpadfs.LstatOrStat(fs, rootFilePath)
		if errors.Is(err, hackpadfs.ErrNotExist) {
			// file really doesn't exist, skip...
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	candidateInfo, err := hackpadfs.LstatOrStat(fs, rootFilePath)
	if err != nil {
		return false, err
	}
	if rootInfo.Mode().Type() != candidateInfo.Mode().Type() {
		// mismatched types, skip...
		return false, nil
	}
	if rootInfo.Mode()&hackpadfs.ModeSymlink != 0 {
		// NOTE: Doesn't check the link's destination.
		// For cleaning up duplicates, this usually isn't a problem. Worst case, the link can be recreated.
		return true, nil
	}
	if !rootInfo.Mode().IsRegular() {
		// unexpected typed file, skip...
		return false, nil
	}
	if rootInfo.Size() != candidateInfo.Size() {
		// sizes differ, skip...
		return false, nil
	}

	rootHash, err := hashFile(fs, rootFilePath)
	if err != nil {
		return false, err
	}
	candidateHash, err := hashFile(fs, candidateFilePath)
	if err != nil {
		return false, err
	}
	return rootHash == candidateHash, nil
}

func hashFile(fs hackpadfs.FS, filePath string) (string, error) {
	file, err := fs.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := md5.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}
	return string(hasher.Sum(nil)), nil
}

type osPathFS interface {
	ToOSPath(fsPath string) (osPath string, err error)
}

func brctlDownload(fs hackpadfs.FS, filePath string) error {
	if fs, ok := fs.(osPathFS); ok {
		var err error
		filePath, err = fs.ToOSPath(filePath)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command("brctl", "download", filePath)
	return cmd.Run()
}
