package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/snabb/webostv"
)

type help struct {
	*tview.TextView
    pointerSocket *webostv.PointerSocket
	updateInfo func(str string)
}

func newHelp() *help {
	w := tview.NewTextView()
	w.SetBorder(true)
	w.SetScrollable(true)
	w.SetTitle("Controls")
	w.SetWrap(false)
	fmt.Fprintln(w, "Tab ⇥ / ⇤ next / prev")
	fmt.Fprintln(w, "V         volume")
	fmt.Fprintln(w, ".         play/pause");
	fmt.Fprintln(w, "C         channels")
	fmt.Fprintln(w, "I         inputs")
	fmt.Fprintln(w, "A         apps")
	fmt.Fprintln(w, "arrows    move")
	fmt.Fprintln(w, "Q / Esc   quit")
	fmt.Fprintln(w, "Ctrl+X    turn off+quit")
	w.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		w.ScrollToBeginning()
		return x + 1, y + 1, width - 2, height - 2
	})

	p, err := tv.NewPointerSocket();
	if err != nil {
		print("Could not connect to pointer socket")
		app.Stop()
		return nil
	}

	return &help{
		TextView: w,
		pointerSocket: p,
	}
}

func (h *help) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return h.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()
		var kr rune
		if key == tcell.KeyRune {
			kr = event.Rune()
		}
		switch {
		case key == tcell.KeyRune && (kr == '+' || kr == '='):
			go tv.AudioVolumeUp()
			// XXX check err
		case key == tcell.KeyRune && (kr == '-' || kr == '_'):
			go tv.AudioVolumeDown()
			// XXX check err
		case key == tcell.KeyRune && kr == '.':
			go func() {
				// This does both play and pause.
				if err := tv.MediaControlsPlay(); err != nil {
					h.updateInfo("Failed to play");
				}
			}()
			
		}
	})
}
