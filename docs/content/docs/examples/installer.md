---
title: Installer
weight: 2
---

A multi-step installer demonstrating all major step types: selection, multi-select, form, confirmation, and long-running task execution.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/SerenaFontaine/gwiz"
)

func main() {
    w := gwiz.New(
        gwiz.WithTitle("Mock Installer"),
        gwiz.WithTheme(gwiz.ThemeDracula),
    )

    w.AddStep("welcome", gwiz.InfoStep{
        BaseStep: gwiz.BaseStep{TitleText: "Welcome"},
        Content:  "This installer will walk you through a {cyan}mock installation{/}.\n\nNo actual software will be installed.",
    })

    w.AddStep("backend", &gwiz.SelectStep{
        BaseStep: gwiz.BaseStep{TitleText: "Backend"},
        Prompt:   "Choose your backend:",
        Options: []gwiz.Option{
            {Label: "vLLM", Value: "vllm", Description: "High-throughput inference engine"},
            {Label: "llama.cpp", Value: "llamacpp", Description: "CPU/GPU inference with GGUF models"},
            {Label: "Ollama", Value: "ollama", Description: "Easy-to-use model runner"},
        },
        ResultKey: "backend",
    })

    w.AddStep("services", &gwiz.MultiSelectStep{
        BaseStep: gwiz.BaseStep{TitleText: "Services"},
        Prompt:   "Select services to install:",
        Options: []gwiz.Option{
            {Label: "API Gateway", Value: "gateway"},
            {Label: "Web UI", Value: "webui"},
            {Label: "Monitoring", Value: "monitoring"},
        },
        ResultKey: "services",
        MinSelect: 1,
    })

    w.AddStep("config", &gwiz.FormStep{
        BaseStep: gwiz.BaseStep{TitleText: "Configuration"},
        Fields: []gwiz.FormField{
            {Label: "API Port", Key: "api_port", Default: "8000"},
            {Label: "Host", Key: "host", Default: "0.0.0.0"},
        },
    })

    w.AddStep("confirm", &gwiz.ConfirmStep{
        BaseStep: gwiz.BaseStep{TitleText: "Review"},
        SummaryFunc: func(state gwiz.State) string {
            return fmt.Sprintf(
                "Backend:  %s\nServices: %v\nAPI Port: %s\nHost:     %s",
                state.GetString("backend"),
                state.GetStringSlice("services"),
                state.GetString("api_port"),
                state.GetString("host"),
            )
        },
        ConfirmLabel: "Install",
        ResultKey:    "confirmed",
    })

    w.AddStep("install", &gwiz.ExecStep{
        BaseStep: gwiz.BaseStep{TitleText: "Installing"},
        StepsFunc: func(state gwiz.State) []gwiz.ExecTask {
            return []gwiz.ExecTask{
                {Label: "Downloading " + state.GetString("backend"), Run: func(ctx context.Context, output chan<- string) error {
                    for i := 0; i < 5; i++ {
                        output <- fmt.Sprintf("Downloading chunk %d/5...", i+1)
                        time.Sleep(200 * time.Millisecond)
                    }
                    return nil
                }},
                {Label: "Installing services", Run: func(ctx context.Context, output chan<- string) error {
                    for _, svc := range state.GetStringSlice("services") {
                        output <- fmt.Sprintf("Installing %s...", svc)
                        time.Sleep(300 * time.Millisecond)
                    }
                    return nil
                }},
                {Label: "Configuring", Run: func(ctx context.Context, output chan<- string) error {
                    output <- fmt.Sprintf("Setting API port to %s", state.GetString("api_port"))
                    output <- fmt.Sprintf("Binding to %s", state.GetString("host"))
                    time.Sleep(200 * time.Millisecond)
                    return nil
                }},
            }
        },
    })

    w.AddStep("done", gwiz.InfoStep{
        BaseStep:      gwiz.BaseStep{TitleText: "Complete"},
        ContinueLabel: "Exit",
        ContentFunc: func(state gwiz.State) string {
            return fmt.Sprintf("{green}Installation complete!{/}\n\nBackend %s is ready on %s:%s",
                state.GetString("backend"),
                state.GetString("host"),
                state.GetString("api_port"),
            )
        },
    })

    result, err := w.Run(context.Background())
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    if result.Aborted {
        fmt.Println("Installation cancelled.")
        os.Exit(0)
    }
    fmt.Println("Installation complete!")
}
```

## Step-by-Step Breakdown

| Step | Type | Purpose |
|------|------|---------|
| welcome | `InfoStep` | Displays a welcome message with colored markup |
| backend | `SelectStep` | Single selection with descriptions |
| services | `MultiSelectStep` | Toggle multiple items, requires at least 1 |
| config | `FormStep` | Two labeled text fields with defaults |
| confirm | `ConfirmStep` | Shows summary from state, waits for confirmation |
| install | `ExecStep` | Runs three sub-tasks with live output |
| done | `InfoStep` | Dynamic completion message from state |
