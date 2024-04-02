package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func batteryStatus(_ context.Context, w io.Writer) error {
	chargeNowBytes, err := os.ReadFile("/sys/class/power_supply/BAT1/charge_now")
	if err != nil {
		return err
	}
	totalChargeBytes, err := os.ReadFile("/sys/class/power_supply/BAT1/charge_full")
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
	statusBytes, err := os.ReadFile("/sys/class/power_supply/BAT1/status")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s %.0f%%", batterySummaryForStatus(string(statusBytes)), chargeNow/totalCharge*100)
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
