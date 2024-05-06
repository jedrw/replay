package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewFooter() *tview.TextView {
	footer := tview.NewTextView()
	footer.SetTitle("REPLAY")
	footer.SetBackgroundColor(tcell.ColorDefault)

	footer.SetText("Navigate <Up|Down>	Select <Enter>	Order <Alt+[1-9]>	Search <ASCII>	Replay <Alt+Enter>")
	footer.SetTextAlign(tview.AlignCenter).SetTextColor(altColour)

	return footer
}
