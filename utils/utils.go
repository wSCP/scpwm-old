package utils

import "strings"

func MatchesAny(s string, ss ...string) bool {
	for _, str := range ss {
		if s == str || strings.Contains(str, s) {
			return true
		}
	}
	return false
}

/*
const (
	MAXINT   = int(^uint(0) >> 1)
	MAXSTATE = int(4)
)

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

*/
