package weather

import (
	"fmt"
	"time"

	"github.com/johnstarich/go/gowerline/internal/icon"
	"github.com/johnstarich/go/gowerline/internal/status"
)

func Status(ctx status.Context) (time.Duration, error) {
	if !ctx.CacheExpired() {
		fmt.Fprint(ctx.Writer, ctx.Cache.Content)
		return 0, nil
	}

	latitude, longitude, err := getCachedCurrentLocation(ctx, time.Now())
	if err != nil {
		fmt.Fprint(ctx.Writer, icon.Globe, icon.Warning, ctx.Cache.Content)
		return 0, err
	}

	if err := getLatestWeather(ctx, latitude, longitude); err != nil {
		fmt.Fprint(ctx.Writer, icon.Internet, icon.Warning, ctx.Cache.Content)
		return 0, err
	}
	return 30 * time.Minute, nil
}
