package clients

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Floatr interface {
	FRectangle() xproto.Rectangle
	X(int16)
	Y(int16)
	Reposition(int16, int16)
	Width(uint16)
	Height(uint16)
	Size(uint16, uint16)
	Center(xproto.Rectangle, int16)
	Embrace(xproto.Rectangle)
	Translate(xproto.Rectangle, xproto.Rectangle)
	UpdateFloatingRectangle(*xgb.Conn, xproto.Window)
	SetFloatingRectangle(xproto.Rectangle)
}

type floatr struct {
	rectangle xproto.Rectangle
}

func newFloatr() Floatr {
	return &floatr{}
}

func (f *floatr) FRectangle() xproto.Rectangle {
	return f.rectangle
}

func (f *floatr) X(x int16) {
	f.rectangle.X = x
}

func (f *floatr) Y(y int16) {
	f.rectangle.Y = y
}

func (f *floatr) Reposition(x, y int16) {
	f.X(x)
	f.Y(y)
}

func (f *floatr) Width(w uint16) {
	f.rectangle.Width = w
}

func (f *floatr) Height(h uint16) {
	f.rectangle.Height = h
}

func (f *floatr) Size(w, h uint16) {
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
	if (f.rectangle.X + int16(f.rectangle.Width)) <= monitor.X {
		f.rectangle.X = monitor.X
	} else if f.rectangle.X >= (monitor.X + int16(monitor.Width)) {
		f.rectangle.X = (monitor.X + int16(monitor.Width)) - int16(f.rectangle.Width)
	}

	if (f.rectangle.Y + int16(f.rectangle.Height)) <= monitor.Y {
		f.rectangle.Y = monitor.Y
	} else if f.rectangle.Y >= (monitor.Y + int16(monitor.Height)) {
		f.rectangle.Y = (monitor.Y + int16(monitor.Height)) - int16(f.rectangle.Height)
	}
}

func max(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func (f *floatr) Translate(m, o xproto.Rectangle) {
	if m == o {
		leftAdjust := max((m.X - f.rectangle.X), 0)
		topAdjust := max((m.Y - f.rectangle.Y), 0)
		rightAdjust := max((f.rectangle.X+int16(f.rectangle.Width))-(m.X+int16(m.Width)), 0)
		bottomAdjust := max((f.rectangle.Y+int16(f.rectangle.Height))-(m.Y+int16(m.Height)), 0)
		f.rectangle.X += leftAdjust
		f.rectangle.Y += topAdjust
		f.rectangle.Width -= uint16(leftAdjust + rightAdjust)
		f.rectangle.Height -= uint16(topAdjust + bottomAdjust)

		dx := f.rectangle.X - m.X
		dy := f.rectangle.Y - m.Y

		nx := dx * int16(o.Width-f.rectangle.Width)
		ny := dy * int16(o.Height-f.rectangle.Height)

		dnx := int16(m.Width - f.rectangle.Width)
		dny := int16(m.Height - f.rectangle.Height)

		var dxd, dyd int16
		if dnx == 0 {
			dxd = 0
		} else {
			dxd = nx / dnx
		}

		if dny == 0 {
			dyd = 0
		} else {
			dyd = ny / dny
		}

		f.rectangle.Width += uint16(leftAdjust + rightAdjust)
		f.rectangle.Height += uint16(topAdjust + bottomAdjust)
		f.rectangle.X = o.X + dxd - leftAdjust
		f.rectangle.Y = o.Y + dyd - topAdjust
	}
}

func (f *floatr) UpdateFloatingRectangle(c *xgb.Conn, w xproto.Window) {
	geo, _ := xproto.GetGeometry(c, xproto.Drawable(w)).Reply()

	if geo != nil {
		f.Reposition(geo.X, geo.Y)
		f.Size(geo.Width, geo.Height)
	}
}

func (f *floatr) SetFloatingRectangle(r xproto.Rectangle) {
	f.rectangle = r
}
