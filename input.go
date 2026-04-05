package gwiz

import (
	"github.com/SerenaFontaine/tui"
)

// InputStep collects a single line of text from the user.
type InputStep struct {
	BaseStep
	Prompt       string
	Placeholder  string
	Default      string
	ResultKey    string
	ValidateFunc func(value string) error

	value     string
	cursorPos int
}

func (s *InputStep) Init(state State) Cmd {
	s.value = s.Default
	s.cursorPos = len(s.value)

	if v, ok := state.Get(s.ResultKey); ok {
		if str, ok := v.(string); ok {
			s.value = str
			s.cursorPos = len(str)
		}
	}
	return nil
}

func (s *InputStep) Update(msg Msg, state State) (Step, Cmd) {
	km, ok := msg.(tui.KeyMsg)
	if !ok {
		return s, nil
	}

	switch km.Type {
	case tui.KeyRune:
		s.value = s.value[:s.cursorPos] + string(km.Rune) + s.value[s.cursorPos:]
		s.cursorPos++
	case tui.KeyBackspace:
		if s.cursorPos > 0 {
			s.value = s.value[:s.cursorPos-1] + s.value[s.cursorPos:]
			s.cursorPos--
		}
	case tui.KeyDelete:
		if s.cursorPos < len(s.value) {
			s.value = s.value[:s.cursorPos] + s.value[s.cursorPos+1:]
		}
	case tui.KeyLeft:
		if s.cursorPos > 0 {
			s.cursorPos--
		}
	case tui.KeyRight:
		if s.cursorPos < len(s.value) {
			s.cursorPos++
		}
	case tui.KeyHome, tui.KeyCtrlA:
		s.cursorPos = 0
	case tui.KeyEnd, tui.KeyCtrlE:
		s.cursorPos = len(s.value)
	case tui.KeyEnter:
		return s, tui.Batch(
			func() Msg { return StepResultMsg{Key: s.ResultKey, Value: s.value} },
			func() Msg { return NextMsg{} },
		)
	case tui.KeyEscape:
		return s, func() Msg { return PrevMsg{} }
	}
	return s, nil
}

func (s *InputStep) Validate(state State) error {
	if s.ValidateFunc != nil {
		return s.ValidateFunc(s.value)
	}
	return nil
}

func (s *InputStep) Render(buf *tui.Buffer, area tui.Rect, state State) {
	y := area.Y
	if s.Prompt != "" {
		buf.SetString(area.X, y, s.Prompt, tui.NewStyle())
		y += 2
	}
	display := s.value
	if display == "" && s.Placeholder != "" {
		buf.SetString(area.X, y, s.Placeholder, tui.NewStyle().Dim(true))
	} else {
		buf.SetString(area.X, y, display, tui.NewStyle())
		if s.cursorPos <= len(display) {
			cursorX := area.X + s.cursorPos
			if cursorX < area.X+area.Width {
				ch := ' '
				if s.cursorPos < len(display) {
					ch = rune(display[s.cursorPos])
				}
				buf.SetChar(cursorX, y, ch, tui.NewStyle().Reverse(true))
			}
		}
	}
	y++
	if y < area.Y+area.Height {
		for x := area.X; x < area.X+area.Width && x < area.X+40; x++ {
			buf.SetChar(x, y, '─', tui.NewStyle().Dim(true))
		}
	}
}

func (s *InputStep) AcceptsTextInput() bool { return true }

func (s *InputStep) KeyHints() []KeyHint {
	return []KeyHint{
		{Key: "Enter", Label: "Confirm"},
		{Key: "Esc", Label: "Back"},
	}
}
