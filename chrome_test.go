package gwiz

import (
	"strings"
	"testing"

	"github.com/SerenaFontaine/tui"
)

func bufString(buf *tui.Buffer) string {
	var sb strings.Builder
	for y := 0; y < buf.Height; y++ {
		for x := 0; x < buf.Width; x++ {
			c := buf.Get(x, y)
			if c.Char == 0 {
				sb.WriteRune(' ')
			} else {
				sb.WriteRune(c.Char)
			}
		}
		if y < buf.Height-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func TestRenderHeader_ContainsTitleAndProgress(t *testing.T) {
	buf := tui.NewBuffer(50, 1)
	area := tui.NewRect(0, 0, 50, 1)
	theme := ThemeNord
	renderHeader(buf, area, "My Wizard", "", 2, 5, theme)

	s := bufString(buf)
	if !strings.Contains(s, "My Wizard") {
		t.Fatalf("header missing title, got: %q", s)
	}
	if !strings.Contains(s, "Step 3 of 5") {
		t.Fatalf("header missing progress, got: %q", s)
	}
}

func TestRenderNavBar_ContainsHints(t *testing.T) {
	buf := tui.NewBuffer(60, 1)
	area := tui.NewRect(0, 0, 60, 1)
	theme := ThemeNord
	hints := []KeyHint{
		{Key: "Enter", Label: "Next"},
		{Key: "Esc", Label: "Back"},
	}
	renderNavBar(buf, area, hints, theme)

	s := bufString(buf)
	if !strings.Contains(s, "Enter") || !strings.Contains(s, "Next") {
		t.Fatalf("nav bar missing hint, got: %q", s)
	}
	if !strings.Contains(s, "Esc") || !strings.Contains(s, "Back") {
		t.Fatalf("nav bar missing hint, got: %q", s)
	}
}

func TestRenderChrome_ReturnsInnerArea(t *testing.T) {
	buf := tui.NewBuffer(60, 20)
	area := tui.NewRect(0, 0, 60, 20)
	theme := ThemeNord
	hints := []KeyHint{{Key: "Enter", Label: "Next"}}
	inner := renderChrome(buf, area, "Title", "", 0, 3, hints, theme)

	if inner.Width >= area.Width || inner.Height >= area.Height {
		t.Fatalf("inner area %v not smaller than outer %v", inner, area)
	}
	if inner.Height <= 0 {
		t.Fatal("inner area has no height")
	}
}
