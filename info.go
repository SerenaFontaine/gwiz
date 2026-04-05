package gwiz

import (
	"strings"

	"github.com/SerenaFontaine/tui"
)

// InfoStep displays formatted text content. Used for welcome screens,
// hardware summaries, warnings, or any read-only informational screen.
type InfoStep struct {
	BaseStep
	Content       string
	ContentFunc   func(state State) string
	ContinueLabel string
}

func (s InfoStep) Update(msg Msg, state State) (Step, Cmd) {
	if km, ok := msg.(tui.KeyMsg); ok {
		switch km.Type {
		case tui.KeyEnter:
			return s, func() Msg { return NextMsg{} }
		case tui.KeyEscape, tui.KeyBackspace:
			return s, func() Msg { return PrevMsg{} }
		}
	}
	return s, nil
}

func (s InfoStep) Render(buf *tui.Buffer, area tui.Rect, state State) {
	content := s.Content
	if s.ContentFunc != nil {
		content = s.ContentFunc(state)
	}

	y := area.Y
	if s.TitleText != "" {
		buf.SetString(area.X, y, s.TitleText, tui.NewStyle().Bold(true))
		y++
	}
	if s.DescriptionText != "" {
		buf.SetString(area.X, y, s.DescriptionText, tui.NewStyle().Dim(true))
		y++
	}
	if s.TitleText != "" || s.DescriptionText != "" {
		y++
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if y >= area.Y+area.Height {
			break
		}
		renderMarkupLine(buf, area.X, y, area.Width, line)
		y++
	}
}

func (s InfoStep) KeyHints() []KeyHint {
	label := "Next"
	if s.ContinueLabel != "" {
		label = s.ContinueLabel
	}
	return []KeyHint{
		{Key: "Enter", Label: label},
		{Key: "Esc", Label: "Back"},
		{Key: "q", Label: "Quit"},
	}
}

// renderMarkupLine renders a single line with simple inline markup.
// Supported: **bold**, *dim*, {green}text{/}, {red}text{/}, {yellow}text{/}, {cyan}text{/}
func renderMarkupLine(buf *tui.Buffer, x, y, maxWidth int, line string) {
	cx := x
	i := 0
	style := tui.NewStyle()

	for i < len(line) && cx < x+maxWidth {
		if i+2 <= len(line) && line[i:i+2] == "**" {
			end := strings.Index(line[i+2:], "**")
			if end >= 0 {
				text := line[i+2 : i+2+end]
				cx += buf.SetString(cx, y, text, style.Bold(true))
				i += 2 + end + 2
				continue
			}
		}

		if i < len(line) && line[i] == '*' && (i+1 < len(line) && line[i+1] != '*') {
			end := strings.Index(line[i+1:], "*")
			if end >= 0 {
				text := line[i+1 : i+1+end]
				cx += buf.SetString(cx, y, text, style.Dim(true))
				i += 1 + end + 1
				continue
			}
		}

		if i < len(line) && line[i] == '{' {
			colorEnd := strings.Index(line[i:], "}")
			if colorEnd > 0 {
				colorName := line[i+1 : i+colorEnd]
				color := colorByName(colorName)
				if !color.IsZero() {
					rest := line[i+colorEnd+1:]
					textEnd := strings.Index(rest, "{/}")
					if textEnd >= 0 {
						text := rest[:textEnd]
						cx += buf.SetString(cx, y, text, style.Fg(color))
						i += colorEnd + 1 + textEnd + 3
						continue
					}
				}
			}
		}

		buf.SetChar(cx, y, rune(line[i]), style)
		cx++
		i++
	}
}

func colorByName(name string) tui.Color {
	switch name {
	case "green":
		return tui.Green
	case "red":
		return tui.Red
	case "yellow":
		return tui.Yellow
	case "cyan":
		return tui.Cyan
	case "blue":
		return tui.Blue
	case "magenta":
		return tui.Magenta
	default:
		return tui.Color{}
	}
}
