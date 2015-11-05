package monitors

import (
	"strings"

	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/selector"
)

func Select(monitors *branch.Branch, sel ...selector.Selector) Monitor {
	var probable []Monitor
	for _, s := range sel {
		if s.Node() == selector.NMonitor {
			switch s.References() {
			case selector.Cycled:
				probable = append(probable, Cycled(monitors, s))
			case selector.Focused:
				probable = append(probable, Focused(monitors))
			case selector.Tagged:
				probable = append(probable, Tagged(monitors, s))
			}
		}
	}
	return nil
}

func anyOf(s string, anyof ...string) bool {
	for _, a := range anyof {
		if s == a || strings.Contains(a, s) {
			return true
		}
	}
	return false
}

func Cycled(monitors *branch.Branch, sel selector.Selector) Monitor {
	cy := sel.Raw()
	for _, cycle := range cy {
		switch {
		case anyOf(cycle, "next", "forward"):
			return Next(monitors)
		case anyOf(cycle, "prev", "previous", "backward"):
			return Prev(monitors)
		}
	}
	return nil
}

func Tagged(monitors *branch.Branch, sel selector.Selector) Monitor {
	t := sel.Raw()
	var fn MatchMonitor
	for _, tagged := range t {
		switch tagged {
		case "name":
			fn = func(Monitor) bool { return false }
		case "id":
			fn = func(Monitor) bool { return false }
		case "index":
			fn = func(Monitor) bool { return false }
		}
	}
	return seek(monitors, fn)
}
