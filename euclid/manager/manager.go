package manager

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
		//propagate down, apply across monitors, desktops unless overrode locally
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
