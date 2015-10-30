package clients

type Constraint interface {
	Min(string) uint16
	Max(string) uint16
	SetConstraint(string, string, uint16)
}

type constraint struct {
	minWidth  uint16
	maxWidth  uint16
	minHeight uint16
	maxHeight uint16
}

func newConstraint() *constraint {
	return &constraint{}
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

func (c *constraint) setHeightDirection(direction string, value uint16) {
	switch direction {
	case "min":
		c.minHeight = value
	case "max":
		c.maxHeight = value
	}
}

func (c *constraint) setWidthDirection(direction string, value uint16) {
	switch direction {
	case "min":
		c.minWidth = value
	case "max":
		c.maxWidth = value
	}
}

func (c *constraint) SetConstraint(dimension, direction string, value uint16) {
	switch dimension {
	case "height":
		c.setHeightDirection(direction, value)
	case "width":
		c.setWidthDirection(direction, value)
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
