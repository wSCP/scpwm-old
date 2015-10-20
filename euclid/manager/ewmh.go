package manager

/*
import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/ewmh"
)

type Ewmh interface {
	DesktopsState
	WmState
}

type extendedWindowManagerHints struct {
	DesktopsState
	WmState
}

func NewEwmh(e ewmh.EWMH) Ewmh {
	return &extendedWindowManagerHints{
		DesktopsState: newDesktopsState(e),
		WmState:       newWmState(e),
	}
}

type DesktopsState interface {
	GetDesktopIndex(*Desktop) uint32
	LocateDesktop() bool
	UpdateCurrentDesktop()
	SetWMDesktop(n *Node, d *Desktop)
	UpdateWMDesktops()
}

type desktopsState struct {
	e ewmh.EWMH
}

func newDesktopsState(e ewmh.EWMH) *desktopsState {
	return &desktopsState{
		e: e,
	}
}

//void ewmh_update_number_of_desktops(void)
//{
//	xcb_ewmh_set_number_of_desktops(ewmh, default_screen, num_desktops);
//}

//uint32_t ewmh_get_desktop_index(desktop_t *d)
func (d *desktopsState) GetDesktopIndex(dsk *Desktop) uint32 {
	//var i uint32
	//for (monitor_t *m = mon_head; m != NULL; m = m->next)
	//	for (desktop_t *cd = m->desk_head; cd != NULL; cd = cd->next, i++)
	//		if (d == cd)
	//			return i;
	return 0
}

//bool ewmh_locate_desktop(uint32_t i, coordinates_t *loc)
func (d *desktopsState) LocateDesktop() bool {
	//for (monitor_t *m = mon_head; m != NULL; m = m->next)
	//	for (desktop_t *d = m->desk_head; d != NULL; d = d->next, i--)
	//		if (i == 0) {
	//			loc->monitor = m;
	//			loc->desktop = d;
	//			loc->node = NULL;
	//			return true;
	//		}
	return false
}

//void ewmh_update_current_desktop(void)
func (d *desktopsState) UpdateCurrentDesktop() {
	//uint32_t i = ewmh_get_desktop_index(mon->desk);
	//xcb_ewmh_set_current_desktop(ewmh, default_screen, i);
}

//void ewmh_set_wm_desktop(node_t *n, desktop_t *d)
func (d *desktopsState) SetWMDesktop(n *Node, dsk *Desktop) {
	//uint32_t i = ewmh_get_desktop_index(d);
	//xcb_ewmh_set_wm_desktop(ewmh, n->client->window, i);
}

//void ewmh_update_wm_desktops(void)
func (d *desktopsState) UpdateWMDesktops() {
	//for (monitor_t *m = mon_head; m != NULL; m = m->next)
	//	for (desktop_t *d = m->desk_head; d != NULL; d = d->next) {
	//		uint32_t i = ewmh_get_desktop_index(d);
	//		for (node_t *n = first_extrema(d->root); n != NULL; n = next_leaf(n, d->root))
	//			xcb_ewmh_set_wm_desktop(ewmh, n->client->window, i);
	//	}
}

//void ewmh_update_desktop_names(void)
func (d *desktopsState) UpdateDesktopNames() {
	//char names[MAXLEN];
	//unsigned int i, j;
	//uint32_t names_len;
	//i = 0;

	//for (monitor_t *m = mon_head; m != NULL; m = m->next)
	//	for (desktop_t *d = m->desk_head; d != NULL; d = d->next) {
	//		for (j = 0; d->name[j] != '\0' && (i + j) < sizeof(names); j++)
	//			names[i + j] = d->name[j];
	//		i += j;
	//		if (i < sizeof(names))
	//			names[i++] = '\0';
	//	}

	//if (i < 1)
	//	return;

	//names_len = i - 1;
	//xcb_ewmh_set_desktop_names(ewmh, default_screen, names_len, names);
}

type WmState interface {
	Set(string, xproto.Window, xproto.Window)
	UpdateActiveWindow()
	UpdateClientList()
	AddClient(*Client, xproto.Atom) bool
	RemoveClient(*Client, xproto.Atom) bool
}

type wmState struct {
	e ewmh.EWMH
}

func newWmState(e ewmh.EWMH) *wmState {
	return &wmState{
		e: e,
	}
}

//void ewmh_set_supporting(xcb_window_t win)
func (w *wmState) Set(name string, root, win xproto.Window) {
	//pid := os.Getpid()
	//e.e.SupportingWmCheckSet(e.root, win)
	//e.e.SupportingWmCheckSet(win, win)
	//e.e.SetWmName(w, name)
	//e.e.SetWmPid(w, pid)
}

//void ewmh_update_active_window(void)
func (w *wmState) UpdateActiveWindow() {
	//xcb_window_t win = (mon->desk->focus == NULL ? XCB_NONE : mon->desk->focus->client->window);
	//xcb_ewmh_set_active_window(ewmh, default_screen, win);
}

//void ewmh_update_client_list(void)
func (w *wmState) UpdateClientList() {
	//if (num_clients == 0) {
	//	xcb_ewmh_set_client_list(ewmh, default_screen, 0, NULL);
	//	xcb_ewmh_set_client_list_stacking(ewmh, default_screen, 0, NULL);
	//	return;
	//}

	//xcb_window_t wins[num_clients];
	//unsigned int i = 0;

	//for (monitor_t *m = mon_head; m != NULL; m = m->next)
	//	for (desktop_t *d = m->desk_head; d != NULL; d = d->next)
	//		for (node_t *n = first_extrema(d->root); n != NULL; n = next_leaf(n, d->root))
	//			wins[i++] = n->client->window;

	//xcb_ewmh_set_client_list(ewmh, default_screen, num_clients, wins);
	//xcb_ewmh_set_client_list_stacking(ewmh, default_screen, num_clients, wins);
}

const MAXSTATE = int(4)

//bool ewmh_wm_state_add(client_t *c, xcb_atom_t state)
func (w *wmState) AddClient(c *Client, state xproto.Atom) bool {
	if c.numStates <= MAXSTATE {
		//for (int i = 0; i < c->num_states; i++)
		//	if (c->wm_state[i] == state)
		//		return false;
		//c->wm_state[c->num_states] = state;
		//c->num_states++;
		//xcb_ewmh_set_wm_state(ewmh, c->window, c->num_states, c->wm_state);
	}
	return false
}

//bool ewmh_wm_state_remove(client_t *c, xcb_atom_t state)
func (w *wmState) RemoveClient(c *Client, state xproto.Atom) bool {
	//for (int i = 0; i < c->num_states; i++)
	//	if (c->wm_state[i] == state)
	//	{
	//		for (int j = i; j < (c->num_states - 1); j++)
	//			c->wm_state[j] = c->wm_state[j + 1];
	//		c->num_states--;
	//		xcb_ewmh_set_wm_state(ewmh, c->window, c->num_states, c->wm_state);
	//		return true;
	//	}
	return false
}
*/
