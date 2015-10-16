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
