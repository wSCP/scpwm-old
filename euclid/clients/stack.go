package clients

type Stackr interface {
	GetStack() int
	SetStack(int)
}

func newStackr() *stackr {
	return &stackr{}
}

type stackr struct {
	stack int
}

func (s *stackr) GetStack() int {
	return s.stack
}

func (s *stackr) SetStack(i int) {
	s.stack = i
}
