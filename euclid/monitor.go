package main

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
)

type Monitors []*Monitor

func NewMonitors(e *Euclid) Monitors {
	ms := make(Monitors, 0)
	e.randr, e.xinerama = ms.Init(e, e.XHandle)
	return ms
}

func (ms Monitors) Init(e *Euclid, x XHandle) (bool, bool) {
	var rr, xn bool
	c := x.Conn()

	err := randr.Init(c)
	if err == nil && ms.Update(e, x, false, false, false) {
		rr = true
		randr.SelectInputChecked(c, x.Root(), randr.NotifyMaskScreenChange)
	} else {
		rr = false
		err = xinerama.Init(c)
		if err == nil {
			xia, err := xinerama.IsActive(c).Reply()
			if xia != nil && err == nil {
				xn = true
				xsq, _ := xinerama.QueryScreens(c).Reply()
				xsi := xsq.ScreenInfo
				for i := 0; i < len(xsi); i++ {
					info := xsi[i]
					rect := xproto.Rectangle{info.XOrg, info.YOrg, info.Width, info.Height}
					ms.Add(e, x, rect)
				}
			} else {
				s := x.Screen()
				rect := xproto.Rectangle{0, 0, s.WidthInPixels, s.HeightInPixels}
				ms.Add(e, x, rect)
			}
		}
	}
	return rr, xn
}

func (ms Monitors) Update(e *Euclid, x XHandle, mergeOverlapping, removeUnplugged, removeDisabled bool) bool {
	c := x.Conn()

	sres, err := randr.GetScreenResources(c, x.Root()).Reply()
	if err != nil {
		return false
	}

	var r []*randr.GetOutputInfoReply
	for _, o := range sres.Outputs {
		rp, _ := randr.GetOutputInfo(c, o, xproto.TimeCurrentTime).Reply()
		r = append(r, rp)
	}

	for _, m := range ms {
		m.wired = false
	}

	for i, info := range r {
		if info != nil {
			if info.Crtc != 0 {
				ir, _ := randr.GetCrtcInfo(c, info.Crtc, xproto.TimeCurrentTime).Reply()
				if ir != nil {
					rect := xproto.Rectangle{ir.X, ir.Y, ir.Width, ir.Height}
					m := ms.fromId(sres.Outputs[i])
					if m != nil {
						m.rectangle = rect
						m.UpdateRoot()
						for _, d := range m.desktops {
							for _, _ = range d.clients {
								//mm.Translate(mm, n.Client)
							}
						}
						//mm.Arrange()
						m.wired = true
					} else {
						m := ms.Add(e, x, rect)
						m.name = string(info.Name)
						m.id = sres.Outputs[i]
					}
				}
			} else if !removeDisabled && info.Connection != randr.ConnectionDisconnected {
				m := ms.fromId(sres.Outputs[i])
				if m != nil {
					m.wired = true
				}
			}
		}
	}

	gpo, _ := randr.GetOutputPrimary(c, x.Root()).Reply()
	if gpo != nil {
		primary := ms.fromId(gpo.Output)
		if primary != nil {
			primary.primary = true
			if ms.Focused() != primary {
				primary.focused = true
			}
			//ewmh_update_current_desktop();
		}
	}

	if mergeOverlapping {
		ms.mergeOverlapping()
	}

	if removeUnplugged {
		ms.removeUnplugged()
	}

	//update_motion_recorder();

	return ms.Number() > 0
}

func (ms Monitors) mergeOverlapping() {
	for _, m := range ms {
		if m.wired {
			for _, mm := range ms {
				if m != mm && mm.wired && contains(m.rectangle, mm.rectangle) {
					m.Merge(mm)
					ms.Remove(mm)
				}
			}
		}
	}
}

func (ms Monitors) removeUnplugged() {
	for _, m := range ms {
		if !m.wired {
			focused := ms.Focused()
			focused.Merge(m)
			ms.Remove(m)
		}
	}
}

func (ms Monitors) Add(e *Euclid, x XHandle, r xproto.Rectangle) *Monitor {
	c := x.Conn()
	w := x.NewWindow()
	n := fmt.Sprintf("%s%d", e.String("DefaultMonitorName"), ms.Number()+1)
	m := newMonitor(e, n, w, r)

	xproto.CreateWindow(
		c,
		xproto.WindowClassCopyFromParent,
		m.root.Window,
		x.Root(),
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

	m.root.Lower()

	if e.Bool("FocusFollowsPointer") {
		m.root.Show()
	}

	ms = append(ms, m)

	m.focused = true

	return m
}

func (ms *Monitors) Select(sel ...Selector) bool {
	return false
}

func (ms Monitors) fromPoint(p xproto.Point) *Monitor {
	for _, m := range ms {
		if m.Contains(p) {
			return m
		}
	}
	return nil
}

func (ms Monitors) fromClient(c *Client) *Monitor {
	p := xproto.Point{c.fRectangle.X, c.fRectangle.Y}
	nearest := ms.fromPoint(p)
	if nearest == nil {
		cr := c.floater.rectangle
		x := ((cr.X + int16(cr.Width)) / 2)
		y := ((cr.Y + int16(cr.Height)) / 2)
		dmin := MAXINT
		for _, m := range ms {
			r := m.rectangle
			d := abs((r.X+int16(r.Width)/2)-x) + abs((r.Y+int16(r.Height)/2)-y)
			if d < dmin {
				dmin = d
				nearest = m
			}
		}
	}
	return nearest
}

func (ms Monitors) fromId(id randr.Output) *Monitor {
	for _, m := range ms {
		if m.id == id {
			return m
		}
	}
	return nil
}

func (ms Monitors) Primary() *Monitor {
	for _, p := range ms {
		if p.primary {
			return p
		}

	}
	return nil
}

func (ms Monitors) Focused() *Monitor {
	for _, m := range ms {
		if m.Focused() {
			return m
		}

	}
	return nil
}

func (ms Monitors) Last() *Monitor {
	for _, m := range ms {
		if m.last {
			return m
		}

	}
	return nil
}

func (ms Monitors) Swap(first, second int) {
	//ewmh_update_wm_desktops();
	//ewmh_update_desktop_names();
	//ewmh_update_current_desktop();
}

func (ms Monitors) Remove(m *Monitor) {
	//prev := m.Pop(Prev)
	//next := m.Pop(Next)
	//if prev == nil {
	//	ms.Monitor = next
	//} else {
	//	prev.next = next
	//	next.prev = prev
	//}
	//ms.Remove(m)
}

func (ms Monitors) Number() int {
	return len(ms)
}

type Pad map[Direction]int

func DefaultPad() Pad {
	ret := make(Pad)
	ret[Up], ret[Down], ret[Left], ret[Right] = 0, 0, 0, 0
	return ret
}

func (p Pad) Get(d Direction) int {
	if v, ok := p[d]; ok {
		return v
	}
	return 0
}

type Monitor struct {
	loc       coordinate
	id        randr.Output
	name      string
	rectangle xproto.Rectangle
	root      *window
	wired     bool
	pad       Pad
	stickys   int
	primary   bool
	focused   bool
	last      bool
	desktops  Desktops
}

func newMonitor(e *Euclid, n string, w *window, r xproto.Rectangle) *Monitor {
	m := &Monitor{
		name:      n,
		root:      w,
		pad:       DefaultPad(),
		rectangle: r,
		wired:     true,
	}
	m.loc = Coordinate(e, m, nil, nil)
	m.desktops = NewDesktops()
	m.desktops.Add(e, m, "")
	return m
}

func (m *Monitor) UpdateRoot() {
	r := m.rectangle
	m.root.MoveResize(r.X, r.Y, r.Width, r.Height)
}

func (m *Monitor) Focus() {
	m.focused = true
	//if m.e.Bool("PointerFollowsMonitor") {
	//center_pointer(m->rectangle)
	//}
	//ewmh_update_current_desktop()
}

func (m *Monitor) Primary() bool {
	return m.primary
}

func (m *Monitor) Focused() bool {
	return m.focused
}

func (m *Monitor) Last() bool {
	return m.last
}

func (m *Monitor) Contains(p xproto.Point) bool {
	r := m.rectangle
	return (r.X <= p.X && p.X < (r.X+int16(r.Width)) && r.Y <= p.Y && p.Y < (r.Y+int16(r.Height)))
}

func (m *Monitor) Merge(dst *Monitor) {
	for _, d := range dst.desktops {
		d.loc.m = m
		for _, c := range d.clients {
			c.loc.m = m
		}
	}
	m.desktops = append(m.desktops, dst.desktops...)
}

func (m *Monitor) Embrace(c *Client) {
	cf := c.floater.rectangle
	if (cf.X + int16(cf.Width)) <= m.rectangle.X {
		//	c.fRectangle.X = m.rectangle.X
	} else if cf.X >= (m.rectangle.X + int16(m.rectangle.Width)) {
		//	c.fRectangle.X = (m.rectangle.X + int16(m.rectangle.Width)) - int16(c.fRectangle.Width)
	}

	if (cf.Y + int16(cf.Height)) <= m.rectangle.Y {
		//	c.fRectangle.Y = m.rectangle.Y
	} else if cf.Y >= (m.rectangle.Y + int16(m.rectangle.Height)) {
		//	c.fRectangle.Y = (m.rectangle.Y + int16(m.rectangle.Height)) - int16(c.fRectangle.Height)
	}
}

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

func (m *Monitor) Delete(c *xgb.Conn) {
	m.desktops = nil
	xproto.DestroyWindow(c, m.root.Window)
	m = nil
}

func selectMonitor(e *Euclid, loc coordinate, sel ...Selector) (coordinate, bool) {
	return loc, false
}

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
