package weather

import (
	"encoding/json"
	"net/http"
	"path"
	"testing"
)

func newWeatherGovHandler(tb testing.TB, prefix string, forecast Forecast) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tb.Fatal("Unexpected Weather.gov call:", r.URL.Path)
	})
	mux.HandleFunc("/points/", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(Point{
			Properties: PointProperties{
				ForecastGridData: path.Join(prefix, "forecast"),
			},
		})
		if err != nil {
			tb.Fatal(err)
		}
	})
	mux.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(forecast)
		if err != nil {
			tb.Fatal(err)
		}
	})
	return mux
}
