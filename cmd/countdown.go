package cmd

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/nsfisis/term-clock/internal/term"
)

func drawCountdown(scr *term.Screen, leftSeconds int, bgStyle, fgStyle term.Style) {
	if leftSeconds <= 0 {
		if int(math.Abs(float64(leftSeconds)))%2 == 0 {
			bgStyle, fgStyle = fgStyle, bgStyle
		}
		leftSeconds = 0
	}

	// Clear the entire screen.
	scr.Clear(bgStyle)

	// Calculate number of digits needed
	numDigits := 1
	if leftSeconds > 0 {
		n := leftSeconds
		numDigits = 0
		for n > 0 {
			numDigits++
			n /= 10
		}
	}

	scrW, scrH := scr.Size()
	// Calculate square size based on number of digits
	squareW := scrW / (numDigits*4 + 2)
	squareH := min(scrH/(5+2), squareW)
	if squareW > squareH*3/2 {
		squareW = squareH * 3 / 2
	}
	xOffset := (scrW - squareW*numDigits*4) / 2
	yOffset := (scrH - squareH*5) / 2

	// Extract and display each digit
	n := leftSeconds
	for i := numDigits - 1; i >= 0; i-- {
		digit := n % 10
		term.DrawNumber(scr, digit, xOffset+squareW*i*4, yOffset, squareW, squareH, fgStyle)
		n /= 10
	}
}

func cmdCountdown(cmd *cobra.Command, args []string) {
	totalSeconds, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if totalSeconds > 9999 {
		log.Fatal("Duration over 9999 seconds is not supported.")
	}

	scr, err := term.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer scr.Close()

	startTime := time.Now()

	calcLeftSeconds := func(t time.Time) int {
		elapsed := int(t.Sub(startTime).Seconds())
		return totalSeconds - elapsed
	}

	drawCountdown(scr, calcLeftSeconds(time.Now()), term.BgStyle, term.FgStyle)

	scr.OnResize(func() bool {
		drawCountdown(scr, calcLeftSeconds(time.Now()), term.BgStyle, term.FgStyle)
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
			drawCountdown(scr, calcLeftSeconds(now), term.BgStyle, term.FgStyle)
			scr.Show()
		}
	}
}

var countdownCmd = &cobra.Command{
	Use:   "countdown",
	Short: "Countdown mode",
	Run:   cmdCountdown,
	Args:  cobra.ExactArgs(1),
}
