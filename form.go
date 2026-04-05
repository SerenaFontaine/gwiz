package gwiz

import (
	"fmt"
	"unicode/utf8"

	"github.com/SerenaFontaine/tui"
)

// FormField describes a single field within a FormStep.
type FormField struct {
	Label    string
	Key      string
	Default  string
	Validate func(value string) error
}

// FormStep presents multiple labeled text fields on a single screen.
// cursors stores rune indices (not byte offsets).
type FormStep struct {
	BaseStep
	Fields []FormField

	values   []string
	cursors  []int // rune indices
	focusIdx int
}

func (s *FormStep) Init(state State) Cmd {
	s.values = make([]string, len(s.Fields))
	s.cursors = make([]int, len(s.Fields))
	s.focusIdx = 0
	for i, f := range s.Fields {
		s.values[i] = f.Default
		if v, ok := state.Get(f.Key); ok {
			if str, ok := v.(string); ok {
				s.values[i] = str
			}
		}
		s.cursors[i] = utf8.RuneCountInString(s.values[i])
	}
	return nil
}

func (s *FormStep) Update(msg Msg, state State) (Step, Cmd) {
	km, ok := msg.(tui.KeyMsg)
	if !ok {
		return s, nil
	}

	switch km.Type {
	case tui.KeyTab:
		if s.focusIdx < len(s.Fields)-1 {
			s.focusIdx++
		}
	case tui.KeyBacktab:
		if s.focusIdx > 0 {
			s.focusIdx--
		}
	case tui.KeyEnter:
		if s.focusIdx < len(s.Fields)-1 {
			s.focusIdx++
		} else {
			cmds := make([]tui.Cmd, 0, len(s.Fields)+1)
			for i, f := range s.Fields {
				key := f.Key
				val := s.values[i]
				cmds = append(cmds, func() Msg { return StepResultMsg{Key: key, Value: val} })
			}
			cmds = append(cmds, func() Msg { return NextMsg{} })
			return s, tui.Batch(cmds...)
		}
	case tui.KeyEscape:
		return s, func() Msg { return PrevMsg{} }
	case tui.KeyRune:
		idx := s.focusIdx
		runes := []rune(s.values[idx])
		runes = append(runes[:s.cursors[idx]], append([]rune{km.Rune}, runes[s.cursors[idx]:]...)...)
		s.values[idx] = string(runes)
		s.cursors[idx]++
	case tui.KeyBackspace:
		idx := s.focusIdx
		if s.cursors[idx] > 0 {
			runes := []rune(s.values[idx])
			runes = append(runes[:s.cursors[idx]-1], runes[s.cursors[idx]:]...)
			s.values[idx] = string(runes)
			s.cursors[idx]--
		}
	case tui.KeyLeft:
		if s.cursors[s.focusIdx] > 0 {
			s.cursors[s.focusIdx]--
		}
	case tui.KeyRight:
		runes := []rune(s.values[s.focusIdx])
		if s.cursors[s.focusIdx] < len(runes) {
			s.cursors[s.focusIdx]++
		}
	}
	return s, nil
}

func (s *FormStep) Validate(state State) error {
	for i, f := range s.Fields {
		if f.Validate != nil {
			if err := f.Validate(s.values[i]); err != nil {
				return fmt.Errorf("%s: %w", f.Label, err)
			}
		}
	}
	return nil
}

func (s *FormStep) Render(buf *tui.Buffer, area tui.Rect, state State) {
	y := area.Y
	maxLabel := 0
	for _, f := range s.Fields {
		if len(f.Label) > maxLabel {
			maxLabel = len(f.Label)
		}
	}
	for i, f := range s.Fields {
		if y >= area.Y+area.Height {
			break
		}
		labelStyle := tui.NewStyle().Dim(true)
		inputStyle := tui.NewStyle()
		if i == s.focusIdx {
			labelStyle = tui.NewStyle().Bold(true)
			inputStyle = tui.NewStyle().Bold(true)
		}
		label := fmt.Sprintf("%-*s: ", maxLabel, f.Label)
		buf.SetString(area.X, y, label, labelStyle)
		valX := area.X + len(label)
		val := s.values[i]
		buf.SetString(valX, y, val, inputStyle)
		if i == s.focusIdx {
			runes := []rune(val)
			cursorX := valX + s.cursors[i]
			if cursorX < area.X+area.Width {
				ch := ' '
				if s.cursors[i] < len(runes) {
					ch = runes[s.cursors[i]]
				}
				buf.SetChar(cursorX, y, ch, tui.NewStyle().Reverse(true))
			}
		}
		y += 2
	}
}

func (s *FormStep) AcceptsTextInput() bool { return true }

func (s *FormStep) KeyHints() []KeyHint {
	return []KeyHint{
		{Key: "Tab", Label: "Next field"},
		{Key: "Shift+Tab", Label: "Prev field"},
		{Key: "Enter", Label: "Confirm"},
		{Key: "Esc", Label: "Back"},
	}
}
