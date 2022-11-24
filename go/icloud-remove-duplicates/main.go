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

	"github.com/fatih/color"
	"github.com/hack-pad/hackpadfs"
	osfs "github.com/hack-pad/hackpadfs/os"
	"github.com/pkg/errors"
)

func main() {
	err := run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %+v\n", color.RedString("Failed:"), err)
		os.Exit(1)
	}
}

type Args struct {
	DryRun    bool
	ErrWriter io.Writer
	FS        hackpadfs.FS
	OutWriter io.Writer
	Root      string
}

func run(osArgs []string, outWriter, errWriter io.Writer) error {
	fs := osfs.NewFS()
	args := Args{
		ErrWriter: errWriter,
		FS:        fs,
		OutWriter: outWriter,
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
		isDuplicate, warningErr, err := isDuplicate(filePath, args.FS)
		if err != nil {
			return err
		}
		if warningErr != nil {
			fmt.Fprintln(args.ErrWriter, color.YellowString("Warning:"), warningErr)
		}
		if isDuplicate {
			duplicateFiles = append(duplicateFiles, filePath)
		}
	}

	for _, filePath := range duplicateFiles {
		if args.DryRun {
			fmt.Fprintln(args.OutWriter, filePath)
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

func isDuplicate(candidateFilePath string, fs hackpadfs.FS) (_ bool, warningErr, err error) {
	defer func() {
		warningErr = errors.Wrapf(warningErr, "checking for duplicate from %q", candidateFilePath)
		err = errors.Wrapf(err, "checking for duplicate from %q", candidateFilePath)
	}()

	name := path.Base(candidateFilePath)
	ext := path.Ext(name)
	nameNoExt := strings.TrimSuffix(name, ext)
	rootName := strings.TrimRightFunc(nameNoExt, func(r rune) bool {
		return unicode.IsDigit(r) || r == ' '
	}) + ext
	rootFilePath := path.Join(path.Dir(candidateFilePath), rootName)

	rootInfo, err := hackpadfs.LstatOrStat(fs, rootFilePath)
	if err != nil && !errors.Is(err, hackpadfs.ErrNotExist) {
		return false, nil, err
	}
	if errors.Is(err, hackpadfs.ErrNotExist) {
		// Failed to find root file, so attempt brctl download.
		err = brctlDownload(fs, rootFilePath)
		if err != nil {
			return false, err, nil
		}
		rootInfo, err = hackpadfs.LstatOrStat(fs, rootFilePath)
		if errors.Is(err, hackpadfs.ErrNotExist) {
			// file really doesn't exist, skip...
			return false, errors.New("looks like an icloud duplicate, but no root file was found"), nil
		}
		if err != nil {
			return false, nil, err
		}
	}
	candidateInfo, err := hackpadfs.LstatOrStat(fs, rootFilePath)
	if err != nil {
		return false, nil, err
	}
	if rootInfo.Mode().Type() != candidateInfo.Mode().Type() {
		// mismatched types, skip...
		return false, errors.Errorf("type bits do not match: %s != %s", rootInfo.Mode().Type().String(), candidateInfo.Mode().Type().String()), nil
	}
	if rootInfo.Mode()&hackpadfs.ModeSymlink != 0 {
		// NOTE: Doesn't check the link's destination.
		// For cleaning up duplicates, this usually isn't a problem. Worst case, the link can be recreated.
		return true, nil, nil
	}
	if !rootInfo.Mode().IsRegular() {
		// unexpected typed file, skip...
		return false, errors.Errorf("unrecognized file type: %s", rootInfo.Mode().Type().String()), nil
	}
	if rootInfo.Size() != candidateInfo.Size() {
		// sizes differ, skip...
		return false, errors.Errorf("file size differs: %d != %d", rootInfo.Size(), candidateInfo.Size()), nil
	}

	rootHash, err := hashFile(fs, rootFilePath)
	if err != nil {
		return false, nil, err
	}
	candidateHash, err := hashFile(fs, candidateFilePath)
	if err != nil {
		return false, nil, err
	}
	return rootHash == candidateHash, nil, nil
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
