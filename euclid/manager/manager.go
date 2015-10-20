package manager

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/thrisp/scpwm/euclid/client"
	"github.com/thrisp/scpwm/euclid/commander"
	"github.com/thrisp/scpwm/euclid/handler"
	"github.com/thrisp/scpwm/euclid/monitor"
	"github.com/thrisp/scpwm/euclid/ruler"
	"github.com/thrisp/scpwm/euclid/settings"
)

type Manager struct {
	*log.Logger
	handler.Handler
	*settings.Settings
	*Loop
	ruler.Ruler
	commander.Commander
	monitors monitor.Monitors
	clients  []client.Clients
	//history  *History
}

func New() *Manager {
	l := newLoop()

	m := &Manager{
		Settings:  settings.DefaultSettings(),
		Ruler:     ruler.New(),
		Commander: commander.New(l.Comm),
		Logger:    log.New(os.Stderr, "[SCPWM] ", log.Ldate|log.Lmicroseconds),
	}

	hndl, err := handler.New("", settings.EwmhSupported)
	if err != nil {
		panic(err)
	}
	m.Handler = hndl

	m.monitors = monitor.NewMonitors(m.Settings, m.Handler)

	m.InitClients()

	m.Loop = l

	return m
}

type Loop struct {
	Pre  chan struct{}
	Post chan struct{}
	Quit chan struct{}
	Comm chan string
}

func newLoop() *Loop {
	return &Loop{
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan string, 0),
	}
}

func (m *Manager) lstn(l *net.UnixListener) {
	for {
		conn, err := l.AcceptUnix()
		if err != nil {
			panic(err)
		}
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		r := bytes.Trim(buf[:n], " ")
		resp := m.Process(r, m.Settings, m.Handler)
		conn.Write(resp)
		conn.Close()
	}
}

func (m *Manager) Looping(l *net.UnixListener) *Loop {
	go func() {
		m.lstn(l)
	}()
	lp := m.Loop
	go func() {
		m.Handler.Handle(lp.Pre, lp.Post, lp.Quit)
	}()
	return lp
}

func (m *Manager) InitClients() {
	m.clients = make([]client.Clients, 0)
	for i, m := range m.monitors.All() {
		name := fmt.Sprintf("%s%d", m.Settings.String("DefaultDesktopName"), i)
		m.clients = append(m.clients, NewClients(m.Id(), name, m.Settings))
		//translate client to monitor
	}
}

/*
func (e *Euclid) Select(item string, description string) (coordinate, bool) {
	var obj Node
	loc := Coordinate(e, nil, nil, nil)
	var sel []Selector
	switch item {
	case "monitor":
		sel = Selectors(find, nMonitor, description)
		obj = nMonitor
	case "desktop":
		sel = Selectors(find, nDesktop, description)
		obj = nDesktop
	case "client":
		sel = Selectors(find, nClient, description)
		obj = nClient
	}
	return e.get(obj, loc, sel...)
}

func (e *Euclid) get(obj Node, loc coordinate, sel ...Selector) (coordinate, bool) {
	switch obj {
	case nMonitor:
		//return selectMonitor(e, loc, sel...)
	case nDesktop:
		//return selectDesktop(e, loc, sel...)
	case nClient:
		//return selectClient(e, loc, sel...)
	}
	return loc, false
}

func (e *Euclid) Match(item, description string, location, reference *coordinate) bool {
	var sel []Selector
	switch item {
	case "monitor":
		sel = Selectors(match, nMonitor, description)
	case "desktop":
		sel = Selectors(match, nDesktop, description)
	case "client":
		sel = Selectors(match, nClient, description)
	}
	return e.match(location, reference, sel...)
}

func (e *Euclid) match(location, reference *coordinate, sel ...Selector) bool {
	return false
}

//func locateWindow(h *Handler, w xproto.Window) (*Client, bool) {
//for _, m := range h.monitors {
//	for _, d := range m.desktops {
//		for _, c := range d.clients {
//			if w == c.Window.Window {
//				return c, true
//			}
//		}
//	}
//}
//return nil, false
//}
*/
