package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andatoshiki/wakafetch/wakafetch/types"
)

func formatDateRange(start, end string) string {
	if start == "" || end == "" {
		return ""
	}
	startParts := strings.Split(start, "T")
	endParts := strings.Split(end, "T")
	if len(startParts) < 1 || len(endParts) < 1 {
		return ""
	}
	startDate := startParts[0]
	endDate := endParts[0]

	const layout = "2006-01-02"
	const outLayout = "Jan 2"

	startTime, err1 := time.Parse(layout, startDate)
	endTime, err2 := time.Parse(layout, endDate)
	if err1 != nil || err2 != nil {
		return ""
	}

	startStr := startTime.Format(outLayout)
	endStr := endTime.Format(outLayout)

	if startDate == endDate {
		return startStr
	}
	return fmt.Sprintf("%s to %s", startStr, endStr)
}

func formatBestDay(dateStr string) string {
	if dateStr == "" {
		return "None"
	}
	dateParts := strings.Split(dateStr, "T")
	if len(dateParts) < 1 {
		return dateStr
	}

	const layout = "2006-01-02"
	const outLayout = "January 2, 2006"
	dateTime, err := time.Parse(layout, dateParts[0])
	if err != nil {
		return dateStr
	}
	return dateTime.Format(outLayout)
}

func formatDailyDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	dateParts := strings.Split(dateStr, "T")
	if len(dateParts) < 1 {
		return dateStr
	}

	const layout = "2006-01-02"
	const outLayout = "Jan 2"
	dateTime, err := time.Parse(layout, dateParts[0])
	if err != nil {
		return dateStr
	}
	return dateTime.Format(outLayout)
}

func timeFmt(seconds float64) string {
	sec := int(seconds)
	if sec < 3600 {
		return fmt.Sprintf("%dm %ds", sec/60, sec%60)
	}
	hours := sec / 3600
	minutes := (sec % 3600) / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func timeFmtPad(seconds, maxSeconds float64) string {
	sec := int(seconds)
	pad := 2
	if maxSeconds > 360000 {
		pad = 3
	}
	if sec < 3600 {
		return fmt.Sprintf("%*dm %2ds", pad, sec/60, sec%60)
	}
	hours := sec / 3600
	minutes := (sec % 3600) / 60
	return fmt.Sprintf("%*dh %2dm", pad, hours, minutes)
}

func topItemName(items []types.StatItem, skipUnknown bool) string {
	if len(items) == 0 {
		return "None"
	}
	top := items[0].Name
	if skipUnknown && top == "unknown" {
		if len(items) > 1 {
			top = items[1].Name
		}
	}
	return top
}

func Errorln(format string, args ...any) {
	fmt.Fprintf(os.Stderr, Clr.Red+format+Clr.Reset+"\n", args...)
	os.Exit(1)
}

func Warnln(format string, args ...any) {
	fmt.Fprintf(os.Stderr, Clr.Yellow+format+Clr.Reset+"\n", args...)
}
