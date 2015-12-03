package layout

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/clients"
	"github.com/thrisp/scpwm/euclid/desktops"
	"github.com/thrisp/scpwm/euclid/monitors"
)

type Layout int

func (l Layout) String() string {
	switch l {
	case Tiled:
		return "tiled"
	case Monocle:
		return "monocle"
	}
	return ""
}

const (
	Tiled Layout = iota
	Monocle
)

var StringLayout map[string]Layout = map[string]Layout{
	"tiled":   Tiled,
	"monocle": Monocle,
}

func Change(m monitors.Monitor, d desktops.Desktop, l Layout) {
	d.Set("layout", l.String())
	Arrange(m, d)
}

func Arrange(m monitors.Monitor, d desktops.Desktop) {
	var original, current Layout
	if v, ok := StringLayout[d.Layout()]; ok {
		original = v
	} else {
		original = Tiled
	}
	if d.Bool("LeafMonocle") && clients.TiledCount(d.Clients()) == 1 {
		current = Monocle
	}
	rect := m.Rectangle()
	var gap int
	if d.Bool("GaplessMonocle") && current == Monocle {
		gap = 0
	} else {
		gap = d.Int("WindowGap")
	}

	mp := m.Pad()
	dp := d.Pad()

	rect.X += int16(mp[3] + dp[3] + gap)
	rect.Y += int16(dp[0] + dp[0] + gap)
	rect.Width -= uint16(mp[3] + dp[3] + dp[1] + mp[1] + gap)
	rect.Height -= uint16(mp[0] + dp[0] + dp[2] + mp[2] + gap)

	focus := clients.Focused(d.Clients())
	Apply(m, d, focus, rect, rect)

	d.Set("layout", original.String())
}

func Apply(m monitors.Monitor, d desktops.Desktop, c clients.Client, r, rr xproto.Rectangle) {

}

/*
	dl := StringLayout[d.Layout()]
	if c != nil {
		if clients.IsTail(d.Clients(), c) {
			var bw uint
			if (d.Bool("BorderlessMonocle") &&
				c.Tiled() &&
				dl == Monocle) ||
				c.Fullscreen() {
				bw = 0
			} else {
				bw = c.BorderWidth()
			}

			var cr xproto.Rectangle
			if c.Tiled() || c.Pseudotiled() {
				var wg int
				if d.Bool("GaplessMonocle") && dl == Monocle {
					wg = 0
				} else {
					wg = d.Int("WindowGap")
				}
				if c.Tiled() {
					cr = r
					var bleed uint16
					bleed = uint16(wg + int(2*bw))
					if bleed < cr.Width {
						cr.Width = cr.Width - bleed
					} else {
						cr.Width = 1
					}
					if bleed < cr.Height {
						cr.Height = cr.Height - bleed
					} else {
						cr.Height = 1
					}
				} else {
					cr = c.FRectangle()
					if d.Bool("CenterPseudoTiled") {
						cr.X = r.X - int16(bw) + int16(r.Width-uint16(wg)-cr.Width)/2
						cr.Y = r.Y - int16(bw) + int16(r.Height-uint16(wg)-cr.Height)/2
					} else {
						cr.X = r.X
						cr.Y = r.Y
					}
				}
				c.SetTiledRectangle(cr)
			} else if c.Floating() {
				cr = c.FRectangle()
			} else {
				cr = m.Rectangle()
			}

			//window_move_resize(n->client->window, r.x, r.y, r.width, r.height);
			//window_border_width(n->client->window, bw);
			//window_draw_border(n, d->focus == n, m == mon);

			//if (frozen_pointer->action == ACTION_NONE) {
			//	put_status(SBSC_MASK_WINDOW_GEOMETRY, "window_geometry %s %s 0x%X %ux%u+%i+%i\n", m->name, d->name, n->client->window, r.width, r.height, r.x, r.y);
			//}

			//if (pointer_follows_focus && mon->desk->focus == n && frozen_pointer->action == ACTION_NONE) {
			//	center_pointer(r);
			//}
		} else {
			var cr xproto.Rectangle
			if dl == Monocle {
				cr = r
			} else {
				var fence uint
				if c.Orientation.String() == "vertical" {
					fence = r.Width * c.Ratio()
				} else {

				}
			}
			/*
				xcb_rectangle_t first_rect;
				xcb_rectangle_t second_rect;

				if (d->layout == LAYOUT_MONOCLE || n->first_child->vacant || n->second_child->vacant) {
					first_rect = second_rect = rect;
				} else {
					unsigned int fence;
					if (n->split_type == TYPE_VERTICAL) {
						fence = rect.width * n->split_ratio;
						first_rect = (xcb_rectangle_t) {rect.x, rect.y, fence, rect.height};
						second_rect = (xcb_rectangle_t) {rect.x + fence, rect.y, rect.width - fence, rect.height};
					} else {
						fence = rect.height * n->split_ratio;
						first_rect = (xcb_rectangle_t) {rect.x, rect.y, rect.width, fence};
						second_rect = (xcb_rectangle_t) {rect.x, rect.y + fence, rect.width, rect.height - fence};
					}
				}

				apply_layout(m, d, n->first_child, first_rect, root_rect);
				apply_layout(m, d, n->second_child, second_rect, root_rect);
		}
	}
*/
