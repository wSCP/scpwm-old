package monitors

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/clients"
)

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

const MAXINT = int(^uint(0) >> 1)

func abs(num int16) int {
	if num < 0 {
		num = -num
	}
	return int(num)
}

func FromClient(monitors *branch.Branch, c clients.Client) Monitor {
	cr := c.FRectangle()
	pt := xproto.Point{cr.X, cr.Y}
	var nearest Monitor
	nearest = fromPoint(monitors, pt)
	if nearest == nil {
		x := ((cr.X + int16(cr.Width)) / 2)
		y := ((cr.Y + int16(cr.Height)) / 2)
		var dmin = MAXINT
		fn := func(m Monitor) bool {
			r := m.Rectangle()
			d := abs((r.X+int16(r.Width)/2)-x) + abs((r.Y+int16(r.Height)/2)-y)
			if d < dmin {
				dmin = d
				nearest = m
			}
			return false
		}
		seek(monitors, fn)
	}
	return nearest
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

func Next(monitors *branch.Branch) Monitor {
	return seekoffset(monitors, isFocused, 1)
}

func Prev(monitors *branch.Branch) Monitor {
	return seekoffset(monitors, isFocused, -1)
}

func seekAny(monitors *branch.Branch, fn MatchMonitor) []Monitor {
	var ret []Monitor
	curr := monitors.Front()
	for curr != nil {
		mon := curr.Value.(Monitor)
		if match := fn(mon); match {
			ret = append(ret, mon)
		}
		curr = curr.Next()
	}
	return ret
}

func All(monitors *branch.Branch) []Monitor {
	fn := func(m Monitor) bool { return true }
	return seekAny(monitors, fn)
}
