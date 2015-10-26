package monitors

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"

	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/desktops"
	"github.com/thrisp/scpwm/euclid/settings"
	"github.com/thrisp/scpwm/euclid/windows"
)

type Monitor interface {
	Id() uint32
	Name() string
	Rectangle() xproto.Rectangle
	SetRectangle(xproto.Rectangle)
	UpdateRoot()
	Wired() bool
	SetWired(bool)
	Focus()
	Primary() bool
	Focused() bool
	Last() bool
	Set(string, bool)
	Contains(xproto.Point) bool
	Desktops() *branch.Branch
	Merge(Monitor)
	Delete()
}

type monitor struct {
	settings.Settings
	c         *xgb.Conn
	id        uint32
	name      string
	rectangle xproto.Rectangle
	root      xproto.Window
	wired     bool
	pad       [4]int
	stickys   int
	primary   bool
	focused   bool
	last      bool
	desktops  *branch.Branch
}

func NewMonitor(id uint32, n string, c *xgb.Conn, root xproto.Window, r xproto.Rectangle, s settings.Settings) Monitor {
	win, err := xproto.NewWindowId(c)
	if err != nil {
		return nil
	}
	xproto.CreateWindow(
		c,
		xproto.WindowClassCopyFromParent,
		win,
		root,
		r.X,
		r.Y,
		r.Width,
		r.Height,
		0,
		xproto.WindowClassInputOnly,
		xproto.WindowClassCopyFromParent,
		xproto.CwEventMask,
		[]uint32{xproto.EventMaskEnterWindow},
	)

	xproto.ConfigureWindowChecked(c, win, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeBelow})

	if s.Bool("FocusFollowPointer") {
		windows.SetVisible(true, c, win, root)
	}

	var pad [4]int
	if p, err := s.Query("DefaultPad"); err != nil {
		pad = [4]int{0, 0, 0, 0}
	} else {
		pad = p.Pad()
	}

	m := &monitor{
		Settings:  s,
		c:         c,
		id:        id,
		name:      n,
		root:      root,
		pad:       pad,
		rectangle: r,
		wired:     true,
		focused:   true,
	}

	m.desktops = desktops.New(m.id, m.Settings)

	return m
}

func (m *monitor) Id() uint32 {
	return m.id
}

func (m *monitor) Name() string {
	return m.name
}

func (m *monitor) Rectangle() xproto.Rectangle {
	return m.rectangle
}

func (m *monitor) SetRectangle(r xproto.Rectangle) {
	m.rectangle = r
}
func (m *monitor) UpdateRoot() {
	r := m.rectangle
	xproto.ConfigureWindowChecked(m.c, m.root, xproto.ConfigWindowX, []uint32{uint32(r.X)})
	xproto.ConfigureWindowChecked(m.c, m.root, xproto.ConfigWindowY, []uint32{uint32(r.Y)})
	xproto.ConfigureWindowChecked(m.c, m.root, xproto.ConfigWindowHeight, []uint32{uint32(r.Height)})
	xproto.ConfigureWindowChecked(m.c, m.root, xproto.ConfigWindowWidth, []uint32{uint32(r.Width)})
}

func (m *monitor) Wired() bool {
	return m.wired
}

func (m *monitor) SetWired(v bool) {
	m.wired = v
}

func (m *monitor) Focus() {
	m.focused = true
	if m.Settings.Bool("PointerFollowsMonitor") {
		//center_pointer(m->rectangle)
	}
	//ewmh_update_current_desktop()
}

func (m *monitor) Primary() bool {
	return m.primary
}

func (m *monitor) Focused() bool {
	return m.focused
}

func (m *monitor) Last() bool {
	return m.last
}

func (m *monitor) Set(k string, v bool) {
	switch k {
	case "primary":
		m.primary = v
	case "focused":
		m.focused = v
	case "last":
		m.last = v
	}
}

func (m *monitor) Contains(p xproto.Point) bool {
	r := m.rectangle
	return (r.X <= p.X && p.X < (r.X+int16(r.Width)) && r.Y <= p.Y && p.Y < (r.Y+int16(r.Height)))
}

func (m *monitor) Desktops() *branch.Branch {
	return m.desktops
}

func (m *monitor) Merge(other Monitor) {
	m.desktops.PushBackBranch(other.Desktops())
	desktops.UpdateDesktopsMonitor(m.desktops, m.id)
	other.Delete()
}

func (m *monitor) Delete() {
	m.desktops = nil
	xproto.DestroyWindow(m.c, m.root)
	m = nil
}
