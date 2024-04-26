package weather

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/hack-pad/hackpadfs"
	"github.com/johnstarich/go/gowerline/internal/httpclient"
	"github.com/johnstarich/go/gowerline/internal/status"
	"github.com/oschwald/maxminddb-golang"
	"github.com/pkg/errors"
)

const maxMindDBFileName = "GeoIP2-latest.mmdb"

type ipCoordinates struct {
	Location ipLocation `maxminddb:"location"`
}

type ipLocation struct {
	Longitude float64 `maxminddb:"longitude"`
	Latitude  float64 `maxminddb:"latitude"`
}

type locationCache struct {
	Latitude  float64
	Longitude float64
	ExpiresAt time.Time
}

func getCachedCurrentLocation(ctx status.Context, maxMindDBURL url.URL, now time.Time) (latitude, longitude float64, _ error) {
	coordinatesCache, coordinatesErr := readCachedLocation(ctx.CacheFS)
	if coordinatesErr == nil && coordinatesCache.ExpiresAt.After(now) {
		return coordinatesCache.Latitude, coordinatesCache.Longitude, nil
	}
	latitude, longitude, err := getCurrentLocation(ctx, maxMindDBURL)
	if err != nil {
		if coordinatesErr == nil {
			return coordinatesCache.Latitude, coordinatesCache.Longitude, nil
		}
		return 0, 0, err
	}
	err = writeCachedLocation(ctx.CacheFS, locationCache{
		ExpiresAt: now.Add(1 * time.Hour),
		Latitude:  longitude,
		Longitude: longitude,
	})
	return latitude, longitude, err
}

const locationCacheFileName = "location.json"

func readCachedLocation(fs fs.FS) (locationCache, error) {
	var cachedCoordinates locationCache
	contents, err := hackpadfs.ReadFile(fs, locationCacheFileName)
	if err != nil {
		return locationCache{}, err
	}
	err = json.Unmarshal(contents, &cachedCoordinates)
	if err != nil {
		return locationCache{}, err
	}
	return cachedCoordinates, nil
}

func writeCachedLocation(fs fs.FS, newLocation locationCache) error {
	newContents, err := json.MarshalIndent(newLocation, "", "    ")
	if err != nil {
		return err
	}
	return hackpadfs.WriteFullFile(fs, locationCacheFileName, newContents, 0o700)
}

func getCurrentLocation(ctx status.Context, maxMindDBURL url.URL) (latitude, longitude float64, _ error) {
	_, statErr := hackpadfs.Stat(ctx.CacheFS, maxMindDBFileName)
	if errors.Is(statErr, hackpadfs.ErrNotExist) {
		err := downloadGeoIPs(ctx.Context, maxMindDBURL, ctx.HTTPClient, ctx.CacheFS, ctx.Now())
		if err != nil {
			return 0, 0, errors.WithMessage(err, "failed to set up geo IP database for weather lookup")
		}
	} else if statErr != nil {
		return 0, 0, errors.WithMessage(statErr, "failed to read geo IP database for weather lookup")
	}

	currentIP, err := ctx.Resolver.LookupIPWithResolverHost(ctx.Context, "resolver1.opendns.com:53", "myip.opendns.com")
	if err != nil {
		return 0, 0, errors.WithMessage(err, "failed to get current IP address for geo IP weather lookup")
	}

	maxMindDBFile, err := ctx.CacheFS.Open(maxMindDBFileName)
	if err != nil {
		return 0, 0, err
	}
	defer maxMindDBFile.Close()
	maxMindDBBytes, err := io.ReadAll(maxMindDBFile)
	if err != nil {
		return 0, 0, err
	}
	reader, err := maxminddb.FromBytes(maxMindDBBytes)
	if err != nil {
		return 0, 0, errors.WithMessage(err, "failed to read geo IP database for weather lookup")
	}
	var coordinates ipCoordinates
	err = reader.Lookup(currentIP, &coordinates)
	if err != nil {
		return 0, 0, err
	}
	if coordinates == (ipCoordinates{}) {
		return 0, 0, errors.New("failed to get valid coordinates for current location")
	}
	return coordinates.Location.Latitude, coordinates.Location.Longitude, nil
}

func downloadGeoIPs(ctx context.Context, maxMindDBURL url.URL, httpClient httpclient.Client, cacheFS fs.FS, now time.Time) error {
	thisMonth := now.Format("2006-01")
	downloadURL := maxMindDBURL
	downloadURL.Path = path.Join(downloadURL.Path, "free", fmt.Sprintf("dbip-city-lite-%s.mmdb.gz", thisMonth))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL.String(), nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.Errorf("failed to fetch geo IP database: %s", string(body))
	}
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	dbFile, err := hackpadfs.OpenFile(cacheFS, maxMindDBFileName, os.O_CREATE|os.O_TRUNC|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return errors.Wrap(err, "failed to create geo IP database file")
	}
	defer dbFile.Close()
	_, err = io.Copy(dbFile.(io.Writer), gzipReader)
	return errors.Wrap(err, "failed to download latest geo IP database")
}
