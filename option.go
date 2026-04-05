package gwiz

// Option represents a selectable item in SelectStep and MultiSelectStep.
type Option struct {
	Label       string
	Value       string
	Description string
	Disabled    bool
	DisabledMsg string
}
