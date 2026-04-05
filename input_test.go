package gwiz

import (
	"errors"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestInputStep_InitRestoresFromState(t *testing.T) {
	s := &InputStep{
		BaseStep:  BaseStep{TitleText: "Path"},
		ResultKey: "path",
		Default:   "/etc/default",
	}
	st := newState()
	st.Set("path", "/opt/custom")
	s.Init(st)

	if s.value != "/opt/custom" {
		t.Fatalf("expected restored value '/opt/custom', got %q", s.value)
	}
}

func TestInputStep_InitUsesDefault(t *testing.T) {
	s := &InputStep{
		BaseStep:  BaseStep{TitleText: "Path"},
		ResultKey: "path",
		Default:   "/etc/default",
	}
	s.Init(newState())

	if s.value != "/etc/default" {
		t.Fatalf("expected default '/etc/default', got %q", s.value)
	}
}

func TestInputStep_TypeCharacters(t *testing.T) {
	s := &InputStep{
		BaseStep:  BaseStep{TitleText: "Name"},
		ResultKey: "name",
	}
	s.Init(newState())

	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'h'}, newState())
	is := step.(*InputStep)
	step, _ = is.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'i'}, newState())
	is = step.(*InputStep)

	if is.value != "hi" {
		t.Fatalf("expected 'hi', got %q", is.value)
	}
}

func TestInputStep_Backspace(t *testing.T) {
	s := &InputStep{
		BaseStep:  BaseStep{TitleText: "Name"},
		ResultKey: "name",
	}
	s.Init(newState())
	s.value = "hello"
	s.cursorPos = 5

	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyBackspace}, newState())
	is := step.(*InputStep)

	if is.value != "hell" {
		t.Fatalf("expected 'hell', got %q", is.value)
	}
}

func TestInputStep_EnterSetsStateAndAdvances(t *testing.T) {
	s := &InputStep{
		BaseStep:  BaseStep{TitleText: "Name"},
		ResultKey: "name",
	}
	s.Init(newState())
	s.value = "test"

	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
}

func TestInputStep_ValidateFunc(t *testing.T) {
	s := &InputStep{
		BaseStep:  BaseStep{TitleText: "Name"},
		ResultKey: "name",
		ValidateFunc: func(value string) error {
			if value == "" {
				return errors.New("required")
			}
			return nil
		},
	}
	s.Init(newState())
	s.value = ""

	err := s.Validate(newState())
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "required" {
		t.Fatalf("expected 'required', got %q", err.Error())
	}

	s.value = "filled"
	err = s.Validate(newState())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
