package monitors

import (
	"fmt"

	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/handler"
	"github.com/thrisp/scpwm/euclid/settings"
)

func New(x handler.Handler, s settings.Settings) *branch.Branch {
	monitors := branch.New("monitors")
	Initialize(monitors, x, s)
	return monitors
}

func Initialize(monitors *branch.Branch, x handler.Handler, s settings.Settings) {
	c := x.Conn()
	root := x.Root()

	err := randr.Init(c)
	if err == nil && Update(monitors, x, s) {
		randr.SelectInputChecked(c, root, randr.NotifyMaskScreenChange)
	} else {
		err = xinerama.Init(c)
		if err == nil {
			xia, err := xinerama.IsActive(c).Reply()
			if xia != nil && err == nil {
				xsq, _ := xinerama.QueryScreens(c).Reply()
				xsi := xsq.ScreenInfo
				for i := 0; i < len(xsi); i++ {
					info := xsi[i]
					rect := xproto.Rectangle{info.XOrg, info.YOrg, info.Width, info.Height}
					nm := NewMonitor(uint32(i), fmt.Sprintf("XMonitor%d", i), c, root, rect, s)
					monitors.PushBack(nm)
				}
			} else {
				scr := x.Screen()
				rect := xproto.Rectangle{0, 0, scr.WidthInPixels, scr.HeightInPixels}
				nm := NewMonitor(1, "SCREEN", c, root, rect, s)
				monitors.PushBack(nm)
			}
		}
	}
}

func Update(monitors *branch.Branch, x handler.Handler, s settings.Settings) bool {
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

	curr := monitors.Front()
	for curr != nil {
		m := curr.Value.(Monitor)
		m.SetWired(false)
		curr = curr.Next()
	}

	for i, info := range r {
		if info != nil {
			if info.Crtc != 0 {
				ir, _ := randr.GetCrtcInfo(c, info.Crtc, xproto.TimeCurrentTime).Reply()
				if ir != nil {
					rect := xproto.Rectangle{ir.X, ir.Y, ir.Width, ir.Height}
					m := fromId(monitors, uint32(sres.Outputs[i]))
					if m != nil {
						m.SetRectangle(rect)
						m.UpdateRoot()
						m.SetWired(true)
					} else {
						nm := NewMonitor(uint32(sres.Outputs[i]), string(info.Name), c, root, rect, s)
						monitors.PushBack(nm)
					}
				}
			} else if !s.Bool("RemoveDisabled") && info.Connection != randr.ConnectionDisconnected {
				m := fromId(monitors, uint32(sres.Outputs[i]))
				if m != nil {
					m.SetWired(true)
				}
			}
		}
	}

	gpo, _ := randr.GetOutputPrimary(c, x.Root()).Reply()
	if gpo != nil {
		pm := fromId(monitors, uint32(gpo.Output))
		if pm != nil {
			pm.Set("primary", true)
			pm.Set("focused", true)
			//ewmh_update_current_desktop();
		}
	}

	if s.Bool("MergeOverlapping") {
		mergeOverlapping(monitors)
	}

	if s.Bool("RemoveUnplugged") {
		removeUnplugged(monitors)
	}

	return monitors.Len() > 0
}

func contains(a, b xproto.Rectangle) bool {
	return (a.X <= b.X &&
		(a.X+int16(a.Width)) >= (b.X+int16(b.Width)) &&
		a.Y <= b.Y && (a.Y+int16(a.Height)) >= (b.Y+int16(b.Height)))
}

func mergeOverlapping(monitors *branch.Branch) {
	var mon1, mon2 Monitor
	m := monitors.Front()
	for m != nil {
		mon1 = m.Value.(Monitor)
		if mon1.Wired() {
			mm := monitors.Front()
			for mm != nil {
				mon2 = mm.Value.(Monitor)
				if mon1.Id() != mon2.Id() && mon2.Wired() && contains(mon1.Rectangle(), mon2.Rectangle()) {
					mon1.Merge(mon2)
				}
				mm = mm.Next()
			}
		}
		m = m.Next()
	}
}

func removeUnplugged(monitors *branch.Branch) {
	fn := func(mon Monitor) error {
		if !mon.Wired() {
			focused := Focused(monitors)
			focused.Merge(mon)
		}
		return nil
	}
	update(monitors, fn)
}

func All(monitors *branch.Branch) []Monitor {
	var ret []Monitor
	curr := monitors.Front()
	for curr != nil {
		mon := curr.Value.(Monitor)
		ret = append(ret, mon)
		curr = curr.Next()
	}
	return ret
}

type MatchMonitor func(Monitor) bool

func seek(monitors *branch.Branch, fn MatchMonitor) Monitor {
	curr := monitors.Front()
	for curr != nil {
		mon := curr.Value.(Monitor)
		if match := fn(mon); match {
			return mon
		}
		curr = curr.Next()
	}
	return nil
}

func seekoffset(monitors *branch.Branch, fn MatchMonitor, offset int) Monitor {
	curr := monitors.Front()
	for curr != nil {
		mon := curr.Value.(Monitor)
		if match := fn(mon); match {
			switch offset {
			case -1:
				mon = curr.Prev().Value.(Monitor)
			case 1:
				mon = curr.Next().Value.(Monitor)
			}
			return mon
		}
		curr = curr.Next()
	}
	return nil

}

func fromPoint(monitors *branch.Branch, p xproto.Point) Monitor {
	fn := func(m Monitor) bool {
		return m.Contains(p)
	}
	return seek(monitors, fn)
}

func fromId(monitors *branch.Branch, id uint32) Monitor {
	fn := func(m Monitor) bool {
		return m.Id() == id
	}
	return seek(monitors, fn)
}

func Primary(monitors *branch.Branch) Monitor {
	fn := func(m Monitor) bool {
		return m.Primary()
	}
	return seek(monitors, fn)
}

func isFocused(m Monitor) bool {
	return m.Focused()
}

func Focused(monitors *branch.Branch) Monitor {
	return seek(monitors, isFocused)
}

func Last(monitors *branch.Branch) Monitor {
	fn := func(m Monitor) bool {
		return m.Last()
	}
	return seek(monitors, fn)
}

func Next(monitors *branch.Branch) Monitor {
	return seekoffset(monitors, isFocused, 1)
}

func Prev(monitors *branch.Branch) Monitor {
	return seekoffset(monitors, isFocused, -1)
}

type UpdateMonitor func(Monitor) error

func update(monitors *branch.Branch, fn UpdateMonitor) error {
	curr := monitors.Front()
	for curr != nil {
		mon := curr.Value.(Monitor)
		if err := fn(mon); err != nil {
			return err
		}
		curr = curr.Next()
	}
	return nil
}

/*
const MAXINT = int(^uint(0) >> 1)

func abs(num int16) int {
	if num < 0 {
		num = -num
	}
	return int(num)
}

func (ms *monitors) fromClient(c *Client) *Monitor {
	r := c.floater.rectangle
	p := xproto.Point{r.X, r.Y}
	nearest := ms.fromPoint(p)
	if nearest == nil {
		cr := c.floater.rectangle
		x := ((cr.X + int16(cr.Width)) / 2)
		y := ((cr.Y + int16(cr.Height)) / 2)
		dmin := MAXINT
		for _, m := range ms.m {
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

func (ms *monitors) Swap(i, j int) {
	ms.m[i], ms.m[j] = ms.m[j], ms.m[i]
	ewmh_update_wm_desktops();
	ewmh_update_desktop_names();
	ewmh_update_current_desktop();
}

func (m *Monitor) Closest(cy Cycle, sel DesktopSelect) *Monitor {
	closest := m.Pop(cy)
	for closest != m {
		loc := m.loc //coordinates(m, m.Desktops.Desktop, nil)
		if MatchDesktop(&loc, &loc, sel) {
			return closest
		}
		closest = closest.Pop(cy)
	}
	return nil
}

func (m *Monitor) Nearest(dir Direction, sel DesktopSelect) *Monitor {
	var dmin int = MAXINT
	var ret *Monitor
	r1 := m.rectangle
	curr := m.Pop(Head)
	for curr != nil {
		if curr != m {
			loc := curr.loc //coordinates(curr, curr.Desktops.Desktop, nil)
			if MatchDesktop(&loc, &loc, sel) {
				r2 := curr.rectangle
				if (dir == Left && r2.X < r1.X) ||
					(dir == Right && r2.X >= (r1.X+int16(r1.Width))) ||
					(dir == Up && r2.Y < r1.Y) ||
					(dir == Down && r2.Y >= (r1.Y+int16(r1.Height))) {
					d := abs((r2.X+int16(r2.Width/2))-(r1.X+int16(r1.Width/2))) + abs((r2.Y+int16(r2.Height/2))-(r1.Y+int16(r1.Height/2)))
					if d < dmin {
						dmin = d
						ret = curr
					}
				}
			}
		}
		curr = curr.Pop(Next)
	}
	return ret
}

func selectMonitor(e *Euclid, loc coordinate, sel ...Selector) (coordinate, bool) {
	return loc, false
}

func MonitorFromIndex(e *Euclid, idx int, locate coordinate) bool {
	if m, ok := e.Monitors[idx]; ok {
		locate.m = m
		locate.d, locate.c = nil, nil
		return true
	}
	return false
}

func MonitorFromDescription(description []string, ref, dst *coordinate) bool {
	sel := desktopselectFromString(description)

	dst.m = nil
	var dir Direction
	var cy Cycle
	var hdir Age
	var idx int

	if directionFromString(description, dir) {
		dst->Monitor = nearest_Monitor(ref->Monitor, dir, sel);
	} else if cycleFromString(description, cy) {
		dst->Monitor = closest_Monitor(ref->Monitor, cyc, sel);
	} else if ageFromString(description, hdir) {
		history_find_Monitor(hdi, ref, dst, sel);
	} else if strings.Contains(description, "last") {
		history_find_Monitor(HISTORY_OLDER, ref, dst, sel);
	} else if strings.Contains(description, "primary") {
		if (pri_mon != NULL) {
			coordinates_t loc = {pri_mon, pri_mon->desk, NULL};
			if (desktop_matches(&loc, ref, sel))
				dst->Monitor = pri_mon;
		}
	} else if strings.Contains(description, "focused") {
		coordinates_t loc = {mon, mon->desk, NULL};
		if (desktop_matches(&loc, ref, sel))
			dst->Monitor = mon;
	} else if indexFromString(description, idx) {
		Monitor_from_index(idx, dst);
	} else {
		locate_Monitor(desc, dst);
	}

	return dst.m != nil

return false
}
*/
