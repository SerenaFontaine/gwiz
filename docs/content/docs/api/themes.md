---
title: Themes
weight: 4
---

`Theme` defines the color palette used to render a wizard.

## Theme Struct

```go
type Theme struct {
    Name       string
    Primary    tui.Color
    Secondary  tui.Color
    Background tui.Color
    Surface    tui.Color
    Text       tui.Color
    TextDim    tui.Color
    Success    tui.Color
    Warning    tui.Color
    Error      tui.Color
    Info       tui.Color
}
```

## Built-in Themes

| Variable | Description |
|----------|-------------|
| `ThemeNord` | Cool blue-tinted Nord palette (default) |
| `ThemeGruvbox` | Warm retro Gruvbox palette |
| `ThemeDracula` | Dark purple Dracula palette |
| `ThemeMonochrome` | Minimal black-and-white |

```go
w := gwiz.New(gwiz.WithTheme(gwiz.ThemeDracula))
```

## Style Helpers

Each theme provides convenience methods that return `tui.Style` values:

| Method | Description |
|--------|-------------|
| `AccentStyle()` | Bold style with primary color |
| `DimStyle()` | Style with dimmed text color |
| `ErrorStyle()` | Style with error color |
| `SuccessStyle()` | Style with success color |
| `WarningStyle()` | Style with warning color |
| `TextStyle()` | Style with primary text color |

## Custom Themes

Define a custom theme by creating a `Theme` value with your colors:

```go
var ThemeCustom = gwiz.Theme{
    Name:       "Custom",
    Primary:    tui.Hex("#FF6B6B"),
    Secondary:  tui.Hex("#4ECDC4"),
    Background: tui.Hex("#1A1A2E"),
    Surface:    tui.Hex("#16213E"),
    Text:       tui.Hex("#E8E8E8"),
    TextDim:    tui.Hex("#555555"),
    Success:    tui.Hex("#95E06C"),
    Warning:    tui.Hex("#FFD93D"),
    Error:      tui.Hex("#FF6B6B"),
    Info:       tui.Hex("#6BCB77"),
}

w := gwiz.New(gwiz.WithTheme(ThemeCustom))
```
