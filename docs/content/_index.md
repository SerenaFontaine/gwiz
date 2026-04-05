---
title: gwiz — Step-by-Step Terminal Wizard Framework for Go
---

A step-by-step terminal wizard framework for Go built on [tui](https://github.com/SerenaFontaine/tui).

## Features

- **Step-by-step flow** — Guided wizard with forward/back navigation and conditional skipping
- **8 built-in step types** — Input, Select, MultiSelect, Form, Info, Confirm, Table, Exec
- **State management** — Typed key-value bag flows between steps automatically
- **Validation** — Per-step validation with inline error display
- **Dynamic content** — Steps can generate options, rows, and content from prior state
- **Long-running tasks** — ExecStep with live output, progress, and retry on failure
- **Themes** — Built-in Nord, Gruvbox, Dracula, Monochrome themes
- **Extensible** — Implement the Step interface for custom step types

## Quick Links

- [Getting Started](/docs/getting-started/) — Installation and first steps
- [Examples](/docs/examples/) — Practical usage examples
- [API Reference](/docs/api/) — Exhaustive technical reference
