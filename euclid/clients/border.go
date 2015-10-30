package clients

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Bordr interface {
	BorderWidth() uint
	SetBorderWidth(uint)
	Draw(Client, bool, bool)
	Color(bool, bool) uint32
}

type bordr struct {
	borderWidth uint
}

func newBordr() *bordr {
	return &bordr{}
}

func (b *bordr) BorderWidth() uint {
	return b.borderWidth
}

func (b *bordr) SetBorderWidth(as uint) {
	b.borderWidth = as
}

func (b *bordr) Draw(c Client, focusedWindow, focusedMonitor bool) {
	if b.borderWidth > 0 {
		//void window_draw_border(client_t *n, bool focused_window, bool focused_monitor);
	}
}

func (b *bordr) Color(focusedWindow, focusedMonitor bool) uint32 {
	/*
		//uint32_t get_border_color(client_t *c, bool focused_window, bool focused_monitor);
		if c != nil {
			var pxl uint32

			if focusedMonitor && focusedWindow {
				//if c.locked {
				//	get_color(focused_locked_border_color, c->window, &pxl);
				//} else if c.sticky {
				//	get_color(focused_sticky_border_color, c->window, &pxl);
				//} else if c.private {
				//	get_color(focused_private_border_color, c->window, &pxl);
				//} else {
				//	get_color(focused_border_color, c->window, &pxl);
				//}
			} else if focusedWindow {
				//if c.urgent {
				//	get_color(urgent_border_color, c->window, &pxl);
				//} else if c.locked {
				//	get_color(active_locked_border_color, c->window, &pxl);
				//} else if c.sticky {
				//	get_color(active_sticky_border_color, c->window, &pxl);
				//} else if c.private {
				//	get_color(active_private_border_color, c->window, &pxl);
				//} else {
				//	get_color(active_border_color, c->window, &pxl);
				//}
			} else {
				//if c.urgent {
				//	get_color(urgent_border_color, c->window, &pxl);
				//} else if c.locked {
				//	get_color(normal_locked_border_color, c->window, &pxl);
				//} else if c.sticky {
				//	get_color(normal_sticky_border_color, c->window, &pxl);
				//} else if c.private {
				//	get_color(normal_private_border_color, c->window, &pxl);
				//} else {
				//	get_color(normal_border_color, c->window, &pxl);
				//}
			}

			return pxl
	*/
	return 0
}

func getColor(c *xgb.Conn, win xproto.Window, color string, pxl uint32) bool {
	/*
		reply := xproto.GetWindowAttributes(win, win, nil)
		if reply != nil {
			cm := reply.Colormap

			if strings.Index(color, "#") == 0 {
				var red, green, blue uint
				if n, err := fmt.Sscanf(color, "%02x%02x%02x", &red, &green, &blue); n == 3 && err == nil {
					red *= 0x101
					green *= 0x101
					blue *= 0x101
					if r := xproto.AllocColorUnchecked(c, cm, red, green, blue); r != nil {
						*pxl = r.Pixel
						return true
					}
				}
			} else {
				if r := xproto.AllocNamedColorUnchecked(c, cm, uint16(len(color)), color); r != nil {
					*pxl = r.Pixel
					return true
				}
			}
		}
		pxl = 0
	*/
	return false
}
