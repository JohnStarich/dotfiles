package weather

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hack-pad/hackpadfs"
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

func getCachedCurrentLocation(ctx status.Context, now time.Time) (latitude, longitude float64, _ error) {
	var cachedCoordinates locationCache
	const locationCacheFileName = "location.json"
	contents, err := hackpadfs.ReadFile(ctx.CacheFS, locationCacheFileName)
	if err == nil {
		err = json.Unmarshal(contents, &cachedCoordinates)
	}
	if err == nil && cachedCoordinates.ExpiresAt.After(now) {
		return cachedCoordinates.Latitude, cachedCoordinates.Longitude, nil
	}
	latitude, longitude, err = getCurrentLocation(ctx)
	if err != nil {
		return 0, 0, err
	}
	newContents, err := json.MarshalIndent(locationCache{
		Latitude:  latitude,
		Longitude: longitude,
		ExpiresAt: now.Add(1 * time.Hour),
	}, "", "    ")
	if err == nil {
		err = hackpadfs.WriteFullFile(ctx.CacheFS, locationCacheFileName, newContents, 0o700)
	}
	return latitude, longitude, err
}

func getCurrentLocation(ctx status.Context) (latitude, longitude float64, _ error) {
	_, statErr := hackpadfs.Stat(ctx.CacheFS, maxMindDBFileName)
	if errors.Is(statErr, hackpadfs.ErrNotExist) {
		err := downloadGeoIPs(ctx.Context, ctx.CacheFS)
		if err != nil {
			return 0, 0, errors.WithMessage(err, "failed to set up geo IP database for weather lookup")
		}
	} else if statErr != nil {
		return 0, 0, errors.WithMessage(statErr, "failed to read geo IP database for weather lookup")
	}

	currentIP, err := currentIP(ctx.Context)
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

func downloadGeoIPs(ctx context.Context, cacheFS fs.FS) error {
	thisMonth := time.Now().Format("2006-01")
	downloadURL := fmt.Sprintf("https://download.db-ip.com/free/dbip-city-lite-%s.mmdb.gz", thisMonth)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
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

func currentIP(ctx context.Context) (net.IP, error) {
	// Equivalent of running: nslookup myip.opendns.com resolver1.opendns.com
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp", "resolver1.opendns.com:53")
		},
	}
	ipAddresses, err := resolver.LookupIPAddr(ctx, "myip.opendns.com")
	if err != nil {
		return nil, err
	}
	if len(ipAddresses) == 0 {
		return nil, errors.New("could not resolve current IP address")
	}
	return ipAddresses[0].IP, nil
}
