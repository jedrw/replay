package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newFooter() *tview.TextView {
	footer := tview.NewTextView()
	footer.SetBackgroundColor(tcell.ColorDefault)

	keybindText := commandSelectControls
	footer.SetText(keybindText)
	footer.SetTextAlign(tview.AlignCenter).SetTextColor(altColour)

	return footer
}
