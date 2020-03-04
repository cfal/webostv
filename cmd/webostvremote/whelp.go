package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/cfal/webostv"
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
	fmt.Fprintln(w, ".         play/pause");
	fmt.Fprintln(w, "arrows    move")
	fmt.Fprintln(w, "Enter     enter")
	fmt.Fprintln(w, "Escape/B  back")
	fmt.Fprintln(w, "H         home")
	fmt.Fprintln(w, "I         info")
	fmt.Fprintln(w, "D         dash")
	fmt.Fprintln(w, "Q         quit")
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
		case key == tcell.KeyEnter:
			go h.pointerSocket.ButtonEnter()
		case key == tcell.KeyLeft:
			go h.pointerSocket.ButtonLeft()
		case key == tcell.KeyRight:
			go h.pointerSocket.ButtonRight()
		case key == tcell.KeyUp:
			go h.pointerSocket.ButtonUp()
		case key == tcell.KeyDown:
			go h.pointerSocket.ButtonDown()
		case key == tcell.KeyESC || (key == tcell.KeyRune && (kr == 'b')):
			go h.pointerSocket.ButtonBack()
		case key == tcell.KeyRune && (kr == 'h'):
			go h.pointerSocket.ButtonHome()
		case key == tcell.KeyRune && (kr == 'i'):
			go h.pointerSocket.ButtonInfo()
		case key == tcell.KeyRune && (kr == 'd'):
			go h.pointerSocket.ButtonDash()
		case key == tcell.KeyRune && (kr == 's'):
			go h.pointerSocket.Button("SETUP")
		}
	})
}
