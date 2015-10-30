package ruler

import (
	"fmt"
	"strconv"
)

type RuleCause int

const (
	ClassName RuleCause = iota
	InstanceName
	AnyAll
)

var stringRuleCause map[string]RuleCause = map[string]RuleCause{
	"class":    ClassName,
	"instance": InstanceName,
	"any":      AnyAll,
	"all":      AnyAll,
}

func (r RuleCause) String() string {
	switch r {
	case ClassName:
		return "class"
	case InstanceName:
		return "instance"
	case AnyAll:
		return "any, all"
	}
	return ""
}

type RuleEffect int

const (
	SetSplitDirection RuleEffect = iota
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
	//IsFollowed
	IsManaged
	IsFocused
	ClientDescription
	MonitorDescription
	DesktopDescription
)

var stringRuleEffect map[string]RuleEffect = map[string]RuleEffect{
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
	//"follow":      IsFollowed,         //bool
	"manage":  IsManaged,          //bool
	"focus":   IsFocused,          //bool
	"client":  ClientDescription,  //string
	"monitor": MonitorDescription, //string
	"desktop": DesktopDescription, //string
}

func (r RuleEffect) String() string {
	switch r {
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
	case MonitorDescription:
		return "monitor"
	case DesktopDescription:
		return "desktop"
	}
	return ""
}

type Rule interface {
	Cause() RuleCause
	Key() string
	Effect() RuleEffect
	Value() string
	Reuse() bool
	String() string
	Implementr
}

type Implementr interface {
	Uint16() uint16
	Float() float64
	Bool() bool
}

type rule struct {
	cause  RuleCause
	key    string
	effect RuleEffect
	value  string
	reuse  bool
}

func newRule(cause, key, effect, value string, reuse bool) Rule {
	if rc, ok := stringRuleCause[cause]; ok {
		if re, ok := stringRuleEffect[effect]; ok {
			return &rule{
				cause:  rc,
				key:    cause,
				effect: re,
				value:  value,
				reuse:  reuse,
			}
		}
	}
	return nil
}

func (r *rule) Cause() RuleCause {
	return r.cause
}

func (r *rule) Key() string {
	return r.key
}

func (r *rule) Effect() RuleEffect {
	return r.effect
}

func (r *rule) Value() string {
	return r.value
}

func (r *rule) String() string {
	return fmt.Sprintf(
		"%s %s %s %s %t",
		r.cause,
		r.key,
		r.effect,
		r.value,
		r.reuse,
	)
}

func (r *rule) Uint16() uint16 {
	if i, err := strconv.ParseUint(r.value, 10, 16); err == nil {
		return uint16(i)
	}
	return 0
}

func (r *rule) Float() float64 {
	if f, err := strconv.ParseFloat(r.value, 64); err == nil {
		return f
	}
	return 0.0
}

func (r *rule) Bool() bool {
	if b, err := strconv.ParseBool(r.value); err == nil {
		return b
	}
	return false
}

func (r *rule) Reuse() bool {
	return r.reuse
}
