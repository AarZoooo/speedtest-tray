package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type narrowProgressBar struct {
	widget.ProgressBar
}

func (n *narrowProgressBar) MinSize() fyne.Size {
	return fyne.NewSize(n.ProgressBar.MinSize().Width, 4) // Force height to 4
}

func newNarrowProgressBar() *narrowProgressBar {
	p := &narrowProgressBar{}
	p.Max = 1000
	p.TextFormatter = func() string { return "" }
	p.ExtendBaseWidget(p)
	return p
}

func isSpinner(text string) bool {
	if text == "--" {
		return true
	}
	for _, f := range spinnerFrames {
		if f == text {
			return true
		}
	}
	return false
}
