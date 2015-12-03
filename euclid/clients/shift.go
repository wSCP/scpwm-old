package clients

import "github.com/thrisp/scpwm/euclid/rules"

type ShiftMode int

const (
	Automatic ShiftMode = iota
	Manual
)

var modeString map[string]ShiftMode = map[string]ShiftMode{
	"automatic": Automatic,
	"manual":    Manual,
}

type ShiftOrientation int

const (
	Horizontal ShiftOrientation = iota
	Vertical
)

var orientationString map[string]ShiftOrientation = map[string]ShiftOrientation{
	"horizontal": Horizontal,
	"vertical":   Vertical,
}

func (s ShiftOrientation) String() string {
	switch s {
	case Horizontal:
		return "horizontal"
	case Vertical:
		return "vertical"
	}
	return ""
}

type ShiftDirection int

const (
	Right ShiftDirection = iota
	Down
	Left
	Up
)

var directionString map[string]ShiftDirection = map[string]ShiftDirection{
	"right": Right,
	"down":  Down,
	"left":  Left,
	"up":    Up,
}

type Shiftr interface {
	Rotate() int
	SetRotate(int)
	Ratio() float64
	SetRatio(float64)
	Mode() ShiftMode
	SetMode(string)
	Orientation() ShiftOrientation
	SetOrientation(string)
	Direction() ShiftDirection
	SetDirection(string)
	SetShift(*rules.Consequence)
}

type shiftr struct {
	rotate      int
	ratio       float64
	mode        ShiftMode
	orientation ShiftOrientation
	direction   ShiftDirection
}

func newShiftr() Shiftr {
	return &shiftr{}
}

func (s *shiftr) Rotate() int {
	return s.rotate
}

func (s *shiftr) SetRotate(v int) {
	s.rotate = v
}

func (s *shiftr) Ratio() float64 {
	return s.ratio
}

func (s *shiftr) SetRatio(v float64) {
	s.ratio = v
}

func (s *shiftr) Mode() ShiftMode {
	return s.mode
}

func (s *shiftr) SetMode(v string) {
	if val, ok := modeString[v]; ok {
		s.mode = val
	}
}

func (s *shiftr) Orientation() ShiftOrientation {
	return s.orientation
}

func (s *shiftr) SetOrientation(v string) {
	if val, ok := orientationString[v]; ok {
		s.orientation = val
	}
}

func (s *shiftr) Direction() ShiftDirection {
	return s.direction
}

func (s *shiftr) SetDirection(v string) {
	if val, ok := directionString[v]; ok {
		s.direction = val
	}
}

func (s *shiftr) SetShift(csq *rules.Consequence) {
	if csq.SplitRatio != 0.0 {
		s.ratio = csq.SplitRatio
	}
	if csq.SplitDirection != "" {
		s.mode = Manual
		if val, ok := directionString[csq.SplitDirection]; ok {
			s.direction = val
		}
	}
}
