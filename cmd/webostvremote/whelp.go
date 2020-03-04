package main

import (
	"fmt"
	"sync"
	"time"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/cfal/webostv"
)

type help struct {
	*tview.TextView
	ps *webostv.PointerSocket
	psUsageTime time.Time
	psMutex sync.Mutex
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

	return &help{
		TextView: w,
	}
}

func (h *help) getPointerSocket() *webostv.PointerSocket {
	h.psMutex.Lock()
	defer h.psMutex.Unlock()

	if h.ps != nil && time.Since(h.psUsageTime).Seconds() > 30 {
		h.ps = nil
	}
	if h.ps == nil {
		p, err := tv.NewPointerSocket()
		if err != nil {
			print("Could not connect to pointer socket")
			app.Stop()
			return nil
		}
		h.ps = p
	}
	h.psUsageTime = time.Now()
	return h.ps
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
			go h.getPointerSocket().ButtonEnter()
		case key == tcell.KeyLeft:
			go h.getPointerSocket().ButtonLeft()
		case key == tcell.KeyRight:
			go h.getPointerSocket().ButtonRight()
		case key == tcell.KeyUp:
			go h.getPointerSocket().ButtonUp()
		case key == tcell.KeyDown:
			go h.getPointerSocket().ButtonDown()
		case key == tcell.KeyESC || (key == tcell.KeyRune && (kr == 'b')):
			go h.getPointerSocket().ButtonBack()
		case key == tcell.KeyRune && (kr == 'h'):
			go h.getPointerSocket().ButtonHome()
		case key == tcell.KeyRune && (kr == 'i'):
			go h.getPointerSocket().ButtonInfo()
		case key == tcell.KeyRune && (kr == 'd'):
			go h.getPointerSocket().ButtonDash()
		case key == tcell.KeyRune && (kr == 's'):
			go h.getPointerSocket().Button("SETUP")
		}
	})
}
