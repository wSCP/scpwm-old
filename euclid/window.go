package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Windower interface {
	NewWindow() *Window
	Schedule(*Euclid, xproto.Window) bool
	Manage(*Euclid, *Window, *Rule) bool
	Unmanage(*Window)
	AdoptOrphans()
}

type windower struct {
	conn *xgb.Conn
	root xproto.Window
}

func NewWindower(c *xgb.Conn, r xproto.Window) Windower {
	return &windower{
		conn: c,
		root: r,
	}
}

func (w *windower) NewWindow() *Window {
	win, err := xproto.NewWindowId(w.conn)
	if err != nil {
		return nil
	}
	return &Window{w.conn, win, w.root, 0, 0}
}

func (w *windower) mkWindow(win xproto.Window) *Window {
	return &Window{w.conn, win, w.root, 0, 0}
}

func (w *windower) Schedule(e *Euclid, win xproto.Window) bool {
	/*var loc coordinate
	var overrideRedirect bool

	wa, _ := xproto.GetWindowAttributes(w.conn, win).Reply()
	if wa != nil {
		overrideRedirect = wa.OverrideRedirect
	}

	nw := w.mkWindow(win)

	if !overrideRedirect || !LocateWindow(e, nw, &loc) {
		pr := e.pending
		for pr != nil {
			if pr.window == nw.Window {
				return false
			}
			pr = pr.Pop(Next)
		}

		csq := NewConsequence()
		e.applyRules(nw, csq)
		if !e.scheduleRules(nw, csq) {
			return w.Manage(e, nw, csq)
		}
	}*/
	return false
}

func (w *windower) Manage(e *Euclid, win *Window, r *Rule) bool {
	//void manage_Window(xcb_Window_t win, rule_consequence_t *csq, int fd);
	/*
		monitor_t *m = mon;
		desktop_t *d = mon->desk;
		client_t *f = mon->desk->focus;

		parse_rule_consequence(fd, csq);

		if (!csq->manage) {
			disable_floating_atom(win);
			Window_show(win);
			return;
		}

		PRINTF("manage %X\n", win);

		if (csq->client_desc[0] != '\0') {
			coordinates_t ref = {m, d, f};
			coordinates_t trg = {NULL, NULL, NULL};
			if (client_from_desc(csq->client_desc, &ref, &trg)) {
				m = trg.monitor;
				d = trg.desktop;
				f = trg.client;
			}
		} else if (csq->desktop_desc[0] != '\0') {
			coordinates_t ref = {m, d, NULL};
			coordinates_t trg = {NULL, NULL, NULL};
			if (desktop_from_desc(csq->desktop_desc, &ref, &trg)) {
				m = trg.monitor;
				d = trg.desktop;
				f = trg.desktop->focus;
			}
		} else if (csq->monitor_desc[0] != '\0') {
			coordinates_t ref = {m, NULL, NULL};
			coordinates_t trg = {NULL, NULL, NULL};
			if (monitor_from_desc(csq->monitor_desc, &ref, &trg)) {
				m = trg.monitor;
				d = trg.monitor->desk;
				f = trg.monitor->desk->focus;
			}
		}

		if (csq->sticky) {
			m = mon;
			d = mon->desk;
			f = mon->desk->focus;
		}

		if (csq->split_dir[0] != '\0' && f != NULL) {
			direction_t dir;
			if (parse_direction(csq->split_dir, &dir)) {
				f->split_mode = MODE_MANUAL;
				f->split_dir = dir;
			}
		}

		if (csq->split_ratio != 0 && f != NULL) {
			f->split_ratio = csq->split_ratio;
		}

		client_t *c = make_client(win, csq->border ? d->border_width : 0);
		update_floating_rectangle(c);
		if (c->floating_rectangle.x == 0 && c->floating_rectangle.y == 0)
			csq->center = true;
		c->min_width = csq->min_width;
		c->max_width = csq->max_width;
		c->min_height = csq->min_height;
		c->max_height = csq->max_height;
		monitor_t *mm = monitor_from_client(c);
		embrace_client(mm, c);
		translate_client(mm, m, c);
		if (csq->center)
			Window_center(m, c);

		snprintf(c->class_name, sizeof(c->class_name), "%s", csq->class_name);
		snprintf(c->instance_name, sizeof(c->instance_name), "%s", csq->instance_name);

		csq->floating = csq->floating || d->floating;

		client_t *n = make_client();
		n->client = c;

		put_status(SBSC_MASK_WINDOW_MANAGE, "Window_manage %s %s 0x%X 0x%X\n", m->name, d->name, f!=NULL?f->client->Window:0, win);
		insert_client(m, d, n, f);

		disable_floating_atom(c->Window);
		set_pseudo_tiled(n, csq->pseudo_tiled);
		set_floating(n, csq->floating);
		set_locked(m, d, n, csq->locked);
		set_sticky(m, d, n, csq->sticky);
		set_private(m, d, n, csq->private);

		if (d->focus != NULL && d->focus->client->fullscreen)
			set_fullscreen(d->focus, false);

		set_fullscreen(n, csq->fullscreen);

		arrange(m, d);

		bool give_focus = (csq->focus && (d == mon->desk || csq->follow));

		if (give_focus)
			focus_client(m, d, n);
		else if (csq->focus)
			pseudo_focus(m, d, n);
		else
			stack(n, STACK_ABOVE);

		uint32_t values[] = {CLIENT_EVENT_MASK | (focus_follows_pointer ? XCB_EVENT_MASK_ENTER_WINDOW : 0)};
		xcb_change_Window_attributes(dpy, c->Window, XCB_CW_EVENT_MASK, values);

		if (visible) {
			if (d == m->desk)
				Window_show(n->client->Window);
			else
				Window_hide(n->client->Window);
		}

		//the same function is already called in `focus_client` but has no effects on unmapped Windows
		if (give_focus)
			xcb_set_input_focus(dpy, XCB_INPUT_FOCUS_POINTER_ROOT, win, XCB_CURRENT_TIME);

		num_clients++;
		ewmh_set_wm_desktop(n, d);
		ewmh_update_client_list();
	*/
	return false
}

func (w *windower) Unmanage(win *Window) {
	//void unmanage_Window(xcb_Window_t win);
}

func (w *windower) AdoptOrphans() {}

type Window struct {
	*xgb.Conn
	xproto.Window
	root      xproto.Window
	wmState   xproto.Atom
	numStates int
}

var NoWindow Window = Window{nil, xproto.WindowNone, xproto.WindowNone, 0, 0}

func (w *Window) Close() {
	//send_client_message(w.Window, ewmh->WM_PROTOCOLS, WM_DELETE_WINDOW);
}

func (w *Window) Kill() {
	xproto.KillClientChecked(w.Conn, uint32(w.Window))
}

func (w *Window) DrawBorder(c *Client, focusedWindow, focusedMonitor bool) {
	//void Window_draw_border(client_t *n, bool focused_Window, bool focused_monitor);

	//if n != nil || n.Client.borderWidth > 0 {
	/*
		xcb_window_t win = n->client->window;
		uint32_t border_color_pxl = get_border_color(n->client, focused_window, focused_monitor);

		if (n->split_mode == MODE_AUTOMATIC) {
			xcb_change_window_attributes(dpy, win, XCB_CW_BORDER_PIXEL, &border_color_pxl);
		} else {
			unsigned int border_width = n->client->border_width;
			uint32_t presel_border_color_pxl;
			get_color(presel_border_color, win, &presel_border_color_pxl);

			xcb_rectangle_t actual_rectangle = get_rectangle(n->client);

			uint16_t width = actual_rectangle.width;
			uint16_t height = actual_rectangle.height;

			uint16_t full_width = width + 2 * border_width;
			uint16_t full_height = height + 2 * border_width;

			xcb_rectangle_t border_rectangles[] =
			{
				{ width, 0, 2 * border_width, height + 2 * border_width },
				{ 0, height, width + 2 * border_width, 2 * border_width }
			};

			xcb_rectangle_t *presel_rectangles;

			uint8_t win_depth = root_depth;
			xcb_get_geometry_reply_t *geo = xcb_get_geometry_reply(dpy, xcb_get_geometry(dpy, win), NULL);
			if (geo != NULL)
				win_depth = geo->depth;
			free(geo);

			xcb_pixmap_t pixmap = xcb_generate_id(dpy);
			xcb_create_pixmap(dpy, win_depth, pixmap, win, full_width, full_height);

			xcb_gcontext_t gc = xcb_generate_id(dpy);
			xcb_create_gc(dpy, gc, pixmap, 0, NULL);

			xcb_change_gc(dpy, gc, XCB_GC_FOREGROUND, &border_color_pxl);
			xcb_poly_fill_rectangle(dpy, pixmap, gc, LENGTH(border_rectangles), border_rectangles);

			uint16_t fence = (int16_t) (n->split_ratio * ((n->split_dir == DIR_UP || n->split_dir == DIR_DOWN) ? height : width));
			presel_rectangles = malloc(2 * sizeof(xcb_rectangle_t));
			switch (n->split_dir) {
				case DIR_UP:
					presel_rectangles[0] = (xcb_rectangle_t) {width, 0, 2 * border_width, fence};
					presel_rectangles[1] = (xcb_rectangle_t) {0, height + border_width, full_width, border_width};
					break;
				case DIR_DOWN:
					presel_rectangles[0] = (xcb_rectangle_t) {width, fence + 1, 2 * border_width, height + border_width - (fence + 1)};
					presel_rectangles[1] = (xcb_rectangle_t) {0, height, full_width, border_width};
					break;
				case DIR_LEFT:
					presel_rectangles[0] = (xcb_rectangle_t) {0, height, fence, 2 * border_width};
					presel_rectangles[1] = (xcb_rectangle_t) {width + border_width, 0, border_width, full_height};
					break;
				case DIR_RIGHT:
					presel_rectangles[0] = (xcb_rectangle_t) {fence + 1, height, width + border_width - (fence + 1), 2 * border_width};
					presel_rectangles[1] = (xcb_rectangle_t) {width, 0, border_width, full_height};
					break;
			}
			xcb_change_gc(dpy, gc, XCB_GC_FOREGROUND, &presel_border_color_pxl);
			xcb_poly_fill_rectangle(dpy, pixmap, gc, 2, presel_rectangles);
			xcb_change_window_attributes(dpy, win, XCB_CW_BORDER_PIXMAP, &pixmap);
			free(presel_rectangles);
			xcb_free_gc(dpy, gc);
			xcb_free_pixmap(dpy, pixmap);
	*/
	//}
}

func (w *Window) BorderWidth(bw uint32) {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowBorderWidth, []uint32{bw})
}

func (w *Window) Move(x, y int16) {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowX, []uint32{uint32(x)})
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowY, []uint32{uint32(y)})
}

func (w *Window) Resize(hght, wdth uint16) {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowHeight, []uint32{uint32(hght)})
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowWidth, []uint32{uint32(wdth)})
}

func (w *Window) MoveResize(x, y int16, hght, wdth uint16) {
	w.Move(x, y)
	w.Resize(hght, wdth)
}

func (w *Window) Raise() {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeAbove})
}

func (w *Window) Lower() {
	xproto.ConfigureWindowChecked(w.Conn, w.Window, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeBelow})
}

func (w *Window) Center(m *Monitor, c *Client) {
	/*r := c.fRectangle
	a := m.rectangle

	if r.Width >= a.Width {
		r.X = a.X
	} else {
		r.X = a.X + (int16(a.Width)-int16(r.Width))/2
	}

	if r.Height >= a.Height {
		r.Y = a.Y
	} else {
		r.Y = a.Y + (int16(a.Height)-int16(r.Height))/2
	}

	r.X -= int16(c.borderWidth)
	r.Y -= int16(c.borderWidth)*/
}

func (w *Window) Stack(o *Window, mode uint32) {
	xproto.ConfigureWindowChecked(
		w.Conn,
		w.Window,
		(xproto.ConfigWindowSibling | xproto.ConfigWindowStackMode),
		[]uint32{uint32(o.Window), mode},
	)
}

func (w *Window) Above(o *Window) {
	w.Stack(o, xproto.StackModeAbove)
}

func (w *Window) Below(o *Window) {
	w.Stack(o, xproto.StackModeBelow)
}

var (
	windowOff = []uint32{RootEventMask, xproto.EventMaskSubstructureNotify} //uint32_t values_off[] = {ROOT_EVENT_MASK & ~XCB_EVENT_MASK_SUBSTRUCTURE_NOTIFY};
	windowOn  = []uint32{RootEventMask}
)

func (w *Window) setVisibility(v bool) {
	xproto.ChangeWindowAttributesChecked(w.Conn, w.root, xproto.CwEventMask, windowOff)
	if v {
		xproto.MapWindow(w.Conn, w.Window)
	} else {
		xproto.UnmapWindow(w.Conn, w.Window)
	}
	xproto.ChangeWindowAttributesChecked(w.Conn, w.root, xproto.CwEventMask, windowOn)
}

func (w *Window) Hide() {
	w.setVisibility(false)
}

func (w *Window) Show() {
	w.setVisibility(true)
}

type WindowSet int

const (
	FullScreen WindowSet = iota
	PseudoTiled
	Floating
	Locked
	Sticky
	Private
	Urgent
)

func (w *Window) Set(ws WindowSet, c *Client, v bool) {
	switch ws {
	case FullScreen:
		w.setFullScreen(c, v)
	case PseudoTiled:
		w.setPseudoTiled(c, v)
	case Floating:
		w.setFloating(c, v)
	case Locked:
		w.setLocked(c, v)
	case Sticky:
		w.setSticky(c, v)
	case Private:
		w.setPrivate(c, v)
	case Urgent:
		w.setUrgent(c, v)
	}
}

func (w *Window) setFullScreen(c *Client, v bool) {
	//if n != nil || n.Client.fullscreen != v {
	//n.Client.fullscreen = v
	//if v {
	//ewmh_wm_state_add(c, ewmh->_NET_WM_STATE_FULLSCREEN);
	//} else {
	//ewmh_wm_state_remove(c, ewmh->_NET_WM_STATE_FULLSCREEN);
	//stack(n, STACK_ABOVE);
	//}
	//}
}

func (w *Window) setPseudoTiled(c *Client, v bool) {
	//if n != nil || n.Client.pseudoTiled != v {
	//	n.Client.pseudoTiled = v
	//}
}

func (w *Window) setFloating(c *Client, v bool) {
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

func (w *Window) setLocked(c *Client, v bool) {
	//void set_locked(monitor_t *m, desktop_t *d, client_t *n, bool value);
	//if n != nil || n.Client.locked != v {
	//client_t *c = n->client;

	//PRINTF("set locked %X: %s\n", c->window, BOOLSTR(value));
	//put_status(SBSC_MASK_WINDOW_STATE, "window_state locked %s 0x%X\n", ONOFFSTR(value), c->window);

	//n.Client.locked = v
	//window_draw_border(n, d->focus == n, m == mon);
	//}
}

func (w *Window) setSticky(c *Client, v bool) {
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

func (w *Window) setPrivate(c *Client, v bool) {
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

func (w *Window) setUrgent(c *Client, v bool) {
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

func (w *Window) getColor(color string, pxl uint32) bool {
	//reply := xproto.GetWindowAttributes(w, w, nil)
	//if reply != nil {
	//	cm := reply.Colormap
	//	if strings.Index(color, "#") == 0 {
	/*
		uint red, green, blue
		if (sscanf(color + 1, "%02x%02x%02x", &red, &green, &blue) == 3) {
			//2**16 - 1 == 0xffff and 0x101 * 0xij == 0xijij
			red *= 0x101;
			green *= 0x101;
			blue *= 0x101;
			xcb_alloc_color_reply_t *reply = xcb_alloc_color_reply(dpy, xcb_alloc_color(dpy, map, red, green, blue), NULL);
			if (reply != NULL) {
				*pxl = reply->pixel;
				free(reply);
				return true;
			}
		}
	*/
	//	} else {
	/*
		xcb_alloc_named_color_reply_t *reply = xcb_alloc_named_color_reply(dpy, xcb_alloc_named_color(dpy, map, strlen(col), col), NULL);
		if (reply != NULL) {
			*pxl = reply->pixel;
			free(reply);
			return true;
		}
	*/
	//	}
	//}
	//pxl = 0
	return false
}

func locateWindow(e *Euclid, sel selector, loc coordinate) bool {
	/*m := e.Pop(Head)
	for m != nil {
		d := m.Desktops.Pop(Head)
		for d != nil {
			n := d.root.rightExtrema()
			for n != nil {
				if n.Client.Window == w {
					loc.m = m
					loc.d = d
					loc.n = n
					return true
				}
				n = nextLeaf(n, d.root)
			}
			d = d.Pop(Next)
		}
		m = m.Pop(Next)
	}*/
	return false
}
