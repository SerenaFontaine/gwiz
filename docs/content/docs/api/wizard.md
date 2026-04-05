---
title: Wizard
weight: 1
---

`Wizard` is the top-level orchestrator that implements `tui.Component`.

## Creating a Wizard

```go
func New(opts ...WizardOption) *Wizard
```

Creates a new wizard with the given options. The default theme is `ThemeNord`.

## Options

```go
type WizardOption func(*Wizard)
```

| Function | Description |
|----------|-------------|
| `WithTitle(string)` | Set the title displayed in the wizard header |
| `WithTheme(Theme)` | Set the color theme used for rendering |
| `WithBanner(string)` | Set subtitle text shown next to the title |

```go
w := gwiz.New(
    gwiz.WithTitle("Setup"),
    gwiz.WithTheme(gwiz.ThemeDracula),
    gwiz.WithBanner("v1.0"),
)
```

## Adding Steps

```go
func (w *Wizard) AddStep(name string, step Step)
```

Registers a named step. Steps are presented in the order they are added.

```go
w.AddStep("name", &gwiz.InputStep{
    BaseStep:  gwiz.BaseStep{TitleText: "Your Name"},
    Prompt:    "Enter your name:",
    ResultKey: "name",
})
```

## Running

```go
func (w *Wizard) Run(ctx context.Context) (*Result, error)
```

Starts the wizard in an alternate-screen TUI and blocks until the user completes or aborts. Returns the accumulated state and whether the wizard was aborted.

```go
result, err := w.Run(context.Background())
if err != nil {
    log.Fatal(err)
}
if result.Aborted {
    // user quit early
}
name := result.State.GetString("name")
```

## Result

```go
type Result struct {
    State   State
    Aborted bool
}
```

| Field | Description |
|-------|-------------|
| `State` | The accumulated key-value state from all completed steps |
| `Aborted` | `true` if the user quit before completing all steps |

## Navigation

The wizard handles navigation automatically:

- **Enter** — Validate current step and advance to next
- **Esc** — Go back to previous step
- **q** / **Ctrl+C** — Open quit confirmation dialog (or quit immediately if no state)

Steps that return `true` from `Skippable(state)` are automatically skipped during forward and backward navigation.

## Quit Dialog

When the user presses q or Ctrl+C and the wizard has accumulated state, a confirmation dialog appears. Tab switches between Cancel and Quit. If no state has been collected yet, the wizard exits immediately.
