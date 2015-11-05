package monitors

import (
	"fmt"

	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/handler"
	"github.com/thrisp/scpwm/euclid/selector"
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
	if err == nil {
		if err := Update(monitors, x, s); err == nil {
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
}

func Focus(monitors *branch.Branch, mon Monitor) {
	focused := Focused(monitors)
	focused.Set("focused", false)
	mon.Focus()
}

func FocusSelect(monitors *branch.Branch, sel selector.Selector) {
	//current := Focused(monitors)
	//var toFocus Monitor
	//switch direction {
	//case "next":
	//	toFocus = Next(monitors)
	//case "previous":
	//	toFocus = Prev(monitors)
	//}
	//toFocus.Focus()
	//current.Set("focus", false)
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
