package gwiz

import (
	"github.com/SerenaFontaine/tui"
)

// Theme defines the color palette used to render a wizard.
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

// AccentStyle returns a bold style using the primary color.
func (t Theme) AccentStyle() tui.Style {
	return tui.NewStyle().Fg(t.Primary).Bold(true)
}

// DimStyle returns a style using the dimmed text color.
func (t Theme) DimStyle() tui.Style {
	return tui.NewStyle().Fg(t.TextDim)
}

// ErrorStyle returns a style using the error color.
func (t Theme) ErrorStyle() tui.Style {
	return tui.NewStyle().Fg(t.Error)
}

// SuccessStyle returns a style using the success color.
func (t Theme) SuccessStyle() tui.Style {
	return tui.NewStyle().Fg(t.Success)
}

// WarningStyle returns a style using the warning color.
func (t Theme) WarningStyle() tui.Style {
	return tui.NewStyle().Fg(t.Warning)
}

// TextStyle returns a style using the primary text color.
func (t Theme) TextStyle() tui.Style {
	return tui.NewStyle().Fg(t.Text)
}

// ThemeNord is a cool blue-tinted theme based on the Nord palette.
var ThemeNord = Theme{
	Name:       "Nord",
	Primary:    tui.Hex("#88C0D0"),
	Secondary:  tui.Hex("#81A1C1"),
	Background: tui.Hex("#2E3440"),
	Surface:    tui.Hex("#3B4252"),
	Text:       tui.Hex("#ECEFF4"),
	TextDim:    tui.Hex("#4C566A"),
	Success:    tui.Hex("#A3BE8C"),
	Warning:    tui.Hex("#EBCB8B"),
	Error:      tui.Hex("#BF616A"),
	Info:       tui.Hex("#88C0D0"),
}

// ThemeGruvbox is a warm retro theme based on the Gruvbox palette.
var ThemeGruvbox = Theme{
	Name:       "Gruvbox",
	Primary:    tui.Hex("#FE8019"),
	Secondary:  tui.Hex("#FABD2F"),
	Background: tui.Hex("#282828"),
	Surface:    tui.Hex("#3C3836"),
	Text:       tui.Hex("#EBDBB2"),
	TextDim:    tui.Hex("#665C54"),
	Success:    tui.Hex("#B8BB26"),
	Warning:    tui.Hex("#FABD2F"),
	Error:      tui.Hex("#FB4934"),
	Info:       tui.Hex("#83A598"),
}

// ThemeDracula is a dark purple theme based on the Dracula palette.
var ThemeDracula = Theme{
	Name:       "Dracula",
	Primary:    tui.Hex("#BD93F9"),
	Secondary:  tui.Hex("#FF79C6"),
	Background: tui.Hex("#282A36"),
	Surface:    tui.Hex("#44475A"),
	Text:       tui.Hex("#F8F8F2"),
	TextDim:    tui.Hex("#6272A4"),
	Success:    tui.Hex("#50FA7B"),
	Warning:    tui.Hex("#F1FA8C"),
	Error:      tui.Hex("#FF5555"),
	Info:       tui.Hex("#8BE9FD"),
}

// ThemeMonochrome is a minimal black-and-white theme.
var ThemeMonochrome = Theme{
	Name:       "Monochrome",
	Primary:    tui.White,
	Secondary:  tui.BrightBlack,
	Background: tui.Black,
	Surface:    tui.BrightBlack,
	Text:       tui.White,
	TextDim:    tui.BrightBlack,
	Success:    tui.White,
	Warning:    tui.White,
	Error:      tui.White,
	Info:       tui.White,
}
