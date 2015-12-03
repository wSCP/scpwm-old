package desktops

import (
	"fmt"

	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/settings"
)

type Desktop interface {
	settings.Settings
	Set(string, interface{})
	Index() int
	Monitor() uint32
	Id() string
	Name() string
	Focus()
	Focused() bool
	Floating() bool
	Layout() string
	Clients() *branch.Branch
}

func NewDesktop(monitor uint32, name string, s settings.Settings) Desktop {
	d := &desktop{
		Settings: s,
	}
	d.Set("name", name)
	d.Set("monitor", monitor)
	return d
}

type desktop struct {
	settings.Settings
	idx      int
	monitor  uint32
	name     string
	focused  bool
	floating bool
	layout   string
	clients  *branch.Branch
}

func (d *desktop) Set(k string, v interface{}) {
	switch k {
	case "index":
		if value, ok := v.(int); ok {
			d.idx = value
		}
	case "monitor":
		if value, ok := v.(uint32); ok {
			d.monitor = value
		}
	case "name":
		if value, ok := v.(string); ok {
			d.name = value
		}
	case "focus":
		if value, ok := v.(bool); ok {
			d.focused = value
		}
	case "floating":
		if value, ok := v.(bool); ok {
			d.floating = value
		}
	case "layout":
		if value, ok := v.(string); ok {
			d.layout = value
		}
	}
}

func (d *desktop) Index() int {
	return d.idx
}

func (d *desktop) Monitor() uint32 {
	return d.monitor
}

func (d *desktop) Id() string {
	return fmt.Sprintf("%s-%d-%d", d.name, d.idx, d.monitor)
}

func (d *desktop) Name() string {
	if d.name == "" {
		return d.defaultName()
	}
	return d.name
}

func (d *desktop) defaultName() string {
	return fmt.Sprintf("%s%d", d.String("DefaultDesktopName"), d.idx)
}

func (d *desktop) Focus() {
	d.focused = true
}

func (d *desktop) Focused() bool {
	return d.focused
}

func (d *desktop) Floating() bool {
	return d.floating
}

func (d *desktop) Layout() string {
	return d.layout
}

func (d *desktop) Clients() *branch.Branch {
	return d.clients
}

/*
import (
	"github.com/thrisp/scpwm/euclid/settings"
	"github.com/thrisp/scpwm/euclid/monitor"
)

type Desktop interface {
	Monitor() monitor.Monitor
	Focused() bool
	Last() bool
	Merge(Desktop)
	Arrange()
	Focus()
	Unfocus()
	Show()
	Hide()
	Urgent() bool
	//Clients() Clients
}

type desktop struct {
	///loc         coordinate
	m           monitor.Monitor
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

func New(m monitor.Monitor, n string) Desktop {
	//d := e.desktopDefault()
	//d.layout = tiled
	//d.pad = DefaultPad()
	//if n == "" {
	//	d.name = d.id
	//} else {
	//	d.name = n
	//}
	//d.loc.m = m
	//d.loc.d = d
	//d.clients = NewClients()
	//return d
	return nil
}

func (d *desktop) Focused() bool {
	return d.focused
}

func (d *desktop) Last() bool {
	return d.last
}

func (d *desktop) Merge(dst Desktop) {
	//for _, c := range dst.Clients() {
	//c.loc.m = d.loc.m
	//c.loc.d = d
	//}
	//d.clients = append(d.clients, dst.clients...)
	// reconfigure clients to this desktop
}



//func (d *desktop) TiledCount() int {
	//var cnt int
	//for _, c := range d.clients {
	//	if c.Tiled() {
	//		cnt++
	//	}
	//}
	//return cnt
//}

func (d *desktop) Focus() {
	//d.loc.m.Focus()
	//if !d.focused {
	//	d.focused = true
	//	d.Show()
	//	//ewmh_update_current_Desktop();
	//}
}

func (d *desktop) Show() {
	//if d.e.Bool("visible") {
	//	n := d.root.rightExtrema()
	//	for n != nil {
	//		//window_show(n->client->window);
	//		n = nextLeaf(n, d.root)
	//	}
	//}
}

func (d *desktop) Unfocus() {
	if d.focused {
		d.focused = false
		d.last = true
		d.Hide()
	}
}

func (d *desktop) Hide() {
	//if d.e.Bool("visible") {
	//	n := d.root.rightExtrema()
	//	for n != nil {
	//		//window_hide(n->client->window);
	//		n = nextLeaf(n, d.root)
	//	}
	//}
}

func (d *desktop) Urgent() bool {
	//for _, c := range d.clients {
	//	if c.urgent {
	//		return true
	//	}
	//}
	return false
}

func (d *desktop) Remove() {
	//if d.clients.Number() <= 0 {
	//	d.loc = NoCoordinate
	//	d.pad = nil
	//	d.clients = nil
	//	d = nil
	//}
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

//func (d *Desktop) Transfer(to *Monitor) {
	//cds := d.loc.m.desktops
	//cds.Pop(d)

	//adjust cds focus etc

	//nds := to.desktops
	//nds.d = append(nds.d, d)

	//d.loc.m = to
	//adjust everything to new monitor
	//ewmh_update_wm_Desktops();
	//ewmh_update_Desktop_names();
	//ewmh_update_current_Desktop();
//}

//func (d *Desktop) Neighbor(sel Selector) *Desktop {
//	return nil
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

//func selectDesktop(e *Euclid, sel ...Selector) bool {
//	return false
//}

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
