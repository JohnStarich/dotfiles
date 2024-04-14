package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oschwald/maxminddb-golang"
	"github.com/pkg/errors"
)

var maxMindDBFile string

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	maxMindDBFile = filepath.Join(configDir, appName, "GeoIP2-latest.mmdb")
}

func weatherStatus(ctx StatusContext) error {
	_, statErr := os.Stat(maxMindDBFile)
	if os.IsNotExist(statErr) {
		err := downloadGeoIPs(ctx.Context)
		if err != nil {
			return errors.WithMessage(err, "failed to set up geo IP database for weather lookup")
		}
	} else if statErr != nil {
		return errors.WithMessage(statErr, "failed to read geo IP database for weather lookup")
	}

	currentIP, err := currentIP(ctx.Context)
	if err != nil {
		return errors.WithMessage(err, "failed to get current IP address for geo IP weather lookup")
	}

	reader, err := maxminddb.Open(maxMindDBFile)
	if err != nil {
		return errors.WithMessage(err, "failed to read geo IP database for weather lookup")
	}
	var coordinates struct {
		Location struct {
			Longitude float64 `maxminddb:"longitude"`
			Latitude  float64 `maxminddb:"latitude"`
		} `maxminddb:"location"`
	}
	err = reader.Lookup(currentIP, &coordinates)
	if err != nil {
		return err
	}

	var weather weatherPoint
	err = doHTTPGet(ctx.Context, fmt.Sprintf("https://api.weather.gov/points/%f,%f", coordinates.Location.Latitude, coordinates.Location.Longitude), &weather)
	if err != nil {
		return err
	}
	var forecast weatherForecast
	err = doHTTPGet(ctx.Context, weather.Properties.ForecastGridData, &forecast)
	if err != nil {
		return err
	}

	now := time.Now()
	_, temp, unit := forecast.Properties.Temperature.RecentMeasurement(now)
	temp, unit = toFahrenheit(temp, unit)

	fmt.Fprintf(ctx.Writer, "🌪  %.f°%s ", temp, unit)
	return nil
}

func toFahrenheit(temperature float64, unit string) (float64, string) {
	if unit != "C" {
		return temperature, unit
	}
	return temperature*9/5 + 32, "F"
}

type weatherPoint struct {
	Properties struct {
		ForecastGridData string
	}
}

type weatherForecast struct {
	Properties struct {
		RelativeHumidity    weatherMeasurements
		ApparentTemperature weatherMeasurements
		Temperature         weatherMeasurements
		WindChill           weatherMeasurements
		WindDirection       weatherMeasurements
		WindSpeed           weatherMeasurements
	}
}

type weatherMeasurements struct {
	UnitOfMeasure string `json:"uom"`
	Values        []weatherMeasurement
}

func (m weatherMeasurements) RecentMeasurement(now time.Time) (t time.Time, value float64, unit string) {
	if len(m.Values) == 0 {
		return time.Time{}, 0, ""
	}
	mostRecent := m.Values[0]
	for _, measurement := range m.Values[1:] {
		if measurement.ValidTime.Time.Before(now) && measurement.ValidTime.Time.After(mostRecent.ValidTime.Time) {
			mostRecent = measurement
		}
	}
	switch m.UnitOfMeasure {
	case "wmoUnit:degC":
		unit = "C"
	case "wmoUnit:degF":
		unit = "F"
	default:
		unit = m.UnitOfMeasure
	}
	return mostRecent.ValidTime.Time, mostRecent.Value, unit
}

type weatherMeasurement struct {
	Value     float64
	ValidTime timeAndDuration
}

type timeAndDuration struct {
	Time time.Time
	// Original value contained a duration. Could unmarshal this in the future if needed.
}

func (i *timeAndDuration) UnmarshalText(bytes []byte) error {
	text := string(bytes)
	if before, _, found := strings.Cut(text, "/"); found {
		text = before
	}

	var err error
	i.Time, err = time.Parse(time.RFC3339, string(text))
	return err
}

func doHTTPGet(ctx context.Context, url string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		return errors.Errorf("failed to fetch from %q: %s", url, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(result)
}

func downloadGeoIPs(ctx context.Context) error {
	err := os.MkdirAll(filepath.Dir(maxMindDBFile), 0o700)
	if err != nil {
		return errors.Wrap(err, "failed to create app directory")
	}
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

	dbFile, err := os.OpenFile(maxMindDBFile, os.O_CREATE|os.O_TRUNC|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return errors.Wrap(err, "failed to create geo IP database file")
	}
	defer dbFile.Close()
	_, err = io.Copy(dbFile, gzipReader)
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
