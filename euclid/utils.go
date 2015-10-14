package main

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
)

const MAXINT = int(^uint(0) >> 1)
const MAXSTATE = int(4)

func abs(num int16) int {
	if num < 0 {
		num = -num
	}
	return int(num)
}

func min() {}

func max(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func fmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

type Node int

const (
	nWM Node = iota
	nMonitor
	nDesktop
	nClient
)

var stringNode map[string]Node = map[string]Node{
	"euclid":  nWM,
	"monitor": nMonitor,
	"desktop": nDesktop,
	"client":  nClient,
}

func (n Node) String() string {
	switch n {
	case nWM:
		return "euclid"
	case nMonitor:
		return "monitor"
	case nDesktop:
		return "desktop"
	case nClient:
		return "client"
	}
	return ""
}

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
)

var stringOrientation map[string]Orientation = map[string]Orientation{
	"closer":     Closer,
	"closest":    Closest,
	"further":    Further,
	"furthest":   Furthest,
	"horizontal": Horizontal,
	"vertical":   Vertical,
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

type Pad map[Direction]int

func DefaultPad() Pad {
	ret := make(Pad)
	ret[Up], ret[Down], ret[Left], ret[Right] = 0, 0, 0, 0
	return ret
}

func PadValue(d string, p Pad) int {
	switch d {
	case "up":
		return p[Up]
	case "down":
		return p[Down]
	case "left":
		return p[Left]
	case "right":
		return p[Right]
	}
	return 0
}

func isAppendable(s string, ss []string) bool {
	for _, x := range ss {
		if x == s {
			return false
		}
	}
	return true
}

func doAdd(s string, ss []string) []string {
	if isAppendable(s, ss) {
		ss = append(ss, s)
	}
	return ss
}

func indexFromString(s string, i int) bool {
	var idx int
	n, err := fmt.Sscanf(s, "%d", idx)
	if err != nil {
		return false
	}
	if n != 1 || idx < 1 {
		return false
	}
	i = idx
	return true
}

func contains(a, b xproto.Rectangle) bool {
	return (a.X <= b.X && (a.X+int16(a.Width)) >= (b.X+int16(b.Width)) &&
		a.Y <= b.Y && (a.Y+int16(a.Height)) >= (b.Y+int16(b.Height)))
}
