package selector

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up
)

var stringDirection = map[string]Direction{
	"right": Right,
	"down":  Down,
	"left":  Left,
	"up":    Up,
}

func (d Direction) String() string {
	switch d {
	case Right:
		return "right"
	case Down:
		return "down"
	case Left:
		return "left"
	case Up:
		return "up"
	}
	return ""
}

func opposite(src Direction) Direction {
	var dst Direction
	switch src {
	case Right:
		dst = Left
	case Down:
		dst = Up
	case Left:
		dst = Right
	case Up:
		dst = Down
	}
	return dst
}

type Cycle int

const (
	First Cycle = iota
	Backward
	Prev
	Next
	Forward
	Last
)

var stringCycle map[string]Cycle = map[string]Cycle{
	"first":    First,
	"backward": Backward,
	"previous": Prev,
	"next":     Next,
	"forward":  Forward,
	"last":     Last,
}

func (c Cycle) String() string {
	switch c {
	case First:
		return "first"
	case Forward:
		return "forward"
	case Prev:
		return "previous"
	case Next:
		return "next"
	case Backward:
		return "backward"
	case Last:
		return "last"
	}
	return ""
}

func reverse(src Cycle) Cycle {
	var dst Cycle
	switch src {
	case First:
		dst = Last
	case Forward:
		dst = Backward
	case Prev:
		dst = Next
	case Next:
		dst = Prev
	case Backward:
		dst = Forward
	case Last:
		dst = First
	}
	return dst
}

type Size int

const (
	Biggest Size = iota
	Smallest
)

var stringSize map[string]Size = map[string]Size{
	"biggest":  Biggest,
	"smallest": Smallest,
}

type Orientation int

const (
	Closer Orientation = iota
	Closest
	Further
	Furthest
	Horizontal
	Vertical
	Above
	Below
)

var stringOrientation map[string]Orientation = map[string]Orientation{
	"closer":     Closer,
	"closest":    Closest,
	"further":    Further,
	"furthest":   Furthest,
	"horizontal": Horizontal,
	"vertical":   Vertical,
	"above":      Above,
	"below":      Below,
}

func (o Orientation) String() string {
	switch o {
	case Closer:
		return "closer"
	case Closest:
		return "closest"
	case Further:
		return "further"
	case Furthest:
		return "furthest"
	case Horizontal:
		return "horizontal"
	case Vertical:
		return "vertical"
	case Above:
		return "above"
	case Below:
		return "below"
	}
	return ""
}

type Age int

const (
	Youngest Age = iota
	Younger
	Older
	Oldest
)

var stringAge map[string]Age = map[string]Age{
	"youngest": Youngest,
	"younger":  Younger,
	"older":    Older,
	"oldest":   Oldest,
}
