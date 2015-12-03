package desktops

import "github.com/thrisp/scpwm/euclid/branch"

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
