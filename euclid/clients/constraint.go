package clients

import (
	"github.com/thrisp/scpwm/euclid/rules"
)

type Constraint interface {
	Min(string) uint16
	Max(string) uint16
	SetConstraint(*rules.Consequence)
}

type constraint struct {
	minWidth  uint16
	maxWidth  uint16
	minHeight uint16
	maxHeight uint16
}

func newConstraint(csq *rules.Consequence) Constraint {
	c := &constraint{}
	c.SetConstraint(csq)
	return c
}

func (c *constraint) Min(dimension string) uint16 {
	switch dimension {
	case "height":
		return c.minHeight
	case "width":
		return c.minWidth
	}
	return 0
}

func (c *constraint) Max(dimension string) uint16 {
	switch dimension {
	case "height":
		return c.maxHeight
	case "width":
		return c.maxWidth
	}
	return 0

}

func (c *constraint) SetConstraint(csq *rules.Consequence) {
	if csq.MinHeight != 0 {
		c.minHeight = csq.MinHeight
	}
	if csq.MaxHeight != 0 {
		c.maxHeight = csq.MaxHeight
	}
	if csq.MinWidth != 0 {
		c.minWidth = csq.MinWidth
	}
	if csq.MaxWidth != 0 {
		c.maxWidth = csq.MaxWidth
	}
}

func RestrainWidth(c Constraint, w uint16) uint16 {
	if w < 1 {
		w = 1
	}
	mw := c.Min("width")
	if mw > 0 && w < mw {
		return mw
	}
	xw := c.Max("width")
	if xw > 0 && w > xw {
		return xw
	}
	return w
}

func RestrainHeight(c Constraint, h uint16) uint16 {
	if h < 1 {
		h = 1
	}
	mh := c.Min("height")
	if mh > 0 && h < mh {
		return mh
	}
	xh := c.Max("height")
	if xh > 0 && h > xh {
		return xh
	}
	return h
}

func RestrainSize(c Constraint, w, h uint16) (uint16, uint16) {
	return RestrainWidth(c, w), RestrainHeight(c, h)
}
