package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/andatoshiki/wakafetch/wakafetch/types"
)

const maxTableBarWidth = 10

func dailyBreakdownStr(dailyData []types.DayData) ([]string, int) {
	if len(dailyData) == 0 {
		return []string{}, 0
	}

	output := make([]string, 0, len(dailyData)+2)

	// sort descending by date
	sortedDays := make([]types.DayData, len(dailyData))
	copy(sortedDays, dailyData)

	sort.Slice(sortedDays, func(i, j int) bool {
		dateI := strings.Split(sortedDays[i].Range.Start, "T")[0]
		dateJ := strings.Split(sortedDays[j].Range.Start, "T")[0]
		return dateI > dateJ
	})

	maxSecs := findMaxDailySeconds(sortedDays)
	cols := calculateDailyColumnWidths(sortedDays)

	output = append(output, dailyHeadersStr(cols))
	output = append(output, dailySeparatorStr(cols))

	dailyRows := dailyRowsStr(sortedDays, cols, maxSecs)
	output = append(output, dailyRows...)

	return output, cols.total
}

func dailyHeadersStr(cols dailyColumns) string {
	headerDate := fmt.Sprintf("%-*s", cols.date, "Date")
	headerTotal := fmt.Sprintf("%-*s", cols.time, "Time")
	headerLang := fmt.Sprintf("%-*s", cols.lang, "Language")
	headerProj := fmt.Sprintf("%-*s", cols.project, "Project")

	return Clr.Blue + headerDate + Clr.Reset + " │ " + Clr.Blue + headerTotal + Clr.Reset + " │ " + Clr.Blue + headerLang + Clr.Reset + " │ " + Clr.Blue + headerProj + Clr.Reset
}

func dailySeparatorStr(cols dailyColumns) string {
	return strings.Repeat("─", cols.date) + "─┼─" +
		strings.Repeat("─", cols.time) + "─┼─" +
		strings.Repeat("─", cols.lang) + "─┼─" +
		strings.Repeat("─", cols.project)
}

func dailyRowsStr(dailyData []types.DayData, cols dailyColumns, maxSecs float64) []string {
	output := make([]string, 0, len(dailyData))

	for _, day := range dailyData {
		if day.GrandTotal.TotalSeconds < 60 {
			continue
		}

		date := fmt.Sprintf("%-*s", cols.date, formatDailyDate(day.Range.Start))

		barLength := int((day.GrandTotal.TotalSeconds / maxSecs) * maxTableBarWidth)
		if barLength < 1 && day.GrandTotal.TotalSeconds > 0 {
			barLength = 1
		}
		bar := strings.Repeat(barChar, barLength)
		timeStr := timeFmtPad(day.GrandTotal.TotalSeconds, maxSecs)
		timeWithBar := timeStr + " " + bar
		totalFormatted := fmt.Sprintf("%-*s", cols.time, timeWithBar)

		topLang := fmt.Sprintf("%-*s", cols.lang, topItemName(day.Languages, false))
		topProj := fmt.Sprintf("%-*s", cols.project, topItemName(day.Projects, true))

		row := date + " │ " + Clr.Green + totalFormatted + Clr.Reset + " │ " + topLang + " │ " + topProj
		output = append(output, row)
	}

	return output
}

func findMaxDailySeconds(dailyData []types.DayData) float64 {
	maxSecs := 0.0
	for _, day := range dailyData {
		if day.GrandTotal.TotalSeconds > maxSecs {
			maxSecs = day.GrandTotal.TotalSeconds
		}
	}
	return maxSecs
}

type dailyColumns struct {
	date    int
	time    int
	lang    int
	project int
	total   int
}

func calculateDailyColumnWidths(dailyData []types.DayData) dailyColumns {
	cols := dailyColumns{
		date:    4, // "Date"
		time:    4, // "Time"
		lang:    8, // "Language"
		project: 7, // "Project"
		total:   0,
	}

	maxSecs := findMaxDailySeconds(dailyData)

	for _, day := range dailyData {
		if day.GrandTotal.TotalSeconds < 60 {
			continue
		}

		dateStr := formatDailyDate(day.Range.Start)
		totalStr := timeFmtPad(day.GrandTotal.TotalSeconds, maxSecs)
		topLang := topItemName(day.Languages, false)
		topProj := topItemName(day.Projects, true)

		if len(dateStr) > cols.date {
			cols.date = len(dateStr)
		}
		if len(totalStr) > cols.time {
			cols.time = len(totalStr)
		}
		if len(topLang) > cols.lang {
			cols.lang = len(topLang)
		}
		if len(topProj) > cols.project {
			cols.project = len(topProj)
		}
	}

	// spaces for bar
	cols.time += maxTableBarWidth + 1
	cols.total = cols.date + 3 + cols.time + 3 + cols.lang + 3 + cols.project
	return cols
}
