package main

import (
	"bytes"
	"fmt"
	"strings"
)

type query struct {
	*bytes.Buffer
}

func Query() *query {
	return &query{new(bytes.Buffer)}
}

//func (q *query) monitors(loc *Coordinate) *query {
//void query_monitors(coordinates_t loc, domain_t dom, FILE *rsp);
//	return q
//}

//func (q *query) desktops(m *Monitor, loc *Coordinate, depth uint, dom domain) *query {
//void query_desktops(monitor_t *m, domain_t dom, coordinates_t loc, unsigned int depth, FILE *rsp);
//	return q
//}

//func (q *query) Node(d *Desktop, n *Node, depth uint) *query {
//void query_tree(desktop_t *d, node_t *n, FILE *rsp, unsigned int depth);
//return q
//}

//func (q *query) History(loc *Coordinate) *query {
//void query_history(coordinates_t loc, FILE *rsp);
//return q
//}

//func (q *query) Stack() *query {
//void query_stack(FILE *rsp);
//return q
//}

//func (q *query) Windows(loc *Coordinate) *query {
//void query_windows(coordinates_t loc, FILE *rsp);
//	return q
//}

//func (q *query) Pointer(ptr *Pointer) *query {
//	return q
//}

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

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up
)

var stringDirection = map[string]Direction{
	"right": Right,
	"down":  Down,
	"left":  Left,
	"up":    Up,
}

func (d Direction) String() string {
	switch d {
	case Right:
		return "right"
	case Down:
		return "down"
	case Left:
		return "left"
	case Up:
		return "up"
	}
	return ""
}

func opposite(src Direction) Direction {
	var dst Direction
	switch src {
	case Right:
		dst = Left
	case Down:
		dst = Up
	case Left:
		dst = Right
	case Up:
		dst = Down
	}
	return dst
}

type Cycle int

const (
	First Cycle = iota
	Backward
	Prev
	Next
	Forward
	Last
)

var stringCycle map[string]Cycle = map[string]Cycle{
	"first":    First,
	"backward": Backward,
	"previous": Prev,
	"next":     Next,
	"forward":  Forward,
	"last":     Last,
}

func (c Cycle) String() string {
	switch c {
	case First:
		return "first"
	case Forward:
		return "forward"
	case Prev:
		return "previous"
	case Next:
		return "next"
	case Backward:
		return "backward"
	case Last:
		return "last"
	}
	return ""
}

func reverse(src Cycle) Cycle {
	var dst Cycle
	switch src {
	case First:
		dst = Last
	case Forward:
		dst = Backward
	case Prev:
		dst = Next
	case Next:
		dst = Prev
	case Backward:
		dst = Forward
	case Last:
		dst = First
	}
	return dst
}

type Size int

const (
	Biggest Size = iota
	Smallest
)

var stringSize map[string]Size = map[string]Size{
	"biggest":  Biggest,
	"smallest": Smallest,
}

type Orientation int

const (
	Closer Orientation = iota
	Closest
	Further
	Furthest
	Horizontal
	Vertical
	Above
	Below
)

var stringOrientation map[string]Orientation = map[string]Orientation{
	"closer":     Closer,
	"closest":    Closest,
	"further":    Further,
	"furthest":   Furthest,
	"horizontal": Horizontal,
	"vertical":   Vertical,
	"above":      Above,
	"below":      Below,
}

func (o Orientation) String() string {
	switch o {
	case Closer:
		return "closer"
	case Closest:
		return "closest"
	case Further:
		return "further"
	case Furthest:
		return "furthest"
	case Horizontal:
		return "horizontal"
	case Vertical:
		return "vertical"
	case Above:
		return "above"
	case Below:
		return "below"
	}
	return ""
}

type Age int

const (
	Youngest Age = iota
	Younger
	Older
	Oldest
)

var stringAge map[string]Age = map[string]Age{
	"youngest": Youngest,
	"younger":  Younger,
	"older":    Older,
	"oldest":   Oldest,
}

type Node int

const (
	nWM Node = iota
	nMonitor
	nDesktop
	nClient
)

var stringNode map[string]Node = map[string]Node{
	"euclid":  nWM,
	"monitor": nMonitor,
	"desktop": nDesktop,
	"client":  nClient,
}

func (n Node) String() string {
	switch n {
	case nWM:
		return "euclid"
	case nMonitor:
		return "monitor"
	case nDesktop:
		return "desktop"
	case nClient:
		return "client"
	}
	return ""
}

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
		return selectMonitor(e, loc, sel...)
	case nDesktop:
		return selectDesktop(e, loc, sel...)
	case nClient:
		return selectClient(e, loc, sel...)
	}
	return loc, false
}

/*
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
*/
