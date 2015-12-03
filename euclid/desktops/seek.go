package desktops

import "github.com/thrisp/scpwm/euclid/branch"

type MatchDesktop func(Desktop) bool

func seek(desktops *branch.Branch, fn MatchDesktop) Desktop {
	curr := desktops.Front()
	for curr != nil {
		desktop := curr.Value.(Desktop)
		if match := fn(desktop); match {
			return desktop
		}
		curr = curr.Next()
	}
	return nil
}

func isFocused(d Desktop) bool {
	if d.Focused() {
		return true
	}
	return false
}

func Focused(desktops *branch.Branch) Desktop {
	return seek(desktops, isFocused)
}

func seekOffset(desktops *branch.Branch, fn MatchDesktop, offset int) Desktop {
	curr := desktops.Front()
	for curr != nil {
		desktop := curr.Value.(Desktop)
		if match := fn(desktop); match {
			switch offset {
			case -1:
				desktop = curr.PrevContinuous().Value.(Desktop)
			case 1:
				desktop = curr.NextContinuous().Value.(Desktop)
			}
			return desktop
		}
		curr = curr.Next()
	}
	return nil
}

func Prev(desktops *branch.Branch) Desktop {
	return seekOffset(desktops, isFocused, -1)
}

func Next(desktops *branch.Branch) Desktop {
	return seekOffset(desktops, isFocused, 1)
}

func seekAny(desktops *branch.Branch, fn MatchDesktop) []Desktop {
	var ret []Desktop
	curr := desktops.Front()
	for curr != nil {
		dsk := curr.Value.(Desktop)
		if match := fn(dsk); match {
			ret = append(ret, dsk)
		}
		curr = curr.Next()
	}
	return ret
}

func All(desktops *branch.Branch) []Desktop {
	fn := func(d Desktop) bool { return true }
	return seekAny(desktops, fn)
}
