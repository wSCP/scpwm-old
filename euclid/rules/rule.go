package rules

import (
	"fmt"
	"strconv"
)

type Apply int

const (
	SetSplitDirection Apply = iota
	SetSplitRatio
	MinWidth
	MaxWidth
	MinHeight
	MaxHeight
	IsPseudoTiled
	IsFloating
	IsLocked
	IsSticky
	IsPrivate
	IsBordered
	IsCentered
	IsManaged
	IsFocused
	ClientSelector
	DesktopSelector
	MonitorSelector
)

var stringApply map[string]Apply = map[string]Apply{
	"direction":   SetSplitDirection, //string
	"ratio":       SetSplitRatio,     //float64
	"minwidth":    MinWidth,          //uint16
	"maxwidth":    MaxWidth,          //uint16
	"minheight":   MinHeight,         //uint16
	"maxheight":   MaxHeight,         //uint16
	"pseudotiled": IsPseudoTiled,     //bool
	"floating":    IsFloating,        //bool
	"lock":        IsLocked,          //bool
	"sticky":      IsSticky,          //bool
	"private":     IsPrivate,         //bool
	"bordered":    IsBordered,        //bool
	"centered":    IsCentered,        //bool
	"manage":      IsManaged,         //bool
	"focus":       IsFocused,         //bool
	"client":      ClientSelector,    //string
	"desktop":     DesktopSelector,   //string
	"monitor":     MonitorSelector,   //string
}

func (a Apply) String() string {
	switch a {
	case SetSplitDirection:
		return "split direction"
	case SetSplitRatio:
		return "split ratio"
	case MinWidth:
		return "minimum width"
	case MaxWidth:
		return "maximum width"
	case MinHeight:
		return "minimum height"
	case MaxHeight:
		return "maximum height"
	case IsPseudoTiled:
		return "pseudo tiled"
	case IsFloating:
		return "floating"
	case IsLocked:
		return "locked"
	case IsSticky:
		return "sticky"
	case IsPrivate:
		return "private"
	case IsBordered:
		return "bordered"
	case IsCentered:
		return "centered"
	case IsManaged:
		return "managed"
	case IsFocused:
		return "focused"
	case ClientSelector:
		return "select client"
	case DesktopSelector:
		return "select desktop"
	case MonitorSelector:
		return "select monitor"
	}
	return ""
}

type Component interface {
	Applies() Apply
	Implementr
}

type Implementr interface {
	String() string
	Uint16() uint16
	Float() float64
	Bool() bool
}

type component struct {
	applies Apply
	raw     string
}

func newComponent(apply, to string) Component {
	if a, ok := stringApply[apply]; ok {
		return &component{
			applies: a,
			raw:     to,
		}
	}
	return nil
}

func (c *component) Applies() Apply {
	return c.applies
}

func (c *component) String() string {
	return c.raw
}

func (c *component) Uint16() uint16 {
	if i, err := strconv.ParseUint(c.raw, 10, 16); err == nil {
		return uint16(i)
	}
	return 0
}

func (c *component) Float() float64 {
	if f, err := strconv.ParseFloat(c.raw, 64); err == nil {
		return f
	}
	return 0.0
}

func (c *component) Bool() bool {
	if b, err := strconv.ParseBool(c.raw); err == nil {
		return b
	}
	return false
}

type Rule interface {
	Cause() Component
	Effect() Component
	String() string
	Reuse() bool
}

type rule struct {
	cause  Component
	effect Component
	reuse  bool
}

func newRule(c, k, e, v string, reuse bool) Rule {
	if cause := newComponent(c, k); cause != nil {
		if effect := newComponent(e, v); effect != nil {
			return &rule{cause, effect, reuse}
		}
	}
	return nil
}

func (r *rule) Cause() Component {
	return r.cause
}

func (r *rule) Effect() Component {
	return r.effect
}

func (r *rule) String() string {
	return fmt.Sprintf("%s %s", r.cause.String(), r.effect.String())
}

func (r *rule) Reuse() bool {
	return r.reuse
}
