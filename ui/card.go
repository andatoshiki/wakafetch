package ui

import "strings"

func cardify(content []string, header string, contentWidth int, rightPad int) ([]string, int) {
	var (
		borderTopLeft     = Clr.MidGray + "╭" + Clr.Reset
		borderTopRight    = Clr.MidGray + "╮" + Clr.Reset
		borderBottomLeft  = Clr.MidGray + "╰" + Clr.Reset
		borderBottomRight = Clr.MidGray + "╯" + Clr.Reset
		borderHorizontal  = Clr.MidGray + "─" + Clr.Reset
		borderVertical    = Clr.MidGray + "│" + Clr.Reset
	)

	if len(content) == 0 {
		return []string{}, 0
	}

	// Base card width is content width plus padding and borders.
	// rightPad is used to make narrower cards match the widest card in a section.
	cardWidth := contentWidth + rightPad + 4

	// If header is longer than content area, expand card to fit header.
	if len(header)+2 > cardWidth { // +2 for corners
		cardWidth = len(header) + 2
	}

	availableSpace := max(0, cardWidth-len(header)-2) // -2 for corner chars

	leftPadding := availableSpace / 2
	rightPadding := availableSpace - leftPadding

	headerLine := borderTopLeft +
		strings.Repeat(borderHorizontal, leftPadding) +
		Clr.Bold + Clr.Yellow + header + Clr.Reset +
		strings.Repeat(borderHorizontal, rightPadding) +
		borderTopRight

	result := make([]string, 0, len(content)+3)
	result = append(result, headerLine)

	// content lines
	for _, line := range content {
		// Pad so that the inner width (between borders) is consistent across cards.
		padding := max(0, cardWidth-contentWidth-4)
		contentLine := borderVertical + " " + line + strings.Repeat(" ", padding) + " " + borderVertical
		result = append(result, contentLine)
	}

	// bottom border
	actualCardWidth := cardWidth
	bottomLine := borderBottomLeft + strings.Repeat(borderHorizontal, actualCardWidth-2) + borderBottomRight
	result = append(result, bottomLine)

	return result, actualCardWidth
}
