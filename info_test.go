package gwiz

import (
	"strings"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestInfoStep_Title(t *testing.T) {
	s := InfoStep{BaseStep: BaseStep{TitleText: "Welcome"}}
	if s.Title() != "Welcome" {
		t.Fatalf("expected 'Welcome', got %q", s.Title())
	}
}

func TestInfoStep_RenderStaticContent(t *testing.T) {
	s := InfoStep{
		BaseStep: BaseStep{TitleText: "Welcome"},
		Content:  "Hello, world!",
	}
	buf := tui.NewBuffer(40, 10)
	area := tui.NewRect(0, 0, 40, 10)
	s.Render(buf, area, newState())

	out := bufString(buf)
	if !strings.Contains(out, "Hello, world!") {
		t.Fatalf("expected content in render, got: %q", out)
	}
}

func TestInfoStep_RenderDynamicContent(t *testing.T) {
	s := InfoStep{
		BaseStep: BaseStep{TitleText: "Dynamic"},
		ContentFunc: func(state State) string {
			return "Value: " + state.GetString("key")
		},
	}
	st := newState()
	st.Set("key", "42")

	buf := tui.NewBuffer(40, 10)
	area := tui.NewRect(0, 0, 40, 10)
	s.Render(buf, area, st)

	out := bufString(buf)
	if !strings.Contains(out, "Value: 42") {
		t.Fatalf("expected dynamic content, got: %q", out)
	}
}

func TestInfoStep_EnterAdvances(t *testing.T) {
	s := InfoStep{BaseStep: BaseStep{TitleText: "Welcome"}, Content: "Hi"}
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil {
		t.Fatal("expected command from Enter key")
	}
	msg := cmd()
	if _, ok := msg.(NextMsg); !ok {
		t.Fatalf("expected NextMsg, got %T", msg)
	}
}

func TestInfoStep_ContentFuncOverridesContent(t *testing.T) {
	s := InfoStep{
		BaseStep:    BaseStep{TitleText: "Test"},
		Content:     "Static",
		ContentFunc: func(state State) string { return "Dynamic" },
	}
	buf := tui.NewBuffer(40, 10)
	area := tui.NewRect(0, 0, 40, 10)
	s.Render(buf, area, newState())

	out := bufString(buf)
	if strings.Contains(out, "Static") {
		t.Fatal("ContentFunc should override Content")
	}
	if !strings.Contains(out, "Dynamic") {
		t.Fatalf("expected dynamic content, got: %q", out)
	}
}

func TestInfoStep_KeyHints(t *testing.T) {
	s := InfoStep{BaseStep: BaseStep{TitleText: "Test"}, ContinueLabel: "Go"}
	hints := s.KeyHints()
	found := false
	for _, h := range hints {
		if h.Label == "Go" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected ContinueLabel in hints")
	}
}
