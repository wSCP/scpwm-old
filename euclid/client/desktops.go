package client

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
