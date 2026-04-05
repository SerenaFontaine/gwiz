// Package gwiz provides a step-by-step terminal wizard framework built on
// [github.com/SerenaFontaine/tui].
//
// A wizard is a sequence of steps that collect user input, display information,
// and optionally run long-running tasks. State flows between steps via a
// typed key-value bag, and steps can be conditionally skipped based on
// prior answers.
//
// Applications create a [Wizard] with [New], add steps with [Wizard.AddStep],
// and run with [Wizard.Run]:
//
//	w := gwiz.New(
//	    gwiz.WithTitle("Setup"),
//	    gwiz.WithTheme(gwiz.ThemeNord),
//	)
//	w.AddStep("name", &gwiz.InputStep{
//	    BaseStep:  gwiz.BaseStep{TitleText: "Your Name"},
//	    Prompt:    "Enter your name:",
//	    ResultKey: "name",
//	})
//	result, err := w.Run(ctx)
//
// Built-in step types include [InputStep] for text input, [SelectStep] and
// [MultiSelectStep] for option selection, [FormStep] for multi-field forms,
// [InfoStep] for read-only content, [ConfirmStep] for confirmation screens,
// [TableStep] for tabular data, and [ExecStep] for long-running tasks with
// live output.
//
// Custom steps can be created by implementing the [Step] interface. Embed
// [BaseStep] to provide sensible defaults and only override the methods
// you need.
package gwiz
