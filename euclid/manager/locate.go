package manager

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/clients"
	"github.com/thrisp/scpwm/euclid/desktops"
	"github.com/thrisp/scpwm/euclid/monitors"
)

func (m *Manager) exists(win xproto.Window) bool {
	_, exists := m.locateClient(win)
	return exists
}

type Location struct {
	m monitors.Monitor
	d desktops.Desktop
	c clients.Client
}

func (m *Manager) Current() Location {
	var mon monitors.Monitor
	var dsk desktops.Desktop
	var clt clients.Client
	mon = monitors.Focused(m.Branch)
	dsk = desktops.Focused(mon.Desktops())
	clt = clients.Focused(dsk.Clients())
	return Location{mon, dsk, clt}
}

//func (m *Manager) Locate(origin Location, sel ...Selector) (Location, bool) {
//curr := &loc
//for _, mon := range monitors.All(m.Branch) {
//curr.m = mon
//if fn(loc) {
//	return loc, true
//}
//for _, desk := range desktops.All(mon.Desktops()) {
//curr.d = desk
//if fn(loc) {
//	return loc, true
//}
//for _, cli := range clients.All(desk.Clients()) {
//curr.c = client
//if fn(loc) {
//	return loc, true
//}
//}
//}
//}
//return origin, false
//}

func (m *Manager) locateClient(win xproto.Window) (clients.Client, bool) {
	for _, c := range m.Clients() {
		if c.XWindow() == win {
			return c, true
		}
	}
	return nil, false
}
