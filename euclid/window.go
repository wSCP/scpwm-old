package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func locateWindow(e *Euclid, w xproto.Window) (*Client, bool) {
	for _, m := range e.monitors {
		for _, d := range m.desktops {
			for _, c := range d.clients {
				if w == c.Window {
					return c, true
				}
			}
		}
	}
	return nil, false
}

type window struct {
	xproto.Window
	*xgb.Conn
	root xproto.Window
}

var NoWindow window = window{nil, xproto.WindowNone, xproto.WindowNone}

func (w *Window) Close() {
	//send_client_message(w.Window, ewmh->WM_PROTOCOLS, WM_DELETE_WINDOW);
}

func (w *Window) Kill() {
	xproto.KillClientChecked(w.Conn, uint32(w.Window))
}

func (w *Window) BorderWidth(bw uint32) {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowBorderWidth, []uint32{bw})
}

func (w *Window) Move(x, y int16) {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowX, []uint32{uint32(x)})
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowY, []uint32{uint32(y)})
}

func (w *Window) Resize(hght, wdth uint16) {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowHeight, []uint32{uint32(hght)})
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowWidth, []uint32{uint32(wdth)})
}

func (w *Window) MoveResize(x, y int16, hght, wdth uint16) {
	w.Move(x, y)
	w.Resize(hght, wdth)
}

func (w *Window) Raise() {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeAbove})
}

func (w *Window) Lower() {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeBelow})
}

func (w *Window) Stack(o *Window, mode uint32) {
	xproto.ConfigureWindowChecked(
		w.Conn,
		w.Window,
		(xproto.ConfigWindowSibling | xproto.ConfigWindowStackMode),
		[]uint32{uint32(o.Window), mode},
	)
}

func (w *Window) Above(o *Window) {
	w.Stack(o, xproto.StackModeAbove)
}

func (w *Window) Below(o *Window) {
	w.Stack(o, xproto.StackModeBelow)
}

var (
	windowOff = []uint32{RootEventMask, xproto.EventMaskSubstructureNotify} //uint32_t values_off[] = {ROOT_EVENT_MASK & ~XCB_EVENT_MASK_SUBSTRUCTURE_NOTIFY};
	windowOn  = []uint32{RootEventMask}
)

func (w *Window) setVisibility(v bool) {
	xproto.ChangeWindowAttributesChecked(w.Conn, w.root, xproto.CwEventMask, windowOff)
	if v {
		xproto.MapWindow(w.Conn, w.Window)
	} else {
		xproto.UnmapWindow(w.Conn, w.Window)
	}
	xproto.ChangeWindowAttributesChecked(w.Conn, w.root, xproto.CwEventMask, windowOn)
}

func (w *Window) Hide() {
	w.setVisibility(false)
}

func (w *Window) Show() {
	w.setVisibility(true)
}
