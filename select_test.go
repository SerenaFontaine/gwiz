package gwiz

import (
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestSelectStep_InitSetsDefaultCursor(t *testing.T) {
	s := &SelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}},
		ResultKey: "choice",
		Default:   "b",
	}
	s.Init(newState())
	if s.cursor != 1 {
		t.Fatalf("expected cursor at 1 (default 'b'), got %d", s.cursor)
	}
}

func TestSelectStep_InitOptionsFuncOverridesOptions(t *testing.T) {
	s := &SelectStep{
		BaseStep: BaseStep{TitleText: "Pick"},
		OptionsFunc: func(state State) []Option {
			return []Option{{Label: "X", Value: "x"}}
		},
		ResultKey: "choice",
	}
	s.Init(newState())
	if len(s.options) != 1 || s.options[0].Value != "x" {
		t.Fatal("OptionsFunc should populate options")
	}
}

func TestSelectStep_NavigateUpDown(t *testing.T) {
	s := &SelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}, {Label: "C", Value: "c"}},
		ResultKey: "choice",
	}
	s.Init(newState())

	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyDown}, newState())
	ss := step.(*SelectStep)
	if ss.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", ss.cursor)
	}

	step, _ = ss.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'j'}, newState())
	ss = step.(*SelectStep)
	if ss.cursor != 2 {
		t.Fatalf("expected cursor 2, got %d", ss.cursor)
	}

	step, _ = ss.Update(tui.KeyMsg{Type: tui.KeyUp}, newState())
	ss = step.(*SelectStep)
	if ss.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", ss.cursor)
	}
}

func TestSelectStep_SkipsDisabledOnEnter(t *testing.T) {
	s := &SelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a", Disabled: true}, {Label: "B", Value: "b"}},
		ResultKey: "choice",
	}
	s.Init(newState())
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd != nil {
		t.Fatal("should not advance on disabled option")
	}
}

func TestSelectStep_EnterSetsStateAndAdvances(t *testing.T) {
	s := &SelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}},
		ResultKey: "choice",
	}
	st := newState()
	s.Init(st)

	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyDown}, st)
	ss := step.(*SelectStep)
	_, cmd := ss.Update(tui.KeyMsg{Type: tui.KeyEnter}, st)
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
}

func TestSelectStep_BoundsCheck(t *testing.T) {
	s := &SelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}},
		ResultKey: "choice",
	}
	s.Init(newState())

	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyUp}, newState())
	ss := step.(*SelectStep)
	if ss.cursor != 0 {
		t.Fatalf("cursor should stay at 0, got %d", ss.cursor)
	}

	step, _ = ss.Update(tui.KeyMsg{Type: tui.KeyDown}, newState())
	ss = step.(*SelectStep)
	if ss.cursor != 0 {
		t.Fatalf("cursor should stay at 0, got %d", ss.cursor)
	}
}
