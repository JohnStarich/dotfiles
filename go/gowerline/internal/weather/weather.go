package weather

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/johnstarich/go/gowerline/internal/status"
)

func getLatestWeather(ctx status.Context, weatherGovURL url.URL, latitude, longitude float64) (string, error) {
	pointURL := weatherGovURL
	pointURL.Path = path.Join(pointURL.Path, "points", fmt.Sprintf("%f,%f", latitude, longitude))
	var point Point
	err := doJSONGet(ctx.Context, ctx.HTTPClient, pointURL.String(), &point)
	if err != nil {
		return "", err
	}
	var forecast Forecast
	err = doJSONGet(ctx.Context, ctx.HTTPClient, point.Properties.ForecastGridData, &forecast)
	if err != nil {
		return "", err
	}

	var weather strings.Builder
	now := time.Now()
	_, state := forecast.Properties.Weather.Recent(now)
	weather.WriteString(state.Icon())

	_, temp, unit, tempOK := forecast.Properties.Temperature.Recent(now)
	temp, unit = toFahrenheit(temp, unit)
	if tempOK {
		weather.WriteString(fmt.Sprintf(" %.f°%s", temp, unit))
	}
	return weather.String(), nil
}

func toFahrenheit(temperature float64, unit string) (float64, string) {
	if unit != "C" {
		return temperature, unit
	}
	return temperature*9/5 + 32, "F"
}

type Point struct {
	Properties PointProperties
}

type PointProperties struct {
	ForecastGridData string
}

type Forecast struct {
	Properties ForecastProperties
}

type ForecastProperties struct {
	RelativeHumidity    Measurements
	ApparentTemperature Measurements
	Temperature         Measurements
	Weather             StateMeasurements
	WindChill           Measurements
	WindDirection       Measurements
	WindSpeed           Measurements
}

type Measurements struct {
	UnitOfMeasure string `json:"uom"`
	Values        []Measurement
}

func (m Measurements) Recent(now time.Time) (t time.Time, value float64, unit string, ok bool) {
	if len(m.Values) == 0 {
		return time.Time{}, 0, "", false
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
	return mostRecent.ValidTime.Time, mostRecent.Value, unit, true
}

type Measurement struct {
	Value     float64
	ValidTime TimeAndDuration
}

type StateMeasurements struct {
	Values []StateMeasurement
}

type StateMeasurement struct {
	ValidTime TimeAndDuration
	Value     []StateValue
}

func (m StateMeasurements) Recent(now time.Time) (time.Time, State) {
	if len(m.Values) == 0 {
		return time.Time{}, stateUnknown
	}
	mostRecent := m.Values[0]
	mostRecentState := stateUnknown
	if len(mostRecent.Value) > 0 && mostRecent.Value[0].Weather != nil {
		mostRecentState = *mostRecent.Value[0].Weather
	}

	for _, measurement := range m.Values[1:] {
		if measurement.ValidTime.Time.Before(now) && measurement.ValidTime.Time.After(mostRecent.ValidTime.Time) {
			for _, value := range measurement.Value {
				if value.Weather != nil {
					mostRecent = measurement
					mostRecentState = *value.Weather
					break
				}
			}
		}
	}
	return mostRecent.ValidTime.Time, mostRecentState
}

type StateValue struct {
	Weather *State
}

type TimeAndDuration struct {
	Time time.Time
	// Original value contained a duration. Could unmarshal this in the future if needed.
}

func (t *TimeAndDuration) MarshalText() ([]byte, error) {
	b, err := t.Time.MarshalText()
	if err != nil {
		return nil, err
	}
	return append(b, []byte("/PT2H")...), nil
}

func (t *TimeAndDuration) UnmarshalText(bytes []byte) error {
	text := string(bytes)
	if before, _, found := strings.Cut(text, "/"); found {
		text = before
	}

	var err error
	t.Time, err = time.Parse(time.RFC3339, string(text))
	return err
}
