package rules

import (
	"fmt"
	"strconv"
)

type Ruler interface {
	Rule(...string) bool
	Unrule(...string) bool
	Applicable(...string) []Rule
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

func (r *ruler) Rule(d ...string) bool {
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

func (r *ruler) Unrule(d ...string) bool {
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

func contains(tags []string, having ...string) bool {
	for _, t := range tags {
		for _, h := range having {
			if t == h {
				return true
			}
		}
	}
	return false
}

func (r *ruler) Applicable(tags ...string) []Rule {
	var ret []Rule
	for _, rule := range r.pending {
		cause := rule.Cause().String()
		if contains(tags, cause) {
			ret = append(ret, rule)
		}
	}
	return nil
}

func (r *ruler) Pending() []Rule {
	return r.pending
}
