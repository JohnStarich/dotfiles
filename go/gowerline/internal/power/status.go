package power

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
	err := writeBatteryStatus(ctx)
	if err != nil {
		fmt.Fprint(ctx.Writer, icon.Warning, ctx.Cache.Content)
		return 0, nil
	}
	return 5 * time.Minute, nil
}
