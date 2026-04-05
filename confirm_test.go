package gwiz

import (
	"strings"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestConfirmStep_TabSwitchesButtons(t *testing.T) {
	s := &ConfirmStep{
		BaseStep:    BaseStep{TitleText: "Confirm"},
		SummaryFunc: func(state State) string { return "Summary" },
		ResultKey:   "confirmed",
	}
	s.Init(newState())
	if s.selected != 0 {
		t.Fatalf("expected initial selection 0, got %d", s.selected)
	}

	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyTab}, newState())
	cs := step.(*ConfirmStep)
	if cs.selected != 1 {
		t.Fatalf("expected selection 1 after tab, got %d", cs.selected)
	}

	step, _ = cs.Update(tui.KeyMsg{Type: tui.KeyTab}, newState())
	cs = step.(*ConfirmStep)
	if cs.selected != 0 {
		t.Fatal("tab should wrap around to 0")
	}
}

func TestConfirmStep_EnterConfirm(t *testing.T) {
	s := &ConfirmStep{
		BaseStep:    BaseStep{TitleText: "Confirm"},
		SummaryFunc: func(state State) string { return "Summary" },
		ResultKey:   "confirmed",
	}
	s.Init(newState())
	s.selected = 0
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
}

func TestConfirmStep_EnterCancel(t *testing.T) {
	s := &ConfirmStep{
		BaseStep:    BaseStep{TitleText: "Confirm"},
		SummaryFunc: func(state State) string { return "Summary" },
		ResultKey:   "confirmed",
	}
	s.Init(newState())
	s.selected = 1
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil {
		t.Fatal("expected command on Enter (back)")
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
