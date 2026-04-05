package gwiz

import (
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestNew_DefaultTheme(t *testing.T) {
	w := New()
	if w.theme.Name == "" {
		t.Fatal("expected default theme")
	}
}

func TestNew_WithOptions(t *testing.T) {
	w := New(
		WithTitle("Test Wizard"),
		WithTheme(ThemeDracula),
		WithBanner("v1.0"),
	)
	if w.title != "Test Wizard" {
		t.Fatalf("expected title 'Test Wizard', got %q", w.title)
	}
	if w.theme.Name != "Dracula" {
		t.Fatalf("expected Dracula theme, got %q", w.theme.Name)
	}
	if w.banner != "v1.0" {
		t.Fatalf("expected banner 'v1.0', got %q", w.banner)
	}
}

func TestAddStep_InOrder(t *testing.T) {
	w := New()
	w.AddStep("a", InfoStep{BaseStep: BaseStep{TitleText: "A"}})
	w.AddStep("b", InfoStep{BaseStep: BaseStep{TitleText: "B"}})
	if len(w.steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(w.steps))
	}
	if w.steps[0].name != "a" || w.steps[1].name != "b" {
		t.Fatal("steps not in insertion order")
	}
}

func TestWizardComponent_Init(t *testing.T) {
	w := New(WithTitle("Test"))
	w.AddStep("a", InfoStep{BaseStep: BaseStep{TitleText: "A"}})
	w.AddStep("b", InfoStep{BaseStep: BaseStep{TitleText: "B"}})

	cmd := w.Init()
	_ = cmd

	if w.current != 0 {
		t.Fatalf("expected current step 0, got %d", w.current)
	}
}

func TestWizardComponent_UpdateNextMsg(t *testing.T) {
	w := New(WithTitle("Test"))
	w.AddStep("a", InfoStep{BaseStep: BaseStep{TitleText: "A"}})
	w.AddStep("b", InfoStep{BaseStep: BaseStep{TitleText: "B"}})
	w.Init()

	comp, _ := w.Update(NextMsg{})
	wiz := comp.(*Wizard)
	if wiz.current != 1 {
		t.Fatalf("expected current step 1, got %d", wiz.current)
	}
}

func TestWizardComponent_UpdatePrevMsg(t *testing.T) {
	w := New(WithTitle("Test"))
	w.AddStep("a", InfoStep{BaseStep: BaseStep{TitleText: "A"}})
	w.AddStep("b", InfoStep{BaseStep: BaseStep{TitleText: "B"}})
	w.Init()
	w.Update(NextMsg{})

	comp, _ := w.Update(PrevMsg{})
	wiz := comp.(*Wizard)
	if wiz.current != 0 {
		t.Fatalf("expected current step 0, got %d", wiz.current)
	}
}

func TestWizardComponent_UpdateStepResultMsg(t *testing.T) {
	w := New(WithTitle("Test"))
	w.AddStep("a", InfoStep{BaseStep: BaseStep{TitleText: "A"}})
	w.Init()

	w.Update(StepResultMsg{Key: "color", Value: "blue"})
	if w.state.GetString("color") != "blue" {
		t.Fatalf("expected state 'color'='blue', got %q", w.state.GetString("color"))
	}
}

func TestWizardComponent_Render(t *testing.T) {
	w := New(WithTitle("My Wizard"))
	w.AddStep("a", InfoStep{BaseStep: BaseStep{TitleText: "Welcome"}, Content: "Hello"})
	w.Init()

	buf := tui.NewBuffer(60, 20)
	area := tui.NewRect(0, 0, 60, 20)
	w.Render(buf, area)
}

func TestWizardComponent_NextPastEnd(t *testing.T) {
	w := New(WithTitle("Test"))
	w.AddStep("a", InfoStep{BaseStep: BaseStep{TitleText: "A"}})
	w.Init()

	comp, cmd := w.Update(NextMsg{})
	_ = comp
	if cmd == nil {
		t.Fatal("expected quit command when advancing past last step")
	}
}
