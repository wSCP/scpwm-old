package desktops

import (
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/settings"
)

func New(monitor uint32, s settings.Settings) *branch.Branch {
	return branch.New("desktops")
}

func Initialize(desktops *branch.Branch, monitor uint32, s settings.Settings) {
	if desktops.Len() == 0 {
		d := NewDesktop(monitor, "", s)
		desktops.PushFront(d)
	}
	UpdateDesktopsIndex(desktops)
}

func Add(desktops *branch.Branch, monitor uint32, name string, s settings.Settings) {
	nd := NewDesktop(monitor, name, s)
	desktops.PushBack(nd)
	UpdateDesktopsIndex(desktops)
}

type MatchDesktop func(Desktop) bool

func seek(desktops *branch.Branch, fn MatchDesktop) Desktop {
	curr := desktops.Front()
	for curr != nil {
		desktop := curr.Value.(Desktop)
		if match := fn(desktop); match {
			return desktop
		}
		curr = curr.Next()
	}
	return nil
}

func seekOffset(desktops *branch.Branch, fn MatchDesktop, offset int) Desktop {
	curr := desktops.Front()
	for curr != nil {
		desktop := curr.Value.(Desktop)
		if match := fn(desktop); match {
			switch offset {
			case -1:
				desktop = curr.PrevContinuous().Value.(Desktop)
			case 1:
				desktop = curr.NextContinuous().Value.(Desktop)
			}
			return desktop
		}
		curr = curr.Next()
	}
	return nil
}

func isFocused(d Desktop) bool {
	if d.Focused() {
		return true
	}
	return false
}

func Focused(desktops *branch.Branch) Desktop {
	return seek(desktops, isFocused)
}

func Prev(desktops *branch.Branch) Desktop {
	return seekOffset(desktops, isFocused, -1)
}

func Next(desktops *branch.Branch) Desktop {
	return seekOffset(desktops, isFocused, 1)
}

func seekAny(desktops *branch.Branch, fn MatchDesktop) []Desktop {
	var ret []Desktop
	curr := desktops.Front()
	for curr != nil {
		dsk := curr.Value.(Desktop)
		if match := fn(dsk); match {
			ret = append(ret, dsk)
		}
		curr = curr.Next()
	}
	return ret
}

func All(desktops *branch.Branch) []Desktop {
	fn := func(d Desktop) bool { return true }
	return seekAny(desktops, fn)
}

type UpdateDesktop func(d Desktop) error

func update(desktops *branch.Branch, fn UpdateDesktop) error {
	curr := desktops.Front()
	for curr != nil {
		d := curr.Value.(Desktop)
		if err := fn(d); err != nil {
			return err
		}
		curr = curr.Next()
	}
	return nil
}

func UpdateDesktopsMonitor(desktops *branch.Branch, id uint32) {
	fn := func(d Desktop) error {
		d.Set("monitor", id)
		return nil
	}
	update(desktops, fn)
}

func UpdateDesktopsIndex(desktops *branch.Branch) {
	idx := 1
	fn := func(d Desktop) error {
		d.Set("index", idx)
		idx++
		return nil
	}
	update(desktops, fn)
}

/*
type Desktops interface {

}

type desktops struct {
	d       []*Desktop
	clients Clients
}

func NewDesktops() Desktops {
	d := &desktops{
		d: make([]Desktop, 0),
	}
	return d
}

func (ds *desktops) Settings(s *settings.Settings) {}

func (ds *desktops) Add(*Desktop) {
	//d := newDesktop(e, m, name)
	//ds = append(ds, d)
	//return d
	//ewmh_update_number_of_Desktops()
	//ewmh_update_Desktop_names()
}

//func (ds Desktops) Select(sel ...Selector) bool {
//	return false
//}

func (ds *desktops) Pop(p *Desktop) *Desktop {
	for i, d := range ds.d {
		if d.id == p.id {
			ds.d = ds.d[:i+copy(ds.d[i:], ds.d[i+1:])]
		}
		return d
	}
	return nil
}

func (ds *desktops) Swap(d1, d2 int) {
	/*
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

func (ds *desktops) Focused() *Desktop {
	for _, d := range ds.d {
		if d.Focused() {
			return d
		}
	}
	return nil
}

func (ds *desktops) Last() *Desktop {
	for _, d := range ds.d {
		if d.Last() {
			return d
		}
	}
	return nil
}

func (ds *desktops) Remove(r *Desktop) {
	for i, d := range ds.d {
		if r.id == d.id {
			ds.d = ds.d[:i+copy(ds.d[i:], ds.d[i+1:])]
			if d.Focused() {
				nd := ds.d[i]
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

func (ds *desktops) Number() int {
	return len(ds.d)
}

func (ds *desktops) Clients() Clients {
	return ds.clients
}

//func (e *Euclid) desktopDefault() *Desktop {
//	e.desktopCount++
//	id := fmt.Sprintf("%s%d", e.String("defaultDesktopName"), e.desktopCount)
//	return &Desktop{
//		loc:         Coordinate(e, nil, nil, nil),
//		id:          id,
//		windowGap:   e.Int("WindowGap"),
//		borderWidth: uint(e.Int("BorderWidth")),
//	}
//}

*/
