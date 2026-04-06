package gwiz

import (
	"strings"

	"github.com/SerenaFontaine/tui"
)

// ConfirmStep displays a summary and asks the user to confirm or go back.
// Enter confirms and advances; Esc goes back.
type ConfirmStep struct {
	BaseStep
	SummaryFunc  func(state State) string
	ConfirmLabel string
	CancelLabel  string
	ResultKey    string
}

func (s *ConfirmStep) Init(state State) Cmd {
	return nil
}

func (s *ConfirmStep) Update(msg Msg, state State) (Step, Cmd) {
	km, ok := msg.(tui.KeyMsg)
	if !ok {
		return s, nil
	}

	switch km.Type {
	case tui.KeyEnter:
		return s, tui.Batch(
			func() Msg { return StepResultMsg{Key: s.ResultKey, Value: true} },
			func() Msg { return NextMsg{} },
		)
	case tui.KeyEscape, tui.KeyBackspace:
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
			if y >= area.Y+area.Height {
				break
			}
			buf.SetString(area.X, y, line, tui.NewStyle())
			y++
		}
	}
}

func (s *ConfirmStep) KeyHints() []KeyHint {
	confirmLabel := s.ConfirmLabel
	if confirmLabel == "" {
		confirmLabel = "Confirm"
	}
	return []KeyHint{
		{Key: "Enter", Label: confirmLabel},
		{Key: "Esc", Label: "Back"},
		{Key: "q", Label: "Quit"},
	}
}
