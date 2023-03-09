package term

import (
	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	scr tcell.Screen
	onResize func() bool
	QuitC chan struct{}
}

func NewScreen() (*Screen, error) {
	scr, err := tcell.NewScreen()
	if err != nil{
		return nil, err
	}
	err = scr.Init()
	if err != nil{
		return nil, err
	}
	return &Screen{
		scr: scr,
		QuitC: make(chan struct{}),
	}, nil
}

func (scr *Screen) Close() {
	scr.scr.Fini()
}

func (scr *Screen) Size() (int, int) {
	return scr.scr.Size()
}

func (scr *Screen) Clear(style Style) {
	scr.scr.SetStyle(tcell.Style(style))
	scr.scr.Clear()
}

func (scr *Screen) OnResize(handler func() bool) {
	scr.onResize = handler
}

func (scr *Screen) Show() {
	scr.scr.Show()
}

func (scr *Screen) DoEventLoop() {
	for {
		scr.scr.Show()

		ev := scr.scr.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			if scr.onResize != nil {
				if quit := scr.onResize(); quit {
					return
				}
			}
			scr.scr.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				close(scr.QuitC)
				return
			}
			scr.scr.Sync()
		}
	}
}
