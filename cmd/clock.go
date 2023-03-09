package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/nsfisis/term-clock/internal/term"
)

func drawClock(scr *term.Screen, now time.Time, bgStyle, fgStyle term.Style) {
	// Clear the entire screen.
	scr.Clear(bgStyle)

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
	squareH := scrH / (5 + 2)
	if squareH > squareW {
		squareH = squareW
	}
	if squareW > squareH*3/2 {
		squareW = squareH * 3 / 2
	}
	xOffset := (scrW - squareW*17) / 2
	yOffset := (scrH - squareH*5) / 2

	// Hour
	hour := now.Hour()
	term.DrawNumber(scr, hour/10, xOffset+squareW*0, yOffset, squareW, squareH, fgStyle)
	term.DrawNumber(scr, hour%10, xOffset+squareW*4, yOffset, squareW, squareH, fgStyle)

	// Colon
	term.DrawSquare(scr, xOffset+squareW*8, yOffset+squareH*1, squareW, squareH, fgStyle)
	term.DrawSquare(scr, xOffset+squareW*8, yOffset+squareH*3, squareW, squareH, fgStyle)

	// Minute
	minute := now.Minute()
	term.DrawNumber(scr, minute/10, xOffset+squareW*10, yOffset, squareW, squareH, fgStyle)
	term.DrawNumber(scr, minute%10, xOffset+squareW*14, yOffset, squareW, squareH, fgStyle)
}

func cmdClock(cmd *cobra.Command, args []string) {
	scr, err := term.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer scr.Close()

	drawClock(scr, time.Now(), term.BgStyle, term.FgStyle)

	scr.OnResize(func() bool {
		drawClock(scr, time.Now(), term.BgStyle, term.FgStyle)
		return false
	})
	go scr.DoEventLoop()

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	prev := time.Now()
	for {
		select {
		case <-scr.QuitC:
			return
		case now := <-t.C:
			if now.Minute() != prev.Minute() {
				drawClock(scr, now, term.BgStyle, term.FgStyle)
				scr.Show()
				prev = now
			}
		}
	}
}

var clockCmd = &cobra.Command{
	Use:   "clock",
	Short: "Clock mode",
	Run:   cmdClock,
}
