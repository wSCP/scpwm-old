package client

import (
	"github.com/thrisp/scpwm/euclid/settings"
)

type Clients interface {
	Id() uint32
	Name() string
	Settings(*settings.Settings)
	Add(Client)
	Pop(Client) Client
	Swap(int, int)
	Merge(Clients)
	Number() int
	Desktop
	Stack
}

type Desktop interface {
	Set(*settings.Settings)
	Gap() int
	Border() uint
	Pad() [4]int
	Focus()
	Focused() bool
	Last() bool
}

type desktop struct {
	//layout      Layout
	pad         [4]int
	windowGap   int
	borderWidth uint
	floating    bool
	focused     bool
	last        bool
}

func (d *desktop) Set(s *settings.Settings) {
	//pad:,
	d.windowGap = s.Int("WindowGap")
	d.borderWidth = uint(s.Int("BorderWidth"))
}

func (d *desktop) Gap() int {
	return d.windowGap
}

func (d *desktop) Border() uint {
	return d.borderWidth
}

func (d *desktop) Pad() [4]int {
	return d.pad
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

type clients struct {
	id   uint32
	name string
	c    []Client
	Desktop
	Stack
}

func NewClients(id uint32, name string, s *settings.Settings) Clients {
	cs := &clients{
		id:      id,
		name:    name,
		c:       make([]Client, 0),
		Desktop: &desktop{},
	}
	cs.Settings(s)
	return cs
}

func (cs *clients) Id() uint32 {
	return cs.id
}

func (cs *clients) Name() string {
	return cs.name
}

func (cs *clients) Settings(s *settings.Settings) {
	cs.Desktop.Set(s)
}

func (cs *clients) Add(c Client) {
	cs.c = append(cs.c, c)
}

func (cs *clients) Pop(p Client) Client {
	for _, _ = range cs.c {
		//if c.Id() == p.Id() {
		//	cs.c = cs.c[:i+copy(cs.c[i:], cs.c[i+1:])]
		//}
		//return c
		// customize comparator of clients
	}
	return nil
}

func (cs *clients) Swap(i, j int) {
	cs.c[i], cs.c[j] = cs.c[j], cs.c[i]
}

func (cs *clients) Merge(o Clients) {}

func (cs *clients) Number() int {
	return len(cs.c)
}

func (cs *clients) Min() int {
	var ret int
	for _, c := range cs.c {
		st := c.GetStack()
		if st <= ret {
			ret = st
		}
	}
	return ret
}

func (cs *clients) Max() int {
	var ret int
	for _, c := range cs.c {
		st := c.GetStack()
		if st <= ret {
			ret = st
		}
	}
	return ret
}

func (cs *clients) Reset() {
	for i, c := range cs.c {
		c.SetStack(i)
	}
}
