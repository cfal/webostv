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
		h.ps.Close()
		h.ps = nil
	}
	if h.ps == nil {
		p, err := tv.NewPointerSocket()
		if err != nil {
			app.Stop()
			print("Could not connect to pointer socket")
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
			if kr >= 'a' && kr <= 'z' {
				kr -= ('a' - 'A')
			}
		}
		switch {
		case key == tcell.KeyRune && (kr == '+' || kr == '='):
			go tv.AudioVolumeUp()
			// XXX check err
		case key == tcell.KeyRune && (kr == '-' || kr == '_'):
			go tv.AudioVolumeDown()
			// XXX check err
		case key == tcell.KeyRune && kr == ' ':
			go func() {
				// This does both play and pause.
				if err := tv.MediaControlsPlay(); err != nil {
					h.updateInfo("Failed to play");
				}
			}()
		}

		// Pointer socket keys
		if time.Since(h.psUsageTime).Milliseconds() < 150 {
			return
		}

		switch {
		case key == tcell.KeyEnter:
			go h.getPointerSocket().ButtonEnter()
		case key == tcell.KeyLeft || (key == tcell.KeyRune && kr == 'J'):
			go h.getPointerSocket().ButtonLeft()
		case key == tcell.KeyRight || (key == tcell.KeyRune && kr == 'L'):
			go h.getPointerSocket().ButtonRight()
		case key == tcell.KeyUp || (key == tcell.KeyRune && kr == 'I'):
			go h.getPointerSocket().ButtonUp()
		case key == tcell.KeyDown || (key == tcell.KeyRune && kr == 'K'):
			go h.getPointerSocket().ButtonDown()
		case key == tcell.KeyESC || key == tcell.KeyBackspace || key == tcell.KeyBackspace2 || (key == tcell.KeyRune && (kr == 'B')):
			go h.getPointerSocket().ButtonBack()
		case key == tcell.KeyRune && (kr == 'H'):
			go h.getPointerSocket().ButtonHome()
		case key == tcell.KeyRune && (kr == 'I'):
			go h.getPointerSocket().ButtonInfo()
		case key == tcell.KeyRune && (kr == 'D'):
			go h.getPointerSocket().ButtonDash()
		case key == tcell.KeyRune && (kr == 'X'):
			go h.getPointerSocket().ButtonExit()
		}
	})
}
