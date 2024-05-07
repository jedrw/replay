package tui

import (
	"fmt"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewFooter() *tview.TextView {
	footer := tview.NewTextView()
	footer.SetTitle("REPLAY")
	footer.SetBackgroundColor(tcell.ColorDefault)

	var replayKeyBind string
	switch runtime.GOOS {
	case "darwin":
		replayKeyBind = "Ctrl+r"
	default:
		replayKeyBind = "Alt+Enter"
	}

	keybindText := fmt.Sprintf("Navigate <Up|Down>	Select/Deselect <Enter>	Order <F[1-9]>	Search <ASCII>	Replay <%s>", replayKeyBind)
	footer.SetText(keybindText)
	footer.SetTextAlign(tview.AlignCenter).SetTextColor(altColour)

	return footer
}
