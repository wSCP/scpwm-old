package ewmh

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"

	"scpwm.local/scpwm/euclid/atomic"
)

type EWMH interface {
	SupportedGet() ([]string, error)
	SupportedSet([]string) error
	SupportingWmCheckGet(xproto.Window) (xproto.Window, error)
	SupportingWmCheckSet(xproto.Window, xproto.Window) error
	DesktopNamesGet() ([]string, error)
	DesktopNamesSet([]string) error
	NumberOfDesktopsGet() (uint, error)
	CurrentDesktopGet() (uint, error)
	ClientListGet() ([]xproto.Window, error)
	ActiveWindowGet() (xproto.Window, error)
	ActiveWindowSet(xproto.Window) error
	WmDesktopGet(xproto.Window) (uint, error)
	WmDesktopSet(xproto.Window, uint) error
	WmStateGet(xproto.Window) ([]string, error)
	WmStateSet(xproto.Window, []string) error
	WmWindowTypeGet(xproto.Window) ([]string, error)
}

type ewmh struct {
	conn *xgb.Conn
	root xproto.Window
	atom atomic.Atomic
}

func New(c *xgb.Conn, r xproto.Window, a atomic.Atomic) EWMH {
	return &ewmh{
		conn: c,
		root: r,
		atom: a,
	}
}

func (e *ewmh) SupportedGet() ([]string, error) {
	reply, err := e.atom.GetProp(e.root, "_NET_SUPPORTED")
	return e.atom.PropValAtoms(reply, err)
}

func (e *ewmh) SupportedSet(atomNames []string) error {
	atoms, err := e.atom.StrToAtoms(atomNames)
	if err != nil {
		return err
	}

	return e.atom.ChangeProp32(e.root, "_NET_SUPPORTED", "ATOM", atoms...)
}

func (e *ewmh) SupportingWmCheckGet(w xproto.Window) (xproto.Window, error) {
	return e.atom.PropValWindow(e.atom.GetProp(w, "_NET_SUPPORTING_WM_CHECK"))
}

func (e *ewmh) SupportingWmCheckSet(w xproto.Window, wmw xproto.Window) error {
	return e.atom.ChangeProp32(w, "_NET_SUPPORTING_WM_CHECK", "WINDOW", uint(wmw))
}

func (e *ewmh) DesktopNamesGet() ([]string, error) {
	return e.atom.PropValStrs(e.atom.GetProp(e.root, "_NET_DESKTOP_NAMES"))
}

func (e *ewmh) DesktopNamesSet(names []string) error {
	nullterm := make([]byte, 0)
	for _, name := range names {
		nullterm = append(nullterm, name...)
		nullterm = append(nullterm, 0)
	}
	return e.atom.ChangeProp(e.root, 8, "_NET_DESKTOP_NAMES", "UTF8_STRING", nullterm)
}

func (e *ewmh) NumberOfDesktopsGet() (uint, error) {
	return e.atom.PropValNum(e.atom.GetProp(e.root, "_NET_NUMBER_OF_DESKTOPS"))
}

func (e *ewmh) CurrentDesktopGet() (uint, error) {
	return e.atom.PropValNum(e.atom.GetProp(e.root, "_NET_CURRENT_DESKTOP"))
}

func (e *ewmh) ClientListGet() ([]xproto.Window, error) {
	return e.atom.PropValWindows(e.atom.GetProp(e.root, "_NET_CLIENT_LIST"))
}

func (e *ewmh) ActiveWindowGet() (xproto.Window, error) {
	return e.atom.PropValWindow(e.atom.GetProp(e.root, "_NET_ACTIVE_WINDOW"))
}

func (e *ewmh) ActiveWindowSet(w xproto.Window) error {
	return e.atom.ChangeProp32(e.root, "_NET_ACTIVE_WINDOW", "WINDOW", uint(w))
}

//func (e *ewmh) CloseWindow(w xproto.Window) error {
//	return e.CloseWindowExtra(w, 0, 2)
//}

//func (e *ewmh) CloseWindowExtra(w xproto.Window, time xproto.Timestamp, source int) error {
//	atm, err := e.atom.Atom("_NET_CLOSE_WINDOW")
//	if err != nil {
//		return err
//	}
//	return ClientEvent(e.conn, e.root, w, atm, int(time), source)
//}

func (e *ewmh) WmDesktopGet(w xproto.Window) (uint, error) {
	return e.atom.PropValNum(e.atom.GetProp(w, "_NET_WM_DESKTOP"))
}

func (e *ewmh) WmDesktopSet(w xproto.Window, desk uint) error {
	return e.atom.ChangeProp32(w, "_NET_WM_DESKTOP", "CARDINAL", uint(desk))
}

func (e *ewmh) WmStateGet(w xproto.Window) ([]string, error) {
	raw, err := e.atom.GetProp(w, "_NET_WM_STATE")
	return e.atom.PropValAtoms(raw, err)
}

func (e *ewmh) WmStateSet(w xproto.Window, atomNames []string) error {
	atoms, err := e.atom.StrToAtoms(atomNames)
	if err != nil {
		return err
	}

	return e.atom.ChangeProp32(w, "_NET_WM_STATE", "ATOM", atoms...)
}

//NET_WM_STATE_FULLSCREEN

//NET_WM_STATE_STICKY

//NET_WM_STATE_DEMANDS_ATTENTION

func (e *ewmh) WmWindowTypeGet(w xproto.Window) ([]string, error) {
	raw, err := e.atom.GetProp(w, "_NET_WM_WINDOW_TYPE")
	return e.atom.PropValAtoms(raw, err)
}

//NET_WM_WINDOW_TYPE_DOCK

//NET_WM_WINDOW_TYPE_DESKTOP

//NET_WM_WINDOW_TYPE_NOTIFICATION

//NET_WM_WINDOW_TYPE_DIALOG

//WM_WINDOW_TYPE_UTILITY

//WM_WINDOW_TYPE_TOOLBAR
