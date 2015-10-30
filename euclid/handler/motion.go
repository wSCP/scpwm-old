package handler

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/scpwm/euclid/window"
)

type Motion interface {
	Enable()
	Disable()
	Renew()
}

type motion struct {
	window.Window
}

func NewMotion(c *xgb.Conn, w, r xproto.Window) Motion {
	return &motion{
		Window: window.New(c, w, r),
	}
}

func (m *motion) Enable() {
	m.Raise()
	m.Show()
}

func (m *motion) Disable() {
	m.Hide()
}

func (m *motion) Renew() {
	geo, _ := xproto.GetGeometry(m.Conn(), xproto.Drawable(m.XRoot())).Reply()

	if geo != nil {
		m.Resize(geo.Width, geo.Height)
	}
}

/*
func mkMotion(s *xproto.ScreenInfo, c *xgb.Conn) (xproto.Window, error) {
	motion, err := xproto.NewWindowId(c)
	if err != nil {
		return motion, err
	}
	xproto.CreateWindow(
		c,
		s.RootDepth,
		motion,
		s.Root,
		0,
		0,
		s.WidthInPixels,
		s.HeightInPixels,
		0,
		xproto.WindowClassInputOnly,
		s.RootVisual,
		xproto.CwEventMask,
		[]uint32{xproto.EventMaskPointerMotion},
	)
	xproto.MapWindow(c, motion)
	return motion, nil
}
*/
