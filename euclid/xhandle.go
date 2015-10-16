package main

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/davecgh/go-spew/spew"

	"scpwm.local/scpwm/euclid/atomic"
	"scpwm.local/scpwm/euclid/ewmh"
)

var (
	RootEventMask   uint32 = (xproto.EventMaskSubstructureNotify | xproto.EventMaskSubstructureRedirect)
	ClientEventMask uint32 = (xproto.EventMaskPropertyChange | xproto.EventMaskFocusChange)
)

type Event struct {
	evt xgb.Event
	err xgb.Error
}

type XHandle interface {
	Pointer() Pointer
	Connectr
	Informr
	Eventr
	Windowr
	Focusr
	atomic.Atomic
	Ewmh
}

type Connectr interface {
	Conn() *xgb.Conn
}

type Informr interface {
	Setup() *xproto.SetupInfo
	Screen() *xproto.ScreenInfo
	Root() xproto.Window
	Meta() xproto.Window
	Motion() xproto.Window
}

type Eventr interface {
	Enqueue(xgb.Event, xgb.Error)
	Dequeue() (xgb.Event, xgb.Error)
	Evt(chan struct{}, chan struct{}, chan struct{})
	Empty() bool
	Quitr
}

type Quitr interface {
	Quit()
	Quitting() bool
}

type Windowr interface {
	New() *Window
	Make(xproto.Window) *Window
	Schedule(*Euclid, xproto.Window) bool
	Manage(*Euclid, *Window, ...Rule) bool
	Unmanage(*Window)
	AdoptOrphans()
}

type Focusr interface {
	SetInputFocus(*Client)
	ClearInputFocus()
}

type xhandle struct {
	conn    *xgb.Conn
	setup   *xproto.SetupInfo
	screen  *xproto.ScreenInfo
	root    xproto.Window
	meta    xproto.Window
	motion  xproto.Window
	pointer Pointer
	Events  []Event
	EvtsLck *sync.RWMutex
	quit    bool
	Windowr
	atomic.Atomic
	Ewmh
}

func mkMeta(s *xproto.ScreenInfo, c *xgb.Conn) (xproto.Window, error) {
	meta, err := xproto.NewWindowId(c)
	if err != nil {
		return meta, err
	}

	xproto.CreateWindow(
		c,
		s.RootDepth,
		meta,
		s.Root,
		-1,
		-1,
		1,
		1,
		0,
		xproto.WindowClassInputOnly,
		s.RootVisual,
		0,
		[]uint32{},
	)
	xproto.MapWindow(c, meta)
	return meta, nil
}

func mkMotion(s *xproto.ScreenInfo, c *xgb.Conn) (xproto.Window, error) {
	motion, err := xproto.NewWindowId(c)
	if err != nil {
		return motion, err
	}
	xproto.CreateWindow(
		c,
		s.RootDepth,
		motion,
		s.Root,
		0,
		0,
		s.WidthInPixels,
		s.HeightInPixels,
		0,
		xproto.WindowClassInputOnly,
		s.RootVisual,
		xproto.CwEventMask,
		[]uint32{xproto.EventMaskPointerMotion},
	)
	xproto.MapWindow(c, motion)
	return motion, nil
}

func NewXHandle(display string, ewhm []string) (*xhandle, error) {
	c, err := xgb.NewConnDisplay(display)
	if err != nil {
		return nil, err
	}

	setup := xproto.Setup(c)
	screen := setup.DefaultScreen(c)
	meta, err := mkMeta(screen, c)
	if err != nil {
		return nil, err
	}
	motion, err := mkMotion(screen, c)
	if err != nil {
		return nil, err
	}

	h := &xhandle{
		conn:    c,
		setup:   setup,
		screen:  screen,
		root:    screen.Root,
		meta:    meta,
		motion:  motion,
		Events:  make([]Event, 0, 1000),
		EvtsLck: &sync.RWMutex{},
	}

	mr := NewMotionRecorder(h.conn, h.root, h.motion)
	h.pointer = NewPointer(mr)

	h.Windowr = NewWindowr(h.conn, h.root)

	h.InputFocus = NewInputFocus(h.conn, h.root)

	h.Atomic = atomic.New(h.conn)

	EWMH := ewmh.New(h.conn, h.root, h.Atomic)
	err = EWMH.SupportedSet(ewhm)
	if err != nil {
		return nil, err
	}

	h.Ewmh = NewEwmh(EWMH)
	//h.Ewmh.Set("string name", h.root, h.meta)

	return h, nil
}

func (h *xhandle) Conn() *xgb.Conn {
	return h.conn
}

func (h *xhandle) Setup() *xproto.SetupInfo {
	return h.setup
}

func (h *xhandle) Screen() *xproto.ScreenInfo {
	return h.screen
}

func (h *xhandle) Root() xproto.Window {
	return h.root
}

func (h *xhandle) Meta() xproto.Window {
	return h.meta
}

func (h *xhandle) Motion() xproto.Window {
	return h.motion
}

func (h *xhandle) Pointer() Pointer {
	return h.pointer
}

func (h *xhandle) Empty() bool {
	h.EvtsLck.Lock()
	defer h.EvtsLck.Unlock()

	return len(h.Events) == 0
}

func (h *xhandle) Quit() {
	h.quit = true
}

func (h *xhandle) Quitting() bool {
	return h.quit
}

func (h *xhandle) Enqueue(evt xgb.Event, err xgb.Error) {
	h.EvtsLck.Lock()
	defer h.EvtsLck.Unlock()

	h.Events = append(h.Events, Event{
		evt: evt,
		err: err,
	})
}

func (h *xhandle) Dequeue() (xgb.Event, xgb.Error) {
	h.EvtsLck.Lock()
	defer h.EvtsLck.Unlock()

	e := h.Events[0]
	h.Events = h.Events[1:]
	return e.evt, e.err
}

func (h *xhandle) Evt(pre, post, quit chan struct{}) {
	for {
		if h.Quitting() {
			if quit != nil {
				quit <- struct{}{}
			}
			break
		}

		read(h)

		process(h, pre, post)
	}
}

func read(h XHandle) {
	ev, err := h.Conn().WaitForEvent()
	if ev == nil && err == nil {
		//Logger.Fatal("BUG: Could not read an event or an error.")
	}
	h.Enqueue(ev, err)
}

func process(h XHandle, pre, post chan struct{}) {
	for !h.Empty() {
		if h.Quitting() {
			return
		}

		pre <- struct{}{}

		ev, err := h.Dequeue()

		if err != nil {
			//Logger.Println(EventError(err.Error()))
			post <- struct{}{}
			continue
		}

		if ev == nil {
			//Logger.Fatal("BUG: Expected an event but got nil.")
		}

		switch event := ev.(type) {
		case xproto.MapRequestEvent:
			spew.Dump(event)
		case xproto.DestroyNotifyEvent:
			spew.Dump(event)
		case xproto.UnmapNotifyEvent:
			spew.Dump(event)
		case xproto.ClientMessageEvent:
			spew.Dump(event)
		case xproto.ConfigureRequestEvent:
			spew.Dump(event)
		case xproto.PropertyNotifyEvent:
			spew.Dump(event)
		case xproto.EnterNotifyEvent:
			spew.Dump(event)
		case xproto.MotionNotifyEvent:
			spew.Dump(event)
		case xproto.FocusInEvent:
			spew.Dump(event)
		case randr.ScreenChangeNotifyEvent:
			spew.Dump(event)
		default:
			fmt.Printf("xhandle: %+v", event)
		}

		post <- struct{}{}
	}
}

func (h *xhandle) SetInputFocus(c *Client) {
	if c == nil {
		i.ClearInputFocus()
	} else {
		if c.icccmFocus {
			//send_client_message(n->client->window, ewmh->WM_PROTOCOLS, WM_TAKE_FOCUS)
		}
		xproto.SetInputFocusChecked(h.conn, xproto.InputFocusPointerRoot, c.window.Window, xproto.TimeCurrentTime)
	}
}

func (h *xhandle) ClearInputFocus() {
	xproto.SetInputFocusChecked(h.conn, xproto.InputFocusPointerRoot, h.root, xproto.TimeCurrentTime)
}

type windowr struct {
	conn *xgb.Conn
	root xproto.Window
}

func NewWindowr(c *xgb.Conn, r xproto.Window) Windowr {
	return &windower{
		conn: c,
		root: r,
	}
}

func (w *windowr) New() *window {
	win, err := xproto.NewWindowId(w.conn)
	if err != nil {
		return nil
	}
	return &window{w.conn, win, w.root}
}

func (w *windowr) Make(win xproto.Window) *window {
	return &window{w.conn, win, w.root}
}

func (w *windowr) Schedule(e *Euclid, win xproto.Window) bool {
	/*
			var loc coordinate
			var overrideRedirect bool

			wa, _ := xproto.GetWindowAttributes(w.conn, win).Reply()
			if wa != nil {
				overrideRedirect = wa.OverrideRedirect
			}

			if !overrideRedirect {
				if _, exists := locateWindow(e, win); !exists {
					// nw := w.MakeWindow(win)
					// rules
					// return w.Manage(e, nw, rule)
				}
		    }

		    return false
	*/
}

func (w *windowr) Manage(e *Euclid, win *window, r ...Rule) bool {
	//void manage_Window(xcb_Window_t win, rule_consequence_t *csq, int fd);
}

func (w *windowr) Unmanage(win *window) {
	//void unmanage_Window(xcb_Window_t win);
}

func (w *windowr) AdoptOrphans() {}

func ClientEvent(c *xgb.Conn, root, w xproto.Window, a xproto.Atom, data ...interface{}) error {
	evMask := (xproto.EventMaskSubstructureNotify | xproto.EventMaskSubstructureRedirect)
	cm, err := mkClientMessage(32, w, a, data...)
	if err != nil {
		return err
	}

	return xproto.SendEventChecked(c, false, root, uint32(evMask), string(cm.Bytes())).Check()
}

func mkClientMessage(format byte, w xproto.Window, t xproto.Atom, data ...interface{}) (*xproto.ClientMessageEvent, error) {
	// Create the client data list first
	var clientData xproto.ClientMessageDataUnion

	// Don't support formats 8 or 16 yet.
	switch format {
	case 8:
		buf := make([]byte, 20)
		for i := 0; i < 20; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = data[i].(byte)
		}
		clientData = xproto.ClientMessageDataUnionData8New(buf)
	case 16:
		buf := make([]uint16, 10)
		for i := 0; i < 10; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = uint16(data[i].(int16))
		}
		clientData = xproto.ClientMessageDataUnionData16New(buf)
	case 32:
		buf := make([]uint32, 5)
		for i := 0; i < 5; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = uint32(data[i].(int))
		}
		clientData = xproto.ClientMessageDataUnionData32New(buf)
	default:
		return nil, fmt.Errorf("mkClientMessage: Unsupported format '%d'.", format)
	}

	return &xproto.ClientMessageEvent{
		Format: format,
		Window: w,
		Type:   t,
		Data:   clientData,
	}, nil
}

/*
void set_floating_atom(xcb_Window_t win, uint32_t value);
void enable_floating_atom(xcb_Window_t win);
void disable_floating_atom(xcb_Window_t win);
void get_atom(char *name, xcb_atom_t *atom);
void set_atom(xcb_Window_t win, xcb_atom_t atom, uint32_t value);
bool has_proto(xcb_atom_t atom, xcb_icccm_get_wm_protocols_reply_t *protocols);
void send_client_message(xcb_Window_t win, xcb_atom_t property, xcb_atom_t value);


*/ //uint32_t values[] = {ROOT_EVENT_MASK};
//xcb_generic_error_t *e = xcb_request_check(dpy, xcb_change_window_attributes_checked(dpy, root, XCB_CW_EVENT_MASK, values));
//var values []uint32 = []uint32{RootEventMask}
//err = xproto.ChangeWindowAttributesChecked(c, screen.Root, xproto.CwEventMask, values).Check()
//if err != nil {
//	spew.Dump(err)
//}
