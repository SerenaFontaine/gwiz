package gwiz

import (
	"strings"

	"github.com/SerenaFontaine/tui"
)

// ConfirmStep displays a summary and asks the user to confirm or go back.
type ConfirmStep struct {
	BaseStep
	SummaryFunc  func(state State) string
	ConfirmLabel string
	CancelLabel  string
	ResultKey    string

	selected int // 0 = confirm, 1 = cancel
}

func (s *ConfirmStep) Init(state State) Cmd {
	s.selected = 0
	return nil
}

func (s *ConfirmStep) Update(msg Msg, state State) (Step, Cmd) {
	km, ok := msg.(tui.KeyMsg)
	if !ok {
		return s, nil
	}

	switch km.Type {
	case tui.KeyTab, tui.KeyRight, tui.KeyLeft:
		s.selected = 1 - s.selected
	case tui.KeyEnter:
		if s.selected == 0 {
			return s, tui.Batch(
				func() Msg { return StepResultMsg{Key: s.ResultKey, Value: true} },
				func() Msg { return NextMsg{} },
			)
		}
		return s, func() Msg { return PrevMsg{} }
	case tui.KeyEscape:
		return s, func() Msg { return PrevMsg{} }
	}
	return s, nil
}

func (s *ConfirmStep) Render(buf *tui.Buffer, area tui.Rect, state State) {
	y := area.Y
	if s.SummaryFunc != nil {
		summary := s.SummaryFunc(state)
		lines := strings.Split(summary, "\n")
		for _, line := range lines {
			if y >= area.Y+area.Height-3 {
				break
			}
			buf.SetString(area.X, y, line, tui.NewStyle())
			y++
		}
	}

	y = area.Y + area.Height - 2
	if y < area.Y {
		y = area.Y
	}

	confirmLabel := s.ConfirmLabel
	if confirmLabel == "" {
		confirmLabel = "Confirm"
	}
	cancelLabel := s.CancelLabel
	if cancelLabel == "" {
		cancelLabel = "Back"
	}

	confirmBtn := "[ " + confirmLabel + " ]"
	cancelBtn := "[ " + cancelLabel + " ]"

	confirmStyle := tui.NewStyle().Dim(true)
	cancelStyle := tui.NewStyle().Dim(true)
	if s.selected == 0 {
		confirmStyle = tui.NewStyle().Bold(true).Reverse(true)
	} else {
		cancelStyle = tui.NewStyle().Bold(true).Reverse(true)
	}

	cx := area.X
	cx += buf.SetString(cx, y, confirmBtn, confirmStyle)
	cx += 2
	buf.SetString(cx, y, cancelBtn, cancelStyle)
}

func (s *ConfirmStep) KeyHints() []KeyHint {
	return []KeyHint{
		{Key: "Tab", Label: "Switch"},
		{Key: "Enter", Label: "Select"},
		{Key: "Esc", Label: "Back"},
		{Key: "q", Label: "Quit"},
	}
}
