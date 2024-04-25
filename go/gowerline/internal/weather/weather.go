package weather

import (
	"fmt"
	"strings"
	"time"

	"github.com/johnstarich/go/gowerline/internal/status"
)

func getLatestWeather(ctx status.Context, latitude, longitude float64) error {
	var weather Point
	err := doJSONGet(ctx.Context, fmt.Sprintf("https://api.weather.gov/points/%f,%f", latitude, longitude), &weather)
	if err != nil {
		return err
	}
	var forecast Forecast
	err = doJSONGet(ctx.Context, weather.Properties.ForecastGridData, &forecast)
	if err != nil {
		return err
	}

	now := time.Now()
	_, temp, unit := forecast.Properties.Temperature.Recent(now)
	temp, unit = toFahrenheit(temp, unit)
	_, w := forecast.Properties.Weather.Recent(now)

	fmt.Fprintf(ctx.Writer, "%s  %.f°%s ", w.Icon(), temp, unit)
	return nil
}

func toFahrenheit(temperature float64, unit string) (float64, string) {
	if unit != "C" {
		return temperature, unit
	}
	return temperature*9/5 + 32, "F"
}

type Point struct {
	Properties struct {
		ForecastGridData string
	}
}

type Forecast struct {
	Properties struct {
		RelativeHumidity    Measurements
		ApparentTemperature Measurements
		Temperature         Measurements
		Weather             StateMeasurements
		WindChill           Measurements
		WindDirection       Measurements
		WindSpeed           Measurements
	}
}

type Measurements struct {
	UnitOfMeasure string `json:"uom"`
	Values        []Measurement
}

func (m Measurements) Recent(now time.Time) (t time.Time, value float64, unit string) {
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
	var mostRecentState State
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

func (t *TimeAndDuration) UnmarshalText(bytes []byte) error {
	text := string(bytes)
	if before, _, found := strings.Cut(text, "/"); found {
		text = before
	}

	var err error
	t.Time, err = time.Parse(time.RFC3339, string(text))
	return err
}
