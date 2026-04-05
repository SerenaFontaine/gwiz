package gwiz

import (
	"fmt"
	"strings"

	"github.com/SerenaFontaine/tui"
)

// renderChrome draws the wizard border, header, and nav bar. Returns the content area Rect.
func renderChrome(buf *tui.Buffer, area tui.Rect, title, banner string, currentStep, totalSteps int, hints []KeyHint, theme Theme) tui.Rect {
	block := tui.Block{
		Border: tui.BorderRounded,
		Style:  tui.NewStyle().Fg(theme.Primary),
	}
	inner := block.Render(buf, area)

	// Split inner area: header (1 row) | content (flex) | nav bar (1 row)
	rows := tui.VSplit(inner, tui.Fixed(1), tui.Flex(1), tui.Fixed(1))
	headerArea := rows[0]
	contentArea := rows[1]
	navArea := rows[2]

	// Draw separator lines
	buf.DrawHLine(inner.X, headerArea.Bottom(), inner.Width, '─', tui.NewStyle().Fg(theme.Primary))
	buf.DrawHLine(inner.X, navArea.Y-1, inner.Width, '─', tui.NewStyle().Fg(theme.Primary))

	// Adjust content area to account for separator lines
	contentArea = tui.NewRect(contentArea.X, contentArea.Y+1, contentArea.Width, contentArea.Height-2)
	if contentArea.Height < 0 {
		contentArea = tui.NewRect(contentArea.X, contentArea.Y, contentArea.Width, 0)
	}

	renderHeader(buf, headerArea, title, banner, currentStep, totalSteps, theme)
	renderNavBar(buf, navArea, hints, theme)

	return contentArea.InnerPadding(0, 1, 0, 1)
}

// renderHeader draws the wizard title (left), optional banner, and step progress (right).
func renderHeader(buf *tui.Buffer, area tui.Rect, title, banner string, currentStep, totalSteps int, theme Theme) {
	titleStyle := theme.AccentStyle()
	cx := buf.SetString(area.X+1, area.Y, title, titleStyle)
	if banner != "" {
		buf.SetString(area.X+1+cx+1, area.Y, banner, theme.DimStyle())
	}

	progress := fmt.Sprintf("Step %d of %d", currentStep+1, totalSteps)
	progressStyle := theme.DimStyle()
	x := area.X + area.Width - len(progress) - 1
	if x > area.X+len(title)+2 {
		buf.SetString(x, area.Y, progress, progressStyle)
	}
}

// renderNavBar draws context-sensitive keybindings.
func renderNavBar(buf *tui.Buffer, area tui.Rect, hints []KeyHint, theme Theme) {
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = h.Key + " " + h.Label
	}
	line := "  " + strings.Join(parts, "   ")

	keyStyle := tui.NewStyle().Fg(theme.Text).Bold(true)
	buf.SetString(area.X, area.Y, line, keyStyle)
}
