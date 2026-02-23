package cmd

import (
	"github.com/nsfisis/term-clock/internal/term"
)

func calcSquareSize(scr *term.Screen) (int, int, int, int) {
	// Calculate square width/height and offset.
	scrW, scrH := scr.Size()
	//         17
	// <--------------->
	// ### ###   ### ### ^
	// # # # # # # # # # |
	// # # # #   # # # # | 5
	// # # # # # # # # # |
	// ### ###   ### ### v
	squareW := scrW / (17 + 2)
	squareH := min(scrH/(5+2), squareW)
	if squareW > squareH*3/2 {
		squareW = squareH * 3 / 2
	}
	xOffset := (scrW - squareW*17) / 2
	yOffset := (scrH - squareH*5) / 2

	return squareW, squareH, xOffset, yOffset
}
