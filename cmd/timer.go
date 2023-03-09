package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/nsfisis/term-clock/internal/term"
)

func drawTimer(scr *term.Screen, leftTime time.Duration, bgStyle, fgStyle term.Style) {
	if leftTime <= 0 {
		leftTime = 0
		bgStyle, fgStyle = fgStyle, bgStyle
	}

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

	// Minute
	minute := leftTime.Minutes()
	term.DrawNumber(scr, int(minute)/10, xOffset+squareW*0, yOffset, squareW, squareH, fgStyle)
	term.DrawNumber(scr, int(minute)%10, xOffset+squareW*4, yOffset, squareW, squareH, fgStyle)

	// Colon
	term.DrawSquare(scr, xOffset+squareW*8, yOffset+squareH*1, squareW, squareH, fgStyle)
	term.DrawSquare(scr, xOffset+squareW*8, yOffset+squareH*3, squareW, squareH, fgStyle)

	// Second
	second := leftTime.Seconds()
	term.DrawNumber(scr, int(second)/10, xOffset+squareW*10, yOffset, squareW, squareH, fgStyle)
	term.DrawNumber(scr, int(second)%10, xOffset+squareW*14, yOffset, squareW, squareH, fgStyle)
}

func cmdTimer(cmd *cobra.Command, args []string) {
	timerTime, err := time.ParseDuration(args[0])
	if err != nil {
		log.Fatalf("%+v", err)
	}

	scr, err := term.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer scr.Close()

	startTime := time.Now()

	drawTimer(scr, (timerTime - time.Now().Sub(startTime)).Round(time.Second), term.BgStyle, term.FgStyle)

	scr.OnResize(func() bool {
		drawTimer(scr, (timerTime - time.Now().Sub(startTime)).Round(time.Second), term.BgStyle, term.FgStyle)
		return false
	})
	go scr.DoEventLoop()

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-scr.QuitC:
			return
		case now := <-t.C:
			drawTimer(scr, (timerTime - now.Sub(startTime)).Round(time.Second), term.BgStyle, term.FgStyle)
			scr.Show()
		}
	}
}

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Timer mode",
	Run:   cmdTimer,
	Args:  cobra.ExactArgs(1),
}
