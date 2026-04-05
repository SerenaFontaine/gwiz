package gwiz

import (
	"context"
	"fmt"
	"strings"

	"github.com/SerenaFontaine/tui"
)

// ExecTask is a single sub-task within a multi-step ExecStep.
type ExecTask struct {
	Label string
	Run   func(ctx context.Context, output chan<- string) error
}

// ExecStep runs a long-running task with live output.
type ExecStep struct {
	BaseStep
	TaskFunc       func(ctx context.Context, state State, output chan<- string) error
	StepsFunc      func(state State) []ExecTask
	FailureActions ExecFailureAction

	// internal state
	running      bool
	completed    bool
	failed       bool
	err          error
	outputLines  []string
	tasks        []ExecTask
	currentTask  int
	tasksDone    int
	failedTask   int

	// failure menu
	failureSelected int

	// references needed for retry
	state  State
	cancel context.CancelFunc
}

func (s *ExecStep) Init(state State) Cmd {
	s.state = state
	s.outputLines = nil
	s.completed = false
	s.failed = false
	s.err = nil
	s.currentTask = 0
	s.tasksDone = 0
	s.failedTask = -1
	s.failureSelected = 0

	if s.FailureActions == 0 {
		s.FailureActions = ExecRetry | ExecRetryFrom | ExecAbort
	}

	if s.StepsFunc != nil {
		s.tasks = s.StepsFunc(state)
	}

	s.running = true
	return s.startExecution(state, 0)
}

func (s *ExecStep) startExecution(state State, fromTask int) Cmd {
	s.running = true
	s.completed = false
	s.failed = false
	s.err = nil
	s.currentTask = fromTask
	s.tasksDone = fromTask
	s.failedTask = -1

	return func() Msg {
		ctx, cancel := context.WithCancel(context.Background())
		s.cancel = cancel
		output := make(chan string, 100)

		go func() {
			defer close(output)
			defer cancel()

			if s.TaskFunc != nil {
				err := s.TaskFunc(ctx, state, output)
				if err != nil {
					output <- fmt.Sprintf("FAILED: %v", err)
				}
				return
			}

			// Multi-step execution
			for i := fromTask; i < len(s.tasks); i++ {
				task := s.tasks[i]
				output <- fmt.Sprintf("--- %s ---", task.Label)
				taskOutput := make(chan string, 100)
				done := make(chan struct{})
				go func() {
					defer close(done)
					for line := range taskOutput {
						output <- line
					}
				}()
				err := task.Run(ctx, taskOutput)
				close(taskOutput)
				<-done // wait for forwarding goroutine
				if err != nil {
					output <- fmt.Sprintf("FAILED: %s: %v", task.Label, err)
					return
				}
				output <- fmt.Sprintf("OK: %s", task.Label)
			}
		}()

		// Collect output
		var lines []string
		for line := range output {
			lines = append(lines, line)
		}

		// Determine result
		var execErr error
		for _, line := range lines {
			if strings.HasPrefix(line, "FAILED:") {
				execErr = fmt.Errorf("%s", line)
				break
			}
		}

		return execBatchResult{lines: lines, err: execErr}
	}
}

// execBatchResult is an internal message carrying collected output.
type execBatchResult struct {
	lines []string
	err   error
}

func (s *ExecStep) Update(msg Msg, state State) (Step, Cmd) {
	switch msg := msg.(type) {
	case execBatchResult:
		s.outputLines = append(s.outputLines, msg.lines...)
		if msg.err != nil {
			s.running = false
			s.failed = true
			s.err = msg.err
		} else {
			s.running = false
			s.completed = true
		}
		return s, nil

	case ExecOutputMsg:
		s.outputLines = append(s.outputLines, msg.Line)
		return s, nil

	case ExecDoneMsg:
		s.running = false
		if msg.Err != nil {
			s.failed = true
			s.err = msg.Err
		} else {
			s.completed = true
		}
		return s, nil

	case tui.KeyMsg:
		if s.running {
			return s, nil
		}
		if s.completed {
			if msg.Type == tui.KeyEnter {
				return s, func() Msg { return NextMsg{} }
			}
		}
		if s.failed {
			return s.handleFailureInput(msg, state)
		}
	}
	return s, nil
}

func (s *ExecStep) handleFailureInput(msg tui.KeyMsg, state State) (Step, Cmd) {
	actions := s.availableActions()
	switch msg.Type {
	case tui.KeyTab, tui.KeyRight:
		s.failureSelected = (s.failureSelected + 1) % len(actions)
	case tui.KeyBacktab, tui.KeyLeft:
		s.failureSelected = (s.failureSelected - 1 + len(actions)) % len(actions)
	case tui.KeyEnter:
		if len(actions) == 0 {
			return s, nil
		}
		action := actions[s.failureSelected]
		switch action {
		case ExecRetry:
			s.outputLines = nil
			return s, s.startExecution(state, 0)
		case ExecRetryFrom:
			return s, s.startExecution(state, s.failedTask)
		case ExecAbort:
			return s, func() Msg { return QuitMsg{} }
		}
	}
	return s, nil
}

func (s *ExecStep) availableActions() []ExecFailureAction {
	var actions []ExecFailureAction
	if s.FailureActions&ExecRetry != 0 {
		actions = append(actions, ExecRetry)
	}
	if s.FailureActions&ExecRetryFrom != 0 && len(s.tasks) > 0 {
		actions = append(actions, ExecRetryFrom)
	}
	if s.FailureActions&ExecAbort != 0 {
		actions = append(actions, ExecAbort)
	}
	return actions
}

func (s *ExecStep) Render(buf *tui.Buffer, area tui.Rect, state State) {
	y := area.Y

	// Progress indicator for multi-step
	if len(s.tasks) > 0 && s.running {
		progress := fmt.Sprintf("Step %d/%d", s.currentTask+1, len(s.tasks))
		if s.currentTask < len(s.tasks) {
			progress += ": " + s.tasks[s.currentTask].Label
		}
		buf.SetString(area.X, y, progress, tui.NewStyle().Bold(true))
		y += 2
	}

	// Spinner if running
	if s.running {
		buf.SetString(area.X, y, "⠋ Running...", tui.NewStyle())
		y++
	}

	// Output lines (show last N lines that fit)
	outputAreaHeight := area.Y + area.Height - y - 3
	if outputAreaHeight < 1 {
		outputAreaHeight = 1
	}
	startLine := 0
	if len(s.outputLines) > outputAreaHeight {
		startLine = len(s.outputLines) - outputAreaHeight
	}
	for i := startLine; i < len(s.outputLines); i++ {
		if y >= area.Y+area.Height-2 {
			break
		}
		buf.SetString(area.X, y, s.outputLines[i], tui.NewStyle().Dim(true))
		y++
	}

	// Status line
	statusY := area.Y + area.Height - 2
	if s.completed {
		buf.SetString(area.X, statusY, "Done! Press Enter to continue.", tui.NewStyle().Bold(true))
	} else if s.failed {
		buf.SetString(area.X, statusY, "Error: "+s.err.Error(), tui.NewStyle().Bold(true))
		btnY := statusY + 1
		actions := s.availableActions()
		cx := area.X
		for i, action := range actions {
			label := failureActionLabel(action)
			style := tui.NewStyle().Dim(true)
			if i == s.failureSelected {
				style = tui.NewStyle().Bold(true).Reverse(true)
			}
			btn := "[ " + label + " ]"
			cx += buf.SetString(cx, btnY, btn, style)
			cx += 2
		}
	}
}

func failureActionLabel(a ExecFailureAction) string {
	switch a {
	case ExecRetry:
		return "Retry All"
	case ExecRetryFrom:
		return "Retry From Failed"
	case ExecAbort:
		return "Abort"
	default:
		return "?"
	}
}

func (s *ExecStep) KeyHints() []KeyHint {
	if s.running {
		return []KeyHint{}
	}
	if s.completed {
		return []KeyHint{{Key: "Enter", Label: "Continue"}}
	}
	return []KeyHint{
		{Key: "Tab", Label: "Switch"},
		{Key: "Enter", Label: "Select"},
	}
}

func (s *ExecStep) Skippable(state State) bool { return false }
