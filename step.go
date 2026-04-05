package gwiz

import (
	"github.com/SerenaFontaine/tui"
)

// Type aliases so consumers don't need to import tui directly.
type Cmd = tui.Cmd
type Msg = tui.Msg
type Buffer = tui.Buffer
type Rect = tui.Rect

// Step is the interface that all wizard steps implement.
type Step interface {
	Title() string
	Description() string
	Init(state State) Cmd
	Update(msg Msg, state State) (Step, Cmd)
	Render(buf *Buffer, area Rect, state State)
	Validate(state State) error
	Skippable(state State) bool
}

// BaseStep provides default implementations for the Step interface.
// Embed this in custom steps to only override what you need.
type BaseStep struct {
	TitleText       string
	DescriptionText string
}

func (b BaseStep) Title() string                           { return b.TitleText }
func (b BaseStep) Description() string                     { return b.DescriptionText }
func (b BaseStep) Init(state State) Cmd                    { return nil }
func (b BaseStep) Update(msg Msg, state State) (Step, Cmd) { return b, nil }
func (b BaseStep) Render(buf *Buffer, area Rect, state State) {}
func (b BaseStep) Validate(state State) error              { return nil }
func (b BaseStep) Skippable(state State) bool              { return false }

// KeyHints returns the default keybindings shown in the nav bar.
func (b BaseStep) KeyHints() []KeyHint {
	return []KeyHint{
		{Key: "Enter", Label: "Next"},
		{Key: "Esc", Label: "Back"},
		{Key: "q", Label: "Quit"},
	}
}

// KeyHint describes a keybinding shown in the nav bar.
type KeyHint struct {
	Key   string
	Label string
}

// WizardOption configures a Wizard.
type WizardOption func(*Wizard)

// Wizard-specific message types.

// NextMsg signals the wizard to advance to the next step.
type NextMsg struct{}

// PrevMsg signals the wizard to go back to the previous step.
type PrevMsg struct{}

// QuitMsg signals the wizard to abort.
type QuitMsg struct{}

// ErrorMsg displays an error on the current step.
type ErrorMsg struct {
	Err error
}

// StepResultMsg is a convenience message to set a state value.
type StepResultMsg struct {
	Key   string
	Value any
}

// ExecOutputMsg appends a line to the exec step's output viewport.
type ExecOutputMsg struct {
	Line string
}

// ExecDoneMsg signals that an exec step's task has finished.
type ExecDoneMsg struct {
	Err error
}

// ExecFailureAction controls what options are shown when an ExecStep fails.
type ExecFailureAction int

const (
	ExecRetry     ExecFailureAction = 1 << iota // Restart all tasks
	ExecRetryFrom                                // Restart from failed task
	ExecAbort                                    // Abort the wizard
)

// Result contains the outcome of a wizard run.
type Result struct {
	State   State
	Aborted bool
}
