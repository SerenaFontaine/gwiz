package gwiz

import "testing"

func stubStep(title string, skippable bool) Step {
	return &skipStub{BaseStep: BaseStep{TitleText: title}, skip: skippable}
}

type skipStub struct {
	BaseStep
	skip bool
}

func (s *skipStub) Skippable(state State) bool { return s.skip }

func TestNextStep_Linear(t *testing.T) {
	steps := []registeredStep{
		{name: "a", step: stubStep("A", false)},
		{name: "b", step: stubStep("B", false)},
		{name: "c", step: stubStep("C", false)},
	}
	s := newState()
	got := nextStep(0, steps, s)
	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestNextStep_SkipsSkippable(t *testing.T) {
	steps := []registeredStep{
		{name: "a", step: stubStep("A", false)},
		{name: "b", step: stubStep("B", true)},
		{name: "c", step: stubStep("C", false)},
	}
	s := newState()
	got := nextStep(0, steps, s)
	if got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestNextStep_AllSkippable(t *testing.T) {
	steps := []registeredStep{
		{name: "a", step: stubStep("A", false)},
		{name: "b", step: stubStep("B", true)},
		{name: "c", step: stubStep("C", true)},
	}
	s := newState()
	got := nextStep(0, steps, s)
	if got != -1 {
		t.Fatalf("expected -1 (done), got %d", got)
	}
}

func TestNextStep_AtEnd(t *testing.T) {
	steps := []registeredStep{
		{name: "a", step: stubStep("A", false)},
	}
	s := newState()
	got := nextStep(0, steps, s)
	if got != -1 {
		t.Fatalf("expected -1 (done), got %d", got)
	}
}

func TestPrevStep_Linear(t *testing.T) {
	steps := []registeredStep{
		{name: "a", step: stubStep("A", false)},
		{name: "b", step: stubStep("B", false)},
		{name: "c", step: stubStep("C", false)},
	}
	s := newState()
	got := prevStep(2, steps, s)
	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestPrevStep_SkipsSkippable(t *testing.T) {
	steps := []registeredStep{
		{name: "a", step: stubStep("A", false)},
		{name: "b", step: stubStep("B", true)},
		{name: "c", step: stubStep("C", false)},
	}
	s := newState()
	got := prevStep(2, steps, s)
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestPrevStep_AtStart(t *testing.T) {
	steps := []registeredStep{
		{name: "a", step: stubStep("A", false)},
	}
	s := newState()
	got := prevStep(0, steps, s)
	if got != -1 {
		t.Fatalf("expected -1, got %d", got)
	}
}

func TestFirstStep_SkipsSkippable(t *testing.T) {
	steps := []registeredStep{
		{name: "a", step: stubStep("A", true)},
		{name: "b", step: stubStep("B", true)},
		{name: "c", step: stubStep("C", false)},
	}
	s := newState()
	got := firstStep(steps, s)
	if got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}
