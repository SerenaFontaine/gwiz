package gwiz

import (
	"unicode/utf8"

	"github.com/SerenaFontaine/tui"
)

// InputStep collects a single line of text from the user.
// cursorPos is a rune index (not a byte offset).
type InputStep struct {
	BaseStep
	Prompt       string
	Placeholder  string
	Default      string
	ResultKey    string
	ValidateFunc func(value string) error

	value     string
	cursorPos int // rune index
}

func (s *InputStep) Init(state State) Cmd {
	s.value = s.Default
	s.cursorPos = utf8.RuneCountInString(s.value)

	if v, ok := state.Get(s.ResultKey); ok {
		if str, ok := v.(string); ok {
			s.value = str
			s.cursorPos = utf8.RuneCountInString(str)
		}
	}
	return nil
}

func (s *InputStep) Update(msg Msg, state State) (Step, Cmd) {
	km, ok := msg.(tui.KeyMsg)
	if !ok {
		return s, nil
	}

	runes := []rune(s.value)

	switch km.Type {
	case tui.KeyRune:
		runes = append(runes[:s.cursorPos], append([]rune{km.Rune}, runes[s.cursorPos:]...)...)
		s.value = string(runes)
		s.cursorPos++
	case tui.KeyBackspace:
		if s.cursorPos > 0 {
			runes = append(runes[:s.cursorPos-1], runes[s.cursorPos:]...)
			s.value = string(runes)
			s.cursorPos--
		}
	case tui.KeyDelete:
		if s.cursorPos < len(runes) {
			runes = append(runes[:s.cursorPos], runes[s.cursorPos+1:]...)
			s.value = string(runes)
		}
	case tui.KeyLeft:
		if s.cursorPos > 0 {
			s.cursorPos--
		}
	case tui.KeyRight:
		if s.cursorPos < len(runes) {
			s.cursorPos++
		}
	case tui.KeyHome, tui.KeyCtrlA:
		s.cursorPos = 0
	case tui.KeyEnd, tui.KeyCtrlE:
		s.cursorPos = len(runes)
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
		runes := []rune(display)
		if s.cursorPos <= len(runes) {
			cursorX := area.X + s.cursorPos
			if cursorX < area.X+area.Width {
				ch := ' '
				if s.cursorPos < len(runes) {
					ch = runes[s.cursorPos]
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
