package clients

type Layer int

const (
	Normal Layer = iota
	Below
	Above
)

type Stackr interface {
	GetStack() Layer
	SetStack(Layer)
}

type stackr struct {
	stack Layer
}

func newStackr() Stackr {
	return &stackr{}
}

func (s *stackr) GetStack() Layer {
	return s.stack
}

func (s *stackr) SetStack(l Layer) {
	s.stack = l
}

func CompareStack(c, o Client) int {
	cl, ol := c.GetStack(), o.GetStack()
	if cl == ol {
		if !c.Floating() && o.Floating() {
			return -1
		} else if c.Floating() && o.Floating() {
			return 1
		} else {
			return 0
		}
	} else {
		if cl == Below {
			return -1
		} else if cl == Above {
			return 1
		} else {
			if ol == Above {
				return -1
			}
		}
	}
	return 1
}
