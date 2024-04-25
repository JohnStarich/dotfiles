package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/johnstarich/go/gowerline/internal/icon"
	"github.com/johnstarich/go/gowerline/internal/status"
)

const powerSupplyFSPrefix = "/sys/class/power_supply/"

func batteryStatus(ctx status.Context) (time.Duration, error) {
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

func writeBatteryStatus(ctx status.Context) error {
	batteryDirectories, err := findBatteryDirectories(ctx.Context)
	if err != nil {
		return err
	}
	if len(batteryDirectories) == 0 {
		return errors.New("no battery detected")
	}

	for index, batteryDir := range batteryDirectories {
		if index > 0 {
			fmt.Fprint(ctx.Writer, " ")
		}
		chargeNowBytes, err := os.ReadFile(batteryDir + "/charge_now")
		if err != nil {
			return err
		}
		totalChargeBytes, err := os.ReadFile(batteryDir + "/charge_full_design")
		if err != nil {
			return err
		}
		chargeNow, err := strconv.ParseFloat(strings.TrimSpace(string(chargeNowBytes)), 64)
		if err != nil {
			return err
		}
		totalCharge, err := strconv.ParseFloat(strings.TrimSpace(string(totalChargeBytes)), 64)
		if err != nil {
			return err
		}
		statusBytes, err := os.ReadFile(batteryDir + "/status")
		if err != nil {
			return err
		}
		chargePercent := chargeNow / totalCharge * 100
		if chargePercent > 100 {
			chargePercent = 100
		}
		fmt.Fprintf(ctx.Writer, "%s %.0f%%", batterySummaryForStatus(string(statusBytes)), chargePercent)
	}
	return nil
}

func batterySummaryForStatus(rawLinuxBatteryStatus string) string {
	linuxBatteryStatus := strings.ToLower(strings.TrimSpace(rawLinuxBatteryStatus))
	switch linuxBatteryStatus {
	// Battery statuses are briefly described here: https://github.com/torvalds/linux/blob/026e680b0a08a62b1d948e5a8ca78700bfac0e6e/Documentation/power/power_supply_class.rst
	// And might use similar names from here: https://github.com/torvalds/linux/blob/026e680b0a08a62b1d948e5a8ca78700bfac0e6e/drivers/acpi/battery.c#L41-L43
	case "discharging":
		return icon.Fire
	case "charging":
		return icon.LightningBolt
	case "critical":
		return icon.Warning
	default:
		return icon.Warning + " " + linuxBatteryStatus
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
