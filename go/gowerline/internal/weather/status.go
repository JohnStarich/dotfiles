package weather

import (
	"fmt"
	"net/url"
	"time"

	"github.com/johnstarich/go/gowerline/internal/icon"
	"github.com/johnstarich/go/gowerline/internal/status"
)

func Status(ctx status.Context) (time.Duration, error) {
	maxMindDBURL, err := url.Parse("https://download.db-ip.com")
	if err != nil {
		return 0, err
	}
	weatherGovAPIURL, err := url.Parse("https://api.weather.gov")
	if err != nil {
		return 0, err
	}
	return statusWithEndpoints(ctx, endpoints{
		MaxMindDBURL:     *maxMindDBURL,
		WeatherGovAPIURL: *weatherGovAPIURL,
	})
}

type endpoints struct {
	MaxMindDBURL     url.URL
	WeatherGovAPIURL url.URL
}

func statusWithEndpoints(ctx status.Context, endpoints endpoints) (time.Duration, error) {
	if !ctx.CacheExpired() {
		fmt.Fprint(ctx.Writer, ctx.Cache.Content)
		return 0, nil
	}

	latitude, longitude, err := getCachedCurrentLocation(ctx, endpoints.MaxMindDBURL, time.Now())
	if err != nil {
		fmt.Fprint(ctx.Writer, icon.Globe, ctx.Cache.Content)
		return 0, err
	}

	weather, err := getLatestWeather(ctx, endpoints.WeatherGovAPIURL, latitude, longitude)
	if err != nil {
		fmt.Fprint(ctx.Writer, icon.Internet, ctx.Cache.Content)
		return 0, err
	}
	fmt.Fprint(ctx.Writer, weather)
	return 30 * time.Minute, nil
}
