package client

import "github.com/BurntSushi/xgb/xproto"

type Tilr interface {
	Area() int
}

type tilr struct {
	rectangle xproto.Rectangle
}

func (t *tilr) Area() int {
	return int(t.rectangle.Width * t.rectangle.Height)
}
