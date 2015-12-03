package manager

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/thrisp/scpwm/euclid/clients"
	"github.com/thrisp/scpwm/euclid/layout"
	"github.com/thrisp/scpwm/euclid/monitors"
	"github.com/thrisp/scpwm/euclid/rules"
	"github.com/thrisp/scpwm/euclid/window"
)

func (m *Manager) SetEventFns() {
	m.SetEventFn("MapRequest", m.MapRequest)
	m.SetEventFn("DestroyNotify", m.DestroyNotify)
	m.SetEventFn("UnmapNotify", m.UnmapNotify)
	m.SetEventFn("ClientMessage", m.ClientMessage)
	m.SetEventFn("ConfigureRequest", m.ConfigureRequest)
	m.SetEventFn("PropertyNotify", m.PropertyNotify)
	m.SetEventFn("EnterNotify", m.EnterNotify)
	m.SetEventFn("MotionNotify", m.MotionNotify)
	m.SetEventFn("FocusInEvent", m.FocusIn)
	m.SetEventFn("ScreenChange", m.ScreenChange)
}

var EventError = Xrror("Event : %+v is not recognized despite being passed to Manager function %s.").Out

func (m *Manager) MapRequest(evt xgb.Event) error {
	if mr, ok := evt.(xproto.MapRequestEvent); ok {
		m.schedule(mr.Window)
	}
	return EventError(evt, "MapRequest")
}

func (m *Manager) schedule(win xproto.Window) {
	var overrideRedirect bool
	wa, _ := xproto.GetWindowAttributes(m.Conn(), win).Reply()
	if wa != nil {
		overrideRedirect = wa.OverrideRedirect
	}

	if !overrideRedirect && !m.exists(win) {
		if ci, err := m.WmClassGet(win); err == nil {
			csq := m.Ruler.Applicable(ci.Class, ci.Instance)
			m.manage(win, csq)
		}
	}
}

func (m *Manager) manage(win xproto.Window, csq *rules.Consequence) {
	if !csq.Manage {
		window.SetVisible(true, m.Conn(), win, m.Root())
		csq = nil
		return
	}

	loc := m.Current()

	if csq.WithClient != nil {
		loc, _ = m.Locate(loc, csq.WithClient)
	}
	if csq.OnDesktop != nil {
		loc, _ = m.Locate(loc, csq.OnDesktop)
	}
	if csq.OnMonitor != nil {
		loc, _ = m.Locate(loc, csq.OnMonitor)
	}

	client := clients.NewClient(m.Conn(), win, m.Root(), csq)

	mr := loc.m.Rectangle()
	mm := monitors.FromClient(m.Branch, client)
	mmr := mm.Rectangle()
	client.Embrace(mmr)
	client.Translate(mmr, mr)

	focus := loc.c

	if focus != nil {
		focus.SetShift(csq)
	}

	clients.Add(loc.cs, client, focus, csq)

	layout.Arrange(loc.m, loc.d)

	// adjust focuses

	// ewmh_set_wm_desktop(n, d);
	// ewmh_update_client_list();
}

func (m *Manager) ConfigureRequest(evt xgb.Event) error {
	if cr, ok := evt.(xproto.ConfigureRequestEvent); ok {
		loc := m.Current()
		client, exists := m.locate(cr.Window)
		var w, h uint16
		if exists && client.Tiled() {
			if (cr.ValueMask & xproto.ConfigWindowX) != 0 {
				client.X(cr.X)
			}
			if (cr.ValueMask & xproto.ConfigWindowY) != 0 {
				client.Y(cr.Y)
			}
			if (cr.ValueMask & xproto.ConfigWindowHeight) != 0 {
				w = cr.Width
			}
			if (cr.ValueMask & xproto.ConfigWindowWidth) != 0 {
				h = cr.Height
			}
			if w != 0 {
				client.Width(clients.RestrainWidth(client, w))
			}
			if h != 0 {
				client.Height(clients.RestrainHeight(client, h))
			}

			var evt xproto.ConfigureNotifyEvent
			var rect xproto.Rectangle
			win := client.XWindow()
			bw := client.BorderWidth()

			if client.Fullscreen() {
				rect = loc.m.Rectangle()
			} else {
				rect = client.TRectangle()
			}

			evt.Event = win
			evt.Window = win
			evt.AboveSibling = xproto.WindowNone
			evt.X = rect.X
			evt.Y = rect.Y
			evt.Width = rect.Width
			evt.Height = rect.Height
			evt.BorderWidth = uint16(bw)
			evt.OverrideRedirect = false

			xproto.SendEvent(m.Conn(), false, win, xproto.EventMaskStructureNotify, evt.String())

			if client.Pseudotiled() {
				layout.Arrange(loc.m, loc.d)
			}
		} else {
			var mask uint16
			var values []uint32
			var i int
			if (cr.ValueMask & xproto.ConfigWindowX) != 0 {
				mask |= xproto.ConfigWindowX
				values[i] = uint32(cr.X)
				if exists {
					client.X(cr.X)
				}
			}
			if (cr.ValueMask & xproto.ConfigWindowY) != 0 {
				mask |= xproto.ConfigWindowY
				i++
				values[i] = uint32(cr.Y)
				if exists {
					client.Y(cr.Y)
				}
			}
			if (cr.ValueMask & xproto.ConfigWindowHeight) != 0 {
				mask |= xproto.ConfigWindowHeight
				w = cr.Width
				if exists {
					client.Width(clients.RestrainWidth(client, w))
				}
				i++
				values[i] = uint32(cr.Height)
			}
			if (cr.ValueMask & xproto.ConfigWindowWidth) != 0 {
				mask |= xproto.ConfigWindowWidth
				h = cr.Height
				if exists {
					client.Height(clients.RestrainHeight(client, h))
				}
				i++
				values[i] = uint32(cr.Width)
			}
			if !exists && (cr.ValueMask&xproto.ConfigWindowBorderWidth) != 0 {
				mask |= xproto.ConfigWindowBorderWidth
				i++
				values[i] = uint32(cr.BorderWidth)
			}
			if (cr.ValueMask & xproto.ConfigWindowSibling) != 0 {
				mask |= xproto.ConfigWindowSibling
				i++
				values[i] = uint32(cr.Sibling)
			}
			if (cr.ValueMask & xproto.ConfigWindowStackMode) != 0 {
				mask |= xproto.ConfigWindowStackMode
				i++
				values[i] = uint32(cr.StackMode)
			}

			xproto.ConfigureWindow(m.Conn(), cr.Window, mask, values)
		}
		if exists {
			mt := monitors.FromClient(m.Branch, client)
			client.Translate(mt.Rectangle(), loc.m.Rectangle())
		}
		return nil
	}
	return EventError(evt, "ConfigureRequest")
}

func (m *Manager) unmanage(win xproto.Window) error {
	if loc, exists := m.LocationWindow(win); exists {
		err := m.RemoveClient(loc)
		// adjust pointer
		layout.Arrange(loc.m, loc.d)
		return err
	}
	return nil
}

var RemoveClientError = Xrror("Unable to remove client, branch or client at Location does not exist.")

func (m *Manager) RemoveClient(loc Location) error {
	if loc.cs != nil && loc.c != nil {
		return clients.Remove(loc.cs, loc.c)
	}
	return RemoveClientError
}

func (m *Manager) DestroyNotify(evt xgb.Event) error {
	if dn, ok := evt.(xproto.DestroyNotifyEvent); ok {
		return m.unmanage(dn.Window)
	}
	return EventError(evt, "DestroyNotify")
}

func (m *Manager) UnmapNotify(evt xgb.Event) error {
	if un, ok := evt.(xproto.UnmapNotifyEvent); ok {
		return m.unmanage(un.Window)
	}
	return EventError(evt, "UnmapNotify")
}

func (m *Manager) PropertyNotify(evt xgb.Event) error {
	if pn, ok := evt.(xproto.PropertyNotifyEvent); ok {
		if pn.Atom == xproto.AtomWmHints || pn.Atom == xproto.AtomWmNormalHints {
			if m.exists(pn.Window) {
				//loc := m.LocationWindow(pn.Window)
				switch pn.Atom {
				case xproto.AtomWmHints:
					if hints, err := m.WmHintsGet(pn.Window); err == nil {
						if (hints.Flags & icccm.HintUrgency) != 0 {
							//set_urgency(loc.monitor, loc.desktop, loc.node, xcb_icccm_wm_hints_get_urgency(&hints));
						}
					}
				case xproto.AtomWmNormalHints:
					if hints, err := m.WmNormalHintsGet(pn.Window); err == nil {
						if (hints.Flags & (icccm.SizeHintPMinSize | icccm.SizeHintPMaxSize)) != 0 {
							/*
								c->min_width = size_hints.min_width;
								c->max_width = size_hints.max_width;
								c->min_height = size_hints.min_height;
								c->max_height = size_hints.max_height;
								int w = c->floating_rectangle.width;
								int h = c->floating_rectangle.height;
								restrain_floating_size(c, &w, &h);
								c->floating_rectangle.width = w;
								c->floating_rectangle.height = h;
								arrange(loc.monitor, loc.desktop);
							*/
						}
					}
				}
				return nil
			}
		}
		return nil
	}
	return EventError(evt, "PropertyNotify")
}

func (m *Manager) ClientMessage(evt xgb.Event) error {
	if cm, ok := evt.(xproto.ClientMessageEvent); ok {
		var a xproto.Atom
		var err error
		if a, err = m.Atom("_NET_CURRENT_DESKTOP"); err == nil {
			if cm.Type == a {
				if _, exists := m.LocationDesktopIndex(int(cm.Data.Data32[0])); exists {
					//focus_node(loc.monitor, loc.desktop, loc.desktop->focus);
					//return nil
				}
			}
		}
		if loc, exists := m.LocationWindow(cm.Window); exists {
			if a, err = m.Atom("_NET_WM_STATE"); cm.Type == a {
				m.HandleState(loc, cm.Data.Data32[1], cm.Data.Data32[0])
				m.HandleState(loc, cm.Data.Data32[2], cm.Data.Data32[0])
			}
			if a, err = m.Atom("_NET_ACTIVE_WINDOW"); cm.Type == a {
				//if m.Bool("IgnoreEwhmFocus") && e->data.data32[0] == XCB_EWMH_CLIENT_SOURCE_TYPE_NORMAL {
				//	return nil
				//}
			}
			if a, err = m.Atom("_NET_WM_DESKTOP"); cm.Type == a {
			}
			if a, err = m.Atom("_NET_CLOSE_WINDOW"); cm.Type == a {
			}
		}
		return err
	}
	return EventError(evt, "ClientMessage")
}

const (
	StateRemove = iota
	StateAdd
	StateToggle
)

func (m *Manager) HandleState(loc Location, state, action uint32) {
	var a xproto.Atom
	st := xproto.Atom(state)
	if a, _ = m.Atom("_NET_WM_STATE_FULLSCREEN"); st == a {
		if action == StateAdd {
			//set_state(m, d, n, STATE_FULLSCREEN)
		} else if action == StateRemove {
			//set_state(m, d, n, n->client->last_state)
		} else if action == StateToggle {
			//set_state(m, d, n, IS_FULLSCREEN(n->client) ? n->client->last_state : STATE_FULLSCREEN)
		}
		layout.Arrange(loc.m, loc.d)
	}
	if a, _ = m.Atom("_NET_WM_STATE_BELOW"); st == a {
		if action == StateAdd {
			//set_layer(m, d, n, LAYER_BELOW)
		} else if action == StateRemove {
			//set_layer(m, d, n, LAYER_NORMAL)
		} else if action == StateToggle {
			//set_layer(m, d, n, n->client->layer == LAYER_BELOW ? n->client->last_layer : LAYER_BELOW);
		}
	}
	if a, _ = m.Atom("_NET_WM_STATE_ABOVE"); st == a {
		if action == StateAdd {
			//set_layer(m, d, n, LAYER_ABOVE)
		} else if action == StateRemove {
			//set_layer(m, d, n, n->client->last_layer)
		} else if action == StateToggle {
			//set_layer(m, d, n, n->client->layer == LAYER_ABOVE ? n->client->last_layer : LAYER_ABOVE)
		}
	}
	if a, _ = m.Atom("_NET_WM_STATE_STICKY"); st == a {
		if action == StateAdd {
			//set_sticky(m, d, n, true)
		} else if action == StateRemove {
			//set_sticky(m, d, n, false)
		} else if action == StateToggle {
			//set_sticky(m, d, n, !n->client->sticky)
		}
	}
	if a, _ = m.Atom("_NET_WM_STATE_DEMANDS_ATTENTION"); st == a {
		if action == StateAdd {
			//set_urgency(m, d, n, true)
		} else if action == StateRemove {
			//set_urgency(m, d, n, false)
		} else if action == StateToggle {
			//set_urgency(m, d, n, !n->client->urgent)
		}
	}
}

func (m *Manager) FocusIn(evt xgb.Event) error {
	if fi, ok := evt.(xproto.FocusInEvent); ok {
		if fi.Mode != xproto.NotifyModeGrab || fi.Mode != xproto.NotifyModeUngrab {
			loc := m.Current()
			if fi.Detail == xproto.NotifyDetailAncestor ||
				fi.Detail == xproto.NotifyDetailInferior ||
				fi.Detail == xproto.NotifyDetailNonlinearVirtual ||
				fi.Detail == xproto.NotifyDetailNonlinear &&
					(loc.c == nil || loc.c.XWindow() != fi.Event) {
				//update input focus
			}
		}
		return nil
	}
	return EventError(evt, "FocusIn")
}

func (m *Manager) EnterNotify(evt xgb.Event) error {
	if _, ok := evt.(xproto.EnterNotifyEvent); ok {
		return nil
	}
	return EventError(evt, "EnterNotify")
}

func (m *Manager) MotionNotify(evt xgb.Event) error {
	if _, ok := evt.(xproto.MotionNotifyEvent); ok {
		return nil
	}
	return EventError(evt, "MotionNotify")
}

func (m *Manager) ScreenChange(evt xgb.Event) error {
	return monitors.Update(m.Branch, m.Handler, m.Settings)
}
