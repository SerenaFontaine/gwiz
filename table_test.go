package gwiz

import (
	"strings"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestTableStep_RenderHeaders(t *testing.T) {
	s := &TableStep{
		BaseStep: BaseStep{TitleText: "GPUs"},
		Headers:  []string{"Name", "VRAM"},
		Rows:     [][]string{{"A100", "80GB"}, {"RTX 4090", "24GB"}},
	}
	s.Init(newState())
	buf := tui.NewBuffer(40, 10)
	area := tui.NewRect(0, 0, 40, 10)
	s.Render(buf, area, newState())
	out := bufString(buf)
	if !strings.Contains(out, "Name") || !strings.Contains(out, "VRAM") {
		t.Fatalf("expected headers in render, got: %q", out)
	}
}

func TestTableStep_RenderRows(t *testing.T) {
	s := &TableStep{
		BaseStep: BaseStep{TitleText: "GPUs"},
		Headers:  []string{"Name", "VRAM"},
		Rows:     [][]string{{"A100", "80GB"}, {"RTX 4090", "24GB"}},
	}
	s.Init(newState())
	buf := tui.NewBuffer(40, 10)
	area := tui.NewRect(0, 0, 40, 10)
	s.Render(buf, area, newState())
	out := bufString(buf)
	if !strings.Contains(out, "A100") || !strings.Contains(out, "RTX 4090") {
		t.Fatalf("expected rows in render, got: %q", out)
	}
}

func TestTableStep_SelectableNavigation(t *testing.T) {
	s := &TableStep{
		BaseStep:   BaseStep{TitleText: "GPUs"},
		Headers:    []string{"Name"},
		Rows:       [][]string{{"A100"}, {"RTX 4090"}, {"RTX 3090"}},
		Selectable: true,
		ResultKey:  "gpu",
	}
	s.Init(newState())
	step, _ := s.Update(tui.KeyMsg{Type: tui.KeyDown}, newState())
	ts := step.(*TableStep)
	if ts.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", ts.cursor)
	}
}

func TestTableStep_SelectableEnter(t *testing.T) {
	s := &TableStep{
		BaseStep:   BaseStep{TitleText: "GPUs"},
		Headers:    []string{"Name"},
		Rows:       [][]string{{"A100"}, {"RTX 4090"}},
		Selectable: true,
		ResultKey:  "gpu",
	}
	s.Init(newState())
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil {
		t.Fatal("expected command on Enter for selectable table")
	}
}

func TestTableStep_NonSelectableEnterAdvances(t *testing.T) {
	s := &TableStep{
		BaseStep: BaseStep{TitleText: "GPUs"},
		Headers:  []string{"Name"},
		Rows:     [][]string{{"A100"}},
	}
	s.Init(newState())
	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
}

func TestTableStep_DynamicData(t *testing.T) {
	s := &TableStep{
		BaseStep:    BaseStep{TitleText: "GPUs"},
		HeadersFunc: func(state State) []string { return []string{"GPU"} },
		RowsFunc:    func(state State) [][]string { return [][]string{{"Dynamic GPU"}} },
	}
	s.Init(newState())
	if len(s.rows) != 1 || s.rows[0][0] != "Dynamic GPU" {
		t.Fatalf("expected dynamic rows, got %v", s.rows)
	}
}
