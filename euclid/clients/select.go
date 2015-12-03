package clients

import (
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/selector"
	"github.com/thrisp/scpwm/utils"
)

func countHigh(d map[Client]int) Client {
	var ret Client
	var count int
	for k, v := range d {
		if v > count {
			count = v
			ret = k
		}
	}
	return ret
}

func addTo(c map[Client]int, cs Client) {
	if cs != nil {
		if cnt, ok := c[cs]; ok {
			cnt++
			c[cs] = cnt
		} else {
			c[cs] = 1
		}
	}
}

func Select(selectors []selector.Selector, cls ...*branch.Branch) Client {
	var c map[Client]int
	for _, clients := range cls {
		for _, sel := range selectors {
			if sel.Node() == selector.NClient {
				switch sel.Category() {
				case selector.Cycled:
					sc := Cycled(clients, sel)
					addTo(c, sc)
				case selector.Focused:
					sc := Focused(clients)
					addTo(c, sc)
				case selector.Tagged:
					sc := Tagged(clients, sel)
					addTo(c, sc)
				}
			}
		}
	}
	return countHigh(c)
}

func Cycled(clients *branch.Branch, sel selector.Selector) Client {
	cy := sel.Modifiers()
	for _, cycle := range cy {
		switch {
		case utils.MatchesAny(cycle, "next", "forward"):
			return Next(clients)
		case utils.MatchesAny(cycle, "prev", "previous", "backward"):
			return Prev(clients)
		}
	}
	return nil
}

func Tagged(clients *branch.Branch, sel selector.Selector) Client {
	/*t := sel.Modifiers()
	var fn MatchClient
	for _, tagged := range t {
		switch {
		case utils.MatchesAny(tagged, "name", "named"):
			fn = func(d Client) bool {
				for _, name := range t {
					if d.Name() == name {
						return true
					}
				}
				return false
			}
		case utils.MatchesAny(tagged, "index", "indexed"):
			fn = func(d Client) bool {
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
			fn = func(d Client) bool {
				for _, id := range t {
					if d.Id() == id {
						return true
					}
				}
				return false
			}
		}
	}
	return seek(clients, fn)*/
	return nil
}
