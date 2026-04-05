package gwiz

import (
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestFullWizardFlow(t *testing.T) {
	w := New(WithTitle("Test Wizard"), WithTheme(ThemeNord))

	w.AddStep("welcome", InfoStep{
		BaseStep: BaseStep{TitleText: "Welcome"},
		Content:  "Hello",
	})

	w.AddStep("choice", &SelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}},
		ResultKey: "choice",
	})

	w.AddStep("confirm", &ConfirmStep{
		BaseStep:    BaseStep{TitleText: "Confirm"},
		SummaryFunc: func(state State) string { return "Choice: " + state.GetString("choice") },
		ResultKey:   "confirmed",
	})

	// Simulate wizard flow
	w.Init()
	if w.current != 0 {
		t.Fatalf("expected step 0, got %d", w.current)
	}

	// Render should work without panic
	buf := tui.NewBuffer(80, 24)
	area := tui.NewRect(0, 0, 80, 24)
	w.Render(buf, area)

	// Advance past welcome
	w.Update(NextMsg{})
	if w.current != 1 {
		t.Fatalf("expected step 1, got %d", w.current)
	}

	// Init select step
	w.steps[w.current].step.Init(w.state)

	// Select option B (down + enter)
	w.steps[w.current].step.Update(tui.KeyMsg{Type: tui.KeyDown}, w.state)
	_, cmd := w.steps[w.current].step.Update(tui.KeyMsg{Type: tui.KeyEnter}, w.state)

	// Process the batch command
	if cmd != nil {
		msg := cmd()
		if batch, ok := msg.(tui.BatchMsg); ok {
			for _, c := range batch {
				m := c()
				w.Update(m)
			}
		}
	}

	if w.state.GetString("choice") != "b" {
		t.Fatalf("expected choice 'b', got %q", w.state.GetString("choice"))
	}

	if w.current != 2 {
		t.Fatalf("expected step 2 (confirm), got %d", w.current)
	}

	// Go back
	w.Update(PrevMsg{})
	if w.current != 1 {
		t.Fatalf("expected step 1 after going back, got %d", w.current)
	}

	// State should be preserved
	if w.state.GetString("choice") != "b" {
		t.Fatal("state should be preserved on back navigation")
	}
}

func TestWizardWithSkippableSteps(t *testing.T) {
	w := New(WithTitle("Skip Test"))

	w.AddStep("a", InfoStep{BaseStep: BaseStep{TitleText: "A"}})
	w.AddStep("skip", &skipStub{BaseStep: BaseStep{TitleText: "Skip"}, skip: true})
	w.AddStep("c", InfoStep{BaseStep: BaseStep{TitleText: "C"}})

	w.Init()
	if w.current != 0 {
		t.Fatalf("expected step 0, got %d", w.current)
	}

	w.Update(NextMsg{})
	if w.current != 2 {
		t.Fatalf("expected step 2 (skipping 1), got %d", w.current)
	}

	w.Update(PrevMsg{})
	if w.current != 0 {
		t.Fatalf("expected step 0 (skipping 1 backwards), got %d", w.current)
	}
}
