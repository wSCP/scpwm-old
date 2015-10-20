package query

import (
	"fmt"
	"strings"
)

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
	special                   // relatively unique property of Node, e.g. any client state
)

type Selector interface {
	Verb() verb
	Adjective() adjective
	Raw() string
	Object() Node
	Matcher
}

type Matcher interface {
	Reference() *coordinate
	Location() *coordinate
	Set(string, *coordinate)
}

func defaultMatcher() Matcher {
	return matcher{ref: &NoCoordinate, loc: &NoCoordinate}
}

type matcher struct {
	ref *coordinate
	loc *coordinate
}

func (m matcher) Reference() *coordinate {
	return m.ref
}

func (m matcher) Location() *coordinate {
	return m.loc
}

func (m matcher) Set(k string, c *coordinate) {
	switch k {
	case "reference":
		m.ref = c
	case "location":
		m.loc = c
	}
}

type selector struct {
	verb      verb
	adjective adjective
	raw       string
	object    Node
	Matcher
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
	sel := selector{object: o, Matcher: defaultMatcher()}
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
