package client

import "github.com/BurntSushi/xgb/xproto"

type ClientState int

const (
	FullScreen ClientState = iota
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

func (s *state) Tiled() bool {
	return (!s.floating && !s.fullscreen)
}

func (s *state) Floating() bool {
	return (s.floating && !s.fullscreen)
}

func (s *state) Set(cs ClientState, v bool) {
	switch cs {
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

func (s *state) setFullScreen(v bool) {
	s.fullscreen = v
	if v {
		//ewmh_wm_state_add(c, ewmh->_NET_WM_STATE_FULLSCREEN);
	} else {
		//ewmh_wm_state_remove(c, ewmh->_NET_WM_STATE_FULLSCREEN);
		//stack(n, STACK_ABOVE);
	}
}

func (s *state) setPseudoTiled(v bool) {
	//if n != nil || n.Client.pseudoTiled != v {
	//	n.Client.pseudoTiled = v
	//}
}

func (s *state) setFloating(v bool) {
	//if n != nil || !n.Client.fullscreen || n.Client.floating != v {
	/*
		client_t *c = n->client;

		PRINTF("floating %X: %s\n", c->window, BOOLSTR(value));
		put_status(SBSC_MASK_WINDOW_STATE, "window_state floating %s 0x%X\n", ONOFFSTR(value), c->window);

		n->split_mode = MODE_AUTOMATIC;
		c->floating = n->vacant = value;
		update_vacant_state(n->parent);

		if (value) {
			enable_floating_atom(c->window);
			unrotate_brother(n);
		} else {
			disable_floating_atom(c->window);
			rotate_brother(n);
		}

		stack(n, STACK_ABOVE);
	*/
	//}
}

func (s *state) setLocked(v bool) {
	//void set_locked(monitor_t *m, desktop_t *d, client_t *n, bool value);
	//if n != nil || n.Client.locked != v {
	//client_t *c = n->client;

	//PRINTF("set locked %X: %s\n", c->window, BOOLSTR(value));
	//put_status(SBSC_MASK_WINDOW_STATE, "window_state locked %s 0x%X\n", ONOFFSTR(value), c->window);

	//n.Client.locked = v
	//window_draw_border(n, d->focus == n, m == mon);
	//}
}

func (s *state) setSticky(v bool) {
	//void set_sticky(monitor_t *m, desktop_t *d, client_t *n, bool value);
	/*
		if (n == NULL || n->client->sticky == value)
			return;

		client_t *c = n->client;

		PRINTF("set sticky %X: %s\n", c->window, BOOLSTR(value));
		put_status(SBSC_MASK_WINDOW_STATE, "window_state sticky %s 0x%X\n", ONOFFSTR(value), c->window);

		if (d != m->desk)
			transfer_client(m, d, n, m, m->desk, m->desk->focus);

		c->sticky = value;
		if (value) {
			ewmh_wm_state_add(c, ewmh->_NET_WM_STATE_STICKY);
			m->num_sticky++;
		} else {
			ewmh_wm_state_remove(c, ewmh->_NET_WM_STATE_STICKY);
			m->num_sticky--;
		}

		window_draw_border(n, d->focus == n, m == mon);
	*/
}

func (s *state) setPrivate(v bool) {
	//void set_private(monitor_t *m, desktop_t *d, client_t *n, bool value);
	/*
		if (n == NULL || n->client->private == value)
			return;

		client_t *c = n->client;

		PRINTF("set private %X: %s\n", c->window, BOOLSTR(value));
		put_status(SBSC_MASK_WINDOW_STATE, "window_state private %s 0x%X\n", ONOFFSTR(value), c->window);

		c->private = value;
		update_privacy_level(n, value);
		window_draw_border(n, d->focus == n, m == mon);
	*/
}

func (s *state) setUrgent(v bool) {
	//void set_urgency(monitor_t *m, desktop_t *d, client_t *n, bool value);
	/*
		if (value && mon->desk->focus == n)
			return;
		n->client->urgent = value;
		window_draw_border(n, d->focus == n, m == mon);

		put_status(SBSC_MASK_WINDOW_STATE, "window_state urgent %s 0x%X\n", ONOFFSTR(value), n->client->window);
		put_status(SBSC_MASK_REPORT);
	*/
}

func (s *state) setVacant(v bool) {}
