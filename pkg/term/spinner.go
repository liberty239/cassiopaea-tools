package term

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/pterm/pterm"
)

func Spinner() *pterm.SpinnerPrinter {
	return &pterm.SpinnerPrinter{
		Sequence:            spinner.CharSets[14],
		Style:               &pterm.ThemeDefault.SpinnerStyle,
		Delay:               time.Millisecond * 200,
		ShowTimer:           true,
		TimerRoundingFactor: time.Second,
		TimerStyle:          &pterm.ThemeDefault.TimerStyle,
		MessageStyle:        &pterm.ThemeDefault.SpinnerTextStyle,
		SuccessPrinter:      &pterm.Success,
		FailPrinter:         &pterm.Error,
		WarningPrinter:      &pterm.Warning,
	}
}
