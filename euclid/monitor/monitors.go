package monitor

import (
	"fmt"

	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/handler"
	"github.com/thrisp/scpwm/euclid/settings"
)

type Monitors interface {
	Init(handler.Handler) (bool, bool)
	Settings(*settings.Settings)
	Update(handler.Handler) bool
	Add(Monitor)
	Pop(Monitor) Monitor
	Swap(int, int)
	Primary() Monitor
	Focused() Monitor
	Last() Monitor
	All() []Monitor
	Number() int
}

type monitors struct {
	randr               bool
	xinerama            bool
	mergeoverlapping    bool
	removeunplugged     bool
	removedisabled      bool
	focusfollowspointer bool
	nameDefault         string
	m                   []Monitor
}

func NewMonitors(s *settings.Settings, x handler.Handler) Monitors {
	ms := &monitors{
		m: make([]Monitor, 0),
	}
	ms.Settings(s)
	ms.randr, ms.xinerama = ms.Init(x)
	return ms
}

func (ms *monitors) Init(x handler.Handler) (bool, bool) {
	var rr, xn bool
	c := x.Conn()
	root := x.Root()

	err := randr.Init(c)
	if err == nil && ms.Update(x) {
		rr = true
		randr.SelectInputChecked(c, root, randr.NotifyMaskScreenChange)
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
					nm := NewMonitor(fmt.Sprintf("XMonitor%d", i), uint32(i), c, root, rect)
					ms.Add(nm)
				}
			} else {
				s := x.Screen()
				rect := xproto.Rectangle{0, 0, s.WidthInPixels, s.HeightInPixels}
				nm := NewMonitor("SCREEN", 1, c, root, rect)
				ms.Add(nm)
			}
		}
	}
	return rr, xn
}

func (ms *monitors) Settings(s *settings.Settings) {
	ms.nameDefault = s.String("DefaultMonitorName")
	ms.mergeoverlapping = s.Bool("MergeOverlapping")
	ms.removeunplugged = s.Bool("RemoveUnplugged")
	ms.removedisabled = s.Bool("RemoveDisabled")
	ms.focusfollowspointer = s.Bool("FocusFollowsPointer")
}

func (ms *monitors) Update(x handler.Handler) bool {
	c := x.Conn()
	root := x.Root()

	sres, err := randr.GetScreenResources(c, root).Reply()
	if err != nil {
		return false
	}

	var r []*randr.GetOutputInfoReply
	for _, o := range sres.Outputs {
		rp, _ := randr.GetOutputInfo(c, o, xproto.TimeCurrentTime).Reply()
		r = append(r, rp)
	}

	for _, m := range ms.m {
		m.SetWired(false)
	}

	for i, info := range r {
		if info != nil {
			if info.Crtc != 0 {
				ir, _ := randr.GetCrtcInfo(c, info.Crtc, xproto.TimeCurrentTime).Reply()
				if ir != nil {
					rect := xproto.Rectangle{ir.X, ir.Y, ir.Width, ir.Height}
					m := ms.fromId(sres.Outputs[i])
					if m != nil {
						m.SetRectangle(rect)
						m.UpdateRoot()
						m.SetWired(true)
					} else {
						nm := NewMonitor(string(info.Name), uint32(sres.Outputs[i]), c, root, rect)
						ms.Add(nm)
					}
				}
			} else if !ms.removedisabled && info.Connection != randr.ConnectionDisconnected {
				m := ms.fromId(sres.Outputs[i])
				if m != nil {
					m.SetWired(true)
				}
			}
		}
	}

	gpo, _ := randr.GetOutputPrimary(c, x.Root()).Reply()
	if gpo != nil {
		primary := ms.fromId(gpo.Output)
		if primary != nil {
			//primary.primary = true
			//if ms.Focused() != primary {
			//	primary.focused = true
			//}
			//ewmh_update_current_desktop();
		}
	}

	if ms.mergeoverlapping {
		ms.mergeOverlapping()
	}

	if ms.removeunplugged {
		ms.removeUnplugged()
	}

	return ms.Number() > 0
}

func contains(a, b xproto.Rectangle) bool {
	return (a.X <= b.X && (a.X+int16(a.Width)) >= (b.X+int16(b.Width)) &&
		a.Y <= b.Y && (a.Y+int16(a.Height)) >= (b.Y+int16(b.Height)))
}

func (ms *monitors) mergeOverlapping() {
	//for _, m := range ms.m {
	//if m.wired {
	//	for _, mm := range ms.m {
	//		if m != mm && mm.wired && contains(m.rectangle, mm.rectangle) {
	//			m.Merge(mm)
	//			ms.Remove(mm)
	//		}
	//	}
	//}
	//}
}

func (ms *monitors) removeUnplugged() {
	for _, m := range ms.m {
		if !m.Wired() {
			//focused := ms.Focused()
			//focused.Merge(m)
			//m := ms.Pop(m)
			//m.Delete()
		}
	}
}

func (ms *monitors) Add(m Monitor) {
	//func (ms *monitors) Add(x Handler, r xproto.Rectangle) *Monitor {
	/*
		c := x.Conn()
		w := x.New()
		n := fmt.Sprintf("%s%d", ms.nameDefault, ms.Number()+1)
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

		if ms.focusFollowsPointer {
			m.root.Show()
		}

		ms = append(ms, m)

		m.focused = true

		return m
	*/
}

func (ms *monitors) Pop(m Monitor) Monitor {
	//for _, m := range ms.m {
	//
	//}
	return nil
}

func (ms *monitors) fromPoint(p xproto.Point) Monitor {
	for _, m := range ms.m {
		if m.Contains(p) {
			return m
		}
	}
	return nil
}

const MAXINT = int(^uint(0) >> 1)

func abs(num int16) int {
	if num < 0 {
		num = -num
	}
	return int(num)
}

//func (ms *monitors) fromClient(c *Client) *Monitor {
//	r := c.floater.rectangle
//	p := xproto.Point{r.X, r.Y}
//	nearest := ms.fromPoint(p)
//	if nearest == nil {
//		cr := c.floater.rectangle
//		x := ((cr.X + int16(cr.Width)) / 2)
//		y := ((cr.Y + int16(cr.Height)) / 2)
//		dmin := MAXINT
//		for _, m := range ms.m {
//			r := m.rectangle
//			d := abs((r.X+int16(r.Width)/2)-x) + abs((r.Y+int16(r.Height)/2)-y)
//			if d < dmin {
//				dmin = d
//				nearest = m
//			}
//		}
//	}
//	return nearest
//}

func (ms *monitors) fromId(id randr.Output) Monitor {
	for _, m := range ms.m {
		if m.Id() == id {
			return m
		}
	}
	return nil
}

func (ms *monitors) Primary() Monitor {
	for _, m := range ms.m {
		if m.Primary() {
			return m
		}

	}
	return nil
}

func (ms *monitors) Focused() Monitor {
	for _, m := range ms.m {
		if m.Focused() {
			return m
		}

	}
	return nil
}

func (ms *monitors) Last() Monitor {
	for _, m := range ms.m {
		if m.Last() {
			return m
		}

	}
	return nil
}

func (ms *monitors) Swap(i, j int) {
	ms.m[i], ms.m[j] = ms.m[j], ms.m[i]
	//ewmh_update_wm_desktops();
	//ewmh_update_desktop_names();
	//ewmh_update_current_desktop();
}

func (ms *monitors) All() []Monitor {
	return ms.m
}

func (ms *monitors) Number() int {
	return len(ms.m)
}
