package main

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/xgb"
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

func getColor(c *xgb.Conn, win xproto.Window, color string, pxl uint32) bool {
	reply := xproto.GetWindowAttributes(win, win, nil)
	if reply != nil {
		cm := reply.Colormap

		if strings.Index(color, "#") == 0 {
			var red, green, blue uint
			if n, err := fmt.Sscanf(color, "%02x%02x%02x", &red, &green, &blue); n == 3 && err == nil {
				red *= 0x101
				green *= 0x101
				blue *= 0x101
				if r := xproto.AllocColorUnchecked(c, cm, red, green, blue); r != nil {
					*pxl = r.Pixel
					return true
				}
			}
		} else {
			if r := xproto.AllocNamedColorUnchecked(c, cm, uint16(len(color)), color); r != nil {
				*pxl = r.Pixel
				return true
			}
		}
	}
	pxl = 0
	return false
}
