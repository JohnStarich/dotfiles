package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/johnstarich/go/gowerline/internal/status"
)

const powerSupplyFSPrefix = "/sys/class/power_supply/"

func batteryStatus(ctx status.Context) (time.Duration, error) {
	if !ctx.CacheExpired() {
		fmt.Fprint(ctx.Writer, ctx.Cache.Content)
		return ctx.CacheDuration(), nil
	}

	batteryDirectories, err := findBatteryDirectories(ctx.Context)
	if err != nil {
		return 0, err
	}
	if len(batteryDirectories) == 0 {
		return 0, errors.New("no battery detected")
	}

	for index, batteryDir := range batteryDirectories {
		if index > 0 {
			fmt.Fprint(ctx.Writer, " ")
		}
		chargeNowBytes, err := os.ReadFile(batteryDir + "/charge_now")
		if err != nil {
			return 0, err
		}
		totalChargeBytes, err := os.ReadFile(batteryDir + "/charge_full_design")
		if err != nil {
			return 0, err
		}
		chargeNow, err := strconv.ParseFloat(strings.TrimSpace(string(chargeNowBytes)), 64)
		if err != nil {
			return 0, err
		}
		totalCharge, err := strconv.ParseFloat(strings.TrimSpace(string(totalChargeBytes)), 64)
		if err != nil {
			return 0, err
		}
		statusBytes, err := os.ReadFile(batteryDir + "/status")
		if err != nil {
			return 0, err
		}
		chargePercent := chargeNow / totalCharge * 100
		if chargePercent > 100 {
			chargePercent = 100
		}
		fmt.Fprintf(ctx.Writer, "%s %.0f%%", batterySummaryForStatus(string(statusBytes)), chargePercent)
	}
	return 1 * time.Minute, nil
}

func batterySummaryForStatus(rawLinuxBatteryStatus string) string {
	linuxBatteryStatus := strings.ToLower(strings.TrimSpace(rawLinuxBatteryStatus))
	const warning = "⚠️"
	switch linuxBatteryStatus {
	// Battery statuses are briefly described here: https://github.com/torvalds/linux/blob/026e680b0a08a62b1d948e5a8ca78700bfac0e6e/Documentation/power/power_supply_class.rst
	// And might use similar names from here: https://github.com/torvalds/linux/blob/026e680b0a08a62b1d948e5a8ca78700bfac0e6e/drivers/acpi/battery.c#L41-L43
	case "discharging":
		return "🔥"
	case "charging":
		return "⚡"
	case "critical":
		return warning
	default:
		return warning + " " + linuxBatteryStatus
	}
}

func findBatteryDirectories(context.Context) ([]string, error) {
	powerDirs, err := os.ReadDir(powerSupplyFSPrefix)
	if err != nil {
		return nil, err
	}
	var batteries []string
	for _, dir := range powerDirs {
		if strings.HasPrefix(dir.Name(), "BAT") {
			batteries = append(batteries, powerSupplyFSPrefix+dir.Name())
		}
	}
	return batteries, nil
}
