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

type Reference int

const (
	NoReference Reference = iota
	Cycled                // Cycle
	Directed              // Direction
	Oriented              // Orientation
	Sized                 // Size
	Aged                  // Age
	Focused               // by focus, urgency
	Tagged                // by name, tag or id
	Special               // relatively unique property of Node, e.g. any client state
)

type Selector interface {
	Raw() []string
	Node() Node
	References() Reference
	Modifiers() []string
}

func New(base string) Selector {
	spl := strings.Split(base, " ")
	node := spl[0]
	ref := spl[1]
	modifiers := spl[1:]
	sel := selector{node: DetermineNode(node)}
	switch ref {
	case "first", "last", "next", "previous", "forward", "backward":
		sel.ref = Cycled
	case "right", "down", "left", "up":
		sel.ref = Directed
	case "closest", "closer", "furthest", "further", "horizontal", "vertical":
		sel.ref = Oriented
	case "biggest", "smallest":
		sel.ref = Sized
	case "youngest", "younger", "oldest", "older":
		sel.ref = Aged
	case "primary", "focused", "unfocused", "current", "urgent":
		sel.ref = Focused
	case "name", "id", "index":
		sel.ref = Tagged
	case "local", "free", "occupied", "tiled", "floating", "like", "unlike", "manual", "automatic":
		sel.ref = Special
	}
	if sel.ref == NoReference {
		return nil
	}
	sel.raw = spl
	sel.mod = modifiers
	return sel
}

type selector struct {
	raw  []string
	node Node
	ref  Reference
	mod  []string
}

func (s selector) Raw() []string {
	return s.raw
}

func (s selector) Node() Node {
	return s.node
}

func (s selector) References() Reference {
	return s.ref
}

func (s selector) Modifiers() []string {
	return s.mod
}
