package clients

type ShiftMode int

const (
	Automatic ShiftMode = iota
	Manual
)

type ShiftOrientation int

const (
	Horizontal ShiftOrientation = iota
	Vertical
)

type ShiftDirection int

const (
	Right ShiftDirection = iota
	Down
	Left
	Up
)

type Shiftr interface {
	Rotate() int
	Ratio() float64
	Mode() ShiftMode
	Orientation() ShiftOrientation
	Direction() ShiftDirection
	SetShift(string, string)
}

type shiftr struct {
	rotate      int
	ratio       float64
	mode        ShiftMode
	orientation ShiftOrientation
	direction   ShiftDirection
}

func newShiftr() *shiftr {
	return &shiftr{}
}

func (s *shiftr) SetShift(k, v string) {
	switch k {
	case "rotate":
		//
	case "ratio":
		//
	case "mode":
		//
	case "orientation":
		//
	case "direction":
		//
	}
}

func (s *shiftr) Rotate() int {
	return s.rotate
}

func (s *shiftr) Ratio() float64 {
	return s.ratio
}

func (s *shiftr) Mode() ShiftMode {
	return s.mode
}

func (s *shiftr) Orientation() ShiftOrientation {
	return s.orientation
}

func (s *shiftr) Direction() ShiftDirection {
	return s.direction
}
