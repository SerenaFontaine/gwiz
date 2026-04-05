# gwiz

[![Docs](https://img.shields.io/badge/docs-github_pages-2ea44f?logo=github)](https://serenafontaine.github.io/gwiz/)
[![License](https://img.shields.io/github/license/SerenaFontaine/gwiz)](https://github.com/SerenaFontaine/gwiz/blob/main/LICENSE)
[![Go Reference](https://img.shields.io/badge/go-reference-00ADD8?logo=go)](https://pkg.go.dev/github.com/SerenaFontaine/gwiz)
[![Docs Deploy](https://img.shields.io/github/actions/workflow/status/SerenaFontaine/gwiz/docs-pages.yml?label=docs%20deploy)](https://github.com/SerenaFontaine/gwiz/actions/workflows/docs-pages.yml)

A step-by-step terminal wizard framework for Go built on
[tui](https://github.com/SerenaFontaine/tui).

## Features

- **Step-by-step flow** — Guided wizard with forward/back navigation and conditional skipping
- **8 built-in step types** — Input, Select, MultiSelect, Form, Info, Confirm, Table, Exec
- **State management** — Typed key-value bag flows between steps automatically
- **Validation** — Per-step validation with inline error display
- **Dynamic content** — Steps can generate options, rows, and content from prior state
- **Long-running tasks** — ExecStep runs tasks with live output, progress, and retry on failure
- **Themes** — Built-in Nord, Gruvbox, Dracula, Monochrome themes
- **Inline markup** — Info steps support **bold**, *dim*, and {color}text{/} markup
- **Quit dialog** — Confirms before discarding progress
- **Extensible** — Implement the Step interface for custom step types

## Installation

```bash
go get github.com/SerenaFontaine/gwiz
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/SerenaFontaine/gwiz"
)

func main() {
    w := gwiz.New(
        gwiz.WithTitle("Hello"),
        gwiz.WithTheme(gwiz.ThemeNord),
    )

    w.AddStep("name", &gwiz.InputStep{
        BaseStep:  gwiz.BaseStep{TitleText: "Your Name"},
        Prompt:    "What is your name?",
        ResultKey: "name",
    })

    w.AddStep("confirm", &gwiz.ConfirmStep{
        BaseStep: gwiz.BaseStep{TitleText: "Confirm"},
        SummaryFunc: func(state gwiz.State) string {
            return fmt.Sprintf("Name: %s", state.GetString("name"))
        },
        ConfirmLabel: "Done",
        ResultKey:    "confirmed",
    })

    result, err := w.Run(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    if result.Aborted {
        fmt.Println("Cancelled.")
        return
    }
    fmt.Printf("Hello, %s!\n", result.State.GetString("name"))
}
```

## Usage Examples

### Text Input

```go
w.AddStep("email", &gwiz.InputStep{
    BaseStep:  gwiz.BaseStep{TitleText: "Email"},
    Prompt:    "Enter your email address:",
    ResultKey: "email",
    ValidateFunc: func(v string) error {
        if !strings.Contains(v, "@") {
            return fmt.Errorf("invalid email")
        }
        return nil
    },
})
```

### Single Select

```go
w.AddStep("backend", &gwiz.SelectStep{
    BaseStep: gwiz.BaseStep{TitleText: "Backend"},
    Prompt:   "Choose your backend:",
    Options: []gwiz.Option{
        {Label: "vLLM", Value: "vllm", Description: "High-throughput inference"},
        {Label: "llama.cpp", Value: "llamacpp", Description: "CPU/GPU with GGUF"},
        {Label: "Ollama", Value: "ollama", Description: "Easy model runner"},
    },
    ResultKey: "backend",
})
```

### Multi Select

```go
w.AddStep("services", &gwiz.MultiSelectStep{
    BaseStep:  gwiz.BaseStep{TitleText: "Services"},
    Prompt:    "Select services to install:",
    Options:   []gwiz.Option{
        {Label: "API Gateway", Value: "gateway"},
        {Label: "Web UI", Value: "webui"},
        {Label: "Monitoring", Value: "monitoring"},
    },
    ResultKey: "services",
    MinSelect: 1,
})
```

### Multi-Field Form

```go
w.AddStep("config", &gwiz.FormStep{
    BaseStep: gwiz.BaseStep{TitleText: "Configuration"},
    Fields: []gwiz.FormField{
        {Label: "API Port", Key: "api_port", Default: "8000"},
        {Label: "Host", Key: "host", Default: "0.0.0.0"},
    },
})
```

### Info Screen with Markup

```go
w.AddStep("welcome", gwiz.InfoStep{
    BaseStep: gwiz.BaseStep{TitleText: "Welcome"},
    Content:  "This is a **bold** message with {green}colored{/} text.",
})
```

### Long-Running Tasks

```go
w.AddStep("install", &gwiz.ExecStep{
    BaseStep: gwiz.BaseStep{TitleText: "Installing"},
    StepsFunc: func(state gwiz.State) []gwiz.ExecTask {
        return []gwiz.ExecTask{
            {Label: "Downloading", Run: func(ctx context.Context, out chan<- string) error {
                out <- "Fetching packages..."
                time.Sleep(2 * time.Second)
                return nil
            }},
        }
    },
})
```

### Dynamic Options from State

```go
w.AddStep("model", &gwiz.SelectStep{
    BaseStep: gwiz.BaseStep{TitleText: "Model"},
    Prompt:   "Select a model:",
    OptionsFunc: func(state gwiz.State) []gwiz.Option {
        switch state.GetString("backend") {
        case "vllm":
            return []gwiz.Option{
                {Label: "Llama 3 70B", Value: "llama3-70b"},
                {Label: "Mistral 7B", Value: "mistral-7b"},
            }
        default:
            return []gwiz.Option{
                {Label: "Default Model", Value: "default"},
            }
        }
    },
    ResultKey: "model",
})
```

### Conditional Skipping

Implement `Skippable` on a custom step to skip it based on prior answers:

```go
type OptionalStep struct {
    gwiz.BaseStep
    // ...
}

func (s *OptionalStep) Skippable(state gwiz.State) bool {
    return state.GetString("backend") != "vllm"
}
```

## Step Types

| Step | Description |
|------|-------------|
| `InputStep` | Single-line text input with cursor, placeholder, and validation |
| `SelectStep` | Single-selection list with arrow/vim keys |
| `MultiSelectStep` | Toggle multiple items with space, min/max constraints |
| `FormStep` | Multiple labeled text fields on one screen |
| `InfoStep` | Read-only formatted text with inline markup |
| `ConfirmStep` | Summary display with confirm/cancel buttons |
| `TableStep` | Tabular data display with optional row selection |
| `ExecStep` | Long-running tasks with live output and retry |

## API Reference

### Core

| Type | Description |
|------|-------------|
| `Wizard` | Top-level orchestrator implementing `tui.Component` |
| `Step` | Interface for wizard steps (7 methods) |
| `BaseStep` | Embeddable default `Step` implementation |
| `State` | Typed key-value bag accumulated across steps |
| `Result` | Outcome of a wizard run (state + aborted flag) |
| `Option` | Selectable item for `SelectStep` and `MultiSelectStep` |
| `KeyHint` | Keybinding label shown in the nav bar |

### Options

| Function | Description |
|----------|-------------|
| `WithTitle(string)` | Set the wizard header title |
| `WithTheme(Theme)` | Set the color theme |
| `WithBanner(string)` | Set subtitle text next to the title |

### Messages

| Type | Description |
|------|-------------|
| `NextMsg` | Advance to next step |
| `PrevMsg` | Go back to previous step |
| `QuitMsg` | Abort the wizard |
| `ErrorMsg` | Display an error on the current step |
| `StepResultMsg` | Store a key-value pair in state |
| `ExecOutputMsg` | Append output line in ExecStep |
| `ExecDoneMsg` | Signal ExecStep completion |

### Themes

| Variable | Description |
|----------|-------------|
| `ThemeNord` | Cool blue-tinted Nord palette |
| `ThemeGruvbox` | Warm retro Gruvbox palette |
| `ThemeDracula` | Dark purple Dracula palette |
| `ThemeMonochrome` | Minimal black-and-white |

## License

[MIT](LICENSE)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## References

- [tui — Terminal User Interface Framework](https://github.com/SerenaFontaine/tui)
