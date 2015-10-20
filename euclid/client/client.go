package client

import (
	"github.com/BurntSushi/xgb/xproto"

	"github.com/thrisp/scpwm/euclid/ruler"
)

type Client interface {
	Class() string
	Instance() string
	Rectangle() xproto.Rectangle
	Center(xproto.Rectangle)
	Adjacent(Client, string) bool
	SideHandle(string) *xproto.Point
	Rules() []ruler.Rule
	Window
	Floatr
	Tilr
	State
	Bordr
	Shiftr
	Stackr
}

type client struct {
	class    string
	instance string
	rules    []ruler.Rule
	*window
	*stackr
	*floatr
	*tilr
	*state
	*bordr
	*shiftr
}

func (c *client) Class() string {
	return c.class
}

func (c *client) Instance() string {
	return c.instance
}

func (c *client) Rules() []ruler.Rule {
	return c.rules
}

func (c *client) Rectangle() xproto.Rectangle {
	if c.Tiled() {
		return c.tilr.rectangle
	}
	return c.floatr.rectangle
}

func (c *client) Center(rect xproto.Rectangle) {
	r := c.floatr.rectangle

	if r.Width >= rect.Width {
		r.X = rect.X
	} else {
		r.X = rect.X + (int16(rect.Width)-int16(r.Width))/2
	}

	if r.Height >= rect.Height {
		r.Y = rect.Y
	} else {
		r.Y = rect.Y + (int16(rect.Height)-int16(r.Height))/2
	}

	r.X -= int16(c.borderWidth)
	r.Y -= int16(c.borderWidth)
}

func (c *client) Adjacent(o Client, direction string) bool {
	cr := c.Rectangle()
	or := o.Rectangle()
	switch direction {
	case "right":
		return (cr.X + int16(cr.Width)) == or.X
	case "down":
		return (cr.Y + int16(cr.Height)) == or.Y
	case "left":
		return (or.X + int16(or.Width)) == cr.X
	case "up":
		return (or.Y + int16(or.Height)) == cr.Y
	}
	return false
}

func (c *client) SideHandle(direction string) *xproto.Point {
	rect := c.Rectangle()
	var p xproto.Point
	switch direction {
	case "right":
		p.X = rect.X + int16(rect.Width)
		p.Y = rect.Y + (int16(rect.Height) / 2)
	case "down":
		p.X = rect.X + (int16(rect.Width) / 2)
		p.Y = rect.Y + int16(rect.Height)
	case "left":
		p.X = rect.X
		p.Y = rect.Y + (int16(rect.Height) / 2)
	case "up":
		p.X = rect.X + (int16(rect.Width) / 2)
		p.Y = rect.Y
	}
	return &p
}

/*
type clientType int

const (
	ctAll clientType = iota
	ctFloating
	ctTiled
)

var stringClientType map[string]clientType = map[string]clientType{
	"all":      ctAll,
	"floating": ctFloating,
	"tiled":    ctTiled,
}

type clientClass int

const (
	ccAll clientClass = iota
	ccEqual
	ccDiffer
)

var stringClientClass map[string]clientClass = map[string]clientClass{
	"all":    ccAll,
	"equal":  ccEqual,
	"differ": ccDiffer,
}

type clientMode int

const (
	cmAll clientMode = iota
	cmAutomatic
	cmManual
)

var stringClientMode map[string]clientMode = map[string]clientMode{
	"all":       cmAll,
	"automatic": cmAutomatic,
	"manual":    cmManual,
}

//func selectClient(e *Euclid, sel ...Selector) bool {
//	return false
//}

/*
func (n *Node) fence(dir Direction) *Node {
	/*if n != nil {
		t := n.top
		for t != nil {
			if (dir == Up && t.splitT == Horizontal && t.Rectangle.Y < n.Rectangle.Y) ||
				(dir == Left && t.splitT == Vertical && t.Rectangle.X < n.Rectangle.X) ||
				(dir == Down && t.splitT == Horizontal && (t.Rectangle.Y+int16(t.Rectangle.Height)) > (n.Rectangle.Y+int16(n.Rectangle.Height))) ||
				(dir == Right && t.splitT == Vertical && (t.Rectangle.X+int16(t.Rectangle.Width)) > (n.Rectangle.X+int16(n.Rectangle.Width))) {
				return t
			}
			t = t.top
		}
	}
	return nil
}

type Node struct {
	idx   int
	loc   coordinate
	focus bool
	next  *Node
	prev  *Node
	*xproto.Rectangle
	*Client
	*shift
}

func newNode(root *Node, idx int) *Node {
	n := &Node{
		idx:   idx,
		shift: &shift{},
	}
	loc := root.loc
	loc.n = n
	n.loc = loc
	return n
}

func (n *Node) peek(cy Cycle) *Node {
	var ret *Node
	switch cy {
	case Next:
		ret = n.next
	case Prev:
		ret = n.prev
	}
	return ret
}

func (n *Node) pop(cy Cycle) *Node {
	ret := n.peek(cy)
	if ret.idx != 0 {
		ret.detach()
		return ret
	}
	return nil
}

func (n *Node) detach() {
	if n.idx != 0 {
		next = n.next
		prev = n.prev
		if next != nil {
			next.prev = prev
		}
		if prev != nil {
			prev.next = next
		}
	}
}

func (n *Node) push(o *Node, cy Cycle) {
	switch cy {
	case Next:
		if n.next != nil {
			nxt = n.next
			o.prev = n
			o.nxt = nxt
			n.next = o
		} else {
			n.next = o
		}
	case Prev:
		//
	}
}

func (n *Node) closest(cy Cycle, sel ClientSelect) *Node {
	/*curr := n.peek(cy)
	if curr != nil {
		ref := n.loc
		for curr != n {
			loc := n.loc
			loc.n = curr
			if MatchNode(&loc, &ref, sel) {
				return curr
			}
			curr = curr.peek(cy)
		}
	}
	return nil
}



func (n *Node) applyLayout(r, rr xproto.Rectangle) {}

func (n *Node) destruct() {
	if n != nil {
		t1 := n.right
		t2 := n.left
		if n.Client != nil {
			n.Client = nil
			NumClients--
		}
		n = nil
		t1.destruct()
		t2.destruct()
	}
}

//func (n *Node) focus() {}

func (n *Node) pseudoFocus() {}

func (n *Node) neighbor(dir Direction, sel ClientSelect) *Node {
	/*if n != nil || !n.Client.fullscreen || (n.loc.d.layout != monocle && !n.IsTiled()) {
		var ret *Node

		if n.loc.e.Bool("HistoryAwareFocus") {
			ret = n.neighborHistory(dir, sel)
		}

		if ret == nil {
			if n.loc.e.Bool("FocusByDistance") {
				ret = n.neighborDistance(dir, sel)
			} else {
				ret = n.neighborTree(dir, sel)
			}
		}

		return ret
	}
	return nil
}

func (n *Node) neighborTree(dir Direction, sel ClientSelect) *Node {
	/*if n != nil {
		fence := n.fence(dir)

		if fence != nil {
			var nearest *Node

			if dir == Up || dir == Left {
				nearest = fence.right.leftExtrema()
			} else if dir == Down || dir == Right {
				nearest = fence.left.rightExtrema()
			}

			ref, loc := n.loc, coordinates(n.loc.e, n.loc.m, n.loc.d, nearest)

			if MatchNode(&loc, &ref, sel) {
				return nearest
			} else {
				return nil
			}
		}
	}
	return nil
}

func (n *Node) neighborHistory(dir Direction, sel ClientSelect) *Node {
	/*if n != nil || n.IsTiled() {
		target := n.fence(dir)
		if target != nil {
			switch dir {
			case Up, Left:
				target = target.right
			case Down, Right:
				target = target.left
			}
			var nearest *Node
			//int min_rank = INT_MAX;
			ref := n.loc
			a := target.rightExtrema()
			for a != nil {
				if !a.vacant || n.isAdjacent(a, dir) || a != n {
					//loc := coordinates(n.e, m, d, a)
					//if NodeMatches(loc, ref, sel) {
					//int rank = history_rank(d, a);
					//if (rank >= 0 && rank < min_rank) {
					//	nearest = a;
					//	min_rank = rank;
					//}
				}
				a = nextLeaf(a, target)
			}
			return nearest
		}
	}
	return nil
}

func (n *Node) neighborDistance(dir Direction, sel ClientSelect) *Node {
	/*if n != nil {
		var target *Node
		if n.IsTiled() {
			target = n.fence(dir)
			if target == nil {
				return nil
			}
			if dir == Up || dir == Left {
				target = target.right
			} else if dir == Down || dir == Right {
				target = target.left
			}
		} else {
			//target = d.root
		}

		var nearest *Node
		var dir2 Direction
		var p1, p2 xproto.Point
		//n.getSideHandle(dir, &p1)
		opposite(dir, dir2)
		//n.getSideHandle(dir, &p1)
		opposite(dir, dir2)
		var ds float64 //double ds = DBL_MAX;
		ref := n.loc

		a := target.rightExtrema()
		for a != nil {
			l := ref
			l.n = a
			loc := l //coordinates(m, d, a)
			if a != n || MatchNode(&loc, &ref, sel) || a.IsTiled() == n.IsTiled() || (a.IsTiled() && n.isAdjacent(a, dir)) {
				//	a.getSideHandle(dir2, &p2)
				//	ds2 := distance(p1, p2)
				//	if ds2 < ds {
				//		ds = ds2
				//		nearest = a
				//	}
			}
			a := nextLeaf(a, target)
		}
		return nearest
	}
	return nil
}

func setVacant(n *Node) {
	curr := n
	for curr != nil {
		curr.vacant = (curr.right.vacant && curr.left.vacant)
		curr = curr.top
	}
}

func setPrivacy(n *Node, as bool) {
	var v int
	if as {
		v = 1
	} else {
		v = -1
	}
	curr := n
	for curr != nil {
		curr.privacy += v
		curr = curr.top
	}
}

func NodeFromDescription(desc []string, ref, dst *coordinate) bool {
	/*
			sel := clientselectFromString(desc)
			dst.m = ref.m
			dst.d = ref.d
			dst.n = nil
			var dir Direction
			var cy Cycle
			var hdir Age
			if directionFromString(desc, dir) {
					//dst->node = nearest_neighbor(ref->Monitor, ref->Desktop, ref->node, dir, sel);
					//if (dst->node == NULL && num_Monitors > 1) {
					//	Monitor_t *m = nearest_Monitor(ref->Monitor, dir, (Desktop_select_t) {DESKTOP_STATUS_ALL, false, false});
					//	if (m != NULL) {
					//		coordinates_t loc = {m, m->desk, m->desk->focus};
					//		if (node_matches(&loc, ref, sel)) {
					//			dst->Monitor = m;
					//			dst->Desktop = m->desk;
					//			dst->node = m->desk->focus;
					//		}
					//	}
					//}
			} else if cycleFromString(desc, cy) {
				//dst->node = closest_node(ref->Monitor, ref->Desktop, ref->node, cy, sel);
			} else if ageFromString(desc, hdir) {
				//history_find_node(hdir, ref, dst, sel);
			} else if strings.Contains(desc, "last") {
				//	history_find_node(HISTORY_OLDER, ref, dst, sel);
			} else if strings.Contains(desc, "biggest") {
				//	dst->node = find_biggest(ref->Monitor, ref->Desktop, ref->node, sel);
			} else if strings.Contains(desc, "focused") {
				//	coordinates_t loc = {mon, mon->desk, mon->desk->focus};
				//	if (node_matches(&loc, ref, sel)) {
				//		dst->Monitor = mon;
				//		dst->Desktop = mon->desk;
				//		dst->node = mon->desk->focus;
				//	}
			} else {
				//	long int wid;
				//	if (parse_window_id(desc, &wid))
				//		locate_window(wid, dst);
			}
		return dst.n != nil
	return false
}

func MatchNode(loc, ref *coordinate, sel ClientSelect) bool {
	if loc.n != nil {
	}
	return false
}

/*
type Nodes struct {
	count int
	*Node
}

func NewNodes(e *Euclid, m *Monitor, d *Desktop) *Nodes {
	n := &Node{idx: 0}
	n.loc = coordinates(e, m, d, n)
	return &Nodes{
		Node: n,
	}
}

func (ns *Nodes) New() *Node {
	ns.count++
	return newNode(ns.Node, ns.count)
}

func (ns *Nodes) Biggest(sel ClientSelect) *Node {
	return nil
}

func (ns *Nodes) Circulate(cy Cycle) {}

//func (ns *Nodes) Swap(n1, n2 *Node) bool {
//	return false
//}

//func SwapNodes(ns1, ns2 *Nodes) bool {
//	return false
//}
*/

/*
func (c *coordinate) resetMode() {
	if c.n != nil {
		c.n.splitM = Automatic
		//window_draw_border(loc->Node, loc->Desktop->focus == loc->Node, mon == loc->Monitor);
	} else if c.d != nil {
		a := c.d.root.rightExtrema()
		for a != nil {
			a.splitM = Automatic
			//window_draw_border(a, loc->Desktop->focus == a, mon == loc->Monitor);
			a = nextLeaf(a, c.d.root)
		}
	}
}
*/
