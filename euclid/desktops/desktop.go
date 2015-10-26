package desktops

import (
	"fmt"

	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/settings"
)

type Desktop interface {
	settings.Settings
	Set(string, interface{})
	Index() int
	Monitor() uint32
	Id() string
	Name() string
	Focus()
	Focused() bool
	Last() bool
	Clients() *branch.Branch
}

func NewDesktop(monitor uint32, name string, s settings.Settings) Desktop {
	d := &desktop{
		Settings: s,
	}
	d.setName(name)
	d.setMonitor(monitor)
	return d
}

type desktop struct {
	settings.Settings
	idx     int
	monitor uint32
	name    string
	focused bool
	last    bool
	clients *branch.Branch
}

func (d *desktop) Set(k string, v interface{}) {
	switch k {
	case "index":
		if value, ok := v.(int); ok {
			d.setIndex(value)
		}
	case "monitor":
		if value, ok := v.(uint32); ok {
			d.setMonitor(value)
		}
	case "name":
		if value, ok := v.(string); ok {
			d.setName(value)
		}
	}
}

func (d *desktop) Index() int {
	return d.idx
}

func (d *desktop) setIndex(idx int) {
	d.idx = idx
}

func (d *desktop) Monitor() uint32 {
	return d.monitor
}

func (d *desktop) setMonitor(as uint32) {
	d.monitor = as
}

func (d *desktop) Id() string {
	return fmt.Sprintf("%s-%d-%d", d.name, d.idx, d.monitor)
}

func (d *desktop) Name() string {
	if d.name == "" {
		return d.defaultName()
	}
	return d.name
}

func (d *desktop) defaultName() string {
	return fmt.Sprintf("%s%d", d.String("DefaultDesktopName"), d.idx)
}

func (d *desktop) setName(name string) {
	d.name = name
}

func (d *desktop) Focus() {
	d.focused = true
}

func (d *desktop) Focused() bool {
	return d.focused
}

func (d *desktop) Last() bool {
	return d.last
}

func (d *desktop) Clients() *branch.Branch {
	return d.clients
}
