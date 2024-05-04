package tui

import (
	"github.com/gdamore/tcell/v2"
)

var (
	selectedStyle   = tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorDefault)
	unselectedStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
)
