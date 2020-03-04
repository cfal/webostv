package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type help struct {
	*tview.TextView
}

func newHelp() *help {
	w := tview.NewTextView()
	w.SetBorder(true)
	w.SetScrollable(true)
	w.SetTitle("Controls")
	w.SetWrap(false)
	fmt.Fprintln(w, "Tab ⇥ / ⇤ next / prev")
	fmt.Fprintln(w, "V         volume")
	fmt.Fprintln(w, "C         channels")
	fmt.Fprintln(w, "I         inputs")
	fmt.Fprintln(w, "A         apps")
	fmt.Fprintln(w, "Enter     select")
	fmt.Fprintln(w, "arrows    move")
	fmt.Fprintln(w, "Q / Esc   quit")
	fmt.Fprintln(w, "Ctrl+X    turn off+quit\n")
	fmt.Fprintln(w, "webostvremote © J.Snabb 2018")
	fmt.Fprint(w, "github.com/snabb/webostv")
	w.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		w.ScrollToBeginning()
		return x + 1, y + 1, width - 2, height - 2
	})
	return &help{w}
}

func (h *help) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return h.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()
		var kr rune
		if key == tcell.KeyRune {
			kr = event.Rune()
		}
		switch {
		case key == tcell.KeyRight || (key == tcell.KeyRune && kr == '+'):
			go tv.AudioVolumeUp()
			// XXX check err
		case key == tcell.KeyLeft || (key == tcell.KeyRune && kr == '-'):
			go tv.AudioVolumeDown()
			// XXX check err
		}
	})
}
