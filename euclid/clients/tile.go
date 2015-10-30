package clients

import "github.com/BurntSushi/xgb/xproto"

type Tilr interface {
	Area() int
	SetTiledRectangle(xproto.Rectangle)
}

type tilr struct {
	rectangle xproto.Rectangle
}

func newTilr() *tilr {
	return &tilr{}
}

func (t *tilr) SetTiledRectangle(r xproto.Rectangle) {
	t.rectangle = r
}

func (t *tilr) Area() int {
	return int(t.rectangle.Width * t.rectangle.Height)
}
