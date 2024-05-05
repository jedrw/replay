package tui

import (
	"github.com/gdamore/tcell/v2"
)

var (
	altColour       = tcell.ColorSteelBlue
	selectedStyle   = tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	unselectedStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
)
