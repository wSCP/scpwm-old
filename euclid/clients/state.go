package clients

import "github.com/BurntSushi/xgb/xproto"

type ClientState int

const (
	Focus ClientState = iota
	FullScreen
	PseudoTiled
	Floating
	Locked
	Sticky
	Private
	Urgent
	Vacant
)

type State interface {
	Tiled() bool
	Floating() bool
	Focused() bool
	Get(ClientState) bool
	Set(ClientState, bool)
}

type state struct {
	wmState     xproto.Atom
	numStates   int
	focus       bool
	icccmFocus  bool
	vacant      bool
	pseudoTiled bool
	floating    bool
	fullscreen  bool
	locked      bool
	sticky      bool
	urgent      bool
	private     bool
}

func newState() *state {
	return &state{}
}

func (s *state) Tiled() bool {
	return (!s.floating && !s.fullscreen)
}

func (s *state) Floating() bool {
	return (s.floating && !s.fullscreen)
}

func (s *state) Focused() bool {
	return s.focus
}

func (s *state) Get(cs ClientState) bool {
	switch cs {
	case Focus:
		return s.focus
	case FullScreen:
		return s.fullscreen
	case PseudoTiled:
		return s.pseudoTiled
	case Floating:
		return s.floating
	case Locked:
		return s.locked
	case Sticky:
		return s.sticky
	case Private:
		return s.private
	case Urgent:
		return s.urgent
	case Vacant:
		return s.vacant
	}
	return false
}

func (s *state) Set(cs ClientState, v bool) {
	switch cs {
	case Focus:
		s.setFocus(v)
	case FullScreen:
		s.setFullScreen(v)
	case PseudoTiled:
		s.setPseudoTiled(v)
	case Floating:
		s.setFloating(v)
	case Locked:
		s.setLocked(v)
	case Sticky:
		s.setSticky(v)
	case Private:
		s.setPrivate(v)
	case Urgent:
		s.setUrgent(v)
	case Vacant:
		s.setVacant(v)
	}
}

func (s *state) setFocus(v bool) {
	s.focus = v
}

func (s *state) setFullScreen(v bool) {
	s.fullscreen = v
	//if v {
	//ewmh_wm_state_add(c, ewmh->_NET_WM_STATE_FULLSCREEN);
	//} else {
	//ewmh_wm_state_remove(c, ewmh->_NET_WM_STATE_FULLSCREEN);
	//stack(n, STACK_ABOVE);
	//}
}

func (s *state) setPseudoTiled(v bool) {
	s.pseudoTiled = v
}

func (s *state) setFloating(v bool) {
	s.floating, s.vacant = v, v
	//update_vacant_state(n->parent);
	//if v {
	//	enable_floating_atom(c->window);
	//	unrotate_brother(n);
	//} else {
	//	disable_floating_atom(c->window);
	//	rotate_brother(n);
	//}
	//stack(n, STACK_ABOVE);
}

func (s *state) setLocked(v bool) {
	s.locked = v
	//window_draw_border(n, d->focus == n, m == mon);
}

func (s *state) setSticky(v bool) {
	s.sticky = v
	//if (value) {
	//		ewmh_wm_state_add(c, ewmh->_NET_WM_STATE_STICKY);
	//		m->num_sticky++;
	//	} else {
	//		ewmh_wm_state_remove(c, ewmh->_NET_WM_STATE_STICKY);
	//		m->num_sticky--;
	//	}
	//	window_draw_border(n, d->focus == n, m == mon);
}

func (s *state) setPrivate(v bool) {
	s.private = v
	//update_privacy_level(n, value);
	//window_draw_border(n, d->focus == n, m == mon);
}

func (s *state) setUrgent(v bool) {
	s.urgent = v
	//window_draw_border(n, d->focus == n, m == mon);
}

func (s *state) setVacant(v bool) {
	s.vacant = v
}
