package selector

import "strings"

type Node int

const (
	NNone Node = iota
	NMonitor
	NDesktop
	NClient
)

var stringNode map[string]Node = map[string]Node{
	"none":    NNone,
	"monitor": NMonitor,
	"desktop": NDesktop,
	"client":  NClient,
}

func (n Node) String() string {
	switch n {
	case NNone:
		return "none"
	case NMonitor:
		return "monitor"
	case NDesktop:
		return "desktop"
	case NClient:
		return "client"
	}
	return ""
}

func DetermineNode(from string) Node {
	switch strings.ToLower(from) {
	case "", "none":
		return NNone
	case "monitor", "screen":
		return NMonitor
	case "desktop", "workspace":
		return NDesktop
	case "window", "client":
		return NClient
	}
	return NNone
}

type Category int

const (
	NoCategory Category = iota
	Cycled              // Cycle
	Directed            // Direction
	Oriented            // Orientation
	Sized               // Size
	Aged                // Age
	Focused             // by focus, urgency
	Tagged              // by name, tag or id
	Capacity            // Capacity(of monitor or desktop state)
	State               // client state
)

type Selector interface {
	Raw() []string
	Node() Node
	Category() Category
	Modifiers() []string
}

func New(base string) Selector {
	spl := strings.Split(base, " ")
	node := spl[0]
	cat := spl[1]
	modifiers := spl[1:]
	sel := selector{node: DetermineNode(node)}
	switch cat {
	case "first", "last", "next", "previous", "forward", "backward":
		sel.cat = Cycled
	case "right", "down", "left", "up":
		sel.cat = Directed
	case "closest", "closer", "furthest", "further", "horizontal", "vertical":
		sel.cat = Oriented
	case "biggest", "smallest":
		sel.cat = Sized
	case "youngest", "younger", "oldest", "older":
		sel.cat = Aged
	case "primary", "focused", "unfocused", "current", "urgent":
		sel.cat = Focused
	case "name", "id", "index", "class", "instance":
		sel.cat = Tagged
	case "free", "empty", "full", "occupied", "tiled", "floating":
		sel.cat = Capacity
	case "local", "like", "unlike", "manual", "automatic":
		sel.cat = State
	}
	if sel.cat == NoCategory {
		return nil
	}
	sel.raw = spl
	sel.mod = modifiers
	return sel
}

type selector struct {
	raw  []string
	node Node
	cat  Category
	mod  []string
}

func (s selector) Raw() []string {
	return s.raw
}

func (s selector) Node() Node {
	return s.node
}

func (s selector) Category() Category {
	return s.cat
}

func (s selector) Modifiers() []string {
	return s.mod
}
