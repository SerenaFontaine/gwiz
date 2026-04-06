package gwiz

import (
	"strings"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestConfirmStep_EnterConfirm(t *testing.T) {
	s := &ConfirmStep{
		BaseStep:    BaseStep{TitleText: "Confirm"},
		SummaryFunc: func(state State) string { return "Summary" },
		ResultKey:   "confirmed",
	}
	s.Init(newState())
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
	// Should produce a batch with StepResultMsg + NextMsg
	msg := cmd()
	if batch, ok := msg.(tui.BatchMsg); ok {
		var hasResult, hasNext bool
		for _, c := range batch {
			m := c()
			if _, ok := m.(StepResultMsg); ok {
				hasResult = true
			}
			if _, ok := m.(NextMsg); ok {
				hasNext = true
			}
		}
		if !hasResult || !hasNext {
			t.Fatal("expected batch with StepResultMsg and NextMsg")
		}
	}
}

func TestConfirmStep_EscGoesBack(t *testing.T) {
	s := &ConfirmStep{
		BaseStep:    BaseStep{TitleText: "Confirm"},
		SummaryFunc: func(state State) string { return "Summary" },
		ResultKey:   "confirmed",
	}
	s.Init(newState())
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEscape}, newState())
	if cmd == nil {
		t.Fatal("expected command on Esc")
	}
	msg := cmd()
	if _, ok := msg.(PrevMsg); !ok {
		t.Fatalf("expected PrevMsg, got %T", msg)
	}
}

func TestConfirmStep_RenderShowsSummary(t *testing.T) {
	s := &ConfirmStep{
		BaseStep:    BaseStep{TitleText: "Review"},
		SummaryFunc: func(state State) string { return "Backend: vllm\nGPU: A100" },
		ResultKey:   "confirmed",
	}
	s.Init(newState())
	buf := tui.NewBuffer(40, 15)
	area := tui.NewRect(0, 0, 40, 15)
	s.Render(buf, area, newState())
	out := bufString(buf)
	if !strings.Contains(out, "Backend: vllm") {
		t.Fatalf("expected summary in render, got: %q", out)
	}
}

func TestConfirmStep_KeyHintsShowConfirmLabel(t *testing.T) {
	s := &ConfirmStep{
		BaseStep:     BaseStep{TitleText: "Confirm"},
		ConfirmLabel: "Install",
		ResultKey:    "confirmed",
	}
	hints := s.KeyHints()
	if hints[0].Label != "Install" {
		t.Fatalf("expected confirm label 'Install', got %q", hints[0].Label)
	}
}
