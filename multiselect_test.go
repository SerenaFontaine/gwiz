package gwiz

import (
	"slices"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestMultiSelectStep_ToggleAndConfirm(t *testing.T) {
	s := &MultiSelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}, {Label: "C", Value: "c"}},
		ResultKey: "items",
	}
	st := newState()
	s.Init(st)

	// Toggle A (index 0)
	step, _ := s.Update(tui.KeyMsg{Type: tui.KeySpace}, st)
	ms := step.(*MultiSelectStep)
	if !ms.selected[0] {
		t.Fatal("expected item 0 to be selected")
	}

	// Move down twice and toggle C (index 2)
	step, _ = ms.Update(tui.KeyMsg{Type: tui.KeyDown}, st)
	ms = step.(*MultiSelectStep)
	step, _ = ms.Update(tui.KeyMsg{Type: tui.KeyDown}, st)
	ms = step.(*MultiSelectStep)
	step, _ = ms.Update(tui.KeyMsg{Type: tui.KeySpace}, st)
	ms = step.(*MultiSelectStep)

	// Confirm
	_, cmd := ms.Update(tui.KeyMsg{Type: tui.KeyEnter}, st)
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
}

func TestMultiSelectStep_ValidateMinSelect(t *testing.T) {
	s := &MultiSelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}},
		ResultKey: "items",
		MinSelect: 1,
	}
	s.Init(newState())

	err := s.Validate(newState())
	if err == nil {
		t.Fatal("expected validation error with 0 selections and MinSelect=1")
	}
}

func TestMultiSelectStep_ValidateMaxSelect(t *testing.T) {
	s := &MultiSelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}, {Label: "C", Value: "c"}},
		ResultKey: "items",
		MaxSelect: 1,
	}
	st := newState()
	s.Init(st)

	// Toggle A and try to toggle B
	step, _ := s.Update(tui.KeyMsg{Type: tui.KeySpace}, st)
	ms := step.(*MultiSelectStep)
	step, _ = ms.Update(tui.KeyMsg{Type: tui.KeyDown}, st)
	ms = step.(*MultiSelectStep)
	step, _ = ms.Update(tui.KeyMsg{Type: tui.KeySpace}, st)
	ms = step.(*MultiSelectStep)

	count := 0
	for _, v := range ms.selected {
		if v {
			count++
		}
	}
	if count > 1 {
		t.Fatal("should not allow more than MaxSelect selections")
	}
}

func TestMultiSelectStep_SkipsDisabled(t *testing.T) {
	s := &MultiSelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a", Disabled: true}, {Label: "B", Value: "b"}},
		ResultKey: "items",
	}
	s.Init(newState())

	step, _ := s.Update(tui.KeyMsg{Type: tui.KeySpace}, newState())
	ms := step.(*MultiSelectStep)
	if ms.selected[0] {
		t.Fatal("should not toggle disabled item")
	}
}

func TestMultiSelectStep_RestoresFromState(t *testing.T) {
	s := &MultiSelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}, {Label: "C", Value: "c"}},
		ResultKey: "items",
	}
	st := newState()
	st.Set("items", []string{"a", "c"})
	s.Init(st)

	if !s.selected[0] || s.selected[1] || !s.selected[2] {
		t.Fatalf("expected [true false true], got %v", s.selected)
	}
}

func TestMultiSelectStep_SelectedValues(t *testing.T) {
	s := &MultiSelectStep{
		BaseStep:  BaseStep{TitleText: "Pick"},
		Options:   []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}, {Label: "C", Value: "c"}},
		ResultKey: "items",
	}
	s.Init(newState())
	s.selected[0] = true
	s.selected[2] = true

	vals := s.selectedValues()
	if !slices.Equal(vals, []string{"a", "c"}) {
		t.Fatalf("expected [a c], got %v", vals)
	}
}
