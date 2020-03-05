package main

import (
	// "fmt"
	// "sync"
	// "time"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	// "github.com/cfal/webostv"
)

type ime struct {
	*tview.TextView
}

func newIme() *ime {
	w := tview.NewTextView()
	w.SetBorder(true)
	w.SetScrollable(true)
	w.SetTitle("IME")
	w.SetWrap(false)
	return &ime{
		TextView: w,
	}
}

func (i *ime) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return i.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyRune:
			kr := event.Rune()
			go tv.ImeInsertText(string(kr), false)
		case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyDelete:
			go tv.ImeDeleteCharacters(1)
		case tcell.KeyEnter:
			tv.ImeSendEnterKey()
			go func() {
				app.changeFocus(nil, app.wHelp)
				app.Draw()
			}()
		}
	})
}
