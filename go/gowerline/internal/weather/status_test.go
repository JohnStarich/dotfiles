package weather

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/johnstarich/go/gowerline/internal/dnsresolver"
	"github.com/johnstarich/go/gowerline/internal/icon"
	"github.com/johnstarich/go/gowerline/internal/status"
)

func TestStatus(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Unrecognized HTTP request:", r.Method, r.URL.String())
	})

	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	testIPAddress, geoIPHandler := newGeoIPHandler(t, now)
	const maxMindDBPrefix = "/maxminddb"
	mux.Handle(maxMindDBPrefix+"/", http.StripPrefix("/maxminddb", geoIPHandler))

	const someTemperature = 20
	oneHourAgo := now.Add(-1 * time.Hour)
	forecast := Forecast{
		Properties: ForecastProperties{
			Temperature: Measurements{
				UnitOfMeasure: "wmoUnit:degF",
				Values: []Measurement{
					{
						ValidTime: TimeAndDuration{Time: oneHourAgo},
						Value:     someTemperature,
					},
				},
			},
			Weather: StateMeasurements{
				Values: []StateMeasurement{
					{
						ValidTime: TimeAndDuration{Time: oneHourAgo},
						Value: []StateValue{
							{Weather: stateThunderstorms.toPointer()},
						},
					},
				},
			},
		},
	}
	const weatherGovPrefix = "/weathergov"
	mux.Handle(weatherGovPrefix+"/", http.StripPrefix(weatherGovPrefix, newWeatherGovHandler(t, weatherGovPrefix, forecast)))

	testCtx := status.NewTestContext(t, status.TestConfig{
		Handler: mux,
		Now:     now,
		ResolvedIPAddresses: []dnsresolver.TestIP{
			{
				ResolverHostPort: "resolver1.opendns.com:53",
				Hostname:         "myip.opendns.com",
				IP:               testIPAddress,
			},
		},
	})
	cacheDuration, err := statusWithEndpoints(testCtx.Context, endpoints{
		MaxMindDBURL:     url.URL{Path: maxMindDBPrefix},
		WeatherGovAPIURL: url.URL{Path: weatherGovPrefix},
	})
	if err != nil {
		t.Fatal(err)
	}
	if cacheDuration != 30*time.Minute {
		t.Error("Expected 30m cache duration, got:", cacheDuration)
	}
	const expectedOutput = icon.StormCloud + " 20°F"
	if output := testCtx.Output(); output != expectedOutput {
		t.Errorf("Expected output to be %q, got: %q", expectedOutput, output)
	}
}
