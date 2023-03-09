package term

import (
	"github.com/gdamore/tcell/v2"
)

type Style tcell.Style;

var (
	BgStyle Style
	FgStyle Style
)

func init() {
	BgStyle = Style(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))
	FgStyle = Style(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorOlive))
}

func DrawSquare(scr *Screen, xOffset, yOffset, w, h int, style Style) {
	for dx := 0; dx < w; dx++ {
		x := xOffset + dx
		for dy := 0; dy < h; dy++ {
			y := yOffset + dy
			scr.scr.SetContent(x, y, ' ', nil, tcell.Style(style))
		}
	}
}

func DrawNumber(scr *Screen, n, xOffset, yOffset, squareW, squareH int, style Style) {
	defs := [...][15]bool{
		{
			true, true, true,
			true, false, true,
			true, false, true,
			true, false, true,
			true, true, true,
		},
		{
			false, false, true,
			false, false, true,
			false, false, true,
			false, false, true,
			false, false, true,
		},
		{
			true, true, true,
			false, false, true,
			true, true, true,
			true, false, false,
			true, true, true,
		},
		{
			true, true, true,
			false, false, true,
			true, true, true,
			false, false, true,
			true, true, true,
		},
		{
			true, false, true,
			true, false, true,
			true, true, true,
			false, false, true,
			false, false, true,
		},
		{
			true, true, true,
			true, false, false,
			true, true, true,
			false, false, true,
			true, true, true,
		},
		{
			true, true, true,
			true, false, false,
			true, true, true,
			true, false, true,
			true, true, true,
		},
		{
			true, true, true,
			false, false, true,
			false, false, true,
			false, false, true,
			false, false, true,
		},
		{
			true, true, true,
			true, false, true,
			true, true, true,
			true, false, true,
			true, true, true,
		},
		{
			true, true, true,
			true, false, true,
			true, true, true,
			false, false, true,
			true, true, true,
		},
	}

	squares := defs[n]
	for i, draw := range squares {
		if !draw {
			continue
		}
		x := i % 3
		y := i / 3
		DrawSquare(scr, xOffset+squareW*x, yOffset+squareH*y, squareW, squareH, style)
	}
}
