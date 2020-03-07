package webostv

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
)

type PointerSocket struct {
	Address string
	ws      *websocket.Conn
	sync.Mutex
}

func (dialer *Dialer) DialPointerSocket(address string) (ps *PointerSocket, err error) {
	wsDialer := dialer.WebsocketDialer
	if wsDialer == nil {
		wsDialer = websocket.DefaultDialer
	}
	ws, resp, err := wsDialer.Dial(address, nil)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &PointerSocket{
		Address: address,
		ws:      ws,
	}, nil
}

func (tv *Tv) NewPointerSocket() (ps *PointerSocket, err error) {
	socketPath, err := tv.GetPointerInputSocket()
	if err != nil {
		return nil, err
	}
	return DefaultDialer.DialPointerSocket(socketPath)
}

func (ps *PointerSocket) MessageHandler() (err error) {
	for {
		_, _, err = ps.ws.ReadMessage()
		if err != nil {
			return err
		}
	}
	// not reached
}

func (ps *PointerSocket) Close() (err error) {
	ps.Lock()
	defer ps.Unlock()
	if ps.ws != nil {
		err = ps.ws.Close()
		ps.ws = nil
	}
	return err
}

func (ps *PointerSocket) writeMessage(messageType int, data []byte) error {
	ps.Lock()
	defer ps.Unlock()
	return ps.ws.WriteMessage(messageType, data)
}

func (ps *PointerSocket) Input(btype, bname string) (err error) {
	msg := "type:" + btype + "\n" + "name:" + bname + "\n\n"
	return ps.writeMessage(websocket.TextMessage, []byte(msg))
}

func (ps *PointerSocket) Move(dx, dy int) (err error) {
	msg := fmt.Sprintf("type:move\ndx:%d\ndy:%d\ndown:0\n\n", dx, dy)
	return ps.writeMessage(websocket.TextMessage, []byte(msg))
}

func (ps *PointerSocket) Scroll(dx, dy int) (err error) {
	msg := fmt.Sprintf("type:scroll\ndx:%d\ndy:%d\ndown:0\n\n", dx, dy)
	return ps.writeMessage(websocket.TextMessage, []byte(msg))
}

func (ps *PointerSocket) Click() (err error) {
	msg := "type: click\n\n"
	return ps.writeMessage(websocket.TextMessage, []byte(msg))
}

// Buttons reference:
// - https://www.openhab.org/addons/bindings/lgwebos/
// - https://github.com/TheRealLink/pylgtv/pull/19/files
// - 

func (ps *PointerSocket) Button(bname string) (err error) {
	return ps.Input("button", bname)
}

func (ps *PointerSocket) ButtonUp() (err error) {
	return ps.Button("UP")
}

func (ps *PointerSocket) ButtonDown() (err error) {
	return ps.Button("DOWN")
}

func (ps *PointerSocket) ButtonLeft() (err error) {
	return ps.Button("LEFT")
}

func (ps *PointerSocket) ButtonRight() (err error) {
	return ps.Button("RIGHT")
}

func (ps *PointerSocket) ButtonHome() (err error) {
	return ps.Button("HOME")
}

func (ps *PointerSocket) ButtonBack() (err error) {
	return ps.Button("BACK")
}

func (ps *PointerSocket) ButtonEnter() (err error) {
	return ps.Button("ENTER")
}

func (ps *PointerSocket) ButtonDash() (err error) {
	return ps.Button("DASH")
}

func (ps *PointerSocket) ButtonInfo() (err error) {
	return ps.Button("INFO")
}

func (ps *PointerSocket) ButtonAsterisk() (err error) {
	return ps.Button("ASTERISK")
}

func (ps *PointerSocket) ButtonMenu() (err error) {
	return ps.Button("MENU")
}

func (ps *PointerSocket) ButtonSAP() (err error) {
	return ps.Button("SAP")
}

func (ps *PointerSocket) ButtonCC() (err error) {
	return ps.Button("CC")
}

func (ps *PointerSocket) ButtonExit() (err error) {
	return ps.Button("EXIT")
}

func (ps *PointerSocket) ButtonMute() (err error) {
	return ps.Button("MUTE")
}

func (ps *PointerSocket) ButtonRed() (err error) {
	return ps.Button("RED")
}

func (ps *PointerSocket) ButtonGreen() (err error) {
	return ps.Button("GREEN")
}

func (ps *PointerSocket) ButtonBlue() (err error) {
	return ps.Button("BLUE")
}

func (ps *PointerSocket) ButtonVolumeUp() (err error) {
	return ps.Button("VOLUMEUP")
}

func (ps *PointerSocket) ButtonVolumeDown() (err error) {
	return ps.Button("VOLUMEDOWN")
}

func (ps *PointerSocket) ButtonChannelUp() (err error) {
	return ps.Button("CHANNELUP")
}

func (ps *PointerSocket) ButtonChannelDown() (err error) {
	return ps.Button("CHANNELDOWN")
}

func (ps *PointerSocket) ButtonNumber(num int) (err error) {
	if num < 0 || num > 9 {
		return errors.New("Invalid number")
	}
	return ps.Button(strconv.Itoa(num))
}
