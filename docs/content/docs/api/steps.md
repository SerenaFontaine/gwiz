---
title: Steps
weight: 2
---

All wizard steps implement the `Step` interface. Embed `BaseStep` for default implementations.

## Step Interface

```go
type Step interface {
    Title() string
    Description() string
    Init(state State) Cmd
    Update(msg Msg, state State) (Step, Cmd)
    Render(buf *Buffer, area Rect, state State)
    Validate(state State) error
    Skippable(state State) bool
}
```

### Methods

| Method | Description |
|--------|-------------|
| `Title()` | Step title shown in the header area |
| `Description()` | Step description (used by some step types) |
| `Init(state)` | Called when the step becomes active; restore state here |
| `Update(msg, state)` | Handle messages (key presses, etc.) |
| `Render(buf, area, state)` | Draw the step into the buffer |
| `Validate(state)` | Validate before advancing; return error to block |
| `Skippable(state)` | Return `true` to skip this step during navigation |

## BaseStep

```go
type BaseStep struct {
    TitleText       string
    DescriptionText string
}
```

Provides no-op defaults for all `Step` methods. Embed it in custom steps to only override what you need.

## Optional Interfaces

Steps can optionally implement these interfaces for additional behavior:

| Interface | Method | Description |
|-----------|--------|-------------|
| `textAcceptor` | `AcceptsTextInput() bool` | Prevents 'q' from being intercepted as quit |
| `hinter` | `KeyHints() []KeyHint` | Custom keybinding hints in the nav bar |

## InputStep

Collects a single line of text from the user.

```go
type InputStep struct {
    BaseStep
    Prompt       string
    Placeholder  string
    Default      string
    ResultKey    string
    ValidateFunc func(value string) error
}
```

| Field | Description |
|-------|-------------|
| `Prompt` | Text displayed above the input field |
| `Placeholder` | Dimmed text shown when the field is empty |
| `Default` | Initial value |
| `ResultKey` | State key where the value is stored |
| `ValidateFunc` | Optional validation called before advancing |

**Keys:** Left/Right to move cursor, Home/End to jump, Enter to confirm, Esc to go back.

## SelectStep

Presents a single-selection list.

```go
type SelectStep struct {
    BaseStep
    Prompt      string
    Options     []Option
    OptionsFunc func(state State) []Option
    ResultKey   string
    Default     string
}
```

| Field | Description |
|-------|-------------|
| `Options` | Static list of options |
| `OptionsFunc` | Dynamic options generated from state (overrides `Options`) |
| `Default` | Value of the initially selected option |

**Keys:** Up/Down or j/k to move, Enter to confirm, Esc/Backspace to go back.

## MultiSelectStep

Presents a list of options where multiple items can be toggled.

```go
type MultiSelectStep struct {
    BaseStep
    Prompt      string
    Options     []Option
    OptionsFunc func(state State) []Option
    ResultKey   string
    MinSelect   int
    MaxSelect   int
}
```

| Field | Description |
|-------|-------------|
| `MinSelect` | Minimum required selections (validated before advancing) |
| `MaxSelect` | Maximum allowed selections (toggles blocked at limit) |

The result is stored as a `[]string` of selected values.

**Keys:** Up/Down or j/k to move, Space to toggle, Enter to confirm, Esc/Backspace to go back.

## FormStep

Presents multiple labeled text fields on a single screen.

```go
type FormStep struct {
    BaseStep
    Fields []FormField
}

type FormField struct {
    Label    string
    Key      string
    Default  string
    Validate func(value string) error
}
```

Each field's value is stored in state under its `Key`. Validation runs per-field before advancing.

**Keys:** Tab/Enter to next field, Shift+Tab to previous field, Esc to go back.

## InfoStep

Displays formatted text content. Used for welcome screens, summaries, warnings, or any read-only information.

```go
type InfoStep struct {
    BaseStep
    Content       string
    ContentFunc   func(state State) string
    ContinueLabel string
}
```

| Field | Description |
|-------|-------------|
| `Content` | Static text with inline markup |
| `ContentFunc` | Dynamic content from state (overrides `Content`) |
| `ContinueLabel` | Custom label for the Enter key hint (default: "Next") |

### Inline Markup

| Syntax | Effect |
|--------|--------|
| `**text**` | Bold |
| `*text*` | Dim |
| `{green}text{/}` | Colored (green, red, yellow, cyan, blue, magenta) |

## ConfirmStep

Displays a summary and asks the user to confirm or go back.

```go
type ConfirmStep struct {
    BaseStep
    SummaryFunc  func(state State) string
    ConfirmLabel string
    CancelLabel  string
    ResultKey    string
}
```

Stores `true` in state when confirmed. Cancel navigates back.

**Keys:** Tab/Left/Right to switch buttons, Enter to select, Esc to go back.

## TableStep

Displays tabular data with optional row selection.

```go
type TableStep struct {
    BaseStep
    Headers     []string
    Rows        [][]string
    HeadersFunc func(state State) []string
    RowsFunc    func(state State) [][]string
    Selectable  bool
    ResultKey   string
}
```

When `Selectable` is true, the user can navigate rows and the selected row is stored in state as a `[]string`.

**Keys:** Up/Down or j/k to move (when selectable), Enter to confirm, Esc/Backspace to go back.

## ExecStep

Runs long-running tasks with live output.

```go
type ExecStep struct {
    BaseStep
    TaskFunc       func(ctx context.Context, state State, output chan<- string) error
    StepsFunc      func(state State) []ExecTask
    FailureActions ExecFailureAction
}

type ExecTask struct {
    Label string
    Run   func(ctx context.Context, output chan<- string) error
}
```

| Field | Description |
|-------|-------------|
| `TaskFunc` | Single task with live output |
| `StepsFunc` | Multiple sub-tasks run sequentially (overrides `TaskFunc`) |
| `FailureActions` | Bitflag controlling retry options on failure |

### Failure Actions

```go
const (
    ExecRetry     ExecFailureAction = 1 << iota // Restart all tasks
    ExecRetryFrom                                // Restart from failed task
    ExecAbort                                    // Abort the wizard
)
```

Default: all three actions are enabled.

## Option

```go
type Option struct {
    Label       string
    Value       string
    Description string
    Disabled    bool
    DisabledMsg string
}
```

Used by `SelectStep` and `MultiSelectStep`. Disabled options are shown but cannot be selected.

## KeyHint

```go
type KeyHint struct {
    Key   string
    Label string
}
```

Describes a keybinding shown in the nav bar.
