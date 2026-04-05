---
title: Getting Started
weight: 10
---

## Installation

Add gwiz to your Go module:

```bash
go get github.com/SerenaFontaine/gwiz
```

## Requirements

- Go 1.23 or later
- Any terminal with ANSI escape sequence support

## Quick Start

### Minimal Wizard

A wizard collects user input across a sequence of steps. Create a `Wizard`, add steps, and call `Run`:

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

### With Selection and Form

```go
w.AddStep("backend", &gwiz.SelectStep{
    BaseStep: gwiz.BaseStep{TitleText: "Backend"},
    Prompt:   "Choose your backend:",
    Options: []gwiz.Option{
        {Label: "vLLM", Value: "vllm", Description: "High-throughput inference"},
        {Label: "Ollama", Value: "ollama", Description: "Easy model runner"},
    },
    ResultKey: "backend",
})

w.AddStep("config", &gwiz.FormStep{
    BaseStep: gwiz.BaseStep{TitleText: "Configuration"},
    Fields: []gwiz.FormField{
        {Label: "API Port", Key: "api_port", Default: "8000"},
        {Label: "Host", Key: "host", Default: "0.0.0.0"},
    },
})
```

## How It Works

gwiz follows a linear step-based flow:

1. **Create** — `gwiz.New()` creates a wizard with options (title, theme, banner)
2. **Add steps** — `AddStep()` registers named steps in order
3. **Run** — `Run()` starts the TUI; steps are presented one at a time
4. **Navigate** — Users press Enter to advance, Esc to go back, q/Ctrl+C to quit
5. **Collect** — Each step stores results in a shared `State` key-value bag
6. **Result** — `Run()` returns the accumulated state and whether the wizard was aborted

Steps can be conditionally skipped based on prior answers by implementing the `Skippable` method. State values set by earlier steps are available to later steps via the `State` interface.
