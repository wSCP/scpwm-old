package clients

import (
	"github.com/BurntSushi/xgb/xproto"
)

type Floatr interface {
	Width(int)
	Height(int)
	Size(int, int)
	Center(xproto.Rectangle)
	Embrace(xproto.Rectangle)
	Translate(xproto.Rectangle, xproto.Rectangle)
	SetFloatingRectangle(xproto.Rectangle)
}

type floatr struct {
	rectangle xproto.Rectangle
}

func newFloatr(r xproto.Rectangle) *floatr {
	return &floatr{rectangle: r}
}

func (f *floatr) SetFloatingRectangle(r xproto.Rectangle) {
	f.rectangle = r
	//void update_floating_rectangle(client_t *c);
}

func (f *floatr) Width(w int) {
	//void restrain_floating_width(client_t *c, int *width);
}

func (f *floatr) Height(h int) {
	//void restrain_floating_height(client_t *c, int *height);
}

func (f *floatr) Size(w, h int) {
	f.Width(w)
	f.Height(h)
}

func (f *floatr) Center(rect xproto.Rectangle, borderWidth int16) {
	r := f.rectangle

	if r.Width >= rect.Width {
		r.X = rect.X
	} else {
		r.X = rect.X + (int16(rect.Width)-int16(r.Width))/2
	}

	if r.Height >= rect.Height {
		r.Y = rect.Y
	} else {
		r.Y = rect.Y + (int16(rect.Height)-int16(r.Height))/2
	}

	r.X -= borderWidth
	r.Y -= borderWidth
}

func (f *floatr) Embrace(monitor xproto.Rectangle) {
	//cf := c.floater.rectangle
	//if (cf.X + int16(cf.Width)) <= m.rectangle.X {
	//	c.fRectangle.X = m.rectangle.X
	//} else if cf.X >= (m.rectangle.X + int16(m.rectangle.Width)) {
	//	c.fRectangle.X = (m.rectangle.X + int16(m.rectangle.Width)) - int16(c.fRectangle.Width)
	//}

	//if (cf.Y + int16(cf.Height)) <= m.rectangle.Y {
	//	c.fRectangle.Y = m.rectangle.Y
	//} else if cf.Y >= (m.rectangle.Y + int16(m.rectangle.Height)) {
	//	c.fRectangle.Y = (m.rectangle.Y + int16(m.rectangle.Height)) - int16(c.fRectangle.Height)
	//}
}

func (f *floatr) Translate(src, dst xproto.Rectangle) {
	//func (m *Monitor) Translate(o *Monitor, c *Client) {
	//if m.e.pointer.action == NoAction || m == o {
	//	leftAdjust := max((m.rectangle.X - c.fRectangle.X), 0)
	//	topAdjust := max((m.rectangle.Y - c.fRectangle.Y), 0)
	//	rightAdjust := max((c.fRectangle.X+int16(c.fRectangle.Width))-(m.rectangle.X+int16(m.rectangle.Width)), 0)
	//	bottomAdjust := max((c.fRectangle.Y+int16(c.fRectangle.Height))-(m.rectangle.Y+int16(m.rectangle.Height)), 0)
	//	c.fRectangle.X += leftAdjust
	//	c.fRectangle.Y += topAdjust
	//	c.fRectangle.Width -= uint16(leftAdjust + rightAdjust)
	//	c.fRectangle.Height -= uint16(topAdjust + bottomAdjust)
	//
	//		dx := c.fRectangle.X - m.rectangle.X
	//		dy := c.fRectangle.Y - m.rectangle.Y
	//
	//		nx := dx * int16(o.rectangle.Width-c.fRectangle.Width)
	//		ny := dy * int16(o.rectangle.Height-c.fRectangle.Height)
	//
	//		dnx := int16(m.rectangle.Width - c.fRectangle.Width)
	//		dny := int16(m.rectangle.Height - c.fRectangle.Height)
	//
	//		var dxd, dyd int16
	//		if dnx == 0 {
	//			dxd = 0
	//		} else {
	//			dxd = nx / dnx
	//		}
	//
	//		if dny == 0 {
	//			dyd = 0
	//		} else {
	//			dyd = ny / dny
	//		}

	//		c.fRectangle.Width += uint16(leftAdjust + rightAdjust)
	//		c.fRectangle.Height += uint16(topAdjust + bottomAdjust)
	//		c.fRectangle.X = o.rectangle.X + dxd - leftAdjust
	//		c.fRectangle.Y = o.rectangle.Y + dyd - topAdjust
	//	}
	//}
}
