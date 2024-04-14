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
	var lineCacheData lineCache
	lineCacheBytes, cacheReadErr := hackpadfs.ReadFile(fs, statusLineCacheFileName)
	if cacheReadErr == nil {
		err := json.Unmarshal(lineCacheBytes, &lineCacheData)
		if err != nil {
			return lineCache{}, err
		}
	} else if !errors.Is(cacheReadErr, hackpadfs.ErrNotExist) {
		return lineCache{}, cacheReadErr
	}

	if lineCacheData.Segments == nil {
		lineCacheData.Segments = make(map[string]SegmentCache)
	}
	return lineCacheData, nil
}

func writeLineCache(fs fs.FS, data lineCache) error {
	lineCacheBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return hackpadfs.WriteFullFile(fs, statusLineCacheFileName, lineCacheBytes, 0o600)
}
