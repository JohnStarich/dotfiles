package main

import (
	"fmt"
	"math"
	"path"
	"sort"
	"strings"

	"github.com/johnstarich/go/datasize"
)

type reportChange struct {
	Path  string
	Delta int64
}

func createReport(result TMUtilCompare, olderBackupPath, newerBackupPath string) string {
	const macHDPrefix = "Macintosh HD - Data"
	olderBackupPath = path.Join(olderBackupPath, macHDPrefix)
	newerBackupPath = path.Join(newerBackupPath, macHDPrefix)

	var adds, removes, changes []reportChange
	for _, change := range result.Changes {
		switch {
		case change.AddedItem.Path != "":
			adds = append(adds, reportChange{
				Path:  strings.TrimPrefix(change.AddedItem.Path, newerBackupPath),
				Delta: change.AddedItem.Size,
			})
		case change.RemovedItem.Path != "":
			removes = append(removes, reportChange{
				Path:  strings.TrimPrefix(change.RemovedItem.Path, olderBackupPath),
				Delta: change.RemovedItem.Size,
			})
		case change.ChangedItem.Path != "":
			changes = append(changes, reportChange{
				Path:  strings.TrimPrefix(change.ChangedItem.Path, olderBackupPath),
				Delta: change.ChangedItem.Size,
			})
		}
	}

	sort.Slice(adds, func(a, b int) bool {
		return adds[a].Delta >= adds[b].Delta
	})
	sort.Slice(removes, func(a, b int) bool {
		return removes[a].Delta >= removes[b].Delta
	})
	sort.Slice(changes, func(a, b int) bool {
		return changes[a].Delta >= changes[b].Delta
	})

	return buildFileMap(adds).String()
}

func formatSize(size datasize.Size) string {
	sign := "+"
	if size.Bytes() <= 0 {
		sign = ""
	}
	value, unit := size.FormatSI()
	value = math.Ceil(value)
	if len(unit) == 1 {
		unit = " " + unit
	}
	return fmt.Sprintf("%s%.0f %s", sign, value, unit)
}

type fileMap struct {
	Change   reportChange
	Children map[string]*fileMap
}

func newFileMap(change reportChange) *fileMap {
	return &fileMap{
		Change:   change,
		Children: make(map[string]*fileMap),
	}
}

func (f *fileMap) ensureValue(key string, change reportChange) *fileMap {
	value, ok := f.Children[key]
	if !ok {
		f.Children[key] = newFileMap(change)
		value = f.Children[key]
	}
	return value
}

func (f *fileMap) setIntermediateDeltas() int64 {
	if len(f.Children) == 0 {
		return f.Change.Delta
	}

	f.Change.Delta = 0
	for _, child := range f.Children {
		f.Change.Delta += child.setIntermediateDeltas()
	}
	return f.Change.Delta
}

func (f *fileMap) format(parentPath string, depth, maxDepth int) string {
	exceededMaxDepth := maxDepth != -1 && depth > maxDepth
	if len(f.Children) == 1 && !exceededMaxDepth {
		// if only one child, collapse this line into the next by prepending a parent path
		for _, child := range f.Children {
			return child.format(path.Join(parentPath, f.Change.Path), depth+1, maxDepth)
		}
	}

	var sb strings.Builder
	const indent = "  "
	prefix := strings.Repeat(indent, depth)
	size := formatSize(datasize.Bytes(f.Change.Delta))
	currentPath := path.Join("/", parentPath, path.Base(f.Change.Path))
	sb.WriteString(fmt.Sprintf("%8s%s %s\n", size, prefix, currentPath))

	if exceededMaxDepth {
		// stop early if a max depth is specified and we've exceeded it
		return sb.String()
	}

	var children []*fileMap
	for _, child := range f.Children {
		children = append(children, child)
	}
	// sort by delta, then by path (in reverse order below, since course-grained sorts go last)
	sort.SliceStable(children, func(aIndex, bIndex int) bool {
		a, b := children[aIndex].Change, children[bIndex].Change
		return a.Path < b.Path
	})
	sort.SliceStable(children, func(aIndex, bIndex int) bool {
		a, b := children[aIndex].Change, children[bIndex].Change
		return a.Delta > b.Delta
	})
	for _, child := range children {
		sb.WriteString(child.format("", depth+1, maxDepth))
	}
	return sb.String()
}

func (f *fileMap) String() string {
	return f.format("", 0, -1)
}

func buildFileMap(changes []reportChange) *fileMap {
	files := newFileMap(reportChange{})
	for _, change := range changes {
		currentFiles := files
		pathComponents := strings.Split(change.Path, "/")
		for i, pathComponent := range pathComponents {
			value := reportChange{Path: pathComponent}
			if i == len(pathComponents)-1 {
				value = change
			}
			currentFiles = currentFiles.ensureValue(pathComponent, value)
		}
	}
	files.setIntermediateDeltas()
	return files
}
