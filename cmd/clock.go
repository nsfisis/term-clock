package cmd

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

func drawSquare(scr tcell.Screen, xOffset, yOffset, w, h int, style tcell.Style) {
	for dx := 0; dx < w; dx++ {
		x := xOffset + dx
		for dy := 0; dy < h; dy++ {
			y := yOffset + dy
			scr.SetContent(x, y, ' ', nil, style)
		}
	}
}

func drawNumber(scr tcell.Screen, n, xOffset, yOffset, squareW, squareH int, style tcell.Style) {
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
		drawSquare(scr, xOffset+squareW*x, yOffset+squareH*y, squareW, squareH, style)
	}
}

func drawClock(scr tcell.Screen, now time.Time, bgStyle, clockStyle tcell.Style) {
	// Clear the entire screen.
	scr.SetStyle(bgStyle)
	scr.Clear()

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
	drawNumber(scr, hour/10, xOffset+squareW*0, yOffset, squareW, squareH, clockStyle)
	drawNumber(scr, hour%10, xOffset+squareW*4, yOffset, squareW, squareH, clockStyle)

	// Colon
	drawSquare(scr, xOffset+squareW*8, yOffset+squareH*1, squareW, squareH, clockStyle)
	drawSquare(scr, xOffset+squareW*8, yOffset+squareH*3, squareW, squareH, clockStyle)

	// Minute
	minute := now.Minute()
	drawNumber(scr, minute/10, xOffset+squareW*10, yOffset, squareW, squareH, clockStyle)
	drawNumber(scr, minute%10, xOffset+squareW*14, yOffset, squareW, squareH, clockStyle)
}

func cmdClock(cmd *cobra.Command, args []string) {
	bgStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	clockStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorOlive)

	scr, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := scr.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer scr.Fini()

	drawClock(scr, time.Now(), bgStyle, clockStyle)

	quitC := make(chan struct{})

	go func() {
		for {
			scr.Show()

			ev := scr.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventResize:
				drawClock(scr, time.Now(), bgStyle, clockStyle)
				scr.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					close(quitC)
					return
				}
			}
		}
	}()

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	prev := time.Now()
	for {
		select {
		case <-quitC:
			return
		case now := <-t.C:
			if now.Minute() != prev.Minute() {
				drawClock(scr, now, bgStyle, clockStyle)
				scr.Show()
				prev = now
			}
		}
	}
}

var clockCmd =&cobra.Command{
	Use: "clock",
	Short: "Clock mode",
	Run: cmdClock,
}
