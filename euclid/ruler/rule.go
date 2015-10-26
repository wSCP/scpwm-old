package ruler

import (
	"fmt"
	"strconv"
)

type ruleCause int

const (
	className ruleCause = iota
	instanceName
	clientDescription
	windowDescription
)

var stringRuleCause map[string]ruleCause = map[string]ruleCause{
	"class":    className,
	"instance": instanceName,
	"client":   clientDescription,
	"window":   windowDescription,
}

func (r ruleCause) String() string {
	switch r {
	case className:
		return "class"
	case instanceName:
		return "instance"
	case clientDescription:
		return "client"
	case windowDescription:
		return "window"
	}
	return ""
}

type ruleEffect int

const (
	setSplitDirection ruleEffect = iota
	setSplitRatio
	minWidth
	maxWidth
	minHeight
	maxHeight
	isPseudoTiled
	isFloating
	isLocked
	isSticky
	isPrivate
	isBordered
	isCentered
	isFollowed
	isManaged
	isFocused
	monitorDescription
	desktopDescription
)

var stringRuleEffect map[string]ruleEffect = map[string]ruleEffect{
	"direction":   setSplitDirection,  //string
	"ratio":       setSplitRatio,      //float64
	"minwidth":    minWidth,           //uint16
	"maxwidth":    maxWidth,           //uint16
	"minheight":   minHeight,          //uint16
	"maxheight":   maxHeight,          //uint16
	"pseudotiled": isPseudoTiled,      //bool
	"floating":    isFloating,         //bool
	"lock":        isLocked,           //bool
	"sticky":      isSticky,           //bool
	"private":     isPrivate,          //bool
	"bordered":    isBordered,         //bool
	"centered":    isCentered,         //bool
	"follow":      isFollowed,         //bool
	"manage":      isManaged,          //bool
	"focus":       isFocused,          //bool
	"monitor":     monitorDescription, //string
	"desktop":     desktopDescription, //string
}

func (r ruleEffect) String() string {
	switch r {
	case setSplitDirection:
		return "split direction"
	case setSplitRatio:
		return "split ratio"
	case minWidth:
		return "minimum width"
	case maxWidth:
		return "maximum width"
	case minHeight:
		return "minimum height"
	case maxHeight:
		return "maximum height"
	case isPseudoTiled:
		return "pseudo tiled"
	case isFloating:
		return "floating"
	case isLocked:
		return "locked"
	case isSticky:
		return "sticky"
	case isPrivate:
		return "private"
	case isBordered:
		return "bordered"
	case isCentered:
		return "centered"
	case isManaged:
		return "managed"
	case isFocused:
		return "focused"
	case monitorDescription:
		return "monitor"
	case desktopDescription:
		return "desktop"
	}
	return ""
}

type Rule interface {
	Cause() ruleCause
	Key() string
	Effect() ruleEffect
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
	cause  ruleCause
	key    string
	effect ruleEffect
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

func (r *rule) Cause() ruleCause {
	return r.cause
}

func (r *rule) Key() string {
	return r.key
}

func (r *rule) Effect() ruleEffect {
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