package gwiz

import (
	"context"

	"github.com/SerenaFontaine/tui"
)

// Wizard is the top-level orchestrator that implements tui.Component.
type Wizard struct {
	title   string
	banner  string
	theme   Theme
	steps   []registeredStep
	state   State
	current int

	showQuitDialog bool
	quitDialog     *quitDialogState

	done    bool
	aborted bool

	err error

	app *tui.App
}

type quitDialogState struct {
	selected int // 0 = cancel, 1 = quit
}

// New creates a Wizard with the given options. The default theme is ThemeNord.
func New(opts ...WizardOption) *Wizard {
	w := &Wizard{
		theme: ThemeNord,
		state: newState(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// WithTitle sets the title displayed in the wizard header.
func WithTitle(title string) WizardOption {
	return func(w *Wizard) { w.title = title }
}

// WithTheme sets the color theme used for rendering.
func WithTheme(theme Theme) WizardOption {
	return func(w *Wizard) { w.theme = theme }
}

// WithBanner sets the subtitle text shown next to the title.
func WithBanner(banner string) WizardOption {
	return func(w *Wizard) { w.banner = banner }
}

// AddStep registers a named step in the wizard. Steps are presented in the
// order they are added.
func (w *Wizard) AddStep(name string, step Step) {
	w.steps = append(w.steps, registeredStep{name: name, step: step})
}

// Run starts the wizard in an alternate-screen TUI and blocks until the user
// completes or aborts. It returns the accumulated state and whether the wizard
// was aborted.
func (w *Wizard) Run(ctx context.Context) (*Result, error) {
	app := tui.NewApp(w, tui.WithAltScreen(true))
	w.app = app
	err := app.Run()
	if err != nil {
		return nil, err
	}
	return &Result{
		State:   w.state,
		Aborted: w.aborted,
	}, nil
}

// Init implements tui.Component.
func (w *Wizard) Init() Cmd {
	w.current = firstStep(w.steps, w.state)
	if w.current < 0 {
		return tui.QuitCmd()
	}
	return w.steps[w.current].step.Init(w.state)
}

// Update implements tui.Component.
func (w *Wizard) Update(msg Msg) (tui.Component, Cmd) {
	if w.showQuitDialog {
		return w.updateQuitDialog(msg)
	}

	switch msg := msg.(type) {
	case tui.KeyMsg:
		return w.handleKey(msg)
	case NextMsg:
		return w.advanceStep()
	case PrevMsg:
		return w.retreatStep()
	case QuitMsg:
		return w.handleQuit()
	case StepResultMsg:
		w.state.Set(msg.Key, msg.Value)
		return w, nil
	case ErrorMsg:
		w.err = msg.Err
		return w, nil
	default:
		if w.current >= 0 && w.current < len(w.steps) {
			_, cmd := w.steps[w.current].step.Update(msg, w.state)
			return w, cmd
		}
		return w, nil
	}
}

func (w *Wizard) handleKey(msg tui.KeyMsg) (tui.Component, Cmd) {
	switch msg.Type {
	case tui.KeyCtrlC:
		return w.handleQuit()
	case tui.KeyRune:
		// Only intercept 'q' for quit if the step doesn't accept text input.
		if msg.Rune == 'q' && !w.activeStepAcceptsText() {
			return w.handleQuit()
		}
	}
	if w.current >= 0 && w.current < len(w.steps) {
		_, cmd := w.steps[w.current].step.Update(msg, w.state)
		return w, cmd
	}
	return w, nil
}

// activeStepAcceptsText returns true if the current step handles free-form text input.
func (w *Wizard) activeStepAcceptsText() bool {
	if w.current < 0 || w.current >= len(w.steps) {
		return false
	}
	type textAcceptor interface {
		AcceptsTextInput() bool
	}
	if ta, ok := w.steps[w.current].step.(textAcceptor); ok {
		return ta.AcceptsTextInput()
	}
	return false
}

func (w *Wizard) handleQuit() (tui.Component, Cmd) {
	if len(w.state.Keys()) == 0 {
		w.aborted = true
		return w, tui.QuitCmd()
	}
	w.showQuitDialog = true
	w.quitDialog = &quitDialogState{selected: 0}
	return w, nil
}

func (w *Wizard) updateQuitDialog(msg Msg) (tui.Component, Cmd) {
	if keyMsg, ok := msg.(tui.KeyMsg); ok {
		switch keyMsg.Type {
		case tui.KeyEscape:
			w.showQuitDialog = false
			w.quitDialog = nil
			return w, nil
		case tui.KeyTab, tui.KeyRight, tui.KeyLeft:
			w.quitDialog.selected = 1 - w.quitDialog.selected
			return w, nil
		case tui.KeyEnter:
			if w.quitDialog.selected == 1 {
				w.aborted = true
				w.showQuitDialog = false
				return w, tui.QuitCmd()
			}
			w.showQuitDialog = false
			w.quitDialog = nil
			return w, nil
		case tui.KeyCtrlC:
			w.aborted = true
			return w, tui.QuitCmd()
		}
	}
	return w, nil
}

func (w *Wizard) advanceStep() (tui.Component, Cmd) {
	if w.current >= 0 && w.current < len(w.steps) {
		if err := w.steps[w.current].step.Validate(w.state); err != nil {
			w.err = err
			return w, nil
		}
	}

	next := nextStep(w.current, w.steps, w.state)
	if next < 0 {
		w.done = true
		return w, tui.QuitCmd()
	}

	w.current = next
	w.err = nil
	return w, w.steps[w.current].step.Init(w.state)
}

func (w *Wizard) retreatStep() (tui.Component, Cmd) {
	prev := prevStep(w.current, w.steps, w.state)
	if prev < 0 {
		return w, nil
	}
	w.current = prev
	w.err = nil
	return w, w.steps[w.current].step.Init(w.state)
}

// Render implements tui.Component.
func (w *Wizard) Render(buf *tui.Buffer, area tui.Rect) {
	if w.current < 0 || w.current >= len(w.steps) {
		return
	}

	step := w.steps[w.current].step
	hints := defaultKeyHints(step)

	contentArea := renderChrome(buf, area, w.title, w.banner, w.current, len(w.steps), hints, w.theme)

	step.Render(buf, contentArea, w.state)

	if w.err != nil {
		renderError(buf, contentArea, w.err, w.theme)
	}

	if w.showQuitDialog {
		renderQuitDialog(buf, area, w.quitDialog, w.theme)
	}
}

func defaultKeyHints(step Step) []KeyHint {
	type hinter interface {
		KeyHints() []KeyHint
	}
	if h, ok := step.(hinter); ok {
		return h.KeyHints()
	}
	return BaseStep{}.KeyHints()
}

func renderError(buf *tui.Buffer, area tui.Rect, err error, theme Theme) {
	errStyle := theme.ErrorStyle()
	y := area.Y + area.Height - 1
	if y >= area.Y {
		buf.SetString(area.X, y, "Error: "+err.Error(), errStyle)
	}
}

func renderQuitDialog(buf *tui.Buffer, area tui.Rect, state *quitDialogState, theme Theme) {
	dialogW := 40
	dialogH := 5
	x := area.X + (area.Width-dialogW)/2
	y := area.Y + (area.Height-dialogH)/2
	dialogArea := tui.NewRect(x, y, dialogW, dialogH)

	buf.Fill(dialogArea, tui.Cell{Char: ' ', Style: tui.NewStyle().Bg(theme.Surface)})

	block := tui.Block{
		Border: tui.BorderRounded,
		Style:  tui.NewStyle().Fg(theme.Warning),
		Title:  "Quit?",
	}
	inner := block.Render(buf, dialogArea)

	buf.SetString(inner.X, inner.Y, "Progress will be lost.", tui.NewStyle().Fg(theme.Text).Bg(theme.Surface))

	cancelLabel := "[ Cancel ]"
	quitLabel := "[ Quit ]"
	cancelStyle := tui.NewStyle().Fg(theme.TextDim).Bg(theme.Surface)
	quitStyle := tui.NewStyle().Fg(theme.TextDim).Bg(theme.Surface)
	if state.selected == 0 {
		cancelStyle = tui.NewStyle().Fg(theme.Text).Bg(theme.Primary).Bold(true)
	} else {
		quitStyle = tui.NewStyle().Fg(theme.Text).Bg(theme.Error).Bold(true)
	}
	btnY := inner.Y + inner.Height - 1
	buf.SetString(inner.X, btnY, cancelLabel, cancelStyle)
	buf.SetString(inner.X+len(cancelLabel)+2, btnY, quitLabel, quitStyle)
}
