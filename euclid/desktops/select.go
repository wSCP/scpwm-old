package desktops

import (
	"strconv"

	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/selector"
	"github.com/thrisp/scpwm/utils"
)

func countHigh(d map[Desktop]int) Desktop {
	var ret Desktop
	var count int
	for k, v := range d {
		if v > count {
			count = v
			ret = k
		}
	}
	return ret
}

func addTo(d map[Desktop]int, dsk Desktop) {
	if dsk != nil {
		if c, ok := d[dsk]; ok {
			c++
			d[dsk] = c
		} else {
			d[dsk] = 1
		}
	}
}

func Select(selectors []selector.Selector, ds ...*branch.Branch) Desktop {
	var d map[Desktop]int
	for _, desktops := range ds {
		for _, sel := range selectors {
			if sel.Node() == selector.NDesktop {
				switch sel.Category() {
				case selector.Cycled:
					desk := Cycled(desktops, sel)
					addTo(d, desk)
				case selector.Focused:
					desk := Focused(desktops)
					addTo(d, desk)
				case selector.Tagged:
					desk := Tagged(desktops, sel)
					addTo(d, desk)
				}
			}
		}
	}
	return countHigh(d)
}

func Cycled(desktops *branch.Branch, sel selector.Selector) Desktop {
	cy := sel.Modifiers()
	for _, cycle := range cy {
		switch {
		case utils.MatchesAny(cycle, "next", "forward"):
			return Next(desktops)
		case utils.MatchesAny(cycle, "prev", "previous", "backward"):
			return Prev(desktops)
		}
	}
	return nil
}

func Tagged(desktops *branch.Branch, sel selector.Selector) Desktop {
	t := sel.Modifiers()
	var fn MatchDesktop
	for _, tagged := range t {
		switch {
		case utils.MatchesAny(tagged, "name", "named"):
			fn = func(d Desktop) bool {
				for _, name := range t {
					if d.Name() == name {
						return true
					}
				}
				return false
			}
		case utils.MatchesAny(tagged, "index", "indexed"):
			fn = func(d Desktop) bool {
				for _, idx := range t {
					if n, err := strconv.Atoi(idx); err == nil {
						if d.Index() == n {
							return true
						}
					}
				}
				return false
			}
		case utils.MatchesAny(tagged, "id"):
			fn = func(d Desktop) bool {
				for _, id := range t {
					if d.Id() == id {
						return true
					}
				}
				return false
			}
		}
	}
	return seek(desktops, fn)
}
