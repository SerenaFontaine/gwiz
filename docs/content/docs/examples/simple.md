---
title: Simple Wizard
weight: 1
---

A minimal wizard with an info screen, text input, and confirmation.

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/SerenaFontaine/gwiz"
)

func main() {
    w := gwiz.New(
        gwiz.WithTitle("Simple Wizard"),
        gwiz.WithTheme(gwiz.ThemeNord),
    )

    w.AddStep("welcome", gwiz.InfoStep{
        BaseStep: gwiz.BaseStep{TitleText: "Welcome"},
        Content:  "This is a **simple** wizard example.\n\nPress {green}Enter{/} to continue.",
    })

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
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    if result.Aborted {
        fmt.Println("Wizard aborted.")
        os.Exit(0)
    }
    fmt.Printf("Hello, %s!\n", result.State.GetString("name"))
}
```

## What's Happening

1. **InfoStep** displays a welcome screen with inline markup (`**bold**` and `{green}colored{/}`)
2. **InputStep** collects a name and stores it in state under the key `"name"`
3. **ConfirmStep** reads state to show a summary and waits for confirmation
4. After the wizard exits, `result.State` contains all collected values
