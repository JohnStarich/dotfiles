package status

import (
	"encoding/json"
	"io/fs"
	"time"

	"github.com/hack-pad/hackpadfs"
	"github.com/pkg/errors"
)

type lineCache struct {
	Segments map[string]SegmentCache
}

type SegmentCache struct {
	Content   string
	ExpiresAt time.Time
}

const statusLineCacheFileName = "status-cache.json"

func readLineCache(fs fs.FS) (lineCache, error) {
	lineCacheBytes, err := hackpadfs.ReadFile(fs, statusLineCacheFileName)
	if err != nil && !errors.Is(err, hackpadfs.ErrNotExist) {
		return lineCache{}, err
	}
	var data lineCache
	err = json.Unmarshal(lineCacheBytes, &data)
	return data, err
}

func writeLineCache(fs fs.FS, data lineCache) error {
	lineCacheBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return hackpadfs.WriteFullFile(fs, statusLineCacheFileName, lineCacheBytes, 0o600)
}
