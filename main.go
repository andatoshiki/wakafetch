package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/andatoshiki/wakafetch/wakafetch/types"
	"github.com/andatoshiki/wakafetch/wakafetch/ui"
)

// okay so there are 2 types of responses/endpoints:
// /summary response:
// - gives summary of EACH day, so more granular data
// - have to aggregate data manually for viewing stats
// - supports custom date ranges
// - so, its called when --days(custom range) or --daily/heatmap(granular daily breakdown) is used

// /stats response:
// - gives summary of the ENTIRE range in a single response
// - no need for aggregation, efficient af
// - doesn't support custom date ranges(only rangeStr)
// - so, its the default unless --days or --daily/heatmap is used

func main() {
	config := parseFlags()
	apiURL, apiKey := loadAPIConfig(config)

	if shouldUseSummaryAPI(config) {
		handleSummaryFlow(config, apiKey, apiURL)
	} else {
		handleStatsFlow(config, apiKey, apiURL)
	}
}

func loadAPIConfig(config Config) (string, string) {
	apiURL, apiKey, err := parseConfig()
	if err != nil {
		ui.Errorln(err.Error())
	}

	if *config.apiKeyFlag != "" {
		apiKey = *config.apiKeyFlag
	}

	return apiURL, apiKey
}

func shouldUseSummaryAPI(config Config) bool {
	// Use summary API if: days flag is set, daily/heatmap flags are set, or range is a year
	_, isYear := parseYear(*config.rangeFlag)
	return *config.daysFlag != 0 || *config.dailyFlag || *config.heatmapFlag || isYear
}

func handleStatsFlow(config Config, apiKey, apiURL string) {
	rangeStr := getRangeStr(*config.rangeFlag)

	data, err := fetchStats(apiKey, apiURL, rangeStr, *config.timeoutFlag)
	if err != nil {
		ui.Errorln(err.Error())
	}

	if *config.jsonFlag {
		outputJSON(data)
		return
	}

	ui.DisplayStats(data, *config.fullFlag, rangeStr)
}

func handleSummaryFlow(config Config, apiKey, apiURL string) {
	var data *types.SummaryResponse
	var err error
	var heading string

	// Check if range flag is a year number
	year, isYear := parseYear(*config.rangeFlag)
	if isYear {
		// Fetch data for the entire year
		startDate := fmt.Sprintf("%d-01-01", year)
		endDate := fmt.Sprintf("%d-12-31", year)
		data, err = fetchSummaryWithDates(apiKey, apiURL, startDate, endDate, *config.timeoutFlag)
		if err != nil {
			ui.Errorln(err.Error())
			return
		}
		heading = fmt.Sprintf("Year %d", year)
	} else {
		// Non-year ranges
		if *config.heatmapFlag {
			// For heatmap without a specific year, always show the most recent 12 months.
			today := time.Now()
			endDate := today.Format("2006-01-02")
			startDate := today.AddDate(-1, 0, 1).Format("2006-01-02")

			data, err = fetchSummaryWithDates(apiKey, apiURL, startDate, endDate, *config.timeoutFlag)
			if err != nil {
				ui.Errorln(err.Error())
				return
			}
			heading = "Last 12 months"
		} else {
			// Use existing logic for preset ranges
			rangeStr := getRangeStr(*config.rangeFlag)
			days := *config.daysFlag
			validRange := true
			if days == 0 {
				days, validRange = map[string]int{
					"today":         1,
					"last_7_days":   7,
					"last_30_days":  30,
					"last_6_months": 183,
					"last_year":     365,
				}[rangeStr]
			}

			if !validRange {
				ui.Errorln("This range isn't supported with `--daily` or `--heatmap` flags. Use `--days` instead")
				return
			}

			data, err = fetchSummary(apiKey, apiURL, days, *config.timeoutFlag)
			if err != nil {
				ui.Errorln(err.Error())
				return
			}

			if *config.daysFlag != 0 {
				if days == 1 {
					heading = "Today"
				} else {
					heading = fmt.Sprintf("Last %d days", days)
				}
			} else {
				headingMap := map[string]string{
					"today":         "Today",
					"last_7_days":   "Last 7 days",
					"last_30_days":  "Last 30 days",
					"last_6_months": "Last 6 months",
					"last_year":     "Last year",
					"all_time":      "All time",
				}
				heading = headingMap[rangeStr]
			}
		}
	}

	if *config.jsonFlag {
		outputJSON(data)
		return
	}

	if *config.dailyFlag {
		ui.DisplayBreakdown(data.Data, heading)
		return
	}

	if *config.heatmapFlag {
		ui.DisplayHeatmap(data.Data, heading)
		return
	}

	ui.DisplaySummary(data, *config.fullFlag, heading)
}

func parseYear(rangeFlag string) (int, bool) {
	// Try to parse as a year number (4 digits, typically 2000-2099)
	year, err := strconv.Atoi(rangeFlag)
	if err != nil {
		return 0, false
	}
	// Validate year is reasonable (1900-2100)
	if year >= 1900 && year <= 2100 {
		return year, true
	}
	return 0, false
}

func getRangeStr(rangeFlag string) string {
	// Check if it's a year first
	_, isYear := parseYear(rangeFlag)
	if isYear {
		// Year is handled separately in handleSummaryFlow
		return ""
	}

	rangeStrMap := map[string]string{
		"today":     "today",
		"yesterday": "yesterday",
		"7d":        "last_7_days",
		"30d":       "last_30_days",
		"6m":        "last_6_months",
		"1y":        "last_year",
		"all":       "all_time",
	}
	rangeStr, exists := rangeStrMap[rangeFlag]
	if !exists {
		ui.Errorln("Invalid range: '%s', must be one of %stoday, 7d, 30d, 6m, 1y, all, or a year (e.g., 2023, 2024)", rangeFlag, ui.Clr.Green)
	}
	return rangeStr
}

func colorsShouldBeEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}

	// tty check
	file, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (file.Mode() & os.ModeCharDevice) != 0
}

func outputJSON(data any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		ui.Errorln("Failed to marshal JSON: %s", err.Error())
	}
	fmt.Println(string(jsonData))
}
