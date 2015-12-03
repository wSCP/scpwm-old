package manager

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/clients"
	"github.com/thrisp/scpwm/euclid/desktops"
	"github.com/thrisp/scpwm/euclid/monitors"
	"github.com/thrisp/scpwm/euclid/selector"
)

type Location struct {
	m  monitors.Monitor
	ms *branch.Branch
	d  desktops.Desktop
	ds *branch.Branch
	c  clients.Client
	cs *branch.Branch
}

func (m *Manager) Current() Location {
	var mon monitors.Monitor
	var dsk desktops.Desktop
	var clt clients.Client
	ms := m.Branch
	mon = monitors.Focused(ms)
	ds := mon.Desktops()
	dsk = desktops.Focused(ds)
	cs := dsk.Clients()
	clt = clients.Focused(cs)
	return Location{mon, ms, dsk, ds, clt, cs}
}

func (m *Manager) Locate(origin Location, sel ...selector.Selector) (Location, bool) {
	curr := &origin
	curr.ms = m.Branch
	for _, mon := range monitors.All(curr.ms) {
		if fmon := monitors.Select(sel, curr.ms); fmon != nil {
			curr.m = fmon
			curr.ds = fmon.Desktops()
			curr.d = desktops.Focused(curr.ds)
			curr.cs = curr.d.Clients()
			curr.c = clients.Focused(curr.cs)
			return *curr, true
		}
		curr.m = mon
		curr.ds = mon.Desktops()
		for _, desk := range desktops.All(curr.ds) {
			if fdesk := desktops.Select(sel, curr.ds); fdesk != nil {
				curr.d = fdesk
				curr.cs = curr.d.Clients()
				curr.c = clients.Focused(curr.cs)
				return *curr, true
			}
			curr.d = desk
			curr.cs = desk.Clients()
			if fclient := clients.Select(sel, curr.cs); fclient != nil {
				curr.c = fclient
				return *curr, true
			}
		}
	}
	curr = nil
	return origin, false
}

func (m *Manager) locate(win xproto.Window) (clients.Client, bool) {
	for _, c := range m.Clients() {
		if c.XWindow() == win {
			return c, true
		}
	}
	return nil, false
}

func (m *Manager) exists(win xproto.Window) bool {
	_, exists := m.locate(win)
	return exists
}

func (m *Manager) LocationWindow(win xproto.Window) (Location, bool) {
	var loc Location
	for _, mon := range monitors.All(m.Branch) {
		loc.m = mon
		loc.ms = m.Branch
		ds := mon.Desktops()
		for _, d := range desktops.All(ds) {
			loc.d = d
			loc.ds = ds
			cs := d.Clients()
			loc.cs = cs
			for _, c := range clients.All(cs) {
				if c.XWindow() == win {
					loc.c = c
					return loc, true
				}
			}
		}
	}
	return loc, false
}

func (m *Manager) LocationDesktopIndex(i int) (Location, bool) {
	var loc Location
	for _, mon := range m.Monitors() {
		for _, d := range desktops.All(mon.Desktops()) {
			if i == d.Index() {
				loc.m = mon
				loc.d = d
				loc.c = nil
				return loc, true
			}
		}
	}
	return loc, false
}
