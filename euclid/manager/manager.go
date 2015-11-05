package manager

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/clients"
	"github.com/thrisp/scpwm/euclid/commander"
	"github.com/thrisp/scpwm/euclid/desktops"
	"github.com/thrisp/scpwm/euclid/handler"
	"github.com/thrisp/scpwm/euclid/monitors"
	"github.com/thrisp/scpwm/euclid/rules"
	"github.com/thrisp/scpwm/euclid/settings"
)

type Manager struct {
	*log.Logger
	handler.Handler
	settings.Settings
	*Loops
	rules.Ruler
	commander.Commander
	*branch.Branch
}

func New() *Manager {
	l := newLoops()

	m := &Manager{
		Settings:  settings.DefaultSettings(),
		Ruler:     rules.New(),
		Commander: commander.New(l.Comm),
		Logger:    log.New(os.Stderr, "[SCPWM] ", log.Ldate|log.Lmicroseconds),
	}

	hndl, err := handler.New("", settings.EwmhSupported, m.Logger)
	if err != nil {
		panic(err)
	}
	m.Handler = hndl
	m.SetEventFns()

	m.Branch = monitors.New(m.Handler, m.Settings)

	m.Loops = l

	//history  *History

	return m
}

type Loops struct {
	Pre  chan struct{}
	Post chan struct{}
	Quit chan struct{}
	Comm chan string
	Sys  chan os.Signal
}

func newLoops() *Loops {
	return &Loops{
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan string, 0),
		make(chan os.Signal, 0),
	}
}

func (m *Manager) Looping(l *net.UnixListener) *Loops {
	lp := m.Loops

	go func() {
		m.Commander.Listen(l, m)
	}()

	go func() {
		m.Handler.Handle(lp.Pre, lp.Post, lp.Quit)
	}()

	signal.Notify(
		lp.Sys,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGCHLD,
		syscall.SIGPIPE,
	)

	return lp
}

func (m *Manager) SignalHandler(sig os.Signal) {
	switch sig {
	case syscall.SIGHUP:
		m.Comm <- fmt.Sprintf("received signal SIGHUP, reconfiguring...")
		cp := m.String("ConfigPath")
		err := m.LoadConfig(cp)
		//propagate down, apply across monitors, desktops
		if err != nil {
			m.Println(err.Error())
		}
	case syscall.SIGINT, syscall.SIGTERM:
		m.Println(sig)
		os.Exit(0)
	case syscall.SIGCHLD, syscall.SIGPIPE:
		m.Println(sig)
	}
}

func (m *Manager) Tree() *branch.Branch {
	return m.Branch
}

func (m *Manager) Monitors() []monitors.Monitor {
	return monitors.All(m.Branch)
}

func (m *Manager) Desktops() []desktops.Desktop {
	var ret []desktops.Desktop
	ms := m.Monitors()
	for _, mon := range ms {
		ret = append(ret, desktops.All(mon.Desktops())...)
	}
	return ret
}

func (m *Manager) Clients() []clients.Client {
	var ret []clients.Client
	ds := m.Desktops()
	for _, d := range ds {
		ret = append(ret, clients.All(d.Clients())...)
	}
	return ret
}

func (m *Manager) SetEventFns() {
	m.SetEventFn("MapRequest", m.MapRequest)
	m.SetEventFn("DestroyNotify", m.DestroyNotify)
	m.SetEventFn("UnmapNotify", m.UnmapNotify)
	m.SetEventFn("ClientMessage", m.ClientMessage)
	m.SetEventFn("ConfigureRequest", m.ConfigureRequest)
	m.SetEventFn("PropertyNotify", m.PropertyNotify)
	m.SetEventFn("EnterNotify", m.EnterNotify)
	m.SetEventFn("MotionNotify", m.MotionNotify)
	m.SetEventFn("FocusInEvent", m.FocusIn)
	m.SetEventFn("ScreenChange", m.ScreenChange)
}

var EventError = Xrror("Event : %+v is not recognized despite being passed to Manager function %s.").Out

func (m *Manager) MapRequest(evt xgb.Event) error {
	if mr, ok := evt.(xproto.MapRequestEvent); ok {
		m.schedule(mr.Window)
	}
	return EventError(evt, "MapRequest")
}

func (m *Manager) schedule(win xproto.Window) {
	var overrideRedirect bool
	wa, _ := xproto.GetWindowAttributes(m.Conn(), win).Reply()
	if wa != nil {
		overrideRedirect = wa.OverrideRedirect
	}

	if !overrideRedirect && !m.exists(win) {
		if class, instance, err := m.WmClassGet(win); err == nil {
			rules := m.Ruler.Applicable(class, instance)
			m.manage(win, rules)
		}
	}
}

func (m *Manager) manage(win xproto.Window, rls []rules.Rule) {
	//loc := m.Current()
	//for _, r := range rls {
	//	ef := r.Effect()
	//	switch ef.Applies() {
	//	case rules.IsManaged:
	//	case rules.ClientSelector:
	//loc, _ = m.Locate(loc, Selectors(ef.String())...)
	//	case rules.DesktopSelector:
	//loc, _ = m.Locate(loc, Selectors(ef.String())...)
	//	case rules.MonitorSelector:
	//		//loc = m.Locate(loc, Selectors(ef.String())...)
	//	case rules.IsSticky:
	//	case rules.SetSplitDirection:
	//	//if f != nil {
	//	//		//f.splitM = manual
	//	//		//f.splitd = r.String() to split direction
	//	//}
	//	case rules.SetSplitRatio:
	//		//if f != nil {
	//		//		//f.splitR = rule.splitRatio
	//		//}
	//	}
	//}

	//client := clients.NewClient(class, instance string, c *xgb.Conn, x xproto.Window, r xproto.Window)
	// update client rectangle, max widths & heights
	// embrace & translate to monitor

	// add client to the appropriate tree on top of focus

	// adjust rule floating to desktop floating

	// set client state based on rules

	// arrange desktop

	// if visible desktop -- show window

	// update ewmh
	// ewmh_set_wm_desktop(n, d);
	// ewmh_update_client_list();
}

func (m *Manager) ConfigureRequest(evt xgb.Event) error {
	if cr, ok := evt.(xproto.ConfigureRequestEvent); ok {
		client, exists := m.locateClient(cr.Window)
		//var w, h uint16
		if exists && client.Tiled() {
			//if (cr.ValueMask & xproto.ConfigWindowX) != 0 {
			//c->floating_rectangle.x = e->x;
			//}
			//if (cr.ValueMask & xproto.ConfigWindowY) != 0 {
			//c->floating_rectangle.y = e->y
			//}
			//if (cr.ValueMask & xproto.ConfigWindowHeight) != 0 {
			//	w = cr.Width
			//}
			//if (cr.ValueMask & xproto.ConfigWindowWidth) != 0 {
			//	h = cr.Height
			//}
			//if w != 0 {
			//restrain_floating_width(c, &w);
			//c->floating_rectangle.width = w;
			//}
			//if h != 0 {
			//restrain_floating_height(c, &h);
			//c->floating_rectangle.height = h;
			//}

			//var evt xproto.ConfigureNotifyEvent
			//var rect xproto.Rectangle
			//win := client.XWindow()
			//bw := client.BorderWidth()

			//if (c->fullscreen)
			//	rect = loc.monitor->rectangle;
			//else
			//	rect = c->tiled_rectangle;

			//evt.Event = win
			//evt.Window = win
			//evt.AboveSibling = xproto.WindowNone
			//evt.X = rect.X
			//evt.Y = rect.Y
			//evt.Width = rect.Width
			//evt.Height = rect.Height
			//evt.BorderWidth = bw
			//evt.OverrideRedirect = false

			//xproto.SendEvent(m.Conn(), false, win, EventMask uint32, evt.String())
			//xcb_send_event(dpy, false, win, XCB_EVENT_MASK_STRUCTURE_NOTIFY, (const char *) &evt);

			//if (c->pseudo_tiled)
			//	arrange(loc.monitor, loc.desktop);
		} else {
			//var mask uint16
			//var value [7]uint32
			//var i int
			//if (cr.ValueMask & xproto.ConfigWindowX) != 0 {
			//mask |= xproto.ConfigWindowX
			//values[i] = cr.X
			//if exists {
			//c->floating_rectangle.x = e->x
			//}
			//}
			//if (cr.ValueMask & xproto.ConfigWindowY) != 0 {
			//mask |= xproto.ConfigWindowY
			//i++
			//value[i] = cr.Y
			//if exists {
			//	c->floating_rectangle.y = e->y;
			//}
			//}
			//if (cr.ValueMask & xproto.ConfigWindowHeight) != 0 {
			//mask |= xproto.ConfigWindowHeight
			//w = cr.Width
			//if exists {
			//restrain_floating_width(c, &w);
			//c->floating_rectangle.width = w;
			//}
			//i++
			//values[i] = cr.Height

			//}
			//if (cr.ValueMask & xproto.ConfigWindowWidth) != 0 {
			//mask |= xproto.ConfigWindowWidth
			//h = cr.Height
			//if exists {
			//restrain_floating_height(c, &h);
			//c->floating_rectangle.height = h;
			//}
			//i++
			//values[i] = cr.Width
			//}
			//if !exists && (cr.ValueMask & xproto.ConfigWindowBorderWidth) != 0 {
			// mask |= xproto.ConfigWindowBorderWidth
			// i++
			// values[i] = cr.BorderWidth
			//}
			//if (cr.ValueMask & xproto.ConfigWindowSibling) != 0 {
			// mask |= xproto.ConfigWindowSibling
			// i++
			//value[i] = cr.WindowSibling
			//}
			//if (cr.ValueMask & xproto.ConfigWindowStackMode) != 0 {
			// mask |= xproto.ConfigWindowStackMode
			// i++
			// values[i] = cr.WindowStackMode
			//}

			//xproto.ConfigureWindow(m.Conn(), cr.Window, mask, values)
		}
		//if exists {
		//translate_client(monitor_from_client(c), loc.monitor, c);
		//}
		return nil
	}
	return EventError(evt, "ConfigureRequest")
}

func (m *Manager) unmanage(win xproto.Window) error {
	return nil
}

func (m *Manager) DestroyNotify(evt xgb.Event) error {
	if dn, ok := evt.(xproto.DestroyNotifyEvent); ok {
		return m.unmanage(dn.Window)
	}
	return EventError(evt, "DestroyNotify")
}

func (m *Manager) UnmapNotify(evt xgb.Event) error {
	if un, ok := evt.(xproto.UnmapNotifyEvent); ok {
		return m.unmanage(un.Window)
	}
	return EventError(evt, "UnmapNotify")
}

func (m *Manager) PropertyNotify(evt xgb.Event) error {
	if pn, ok := evt.(xproto.PropertyNotifyEvent); ok {
		if pn.Atom == xproto.AtomWmHints || pn.Atom == xproto.AtomWmNormalHints {
			if m.exists(pn.Window) {
				switch pn.Atom {
				case xproto.AtomWmHints:
					//
				case xproto.AtomWmNormalHints:
					//
				}
				return nil
			}
		}
		return nil
	}
	return EventError(evt, "PropertyNotify")
}

func (m *Manager) ClientMessage(evt xgb.Event) error {
	//if cm, ok := evt.(xproto.ClientMessageEvent); ok {
	//var a xproto.Atom
	//var err error
	//if a, err = m.Atomic.Atom("_NET_CURRENT_DESKTOP"); err == nil {
	//if cm.Type == a {
	//}
	//}
	//if exists = m.locate(cm.Window); exists {
	//if a, err = m.Atomic.Atom("_NET_WM_STATE"); e.Type == a {
	//}
	//if a, err = m.Atomic.Atom("_NET_ACTIVE_WINDOW"); e.Type == a {
	//}
	//if a, err = m.Atomic.Atom("_NET_WM_DESKTOP"); e.Type == a  {
	//}
	//if a, err = m.Atomic.Atom("_NET_CLOSE_WINDOW"); e.Type == a {
	//}
	//}
	//return err
	//}
	return EventError(evt, "ClientMessage")
}

func (m *Manager) FocusIn(evt xgb.Event) error {
	//if fi, ok := evt.(xproto.FocusInEvent); ok {
	//	return nil
	//}
	return EventError(evt, "FocusIn")
}

func (m *Manager) EnterNotify(evt xgb.Event) error {
	//if en, ok := evt.(xproto.EnterNotifyEvent); ok {
	//	return nil
	//}
	return EventError(evt, "EnterNotify")
}

func (m *Manager) MotionNotify(evt xgb.Event) error {
	//if mn, ok := evt.(xproto.MotionNotifyEvent); ok {
	//	return nil
	//}
	return EventError(evt, "MotionNotify")
}

func (m *Manager) ScreenChange(evt xgb.Event) error {
	return monitors.Update(m.Branch, m.Handler, m.Settings)
}
