//go:build darwin

package power

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/johnstarich/go/gowerline/internal/icon"
	"github.com/johnstarich/go/gowerline/internal/status"
	"github.com/johnstarich/go/regext"
	"github.com/pkg/errors"
)

var (
	powerRegex = regext.MustCompile(`
    (?P<percentage> \d+)%; \s
    (?P<status> [^;]+); \s
    (?P<time> (?: \d+:\d+ \s)? .+) \s
    present .*
`)
	powerPercentageSubexp = powerRegex.SubexpIndex("percentage")
	powerStatusSubexp     = powerRegex.SubexpIndex("status")
	powerTimeSubexp       = powerRegex.SubexpIndex("time")
)

func writeBatteryStatus(ctx status.Context) error {
	cmd := exec.CommandContext(ctx.Context, "pmset", "-g", "ps")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.Errorf("failed to run pmset: %w %q", err, stderr.String())
	}

	matches := powerRegex.FindStringSubmatch(stdout.String())
	chargePercentage, err := strconv.Atoi(matches[powerPercentageSubexp])
	if err != nil {
		return errors.Errorf("failed to parse battery percentage: %q %w", matches[powerPercentageSubexp], err)
	}
	status, timeRemaining := batterySummaryForStatus(matches[powerStatusSubexp], matches[powerTimeSubexp])
	fmt.Fprintf(ctx.Writer, "%s %d%% (%s)", status, chargePercentage, timeRemaining)
	return nil
}

func batterySummaryForStatus(rawDarwinBatteryStatus, rawDarwinTimeRemaining string) (status, timeRemaining string) {
	timeRemaining = strings.ReplaceAll(rawDarwinTimeRemaining, " remaining", "")
	if rawDarwinBatteryStatus == "charged" {
		timeRemaining = "full"
	} else if timeRemaining == "0:00" || timeRemaining == "(no estimate)" || timeRemaining == "not charging" {
		timeRemaining = "-:--"
	}
	status = icon.Warning
	switch rawDarwinBatteryStatus {
	case "discharging":
		status = icon.Fire
	case "charging":
		status = icon.LightningBolt
	case "finishing charge":
		status = icon.FullBattery
	case "charged":
		status = icon.FullBattery
	case "AC attached":
		status = icon.Plug
	}
	return
}
