package cmd

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

func drawTimer(scr tcell.Screen, leftTime time.Duration, bgStyle, clockStyle tcell.Style) {
	if leftTime<=0{
		leftTime=0
		bgStyle, clockStyle = clockStyle, bgStyle
	}

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

	// Minute
	minute := leftTime.Minutes()
	drawNumber(scr, int(minute)/10, xOffset+squareW*0, yOffset, squareW, squareH, clockStyle)
	drawNumber(scr, int(minute)%10, xOffset+squareW*4, yOffset, squareW, squareH, clockStyle)

	// Colon
	drawSquare(scr, xOffset+squareW*8, yOffset+squareH*1, squareW, squareH, clockStyle)
	drawSquare(scr, xOffset+squareW*8, yOffset+squareH*3, squareW, squareH, clockStyle)

	// Second
	second := leftTime.Seconds()
	drawNumber(scr, int(second)/10, xOffset+squareW*10, yOffset, squareW, squareH, clockStyle)
	drawNumber(scr, int(second)%10, xOffset+squareW*14, yOffset, squareW, squareH, clockStyle)
}

func cmdTimer(cmd *cobra.Command, args []string) {
	timerTime, err := time.ParseDuration(args[0])
	if err != nil {
	log.Fatalf("%+v", err)
	}

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

	startTime := time.Now()

	drawTimer(scr, (timerTime - time.Now().Sub(startTime)).Round(time.Second), bgStyle, clockStyle)

	quitC := make(chan struct{})

	go func() {
		for {
			scr.Show()

			ev := scr.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventResize:
				drawTimer(scr, (timerTime - time.Now().Sub(startTime)).Round(time.Second), bgStyle, clockStyle)
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

	for {
		select {
		case <-quitC:
			return
		case now := <-t.C:
			drawTimer(scr, (timerTime - now.Sub(startTime)).Round(time.Second), bgStyle, clockStyle)
			scr.Show()
		}
	}
}

var timerCmd =&cobra.Command{
	Use: "timer",
	Short: "Timer mode",
	Run: cmdTimer,
	Args: cobra.ExactArgs(1),
}
