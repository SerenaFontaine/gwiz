package gwiz

import (
	"fmt"

	"github.com/SerenaFontaine/tui"
)

// TableStep displays tabular data with optional row selection.
type TableStep struct {
	BaseStep
	Headers     []string
	Rows        [][]string
	HeadersFunc func(state State) []string
	RowsFunc    func(state State) [][]string
	Selectable  bool
	ResultKey   string

	headers []string
	rows    [][]string
	cursor  int
}

func (s *TableStep) Init(state State) Cmd {
	if s.HeadersFunc != nil {
		s.headers = s.HeadersFunc(state)
	} else {
		s.headers = s.Headers
	}
	if s.RowsFunc != nil {
		s.rows = s.RowsFunc(state)
	} else {
		s.rows = s.Rows
	}
	s.cursor = 0
	return nil
}

func (s *TableStep) Update(msg Msg, state State) (Step, Cmd) {
	km, ok := msg.(tui.KeyMsg)
	if !ok {
		return s, nil
	}

	switch km.Type {
	case tui.KeyDown:
		if s.Selectable && s.cursor < len(s.rows)-1 {
			s.cursor++
		}
	case tui.KeyUp:
		if s.Selectable && s.cursor > 0 {
			s.cursor--
		}
	case tui.KeyRune:
		if s.Selectable {
			switch km.Rune {
			case 'j':
				if s.cursor < len(s.rows)-1 {
					s.cursor++
				}
			case 'k':
				if s.cursor > 0 {
					s.cursor--
				}
			}
		}
	case tui.KeyEnter:
		if s.Selectable && len(s.rows) > 0 {
			row := s.rows[s.cursor]
			return s, tui.Batch(
				func() Msg { return StepResultMsg{Key: s.ResultKey, Value: row} },
				func() Msg { return NextMsg{} },
			)
		}
		return s, func() Msg { return NextMsg{} }
	case tui.KeyEscape, tui.KeyBackspace:
		return s, func() Msg { return PrevMsg{} }
	}
	return s, nil
}

func (s *TableStep) Render(buf *tui.Buffer, area tui.Rect, state State) {
	if len(s.headers) == 0 {
		return
	}

	colWidths := make([]int, len(s.headers))
	for i, h := range s.headers {
		colWidths[i] = len(h)
	}
	for _, row := range s.rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	y := area.Y
	x := area.X
	headerStyle := tui.NewStyle().Bold(true)
	for i, h := range s.headers {
		label := fmt.Sprintf("%-*s", colWidths[i]+2, h)
		buf.SetString(x, y, label, headerStyle)
		x += colWidths[i] + 2
	}
	y++

	tw := totalWidth(colWidths)
	for x := area.X; x < area.X+area.Width && x < area.X+tw; x++ {
		buf.SetChar(x, y, '─', tui.NewStyle().Dim(true))
	}
	y++

	for ri, row := range s.rows {
		if y >= area.Y+area.Height {
			break
		}
		x := area.X
		style := tui.NewStyle()
		if s.Selectable && ri == s.cursor {
			style = style.Bold(true).Reverse(true)
		}
		for ci, cell := range row {
			if ci < len(colWidths) {
				label := fmt.Sprintf("%-*s", colWidths[ci]+2, cell)
				buf.SetString(x, y, label, style)
				x += colWidths[ci] + 2
			}
		}
		y++
	}
}

func totalWidth(widths []int) int {
	total := 0
	for _, w := range widths {
		total += w + 2
	}
	return total
}

func (s *TableStep) KeyHints() []KeyHint {
	if s.Selectable {
		return []KeyHint{
			{Key: "↑/↓", Label: "Select"},
			{Key: "Enter", Label: "Confirm"},
			{Key: "Esc", Label: "Back"},
			{Key: "q", Label: "Quit"},
		}
	}
	return []KeyHint{
		{Key: "Enter", Label: "Next"},
		{Key: "Esc", Label: "Back"},
		{Key: "q", Label: "Quit"},
	}
}
