package query

import (
	"fmt"
	"strings"
)

type Adjective int

const (
	cycled   Adjective = iota // Cycle
	directed                  // Direction
	oriented                  // Orientation
	sized                     // Size
	aged                      // Age
	focused                   // by focus, urgency
	tagged                    // by name, tag or id
	special                   // relatively unique property of Node, e.g. any client state
)

type Selector interface {
	Raw() string
	Adjective() Adjective
	Object() Node
}

type selector struct {
	raw       string
	adjective Adjective
	object    Node
}

func (s selector) Adjective() Adjective {
	return s.adjective
}

func (s selector) Raw() string {
	return s.raw
}

func (s selector) Object() Node {
	return s.object
}

func NewSelector(o Node, adj, xtra string) Selector {
	sel := selector{object: o}
	switch adj {
	case "first", "last", "next", "previous", "forward", "backward":
		sel.adjective = cycled
		sel.raw = adj
	case "right", "down", "left", "up":
		sel.adjective = directed
		sel.raw = adj
	case "closest", "closer", "furthest", "further", "horizontal", "vertical":
		sel.adjective = oriented
		sel.raw = adj
	case "biggest", "smallest":
		sel.adjective = sized
		sel.raw = adj
	case "youngest", "younger", "oldest", "older":
		sel.adjective = aged
		sel.raw = adj
	case "primary", "focused", "unfocused", "current", "urgent":
		sel.adjective = focused
		sel.raw = adj
	case "name", "id", "index":
		sel.adjective = tagged
		sel.raw = fmt.Sprintf("%s %s", adj, xtra)
	case "local", "free", "occupied", "tiled", "floating", "like", "unlike", "manual", "automatic":
		sel.adjective = special
		sel.raw = adj
	}
	return sel
}

func Selectors(o Node, description string) []Selector {
	var ret []Selector
	descs := strings.Split(description, " ")
	for i, desc := range descs {
		if desc == "name" || desc == "id" || desc == "index" {
			ret = append(ret, NewSelector(o, desc, descs[i+1]))
		} else {
			ret = append(ret, NewSelector(o, desc, ""))
		}
	}
	return ret
}

/*
func (e *Euclid) Select(item string, description string) (coordinate, bool) {
	var obj Node
	loc := Coordinate(e, nil, nil, nil)
	var sel []Selector
	switch item {
	case "monitor":
		sel = Selectors(find, nMonitor, description)
		obj = nMonitor
	case "desktop":
		sel = Selectors(find, nDesktop, description)
		obj = nDesktop
	case "client":
		sel = Selectors(find, nClient, description)
		obj = nClient
	}
	return e.get(obj, loc, sel...)
}

func (e *Euclid) get(obj Node, loc coordinate, sel ...Selector) (coordinate, bool) {
	switch obj {
	case nMonitor:
		//return selectMonitor(e, loc, sel...)
	case nDesktop:
		//return selectDesktop(e, loc, sel...)
	case nClient:
		//return selectClient(e, loc, sel...)
	}
	return loc, false
}

func (e *Euclid) Match(item, description string, location, reference *coordinate) bool {
	var sel []Selector
	switch item {
	case "monitor":
		sel = Selectors(match, nMonitor, description)
	case "desktop":
		sel = Selectors(match, nDesktop, description)
	case "client":
		sel = Selectors(match, nClient, description)
	}
	return e.match(location, reference, sel...)
}

func (e *Euclid) match(location, reference *coordinate, sel ...Selector) bool {
	return false
}

//func locateWindow(h *Handler, w xproto.Window) (*Client, bool) {
//for _, m := range h.monitors {
//	for _, d := range m.desktops {
//		for _, c := range d.clients {
//			if w == c.Window.Window {
//				return c, true
//			}
//		}
//	}
//}
//return nil, false
//}
*/
