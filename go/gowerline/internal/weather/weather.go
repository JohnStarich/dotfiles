package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/johnstarich/go/gowerline/internal/status"
	"github.com/pkg/errors"
)

func getLatestWeather(ctx status.Context, latitude, longitude float64) error {
	var weather weatherPoint
	err := doHTTPGet(ctx.Context, fmt.Sprintf("https://api.weather.gov/points/%f,%f", latitude, longitude), &weather)
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
	_, w := forecast.Properties.Weather.RecentMeasurement(now)

	fmt.Fprintf(ctx.Writer, "%s  %.f°%s ", w.Icon(), temp, unit)
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
		Weather             weatherValueMeasurements
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

type weatherValueMeasurements struct {
	Values []weatherValueMeasurement
}

type weatherValueMeasurement struct {
	ValidTime timeAndDuration
	Value     []weatherValue
}

func (m weatherValueMeasurements) RecentMeasurement(now time.Time) (t time.Time, value weatherEnum) {
	if len(m.Values) == 0 {
		return time.Time{}, weatherUnknown
	}
	mostRecent := m.Values[0]
	var mostRecentWeather weatherEnum
	for _, measurement := range m.Values[1:] {
		if measurement.ValidTime.Time.Before(now) && measurement.ValidTime.Time.After(mostRecent.ValidTime.Time) {
			for _, value := range measurement.Value {
				if value.Weather != nil {
					mostRecent = measurement
					mostRecentWeather = *value.Weather
					break
				}
			}
		}
	}
	return mostRecent.ValidTime.Time, mostRecentWeather
}

type weatherValue struct {
	Weather *weatherEnum
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
