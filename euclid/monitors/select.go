package monitors

import (
	"strconv"

	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/selector"
	"github.com/thrisp/scpwm/utils"
)

func countHigh(m map[Monitor]int) Monitor {
	var ret Monitor
	var count int
	for k, v := range m {
		if v > count {
			count = v
			ret = k
		}
	}
	return ret
}

func addTo(m map[Monitor]int, mon Monitor) {
	if mon != nil {
		if c, ok := m[mon]; ok {
			c++
			m[mon] = c
		} else {
			m[mon] = 1
		}
	}
}

func Select(sel []selector.Selector, ms ...*branch.Branch) Monitor {
	var m map[Monitor]int
	for _, monitors := range ms {
		for _, s := range sel {
			if s.Node() == selector.NMonitor {
				switch s.Category() {
				case selector.Cycled:
					mon := Cycled(monitors, s)
					addTo(m, mon)
				case selector.Focused:
					mon := Focused(monitors)
					addTo(m, mon)
				case selector.Tagged:
					mon := Tagged(monitors, s)
					addTo(m, mon)
				}
			}
		}
	}
	return countHigh(m)
}

func Cycled(monitors *branch.Branch, sel selector.Selector) Monitor {
	cy := sel.Modifiers()
	for _, cycle := range cy {
		switch {
		case utils.MatchesAny(cycle, "next", "forward"):
			return Next(monitors)
		case utils.MatchesAny(cycle, "prev", "previous", "backward"):
			return Prev(monitors)
		}
	}
	return nil
}

func Tagged(monitors *branch.Branch, sel selector.Selector) Monitor {
	t := sel.Modifiers()
	var fn MatchMonitor
	for _, tagged := range t {
		switch tagged {
		case "name":
			fn = func(m Monitor) bool {
				for _, name := range t {
					if m.Name() == name {
						return true
					}
				}
				return false
			}
		case "id":
			fn = func(m Monitor) bool {
				for _, id := range t {
					if nid, err := strconv.Atoi(id); err == nil {
						if m.Id() == uint32(nid) {
							return true
						}
					}
				}
				return false
			}
		}
	}
	return seek(monitors, fn)
}
