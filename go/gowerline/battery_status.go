package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const powerSupplyFSPrefix = "/sys/class/power_supply/"

func batteryStatus(ctx context.Context, w io.Writer) error {
	batteryDirectories, err := findBatteryDirectories(ctx)
	if err != nil {
		return err
	}
	if len(batteryDirectories) == 0 {
		return errors.New("no battery detected")
	}

	for index, batteryDir := range batteryDirectories {
		if index > 0 {
			fmt.Fprint(w, " ")
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
		fmt.Fprintf(w, "%s %.0f%%", batterySummaryForStatus(string(statusBytes)), chargePercent)
	}
	return nil
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
