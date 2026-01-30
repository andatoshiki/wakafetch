package ui

import (
	"fmt"
	"strings"

	"github.com/andatoshiki/wakafetch/wakafetch/types"
)

func graphStr(items []types.StatItem, limit int) ([]string, int) {
	if len(items) == 0 {
		return []string{}, 0
	}
	count := len(items)
	if limit > 0 && limit < count {
		count = limit
	}
	visibleItems := items[:count]

	if len(visibleItems) == 0 || visibleItems[0].TotalSeconds == 0 {
		return []string{}, 0
	}

	output := make([]string, 0, len(visibleItems))

	maxNameLength := 0
	maxSeconds := 0.0
	for _, item := range visibleItems {
		maxNameLength = max(maxNameLength, len(item.Name))
		maxSeconds = max(maxSeconds, item.TotalSeconds)
	}
	for _, item := range visibleItems {
		if item.TotalSeconds < 60 {
			continue
		}
		barLength := int((item.TotalSeconds / visibleItems[0].TotalSeconds) * float64(barWidth))
		secondBarLength := barWidth - barLength
		if barLength < 1 {
			barLength = 1
			secondBarLength = barWidth - 1
		}
		bar := strings.Repeat(barChar, barLength)
		secondBar := strings.Repeat(barChar, secondBarLength)
		label := fmt.Sprintf("%-*s ", maxNameLength, item.Name)

		if Clr.Gray == "" {
			secondBar = strings.Repeat(" ", secondBarLength)
		}
		line := label +
			Clr.Green + bar + Clr.Reset +
			Clr.Gray + secondBar + Clr.Reset + " " +
			Clr.Green + timeFmtPad(item.TotalSeconds, maxSeconds) + Clr.Reset
		output = append(output, line)
	}
	graphWidth := maxNameLength + 1 + barWidth + 1 + len(timeFmtPad(maxSeconds, maxSeconds))
	return output, graphWidth
}
