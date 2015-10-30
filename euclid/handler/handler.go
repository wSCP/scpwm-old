package handler

import (
	"log"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/atomic"
	"github.com/thrisp/scpwm/euclid/ewmh"
)

type Handler interface {
	Connectr
	Informr
	Eventr
	Callr
	atomic.Atomic
	ewmh.Ewmh
}

type Connectr interface {
	Conn() *xgb.Conn
}

type Informr interface {
	Setup() *xproto.SetupInfo
	Screen() *xproto.ScreenInfo
	Root() xproto.Window
}

type Eventr interface {
	Enqueue(xgb.Event, xgb.Error)
	Dequeue() (xgb.Event, xgb.Error)
	Handle(chan struct{}, chan struct{}, chan struct{})
	Empty() bool
	Endr
}

type Endr interface {
	End()
	Ending() bool
}

type Call func(xgb.Event) error

type Callr interface {
	SetEventFn(string, Call)
	Call(string, xgb.Event)
}

type handler struct {
	*log.Logger
	conn    *xgb.Conn
	setup   *xproto.SetupInfo
	screen  *xproto.ScreenInfo
	root    xproto.Window
	events  []evnt
	evtsLck *sync.RWMutex
	call    map[string]Call
	callLck *sync.RWMutex
	end     bool
	atomic.Atomic
	ewmh.Ewmh
}

func New(display string, ewhm []string, logr *log.Logger) (Handler, error) {
	c, err := xgb.NewConnDisplay(display)
	if err != nil {
		return nil, err
	}

	setup := xproto.Setup(c)
	screen := setup.DefaultScreen(c)

	h := &handler{
		Logger:  logr,
		conn:    c,
		setup:   setup,
		screen:  screen,
		root:    screen.Root,
		events:  make([]evnt, 0, 1000),
		evtsLck: &sync.RWMutex{},
		call:    make(map[string]Call),
		callLck: &sync.RWMutex{},
	}

	h.Atomic = atomic.New(h.conn)
	h.Atomic.Atom("WM_DELETE_WINDOW")
	h.Atomic.Atom("WM_TAKE_FOCUS")
	h.Atomic.Atom("_SCPWM_FLOATING_WINDOW")

	e := ewmh.New(h.conn, h.root, h.Atomic)
	//err = EWMH.SupportedSet(ewhm)
	//if err != nil {
	//	return nil, err
	//}
	//h.Ewmh.Set("string name", h.root, h.meta)
	h.Ewmh = e

	return h, nil
}

func (h *handler) Conn() *xgb.Conn {
	return h.conn
}

func (h *handler) Setup() *xproto.SetupInfo {
	return h.setup
}

func (h *handler) Screen() *xproto.ScreenInfo {
	return h.screen
}

func (h *handler) Root() xproto.Window {
	return h.root
}

func (h *handler) Empty() bool {
	h.evtsLck.Lock()
	defer h.evtsLck.Unlock()

	return len(h.events) == 0
}

func (h *handler) End() {
	h.end = true
}

func (h *handler) Ending() bool {
	return h.end
}

type evnt struct {
	evt xgb.Event
	err xgb.Error
}

func (h *handler) Enqueue(evt xgb.Event, err xgb.Error) {
	h.evtsLck.Lock()
	defer h.evtsLck.Unlock()

	h.events = append(h.events, evnt{
		evt: evt,
		err: err,
	})
}

func (h *handler) Dequeue() (xgb.Event, xgb.Error) {
	h.evtsLck.Lock()
	defer h.evtsLck.Unlock()

	e := h.events[0]
	h.events = h.events[1:]
	return e.evt, e.err
}

func (h *handler) Handle(pre, post, quit chan struct{}) {
	for {
		if h.Ending() {
			if quit != nil {
				quit <- struct{}{}
			}
			break
		}

		h.read()

		h.process(pre, post)
	}
}

func (h *handler) read() {
	ev, err := h.Conn().WaitForEvent()
	if ev == nil && err == nil {
		h.Fatal("BUG: Could not read an event or an error.")
	}
	h.Enqueue(ev, err)
}

func (h *handler) process(pre, post chan struct{}) {
	for !h.Empty() {
		if h.Ending() {
			return
		}

		pre <- struct{}{}

		ev, err := h.Dequeue()

		if err != nil {
			h.Println(err.Error())
			post <- struct{}{}
			continue
		}

		if ev == nil {
			h.Fatal("BUG: Expected an event but got nil.")
		}

		var tag string
		switch ev.(type) {
		case xproto.MapRequestEvent:
			tag = "MapRequest"
		case xproto.DestroyNotifyEvent:
			tag = "DestroyNotify"
		case xproto.UnmapNotifyEvent:
			tag = "UnmapNotify"
		case xproto.ClientMessageEvent:
			tag = "ClientMessage"
		case xproto.ConfigureRequestEvent:
			tag = "ConfigureRequest"
		case xproto.PropertyNotifyEvent:
			tag = "PropertyNotify"
		case xproto.EnterNotifyEvent:
			tag = "EnterNotify"
		case xproto.MotionNotifyEvent:
			tag = "MotionNotify"
		case xproto.FocusInEvent:
			tag = "FocusIn"
		case randr.ScreenChangeNotifyEvent:
			tag = "ScreenChange"
		}

		if tag != "" {
			h.Call(tag, ev)
		}

		post <- struct{}{}
	}
}

func (h *handler) SetEventFn(tag string, fn Call) {
	h.callLck.Lock()
	defer h.callLck.Unlock()
	h.call[tag] = fn
}

func (h *handler) Call(tag string, evt xgb.Event) {
	h.callLck.Lock()
	defer h.callLck.Unlock()
	if fn, ok := h.call[tag]; ok {
		err := fn(evt)
		if err != nil {
			h.Println("ERROR: %s", err.Error())
		}
	}
}

/*
//mr := NewMotionRecorder(h.conn, h.root, h.motion)
//h.Pointer = NewPointer(mr)

//h.Windowr = NewWindowr(h.conn, h.root)

//h.Monitors = NewMonitors(h)
*/

/*
//func (h *handler) Meta() xproto.Window {
//	return h.meta
//}

//func (h *handler) Motion() xproto.Window {
//	return h.motion
//}

//meta:    meta,
//motion:  motion,
//meta, err := mkMeta(screen, c)
//if err != nil {
//	return nil, err
//}
//motion, err := mkMotion(screen, c)
//if err != nil {
//	return nil, err
//}

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


//var (
//	RootEventMask   uint32 = (xproto.EventMaskSubstructureNotify | xproto.EventMaskSubstructureRedirect)
//	ClientEventMask uint32 = (xproto.EventMaskPropertyChange | xproto.EventMaskFocusChange)
//)

//func (h *handler) SetInputFocus(c *Client) {
//	if c == nil {
//		h.ClearInputFocus()
//	} else {
//		if c.icccmFocus {
//			//send_client_message(n->client->window, ewmh->WM_PROTOCOLS, WM_TAKE_FOCUS)
//		}
//		xproto.SetInputFocusChecked(h.conn, xproto.InputFocusPointerRoot, c.Window.Window, xproto.TimeCurrentTime)
//	}
//}

//func (h *handler) ClearInputFocus() {
//	xproto.SetInputFocusChecked(h.conn, xproto.InputFocusPointerRoot, h.root, xproto.TimeCurrentTime)
//}

//type windowr struct {
//	conn *xgb.Conn
//	root xproto.Window
//}

//func NewWindowr(c *xgb.Conn, r xproto.Window) Windowr {
//	return &windowr{
//		conn: c,
//		root: r,
//	}
//}

//func (w *windowr) New() *Window {
//	win, err := xproto.NewWindowId(w.conn)
//	if err != nil {
//		return nil
//	}
//	return &Window{w.conn, win, w.root}
//}

//func (w *windowr) Make(win xproto.Window) *Window {
//	return &Window{w.conn, win, w.root}
//}

//func (w *windowr) Schedule(r ruler.Ruler, win xproto.Window) bool {
/*
		var loc coordinate
		var overrideRedirect bool

		wa, _ := xproto.GetWindowAttributes(w.conn, win).Reply()
		if wa != nil {
			overrideRedirect = wa.OverrideRedirect
		}

		if !overrideRedirect {
			if _, exists := locateWindow(e, win); !exists {
				// nw := w.Make(win)
				// rules
				// return w.Manage(e, nw, rule)
			}
	    }

	    return false
*/
//	return false
//}

//func (w *windowr) Manage(win *Window, r ...ruler.Rule) bool {
//void manage_Window(xcb_Window_t win, rule_consequence_t *csq, int fd);
//	return false
//}

//func (w *windowr) Unmanage(win *Window) {
//void unmanage_Window(xcb_Window_t win);
//}

//func (w *windowr) AdoptOrphans() {
//	if qtr, _ := xproto.QueryTree(w.conn, w.root).Reply(); qtr != nil {
//int len = xcb_query_tree_children_length(qtr);
//xcb_window_t *wins = xcb_query_tree_children(qtr);
//for (int i = 0; i < len; i++) {
//	uint32_t idx;
//	xcb_window_t win = wins[i];
//	if (xcb_ewmh_get_wm_desktop_reply(ewmh, xcb_ewmh_get_wm_desktop(ewmh, win), &idx, NULL) == 1)
//		schedule_window(win);
//}
//	}
//}

/*
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

void set_floating_atom(xcb_Window_t win, uint32_t value);
void enable_floating_atom(xcb_Window_t win);
void disable_floating_atom(xcb_Window_t win);
void get_atom(char *name, xcb_atom_t *atom);
void set_atom(xcb_Window_t win, xcb_atom_t atom, uint32_t value);
bool has_proto(xcb_atom_t atom, xcb_icccm_get_wm_protocols_reply_t *protocols);
void send_client_message(xcb_Window_t win, xcb_atom_t property, xcb_atom_t value);


//uint32_t values[] = {ROOT_EVENT_MASK};
//xcb_generic_error_t *e = xcb_request_check(dpy, xcb_change_window_attributes_checked(dpy, root, XCB_CW_EVENT_MASK, values));
//var values []uint32 = []uint32{RootEventMask}
//err = xproto.ChangeWindowAttributesChecked(c, screen.Root, xproto.CwEventMask, values).Check()
//if err != nil {
//	spew.Dump(err)
//}

*/
//type Windowr interface {
//	New() *Window
//	Make(xproto.Window) *Window
//	Schedule(ruler.Ruler, xproto.Window) bool
//	Manage(*Window, ...ruler.Rule) bool
//	Unmanage(*Window)
//	AdoptOrphans()
//}

//type Focusr interface {
//	SetInputFocus(*Client)
//	ClearInputFocus()
//}
