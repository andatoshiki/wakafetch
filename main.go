package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	if *config.updateFlag {
		runUpdateCheck(*config.timeoutFlag)
		return
	}
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
			// Heatmap default: last 12 months when no explicit --range. Otherwise respect --range (7d, 30d, etc.).
			heatmapRange := *config.rangeFlag
			if heatmapRange == "yesterday" {
				heatmapRange = "1y"
			}
			startDate, endDate, head, valid := getSummaryRange(heatmapRange)
			if !valid {
				ui.Errorln("Invalid range for heatmap: use today, yesterday, 7d, 30d, 6m, 1y, or a year (e.g. 2024)")
				return
			}
			heading = head
			data, err = fetchSummaryWithDates(apiKey, apiURL, startDate, endDate, *config.timeoutFlag)
			if err != nil {
				ui.Errorln(err.Error())
				return
			}
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

// getSummaryRange returns (startDate, endDate, heading, valid) for summary/heatmap by range flag.
// Used when --heatmap is set and range is not a year.
func getSummaryRange(rangeFlag string) (startDate, endDate, heading string, valid bool) {
	today := time.Now()
	switch rangeFlag {
	case "today":
		d := today.Format("2006-01-02")
		return d, d, "Today", true
	case "yesterday":
		d := today.AddDate(0, 0, -1).Format("2006-01-02")
		return d, d, "Yesterday", true
	case "7d":
		end := today.Format("2006-01-02")
		start := today.AddDate(0, 0, -6).Format("2006-01-02")
		return start, end, "Last 7 days", true
	case "30d":
		end := today.Format("2006-01-02")
		start := today.AddDate(0, 0, -29).Format("2006-01-02")
		return start, end, "Last 30 days", true
	case "6m":
		end := today.Format("2006-01-02")
		start := today.AddDate(0, -6, 0).Format("2006-01-02")
		return start, end, "Last 6 months", true
	case "1y", "all":
		// For heatmap, "all" is capped at last 12 months
		end := today.Format("2006-01-02")
		start := today.AddDate(-1, 0, 1).Format("2006-01-02")
		if rangeFlag == "all" {
			return start, end, "Last 12 months", true
		}
		return start, end, "Last 12 months", true
	default:
		return "", "", "", false
	}
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

const updateCheckURL = "https://api.github.com/repos/andatoshiki/wakafetch/releases/latest"
const installScriptURL = "https://raw.githubusercontent.com/andatoshiki/wakafetch/master/scripts/install.sh"

func runUpdateCheck(timeoutSeconds int) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 10
	}
	client := &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	req, err := http.NewRequest("GET", updateCheckURL, nil)
	if err != nil {
		ui.Errorln("Update check failed: %s", err.Error())
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "wakafetch")
	resp, err := client.Do(req)
	if err != nil {
		ui.Errorln("Update check failed: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ui.Errorln("Update check failed: HTTP %s", resp.Status)
		return
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		ui.Errorln("Update check failed: %s", err.Error())
		return
	}
	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(Version, "v")
	if versionGreater(latest, current) {
		fmt.Println(ui.Clr.Green + "New version " + release.TagName + " available." + ui.Clr.Reset)
		fmt.Println("Install: curl -fsSL " + installScriptURL + " | sh")
	} else {
		fmt.Println(ui.Clr.Green + "Already up to date." + ui.Clr.Reset + " (wakafetch @" + Version + ")")
	}
}

// versionGreater returns true if a > b (semver-style: major.minor.patch).
func versionGreater(a, b string) bool {
	parse := func(s string) (major, minor, patch int) {
		parts := strings.Split(s, ".")
		if len(parts) >= 1 {
			major, _ = strconv.Atoi(parts[0])
		}
		if len(parts) >= 2 {
			minor, _ = strconv.Atoi(parts[1])
		}
		if len(parts) >= 3 {
			patch, _ = strconv.Atoi(strings.TrimSuffix(parts[2], "-next"))
		}
		return major, minor, patch
	}
	ma, mi, pa := parse(a)
	mb, mj, pb := parse(b)
	if ma != mb {
		return ma > mb
	}
	if mi != mj {
		return mi > mj
	}
	return pa > pb
}
