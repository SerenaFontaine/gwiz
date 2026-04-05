package gwiz

// registeredStep pairs a name with a step, in insertion order.
type registeredStep struct {
	name string
	step Step
}

// nextStep returns the index of the next non-skippable step after current,
// or -1 if the wizard is complete.
func nextStep(current int, steps []registeredStep, state State) int {
	for i := current + 1; i < len(steps); i++ {
		if !steps[i].step.Skippable(state) {
			return i
		}
	}
	return -1
}

// prevStep returns the index of the previous non-skippable step before current,
// or -1 if at the first step.
func prevStep(current int, steps []registeredStep, state State) int {
	for i := current - 1; i >= 0; i-- {
		if !steps[i].step.Skippable(state) {
			return i
		}
	}
	return -1
}

// firstStep returns the index of the first non-skippable step, or -1 if all are skippable.
func firstStep(steps []registeredStep, state State) int {
	for i := 0; i < len(steps); i++ {
		if !steps[i].step.Skippable(state) {
			return i
		}
	}
	return -1
}
