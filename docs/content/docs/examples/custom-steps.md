---
title: Custom Steps
weight: 3
---

Create custom step types by implementing the `Step` interface. Embed `BaseStep` to get sensible defaults and only override what you need.

## Basic Custom Step

```go
type LicenseStep struct {
    gwiz.BaseStep
    LicenseText string
    ResultKey   string
    accepted    bool
}

func (s *LicenseStep) Init(state gwiz.State) gwiz.Cmd {
    s.accepted = false
    return nil
}

func (s *LicenseStep) Update(msg gwiz.Msg, state gwiz.State) (gwiz.Step, gwiz.Cmd) {
    if km, ok := msg.(tui.KeyMsg); ok {
        switch km.Type {
        case tui.KeyRune:
            if km.Rune == 'a' {
                s.accepted = true
                return s, tui.Batch(
                    func() gwiz.Msg { return gwiz.StepResultMsg{Key: s.ResultKey, Value: true} },
                    func() gwiz.Msg { return gwiz.NextMsg{} },
                )
            }
        case tui.KeyEscape:
            return s, func() gwiz.Msg { return gwiz.PrevMsg{} }
        }
    }
    return s, nil
}

func (s *LicenseStep) Render(buf *gwiz.Buffer, area gwiz.Rect, state gwiz.State) {
    buf.SetString(area.X, area.Y, s.LicenseText, tui.NewStyle())
    buf.SetString(area.X, area.Y+3, "Press 'a' to accept", tui.NewStyle().Bold(true))
}

func (s *LicenseStep) KeyHints() []gwiz.KeyHint {
    return []gwiz.KeyHint{
        {Key: "a", Label: "Accept"},
        {Key: "Esc", Label: "Back"},
    }
}
```

## Conditional Skipping

Override `Skippable` to skip a step based on prior state:

```go
type VLLMConfigStep struct {
    gwiz.BaseStep
    // ...
}

func (s *VLLMConfigStep) Skippable(state gwiz.State) bool {
    return state.GetString("backend") != "vllm"
}
```

This step is only shown when the user selected "vllm" in a previous step.

## Text Input Detection

If your custom step accepts free-form text input, implement `AcceptsTextInput` to prevent the wizard from intercepting the 'q' key as a quit shortcut:

```go
func (s *LicenseStep) AcceptsTextInput() bool { return false }
```
