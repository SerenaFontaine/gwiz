package gwiz

import (
	"testing"

	"github.com/SerenaFontaine/tui"
)

func TestBuiltInThemes_HaveNames(t *testing.T) {
	themes := []Theme{ThemeNord, ThemeGruvbox, ThemeDracula, ThemeMonochrome}
	for _, th := range themes {
		if th.Name == "" {
			t.Fatal("theme has no name")
		}
	}
}

func TestBuiltInThemes_HaveNonZeroPrimary(t *testing.T) {
	themes := []Theme{ThemeNord, ThemeGruvbox, ThemeDracula, ThemeMonochrome}
	for _, th := range themes {
		if th.Primary == (tui.Color{}) {
			t.Fatalf("theme %q has zero Primary color", th.Name)
		}
	}
}

func TestTheme_Styles(t *testing.T) {
	th := ThemeNord
	_ = th.AccentStyle()
	_ = th.DimStyle()
	_ = th.ErrorStyle()
	_ = th.SuccessStyle()
	_ = th.WarningStyle()
}
