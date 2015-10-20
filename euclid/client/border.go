package client

type Bordr interface {
	Draw(Client, bool, bool)
	Color(bool, bool) uint32
}

type bordr struct {
	borderWidth uint
	minWidth    uint16
	maxWidth    uint16
	minHeight   uint16
	maxHeight   uint16
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
