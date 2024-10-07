package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/nsfisis/term-clock/internal/term"
)

func drawAlarm(scr *term.Screen, now time.Time, alarmTime time.Time, bgStyle, fgStyle term.Style) {
	h1, m1, s1 := now.Clock()
	h2, m2, s2 := alarmTime.Clock()

	if h1*3600+m1*60+s1 >= h2*3600+m2*60+s2 {
		bgStyle, fgStyle = fgStyle, bgStyle
	}

	drawClock(scr, now, bgStyle, fgStyle)
}

func cmdAlarm(cmd *cobra.Command, args []string) {
	alarmTime, err := time.Parse("15:04", args[0])
	if err != nil {
		log.Fatalf("%+v", err)
	}

	scr, err := term.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer scr.Close()

	drawAlarm(scr, time.Now(), alarmTime, term.BgStyle, term.FgStyle)

	scr.OnResize(func() bool {
		drawAlarm(scr, time.Now(), alarmTime, term.BgStyle, term.FgStyle)
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
			drawAlarm(scr, now, alarmTime, term.BgStyle, term.FgStyle)
			scr.Show()
		}
	}
}

var alarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "Alarm mode",
	Run:   cmdAlarm,
	Args:  cobra.ExactArgs(1),
}
