package main

import (
	"fmt"
	"strings"
)

type coordinate struct {
	e *Euclid
	m *Monitor
	d *Desktop
	c *Client
}

var NoCoordinate = coordinate{nil, nil, nil, nil}

func Coordinate(e *Euclid, m *Monitor, d *Desktop, c *Client) coordinate {
	return coordinate{
		e: e,
		m: m,
		d: d,
		c: c,
	}
}

type verb int

const (
	find verb = iota
	match
)

type adjective int

const (
	cycled   adjective = iota // Cycle
	directed                  // Direction
	oriented                  // Orientation
	sized                     // Size
	aged                      // Age
	focused                   // by focus, urgency
	tagged                    // by name, tag or id
	special                   // relatively unique
)

type Selector interface {
	Verb() verb
	Adjective() adjective
	Raw() string
	Object() Node
}

type selector struct {
	verb      verb
	adjective adjective
	raw       string
	object    Node
}

func (s selector) Verb() verb {
	return s.verb
}

func (s selector) Adjective() adjective {
	return s.adjective
}

func (s selector) Raw() string {
	return s.raw
}

func (s selector) Object() Node {
	return s.object
}

func NewSelector(v verb, o Node, adj, xtra string) Selector {
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

func Selectors(v verb, o Node, description string) []Selector {
	var ret []Selector
	descs := strings.Split(description, " ")
	for i, desc := range descs {
		if desc == "name" || desc == "id" || desc == "index" {
			ret = append(ret, NewSelector(v, o, desc, descs[i+1]))
		} else {
			ret = append(ret, NewSelector(v, o, desc, ""))
		}
	}
	return ret
}

func (e *Euclid) Locate(item string, description string) (coordinate, bool) {
	loc := Coordinate(e, nil, nil, nil)
	var sel []Selector
	var obj Node
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
	return e.locate(obj, loc, sel...)
}

func (e *Euclid) locate(obj Node, loc coordinate, sel ...Selector) (coordinate, bool) {
	switch obj {
	case nMonitor:
		return locateMonitor(e, loc, sel...)
	case nDesktop:
		return locateDesktop(e, loc, sel...)
	case nClient:
		return locateClient(e, loc, sel...)
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
