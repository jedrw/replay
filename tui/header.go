package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewHeader() *tview.TextView {
	header := tview.NewTextView()
	header.SetTitle("REPLAY")
	header.SetBackgroundColor(tcell.ColorDefault)

	header.SetText("Move <Up|Down|Left|Right>	Select <Enter>	Order <1-9>	Replay <Alt+Enter>")
	header.SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorSteelBlue)

	return header
}
