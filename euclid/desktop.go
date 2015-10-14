package main

import "fmt"

type Desktops []*Desktop

func NewDesktops(e *Euclid, m *Monitor) Desktops {
	return make(Desktops, 0)
}

func (ds Desktops) Number() int {
	return len(ds)
}

func (ds Desktops) Add(e *Euclid, m *Monitor, name string) *Desktop {
	d := newDesktop(e, m, name)
	ds = append(ds, d)
	return d
	//ewmh_update_number_of_Desktops()
	//ewmh_update_Desktop_names()
}

func (ds Desktops) Find(sel Selector, loc coordinate) bool {
	return false
}

func (ds Desktops) Match(sel Selector, ref, loc *coordinate) bool {
	return false
}

func (ds Desktops) Remove(r *Desktop) {
	for i, d := range ds {
		if r.id == d.id {
			ds = ds[:i+copy(ds[i:], ds[i+1:])]
			if d.Focused() {
				nd := ds[i]
				//shift all d Clients to nd
				nd.Focus()
			}
			d.Remove()
		}
	}
	//ewmh_update_current_Desktop()
	//ewmh_update_number_of_Desktops()
	//ewmh_update_Desktop_names()
}

func (ds Desktops) Focused() *Desktop {
	for _, d := range ds {
		if d.Focused() {
			return d
		}
	}
	return nil
}

func (ds Desktops) Last() *Desktop {
	for _, d := range ds {
		if d.Last() {
			return d
		}
	}
	return nil
}

func (ds Desktops) Pop(p *Desktop) *Desktop {
	for i, d := range ds {
		if d.id == p.id {
			ds = ds[:i+copy(ds[i:], ds[i+1:])]
		}
		return d
	}
	return nil
}

func (ds Desktops) Swap(d1, d2 *Desktop) {
	var i1, i2 int
	for i, d := range ds {
		if d.id == d1.id {
			i1 = i
		}
		if d.id == d2.id {
			i2 = i
		}
	}
	if i1 != i2 {
		ds[i2] = d1
		ds[i1] = d2
	}
	// update_input_focus();
	// ewmh_update_wm_Desktops();
	// ewmh_update_Desktop_names();
	// ewmh_update_current_Desktop();
}

func (e *Euclid) desktopDefault() *Desktop {
	e.desktopCount++
	id := fmt.Sprintf("%s%d", e.String("defaultDesktopName"), e.desktopCount)
	return &Desktop{
		loc:         Coordinate(e, nil, nil, nil),
		id:          id,
		windowGap:   e.Int("WindowGap"),
		borderWidth: uint(e.Int("BorderWidth")),
	}
}

type Desktop struct {
	loc         coordinate
	id          string
	name        string
	layout      Layout
	pad         Pad
	windowGap   int
	borderWidth uint
	floating    bool
	focused     bool
	last        bool
	clients     Clients
}

func newDesktop(e *Euclid, m *Monitor, n string) *Desktop {
	d := e.desktopDefault()
	d.layout = tiled
	d.pad = DefaultPad()
	if n == "" {
		d.name = d.id
	} else {
		d.name = n
	}
	d.loc.m = m
	d.loc.d = d
	d.clients = NewClients()
	return d
}

func (d *Desktop) Focused() bool {
	return d.focused
}

func (d *Desktop) Last() bool {
	return d.last
}

func (d *Desktop) Merge(dst *Desktop) {
	for _, c := range dst.clients {
		c.loc.m = d.loc.m
		c.loc.d = d
	}
	d.clients = append(d.clients, dst.clients...)
	// reconfigure clients to this desktop
}

func (d *Desktop) Arrange() {
	/*
		if d.root != nil {
			setLayout := d.layout

			if d.e.Bool("LeafMonocle") && d.TiledCount() == 1 {
				d.layout = monocle
			}

			rect := d.m.rectangle
			var gap int
			if d.e.Bool("GaplessMonocle") && d.layout == monocle {
				gap = 0
			} else {
				gap = d.windowGap
			}

			rect.X += int16(d.m.pad[Left] + d.pad[Left] + gap)
			rect.Y += int16(d.m.pad[Up] + d.pad[Up] + gap)
			rect.Width -= uint16(d.m.pad[Left] + d.pad[Left] + d.pad[Right] + d.m.pad[Right] + gap)
			rect.Height -= uint16(d.m.pad[Up] + d.pad[Up] + d.pad[Down] + d.m.pad[Down] + gap)
			//ApplyLayout(d.m, d, d.root, rect, rect)

			d.layout = setLayout
		}
	*/
}

func (d *Desktop) TiledCount() int {
	var cnt int
	for _, c := range d.clients {
		if c.IsTiled() {
			cnt++
		}
	}
	return cnt
}

func (d *Desktop) Focus() {
	d.loc.m.Focus()
	if !d.focused {
		d.focused = true
		d.Show()
		//ewmh_update_current_Desktop();
	}
}

func (d *Desktop) Show() {
	//if d.e.Bool("visible") {
	//	n := d.root.rightExtrema()
	//	for n != nil {
	//		//window_show(n->client->window);
	//		n = nextLeaf(n, d.root)
	//	}
	//}
}

func (d *Desktop) UnFocus() {
	if d.focused {
		d.focused = false
		d.last = true
		d.Hide()
	}
}

func (d *Desktop) Hide() {
	//if d.e.Bool("visible") {
	//	n := d.root.rightExtrema()
	//	for n != nil {
	//		//window_hide(n->client->window);
	//		n = nextLeaf(n, d.root)
	//	}
	//}
}

func (d *Desktop) Urgent() bool {
	for _, c := range d.clients {
		if c.urgent {
			return true
		}
	}
	return false
}

func (d *Desktop) Remove() {
	if d.clients.Number() <= 0 {
		d.loc = NoCoordinate
		d.pad = nil
		d.clients = nil
		d = nil
	}
	// raise a notice of some sort
}

type Layout int

const (
	tiled Layout = iota
	monocle
)

func (d *Desktop) Change(l Layout) {
	d.layout = l
	d.Arrange()
}

//func (d *Desktop) Transfer(src, dst *Monitor) {
/*if src != dst {
	prev := d.prev
	next := d.next
	dsrc := src.Desktops

	if d.prev == nil {
		dsrc.Desktop = d.next
	} else {
		prev.next = next
		next.prev = prev
	}

	if d.IsCurrent() {
		d.UnFocus()
		last := dsrc.Last()
		dsrc.SetCurrent(last)
	}

	d.m = dst
	dst.Desktops.Tail().next = d
	d.prev = dst.Desktops.Tail()

	n := d.root.rightExtrema()
	for n != nil {
		dst.Translate(src, n.Client)
		n = nextLeaf(n, d.root)
	}

	d.Arrange()

	d.e.history.transferDesktop(dst, d)

	//ewmh_update_wm_Desktops();
	//ewmh_update_Desktop_names();
	//ewmh_update_current_Desktop();
}
*/
//}

//func (d *Desktop) Closest(cy Cycle, sel DesktopSelect) *Desktop {
//curr := d.Pop(cy)
//for curr != d {
//	if MatchDesktop(&curr.loc, &curr.loc, sel) {
//		return curr
//	}
//	curr = curr.Pop(cy)
//}
//return nil
//}

type desktopStatus int

const (
	dsAll desktopStatus = iota
	dsFree
	dsOccupied
)

var stringDesktopStatus map[string]desktopStatus = map[string]desktopStatus{
	"all":      dsAll,
	"free":     dsFree,
	"occupied": dsOccupied,
}

type desktopUrgency int

const (
	duAll desktopUrgency = iota
	duOn
	duOff
)

var stringDesktopUrgency map[string]desktopUrgency = map[string]desktopUrgency{
	"all": duAll,
	"on":  duOn,
	"off": duOff,
}

func locateDesktop(e *Euclid, loc coordinate, sel ...Selector) (coordinate, bool) {
	return loc, false
}

/*
func DesktopFromDescription(desc []string, ref, dst *coordinate) bool {
	//sel := DesktopSelectFromString(desc)
	//dst.d = nil
		cycle_dir_t cyc;
		history_dir_t hdi;
		char *colon;
		int idx;
		if (parse_cycle_direction(desc, &cyc)) {
			dst->monitor = ref->monitor;
			dst->Desktop = closest_Desktop(ref->monitor, ref->Desktop, cyc, sel);
		} else if (parse_history_direction(desc, &hdi)) {
			history_find_Desktop(hdi, ref, dst, sel);
		} else if (streq("last", desc)) {
			history_find_Desktop(HISTORY_OLDER, ref, dst, sel);
		} else if (streq("focused", desc)) {
			coordinates_t loc = {mon, mon->desk, NULL};
			if (Desktop_matches(&loc, ref, sel)) {
				dst->monitor = mon;
				dst->Desktop = mon->desk;
			}
		} else if ((colon = strchr(desc, ':')) != NULL) {
			*colon = '\0';
			if (monitor_from_desc(desc, ref, dst)) {
				if (streq("focused", colon + 1)) {
					dst->Desktop = dst->monitor->desk;
				} else if (parse_index(colon + 1, &idx)) {
					Desktop_from_index(idx, dst, dst->monitor);
				}
			}
		} else if (parse_index(desc, &idx)) {
			Desktop_from_index(idx, dst, NULL);
		} else {
			locate_Desktop(desc, dst);
		}

		return (dst->Desktop != NULL);
		}
	return false
}

func DesktopFromIndex(idx int, loc *coordinate) bool {
	mon := m.Pop(Head)
	for mon != nil {
		if m != nil && m != mon {
			d := mon.Desktop.Pop(Head)
			for d != nil {
				if idx == 1 {
					loc.m = mon
					loc.d = d
					loc.n = nil
					return true
				}
				d = d.Pop(Next)
				idx--
			}
		}
		mon = mon.Pop(Next)
	}
	return false
}

func MatchDesktop(loc, ref *coordinate, sel DesktopSelect) bool {
	if sel.status != dsAll && loc.d.root == nil {
		if sel.status == dsOccupied || sel.status == dsFree {
			return false
		}
	}

	if sel.urgent && !loc.d.Urgent() {
		return false
	}

	if sel.local && ref.m != loc.m {
		return false
	}

	return true
}
*/
