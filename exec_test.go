package gwiz

import (
	"context"
	"errors"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestExecStep_InitStartsTask(t *testing.T) {
	ran := false
	s := &ExecStep{
		BaseStep: BaseStep{TitleText: "Running"},
		TaskFunc: func(ctx context.Context, state State, output chan<- string) error {
			ran = true
			output <- "done"
			return nil
		},
	}
	cmd := s.Init(newState())
	if cmd == nil {
		t.Fatal("expected Init to return a command")
	}
	if !s.running {
		t.Fatal("expected running=true after Init")
	}
	_ = ran
}

func TestExecStep_ExecDoneMsgSuccess(t *testing.T) {
	s := &ExecStep{
		BaseStep: BaseStep{TitleText: "Running"},
		TaskFunc: func(ctx context.Context, state State, output chan<- string) error {
			return nil
		},
	}
	s.Init(newState())

	step, _ := s.Update(ExecDoneMsg{Err: nil}, newState())
	es := step.(*ExecStep)
	if es.running {
		t.Fatal("should not be running after done")
	}
	if es.failed {
		t.Fatal("should not be failed on success")
	}
}

func TestExecStep_ExecDoneMsgFailure(t *testing.T) {
	s := &ExecStep{
		BaseStep: BaseStep{TitleText: "Running"},
		TaskFunc: func(ctx context.Context, state State, output chan<- string) error {
			return errors.New("boom")
		},
	}
	s.Init(newState())

	step, _ := s.Update(ExecDoneMsg{Err: errors.New("boom")}, newState())
	es := step.(*ExecStep)
	if es.running {
		t.Fatal("should not be running after done")
	}
	if !es.failed {
		t.Fatal("should be failed")
	}
	if es.err == nil || es.err.Error() != "boom" {
		t.Fatalf("expected error 'boom', got %v", es.err)
	}
}

func TestExecStep_ExecOutputMsg(t *testing.T) {
	s := &ExecStep{
		BaseStep: BaseStep{TitleText: "Running"},
		TaskFunc: func(ctx context.Context, state State, output chan<- string) error {
			return nil
		},
	}
	s.Init(newState())

	step, _ := s.Update(ExecOutputMsg{Line: "Installing..."}, newState())
	es := step.(*ExecStep)
	if len(es.outputLines) != 1 || es.outputLines[0] != "Installing..." {
		t.Fatalf("expected output line, got %v", es.outputLines)
	}
}

func TestExecStep_DefaultFailureActions(t *testing.T) {
	s := &ExecStep{
		BaseStep: BaseStep{TitleText: "Running"},
		TaskFunc: func(ctx context.Context, state State, output chan<- string) error {
			return nil
		},
	}
	s.Init(newState())

	expected := ExecRetry | ExecRetryFrom | ExecAbort
	if s.FailureActions != expected {
		t.Fatalf("expected default failure actions %d, got %d", expected, s.FailureActions)
	}
}

func TestExecStep_MultiStepProgress(t *testing.T) {
	s := &ExecStep{
		BaseStep: BaseStep{TitleText: "Installing"},
		StepsFunc: func(state State) []ExecTask {
			return []ExecTask{
				{Label: "Step 1", Run: func(ctx context.Context, output chan<- string) error {
					output <- "step 1 done"
					return nil
				}},
				{Label: "Step 2", Run: func(ctx context.Context, output chan<- string) error {
					output <- "step 2 done"
					return nil
				}},
			}
		},
	}
	cmd := s.Init(newState())
	if cmd == nil {
		t.Fatal("expected Init to return a command")
	}
	if len(s.tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(s.tasks))
	}
}

func TestExecStep_EnterAfterSuccess(t *testing.T) {
	s := &ExecStep{
		BaseStep: BaseStep{TitleText: "Running"},
		TaskFunc: func(ctx context.Context, state State, output chan<- string) error {
			return nil
		},
	}
	s.Init(newState())
	s.running = false
	s.completed = true

	_, cmd := s.Update(tui.KeyMsg{Type: tui.KeyEnter}, newState())
	if cmd == nil {
		t.Fatal("expected NextMsg command after success + Enter")
	}
}
