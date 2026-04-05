---
title: State
weight: 3
---

`State` is a typed key-value bag that accumulates across wizard steps.

## Interface

```go
type State interface {
    Get(key string) (any, bool)
    GetString(key string) string
    GetBool(key string) bool
    GetInt(key string) int
    GetStringSlice(key string) []string
    Set(key string, value any)
    Keys() []string
}
```

## Methods

| Method | Return | Description |
|--------|--------|-------------|
| `Get(key)` | `(any, bool)` | Raw value lookup with existence check |
| `GetString(key)` | `string` | Returns `""` if missing or not a string |
| `GetBool(key)` | `bool` | Returns `false` if missing or not a bool |
| `GetInt(key)` | `int` | Returns `0` if missing or not an int |
| `GetStringSlice(key)` | `[]string` | Returns `nil` if missing or not a string slice |
| `Set(key, value)` | — | Store any value |
| `Keys()` | `[]string` | All stored keys |

## How State Flows

1. When a step sends a `StepResultMsg`, the wizard calls `Set(key, value)` on the state
2. When a step becomes active, `Init(state)` is called — steps typically restore their previous value from state
3. The `SummaryFunc`, `ContentFunc`, `OptionsFunc`, and `RowsFunc` callbacks all receive state to generate dynamic content
4. After the wizard completes, `Result.State` contains all accumulated values

## Usage in Steps

### Reading State

```go
func (s *SelectStep) Init(state State) Cmd {
    // Restore previous selection
    if v, ok := state.Get(s.ResultKey); ok {
        if val, ok := v.(string); ok {
            for i, opt := range s.options {
                if opt.Value == val {
                    s.cursor = i
                    break
                }
            }
        }
    }
    return nil
}
```

### Writing State

Steps write to state by returning a `StepResultMsg` command:

```go
return s, func() Msg {
    return StepResultMsg{Key: s.ResultKey, Value: s.value}
}
```

### Dynamic Content from State

```go
w.AddStep("confirm", &gwiz.ConfirmStep{
    BaseStep: gwiz.BaseStep{TitleText: "Review"},
    SummaryFunc: func(state gwiz.State) string {
        return fmt.Sprintf("Backend: %s\nPort: %s",
            state.GetString("backend"),
            state.GetString("api_port"),
        )
    },
    ResultKey: "confirmed",
})
```
