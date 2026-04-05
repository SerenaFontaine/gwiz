---
title: Messages
weight: 5
---

Messages drive the wizard's state machine. Steps emit them as commands to trigger navigation and state changes.

## Type Aliases

These types are re-exported from `tui` so consumers don't need to import it directly:

```go
type Cmd = tui.Cmd
type Msg = tui.Msg
type Buffer = tui.Buffer
type Rect = tui.Rect
```

## Navigation Messages

| Type | Description |
|------|-------------|
| `NextMsg` | Advance to the next step (triggers validation first) |
| `PrevMsg` | Go back to the previous step |
| `QuitMsg` | Abort the wizard |

### Usage in Steps

```go
// Advance on Enter
case tui.KeyEnter:
    return s, tui.Batch(
        func() Msg { return StepResultMsg{Key: s.ResultKey, Value: s.value} },
        func() Msg { return NextMsg{} },
    )

// Go back on Esc
case tui.KeyEscape:
    return s, func() Msg { return PrevMsg{} }
```

## State Messages

```go
type StepResultMsg struct {
    Key   string
    Value any
}
```

Stores a key-value pair in the wizard's state. Typically sent alongside `NextMsg` in a `tui.Batch`.

## Error Messages

```go
type ErrorMsg struct {
    Err error
}
```

Displays an error on the current step without navigating.

## Exec Messages

| Type | Description |
|------|-------------|
| `ExecOutputMsg` | Appends a line to the exec step's output viewport |
| `ExecDoneMsg` | Signals that an exec step's task has finished |

```go
type ExecOutputMsg struct {
    Line string
}

type ExecDoneMsg struct {
    Err error
}
```

## Commands

Commands are functions that return a `Msg`. The framework runs them asynchronously and feeds the result back through `Update`.

```go
// Single command
return s, func() Msg { return NextMsg{} }

// Multiple commands (run concurrently)
return s, tui.Batch(
    func() Msg { return StepResultMsg{Key: "name", Value: "Alice"} },
    func() Msg { return NextMsg{} },
)
```
