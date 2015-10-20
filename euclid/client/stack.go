package client

type Stack interface {
	Min() int
	Max() int
	Reset()
}

type Stackr interface {
	GetStack() int
	SetStack(int)
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
