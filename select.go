package gwiz

import (
	"github.com/SerenaFontaine/tui"
)

// SelectStep presents a single-selection list.
type SelectStep struct {
	BaseStep
	Prompt      string
	Options     []Option
	OptionsFunc func(state State) []Option
	ResultKey   string
	Default     string

	options []Option
	cursor  int
}

func (s *SelectStep) Init(state State) Cmd {
	if s.OptionsFunc != nil {
		s.options = s.OptionsFunc(state)
	} else {
		s.options = s.Options
	}
	s.cursor = 0
	if s.Default != "" {
		for i, opt := range s.options {
			if opt.Value == s.Default {
				s.cursor = i
				break
			}
		}
	}
	if v, ok := state.Get(s.ResultKey); ok {
		if val, ok := v.(string); ok {
			for i, opt := range s.options {
				if opt.Value == val {
					s.cursor = i
					break
				}
			}
		}
	}
	return nil
}

func (s *SelectStep) Update(msg Msg, state State) (Step, Cmd) {
	km, ok := msg.(tui.KeyMsg)
	if !ok {
		return s, nil
	}

	switch km.Type {
	case tui.KeyDown:
		s.moveDown()
	case tui.KeyUp:
		s.moveUp()
	case tui.KeyRune:
		switch km.Rune {
		case 'j':
			s.moveDown()
		case 'k':
			s.moveUp()
		}
	case tui.KeyEnter:
		if len(s.options) == 0 {
			return s, nil
		}
		opt := s.options[s.cursor]
		if opt.Disabled {
			return s, nil
		}
		return s, tui.Batch(
			func() Msg { return StepResultMsg{Key: s.ResultKey, Value: opt.Value} },
			func() Msg { return NextMsg{} },
		)
	case tui.KeyEscape, tui.KeyBackspace:
		return s, func() Msg { return PrevMsg{} }
	}
	return s, nil
}

func (s *SelectStep) moveDown() {
	if s.cursor < len(s.options)-1 {
		s.cursor++
	}
}

func (s *SelectStep) moveUp() {
	if s.cursor > 0 {
		s.cursor--
	}
}

func (s *SelectStep) Render(buf *tui.Buffer, area tui.Rect, state State) {
	y := area.Y

	if s.Prompt != "" {
		buf.SetString(area.X, y, s.Prompt, tui.NewStyle())
		y += 2
	}

	for i, opt := range s.options {
		if y >= area.Y+area.Height {
			break
		}

		prefix := "  "
		style := tui.NewStyle()
		if i == s.cursor {
			prefix = "> "
			style = style.Bold(true)
		}
		if opt.Disabled {
			style = style.Dim(true)
		}

		label := prefix + opt.Label
		buf.SetString(area.X, y, label, style)
		y++

		if opt.Description != "" && y < area.Y+area.Height {
			desc := "    " + opt.Description
			if opt.Disabled && opt.DisabledMsg != "" {
				desc = "    " + opt.DisabledMsg
			}
			buf.SetString(area.X, y, desc, tui.NewStyle().Dim(true))
			y++
		}
	}
}

func (s *SelectStep) KeyHints() []KeyHint {
	return []KeyHint{
		{Key: "↑/↓", Label: "Select"},
		{Key: "Enter", Label: "Confirm"},
		{Key: "Esc", Label: "Back"},
		{Key: "q", Label: "Quit"},
	}
}
