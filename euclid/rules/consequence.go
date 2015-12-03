package rules

import (
	"github.com/thrisp/scpwm/euclid/selector"
)

type Consequence struct {
	Manage         bool
	Class          string
	Instance       string
	OnMonitor      selector.Selector
	OnDesktop      selector.Selector
	WithClient     selector.Selector
	LayoutState    string
	Layer          string
	SplitDirection string
	SplitRatio     float64
	MinWidth       uint16
	MaxWidth       uint16
	MinHeight      uint16
	MaxHeight      uint16
	Locked         bool
	Sticky         bool
	Center         bool
	Focus          bool
	Border         bool
}

func NewConsequence(class, instance string, rules ...Rule) *Consequence {
	csq := &Consequence{
		Manage:   true,
		Class:    class,
		Instance: instance,
		Focus:    true,
		Border:   true,
	}
	for _, rule := range rules {
		eff := rule.Effect()
		switch eff.Applies() {
		case IsManaged:
			if !eff.Bool() {
				csq.Manage = false
			}
		case MonitorSelector:
			csq.OnMonitor = selector.New(eff.String())
		case DesktopSelector:
			csq.OnDesktop = selector.New(eff.String())
		case ClientSelector:
			csq.WithClient = selector.New(eff.String())
		case Layout:
			csq.LayoutState = eff.String()
		case Layer:
			csq.Layer = eff.String()
		case SplitDirection:
			csq.SplitDirection = eff.String()
		case SplitRatio:
			csq.SplitRatio = eff.Float()
		case MinWidth:
			csq.MinWidth = eff.Uint16()
		case MaxWidth:
			csq.MinWidth = eff.Uint16()
		case MinHeight:
			csq.MinHeight = eff.Uint16()
		case MaxHeight:
			csq.MaxHeight = eff.Uint16()
		case IsLocked:
			csq.Locked = eff.Bool()
		case IsSticky:
			csq.Sticky = eff.Bool()
		case IsCentered:
			csq.Center = eff.Bool()
		case IsFocused:
			csq.Focus = eff.Bool()
		case IsBordered:
			csq.Border = eff.Bool()
		}
	}
	return csq
}
