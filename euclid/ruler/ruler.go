package ruler

import (
	"fmt"
	"strconv"

	"github.com/BurntSushi/xgb/xproto"
)

type Ruler interface {
	Add(...string) bool
	Remove(...string) bool
	Applicable(xproto.Window) ([]Rule, bool)
	Pending() []Rule
}

type ruler struct {
	pending []Rule
}

func New() Ruler {
	return &ruler{
		pending: make([]Rule, 0),
	}
}

func (r *ruler) Add(d ...string) bool {
	var once bool
	switch len(d) {
	case 4:
		once = false
	case 5:
		b, err := strconv.ParseBool(d[4])
		if err == nil {
			once = b
		}
	}
	return r.add(once, d[0], d[1], d[2], d[3])
}

func (r *ruler) add(once bool, d ...string) bool {
	nr := newRule(d[0], d[1], d[2], d[3], once)
	if nr != nil {
		r.pending = append(r.pending, nr)
		return true
	}
	return false
}

func (r *ruler) Remove(d ...string) bool {
	var once string
	switch len(d) {
	case 4:
		once = "false"
	case 5:
		once = d[4]
	}
	return r.remove(fmt.Sprintf("%s %s %s %s %s", d[0], d[1], d[2], d[3], once))
}

func (r *ruler) remove(d string) bool {
	for i, p := range r.pending {
		if p.String() == d {
			r.pending, r.pending[len(r.pending)-1] = append(r.pending[:i], r.pending[i+1:]...), nil
			return true
		}
	}
	return false
}

func (r *ruler) Applicable(w xproto.Window) ([]Rule, bool) {
	return nil, false
}

func (r *ruler) Pending() []Rule {
	return r.pending
}

//var defaultWindowRuleset = []Rule{
//	newRule("window", "manage", "true", true),
//	newRule("window", "focus", "true", true),
//	newRule("window", "bordered", "true", true),
//}

//func (e *Euclid) applyRules(win *Window, csq *Consequence) {}

//func (e *Euclid) scheduleRules(win *Window, csq *Consequence) bool {
//	return false
//}
