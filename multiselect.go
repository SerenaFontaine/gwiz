package gwiz

import (
	"fmt"

	"github.com/SerenaFontaine/tui"
)

// MultiSelectStep presents a list of options where multiple items can be toggled.
type MultiSelectStep struct {
	BaseStep
	Prompt      string
	Options     []Option
	OptionsFunc func(state State) []Option
	ResultKey   string
	MinSelect   int
	MaxSelect   int

	options  []Option
	cursor   int
	selected []bool
}

func (s *MultiSelectStep) Init(state State) Cmd {
	if s.OptionsFunc != nil {
		s.options = s.OptionsFunc(state)
	} else {
		s.options = s.Options
	}
	s.selected = make([]bool, len(s.options))
	s.cursor = 0

	if vals := state.GetStringSlice(s.ResultKey); vals != nil {
		valSet := make(map[string]bool, len(vals))
		for _, v := range vals {
			valSet[v] = true
		}
		for i, opt := range s.options {
			s.selected[i] = valSet[opt.Value]
		}
	}
	return nil
}

func (s *MultiSelectStep) Update(msg Msg, state State) (Step, Cmd) {
	km, ok := msg.(tui.KeyMsg)
	if !ok {
		return s, nil
	}

	switch km.Type {
	case tui.KeyDown:
		if s.cursor < len(s.options)-1 {
			s.cursor++
		}
	case tui.KeyUp:
		if s.cursor > 0 {
			s.cursor--
		}
	case tui.KeyRune:
		switch km.Rune {
		case 'j':
			if s.cursor < len(s.options)-1 {
				s.cursor++
			}
		case 'k':
			if s.cursor > 0 {
				s.cursor--
			}
		}
	case tui.KeySpace:
		if len(s.options) > 0 && !s.options[s.cursor].Disabled {
			if s.selected[s.cursor] {
				s.selected[s.cursor] = false
			} else {
				if s.MaxSelect > 0 && s.selectedCount() >= s.MaxSelect {
					return s, nil
				}
				s.selected[s.cursor] = true
			}
		}
	case tui.KeyEnter:
		vals := s.selectedValues()
		return s, tui.Batch(
			func() Msg { return StepResultMsg{Key: s.ResultKey, Value: vals} },
			func() Msg { return NextMsg{} },
		)
	case tui.KeyEscape, tui.KeyBackspace:
		return s, func() Msg { return PrevMsg{} }
	}
	return s, nil
}

func (s *MultiSelectStep) Validate(state State) error {
	count := s.selectedCount()
	if s.MinSelect > 0 && count < s.MinSelect {
		return fmt.Errorf("select at least %d item(s)", s.MinSelect)
	}
	if s.MaxSelect > 0 && count > s.MaxSelect {
		return fmt.Errorf("select at most %d item(s)", s.MaxSelect)
	}
	return nil
}

func (s *MultiSelectStep) selectedCount() int {
	count := 0
	for _, v := range s.selected {
		if v {
			count++
		}
	}
	return count
}

func (s *MultiSelectStep) selectedValues() []string {
	var vals []string
	for i, opt := range s.options {
		if s.selected[i] {
			vals = append(vals, opt.Value)
		}
	}
	return vals
}

func (s *MultiSelectStep) Render(buf *tui.Buffer, area tui.Rect, state State) {
	y := area.Y
	if s.Prompt != "" {
		buf.SetString(area.X, y, s.Prompt, tui.NewStyle())
		y += 2
	}
	for i, opt := range s.options {
		if y >= area.Y+area.Height {
			break
		}
		check := "[ ] "
		if s.selected[i] {
			check = "[x] "
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
		line := prefix + check + opt.Label
		buf.SetString(area.X, y, line, style)
		y++
		if opt.Description != "" && y < area.Y+area.Height {
			desc := "      " + opt.Description
			if opt.Disabled && opt.DisabledMsg != "" {
				desc = "      " + opt.DisabledMsg
			}
			buf.SetString(area.X, y, desc, tui.NewStyle().Dim(true))
			y++
		}
	}
}

func (s *MultiSelectStep) KeyHints() []KeyHint {
	return []KeyHint{
		{Key: "↑/↓", Label: "Move"},
		{Key: "Space", Label: "Toggle"},
		{Key: "Enter", Label: "Confirm"},
		{Key: "Esc", Label: "Back"},
		{Key: "q", Label: "Quit"},
	}
}
