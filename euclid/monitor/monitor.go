package monitor

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
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
	Contains(xproto.Point) bool
	Delete()
}

type monitor struct {
	c         *xgb.Conn
	id        uint32 //randr.Output
	name      string
	rectangle xproto.Rectangle
	root      xproto.Window
	wired     bool
	pad       [4]int
	stickys   int
	primary   bool
	focused   bool
	last      bool
}

func NewMonitor(n string, id uint32, c *xgb.Conn, root xproto.Window, r xproto.Rectangle) Monitor {
	return &monitor{
		c:         c,
		id:        id,
		name:      n,
		root:      root,
		pad:       [4]int{},
		rectangle: r,
		wired:     true,
	}
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
	//if m.e.Bool("PointerFollowsMonitor") {
	//center_pointer(m->rectangle)
	//}
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

func (m *monitor) Contains(p xproto.Point) bool {
	r := m.rectangle
	return (r.X <= p.X && p.X < (r.X+int16(r.Width)) && r.Y <= p.Y && p.Y < (r.Y+int16(r.Height)))
}

//func (m *monitor) Merge(dst Monitor) {
//for _, d := range dst.desktops {
//d.loc.m = m
//for _, c := range d.clients {
//c.loc.m = m
//}
//}
//m.desktops = append(m.desktops, dst.desktops...)
//}

func (m *monitor) Delete() {
	//m.desktops = nil
	//xproto.DestroyWindow(c, m.root)
	//m = nil
}

//func (m *Monitor) Embrace(c *Client) {
//cf := c.floater.rectangle
//if (cf.X + int16(cf.Width)) <= m.rectangle.X {
//	c.fRectangle.X = m.rectangle.X
//} else if cf.X >= (m.rectangle.X + int16(m.rectangle.Width)) {
//	c.fRectangle.X = (m.rectangle.X + int16(m.rectangle.Width)) - int16(c.fRectangle.Width)
//}

//if (cf.Y + int16(cf.Height)) <= m.rectangle.Y {
//	c.fRectangle.Y = m.rectangle.Y
//} else if cf.Y >= (m.rectangle.Y + int16(m.rectangle.Height)) {
//	c.fRectangle.Y = (m.rectangle.Y + int16(m.rectangle.Height)) - int16(c.fRectangle.Height)
//}
//}

//func (m *Monitor) Translate(o *Monitor, c *Client) {
//if m.e.pointer.action == NoAction || m == o {
//	leftAdjust := max((m.rectangle.X - c.fRectangle.X), 0)
//	topAdjust := max((m.rectangle.Y - c.fRectangle.Y), 0)
//	rightAdjust := max((c.fRectangle.X+int16(c.fRectangle.Width))-(m.rectangle.X+int16(m.rectangle.Width)), 0)
//	bottomAdjust := max((c.fRectangle.Y+int16(c.fRectangle.Height))-(m.rectangle.Y+int16(m.rectangle.Height)), 0)
//	c.fRectangle.X += leftAdjust
//	c.fRectangle.Y += topAdjust
//	c.fRectangle.Width -= uint16(leftAdjust + rightAdjust)
//	c.fRectangle.Height -= uint16(topAdjust + bottomAdjust)
//
//		dx := c.fRectangle.X - m.rectangle.X
//		dy := c.fRectangle.Y - m.rectangle.Y
//
//		nx := dx * int16(o.rectangle.Width-c.fRectangle.Width)
//		ny := dy * int16(o.rectangle.Height-c.fRectangle.Height)
//
//		dnx := int16(m.rectangle.Width - c.fRectangle.Width)
//		dny := int16(m.rectangle.Height - c.fRectangle.Height)
//
//		var dxd, dyd int16
//		if dnx == 0 {
//			dxd = 0
//		} else {
//			dxd = nx / dnx
//		}
//
//		if dny == 0 {
//			dyd = 0
//		} else {
//			dyd = ny / dny
//		}

//		c.fRectangle.Width += uint16(leftAdjust + rightAdjust)
//		c.fRectangle.Height += uint16(topAdjust + bottomAdjust)
//		c.fRectangle.X = o.rectangle.X + dxd - leftAdjust
//		c.fRectangle.Y = o.rectangle.Y + dyd - topAdjust
//	}
//}

//func (m *Monitor) Closest(cy Cycle, sel DesktopSelect) *Monitor {
//	closest := m.Pop(cy)
//	for closest != m {
//		loc := m.loc //coordinates(m, m.Desktops.Desktop, nil)
//		if MatchDesktop(&loc, &loc, sel) {
//			return closest
//		}
//		closest = closest.Pop(cy)
//	}
//	return nil
//}

//func (m *Monitor) Nearest(dir Direction, sel DesktopSelect) *Monitor {
//	var dmin int = MAXINT
///	var ret *Monitor
//	r1 := m.rectangle
//	curr := m.Pop(Head)
//	for curr != nil {
//		if curr != m {
//			loc := curr.loc //coordinates(curr, curr.Desktops.Desktop, nil)
//			if MatchDesktop(&loc, &loc, sel) {
//				r2 := curr.rectangle
//				if (dir == Left && r2.X < r1.X) ||
//					(dir == Right && r2.X >= (r1.X+int16(r1.Width))) ||
//					(dir == Up && r2.Y < r1.Y) ||
//					(dir == Down && r2.Y >= (r1.Y+int16(r1.Height))) {
//					d := abs((r2.X+int16(r2.Width/2))-(r1.X+int16(r1.Width/2))) + abs((r2.Y+int16(r2.Height/2))-(r1.Y+int16(r1.Height/2)))
//					if d < dmin {
//						dmin = d
//						ret = curr
//					}
//				}
//			}
//		}
//		curr = curr.Pop(Next)
//	}
//	return ret
//}

//func selectMonitor(e *Euclid, loc coordinate, sel ...Selector) (coordinate, bool) {
//	return loc, false
//}

//func MonitorFromIndex(e *Euclid, idx int, locate coordinate) bool {
//	if m, ok := e.Monitors[idx]; ok {
//		locate.m = m
//		locate.d, locate.c = nil, nil
//		return true
//	}
//	return false
//}

//func MonitorFromDescription(description []string, ref, dst *coordinate) bool {
/*
	//sel := desktopselectFromString(description)

	dst.m = nil
	var dir Direction
	var cy Cycle
	var hdir Age
	var idx int

	if directionFromString(description, dir) {
		//dst->Monitor = nearest_Monitor(ref->Monitor, dir, sel);
	} else if cycleFromString(description, cy) {
		//dst->Monitor = closest_Monitor(ref->Monitor, cyc, sel);
	} else if ageFromString(description, hdir) {
		//history_find_Monitor(hdi, ref, dst, sel);
	} else if strings.Contains(description, "last") {
		//history_find_Monitor(HISTORY_OLDER, ref, dst, sel);
	} else if strings.Contains(description, "primary") {
		//if (pri_mon != NULL) {
		//	coordinates_t loc = {pri_mon, pri_mon->desk, NULL};
		//	if (desktop_matches(&loc, ref, sel))
		//		dst->Monitor = pri_mon;
		//}
	} else if strings.Contains(description, "focused") {
		//coordinates_t loc = {mon, mon->desk, NULL};
		//if (desktop_matches(&loc, ref, sel))
		//	dst->Monitor = mon;
	} else if indexFromString(description, idx) {
		//Monitor_from_index(idx, dst);
	} else {
		//locate_Monitor(desc, dst);
	}

	return dst.m != nil
*/
//return false
//}
