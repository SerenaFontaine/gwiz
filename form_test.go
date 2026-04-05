package gwiz

import (
	"errors"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestFormStep_Init(t *testing.T) {
	s := &FormStep{
		BaseStep: BaseStep{TitleText: "Config"},
		Fields: []FormField{
			{Label: "Port", Key: "port", Default: "8080"},
			{Label: "Host", Key: "host", Default: "localhost"},
		},
	}
	s.Init(newState())
	if len(s.values) != 2 { t.Fatalf("expected 2 values, got %d", len(s.values)) }
	if s.values[0] != "8080" { t.Fatalf("expected default '8080', got %q", s.values[0]) }
}

func TestFormStep_InitRestoresFromState(t *testing.T) {
	s := &FormStep{
		BaseStep: BaseStep{TitleText: "Config"},
		Fields: []FormField{
			{Label: "Port", Key: "port", Default: "8080"},
			{Label: "Host", Key: "host", Default: "localhost"},
		},
	}
	st := newState()
	st.Set("port", "9090")
	s.Init(st)
	if s.values[0] != "9090" { t.Fatalf("expected restored '9090', got %q", s.values[0]) }
	if s.values[1] != "localhost" { t.Fatalf("expected default 'localhost', got %q", s.values[1]) }
}

func TestFormStep_TabNavigation(t *testing.T) {
	s := &FormStep{
		BaseStep: BaseStep{TitleText: "Config"},
		Fields: []FormField{{Label: "A", Key: "a"}, {Label: "B", Key: "b"}, {Label: "C", Key: "c"}},
	}
	s.Init(newState())
	if s.focusIdx != 0 { t.Fatalf("expected focus at 0, got %d", s.focusIdx) }

	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyTab}, newState())
	fs := step.(*FormStep)
	if fs.focusIdx != 1 { t.Fatalf("expected focus at 1, got %d", fs.focusIdx) }

	step, _ = fs.Update(tui.KeyMsg{Type: tui.KeyBacktab}, newState())
	fs = step.(*FormStep)
	if fs.focusIdx != 0 { t.Fatalf("expected focus at 0, got %d", fs.focusIdx) }
}

func TestFormStep_EnterOnLastFieldAdvances(t *testing.T) {
	s := &FormStep{
		BaseStep: BaseStep{TitleText: "Config"},
		Fields: []FormField{{Label: "A", Key: "a"}, {Label: "B", Key: "b"}},
	}
	s.Init(newState())
	s.focusIdx = 1
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil { t.Fatal("expected command on Enter at last field") }
}

func TestFormStep_EnterOnNonLastFieldTabsForward(t *testing.T) {
	s := &FormStep{
		BaseStep: BaseStep{TitleText: "Config"},
		Fields: []FormField{{Label: "A", Key: "a"}, {Label: "B", Key: "b"}},
	}
	s.Init(newState())
	s.focusIdx = 0
	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	fs := step.(*FormStep)
	if fs.focusIdx != 1 { t.Fatalf("expected focus at 1, got %d", fs.focusIdx) }
}

func TestFormStep_ValidateFields(t *testing.T) {
	s := &FormStep{
		BaseStep: BaseStep{TitleText: "Config"},
		Fields: []FormField{
			{Label: "Port", Key: "port", Validate: func(v string) error {
				if v == "" { return errors.New("required") }
				return nil
			}},
		},
	}
	s.Init(newState())
	s.values[0] = ""
	err := s.Validate(newState())
	if err == nil { t.Fatal("expected validation error") }
}
