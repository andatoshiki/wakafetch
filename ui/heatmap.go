package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/andatoshiki/wakafetch/wakafetch/types"
)

func heatmap(days []types.DayData) ([]string, int) {
	const heatmapChar = "■"               // █ ❐ ▪ ◼ 🙩 🙫 ⛝ ⏹ 🞕 🞔 🞖
	const highlight = "\x1b[38;2;0;%v;0m" // \x1b[38;2;R;G;Bm
	const numLevels = 6                   // 6 intensity levels matching GitHub style

	// Green RGB values for each intensity level (0-5)
	// Level 0: no activity (black), Level 5: maximum activity (light green)
	greenLevels := [numLevels]int{0, 30, 60, 100, 140, 220}

	if len(days) == 0 {
		return []string{}, 0
	}
	startDay, err := time.Parse("2006-01-02", strings.Split(days[0].Range.Start, "T")[0])
	if err != nil {
		return []string{}, 0
	}
	endDay, err := time.Parse("2006-01-02", strings.Split(days[len(days)-1].Range.Start, "T")[0])
	if err != nil {
		return []string{}, 0
	}

	height := 4
	numOfDays := int(endDay.Sub(startDay).Hours()/24) + 1
	cols := getTerminalCols()
	width := 2*int((numOfDays+height-1)/height) - 1 // 2*ciel -1
	for width+4 > cols {
		height++
		width = 2*int((numOfDays+height-1)/height) - 1
	}

	output := make([]string, height)
	dataIndex := 0
	i := 0
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		level := 0 // Default to level 0 (no activity)
		if dataIndex < len(days) && strings.Split(days[dataIndex].Range.Start, "T")[0] == d.Format("2006-01-02") {
			day := days[dataIndex]
			dataIndex++
			// Convert seconds to hours
			hours := day.GrandTotal.TotalSeconds / 3600.0
			// Map hours to intensity level (absolute thresholds)
			// Level 0: 0 hours (black)
			// Level 1: 0-2 hours
			// Level 2: 2-4 hours
			// Level 3: 4-6 hours
			// Level 4: 6-8 hours
			// Level 5: 8+ hours
			if hours >= 8 {
				level = 5
			} else if hours >= 6 {
				level = 4
			} else if hours >= 4 {
				level = 3
			} else if hours >= 2 {
				level = 2
			} else if hours > 0 {
				level = 1
			}
			// hours == 0 remains level 0 (black)
		}
		greenValue := greenLevels[level]
		char := fmt.Sprintf(highlight, greenValue) + heatmapChar + "\x1b[0m"
		output[i%height] += char + " "
		i++
	}
	// trim trailing spaces
	for j := range output {
		output[j] = strings.TrimRight(output[j], " ")
	}

	// ensure same width for all lines
	for i%height != 0 {
		if i/height < 1 {
			output[i%height] += " "
		} else {
			output[i%height] += "  "
		}
		i++
	}
	return output, width
}
