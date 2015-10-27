package desktops

import (
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/settings"
)

func New(monitor uint32, s settings.Settings) *branch.Branch {
	return branch.New("desktops")
}

func Initialize(desktops *branch.Branch, monitor uint32, s settings.Settings) {
	if desktops.Len() == 0 {
		d := NewDesktop(monitor, "", s)
		desktops.PushFront(d)
	}
	UpdateDesktopsIndex(desktops)
}

func All(desktops *branch.Branch) []Desktop {
	var ret []Desktop
	curr := desktops.Front()
	for curr != nil {
		dsk := curr.Value.(Desktop)
		ret = append(ret, dsk)
		curr = curr.Next()
	}
	return ret
}

func Add(desktops *branch.Branch, monitor uint32, name string, s settings.Settings) {
	nd := NewDesktop(monitor, name, s)
	desktops.PushBack(nd)
	UpdateDesktopsIndex(desktops)
}

type SelectDesktop func(Desktop) bool

func seek(desktops *branch.Branch, fn SelectDesktop) Desktop {
	curr := desktops.Front()
	for curr != nil {
		desktop := curr.Value.(Desktop)
		if found := fn(desktop); found {
			return desktop
		}
		curr = curr.Next()
	}
	return nil
}

func seekOffset(desktops *branch.Branch, fn SelectDesktop, offset int) Desktop {
	curr := desktops.Front()
	for curr != nil {
		desktop := curr.Value.(Desktop)
		if found := fn(desktop); found {
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

func isFocused(d Desktop) bool {
	if d.Focused() {
		return true
	}
	return false
}

func Focused(desktops *branch.Branch) Desktop {
	return seek(desktops, isFocused)
}

func Prev(desktops *branch.Branch) Desktop {
	return seekOffset(desktops, isFocused, -1)
}

func Next(desktops *branch.Branch) Desktop {
	return seekOffset(desktops, isFocused, 1)
}

type UpdateDesktop func(d Desktop) error

func update(desktops *branch.Branch, fn UpdateDesktop) error {
	curr := desktops.Front()
	for curr != nil {
		d := curr.Value.(Desktop)
		if err := fn(d); err != nil {
			return err
		}
		curr = curr.Next()
	}
	return nil
}

func UpdateDesktopsMonitor(desktops *branch.Branch, id uint32) {
	fn := func(d Desktop) error {
		d.Set("monitor", id)
		return nil
	}
	update(desktops, fn)
}

func UpdateDesktopsIndex(desktops *branch.Branch) {
	idx := 1
	fn := func(d Desktop) error {
		d.Set("index", idx)
		idx++
		return nil
	}
	update(desktops, fn)
}
