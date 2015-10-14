package main

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

var KeyParseError = Krror("Could not find a valid keycode in the string '%s'. Key binding failed.").Out

func ParseKeyInput(k *Keyboard, s string) (uint16, []xproto.Keycode, error) {
	mods, kcs := uint16(0), []xproto.Keycode{}
	for _, part := range strings.Split(s, "+") {
		switch strings.ToLower(part) {
		case "shift":
			mods |= xproto.ModMaskShift
		case "lock":
			mods |= xproto.ModMaskLock
		case "control", "ctrl":
			mods |= xproto.ModMaskControl
		case "mod1", "alt":
			mods |= xproto.ModMask1
		case "mod2":
			mods |= xproto.ModMask2
		case "mod3":
			mods |= xproto.ModMask3
		case "mod4", "super":
			mods |= xproto.ModMask4
		case "mod5":
			mods |= xproto.ModMask5
		case "any":
			mods |= xproto.ModMaskAny
		default:
			if len(kcs) == 0 {
				kcs = k.StrToKeycode(part)
			}
		}
	}

	if len(kcs) == 0 {
		return 0, nil, KeyParseError(s)
	}

	return mods, kcs, nil
}

var ButtonParseError = Krror("Could not find a valid mouse button in the string '%s'. Button binding failed.").Out

func ParseMouseInput(s string) (uint16, xproto.Button, error) {
	mods, button := uint16(0), xproto.Button(0)
	for _, part := range strings.Split(s, "+") {
		switch strings.ToLower(part) {
		case "shift":
			mods |= xproto.ModMaskShift
		case "lock":
			mods |= xproto.ModMaskLock
		case "control", "ctrl":
			mods |= xproto.ModMaskControl
		case "mod1", "alt":
			mods |= xproto.ModMask1
		case "mod2":
			mods |= xproto.ModMask2
		case "mod3":
			mods |= xproto.ModMask3
		case "mod4", "super":
			mods |= xproto.ModMask4
		case "mod5":
			mods |= xproto.ModMask5
		default:
			switch part {
			case "button1":
				button = xproto.Button(1)
			case "button2":
				button = xproto.Button(2)
			case "button3":
				button = xproto.Button(3)
			case "button4":
				button = xproto.Button(4)
			case "button5":
				button = xproto.Button(5)
			case "buttonany":
				button = xproto.ButtonIndexAny
			}
		}
	}
	if button == 0 {
		return 0, button, ButtonParseError(s)
	}
	return mods, button, nil
}

var UnMappable = Krror("could not get %s mapping: %v, unrecoverable.").Out

func minMaxKeycodeGet(s *xproto.SetupInfo) (xproto.Keycode, xproto.Keycode) {
	return s.MinKeycode, s.MaxKeycode
}

type Keyboard struct {
	Min xproto.Keycode
	Max xproto.Keycode
	*xproto.GetKeyboardMappingReply
}

func (k *Keyboard) Keycodes(keysym xproto.Keysym) []xproto.Keycode {
	var c byte
	var keycode xproto.Keycode
	keycodes := make([]xproto.Keycode, 0)
	set := make(map[xproto.Keycode]bool, 0)

	for kc := int(k.Min); kc <= int(k.Max); kc++ {
		keycode = xproto.Keycode(kc)
		for c = 0; c < k.KeysymsPerKeycode; c++ {
			if keysym == k.Keysymget(keycode, c) && !set[keycode] {
				keycodes = append(keycodes, keycode)
				set[keycode] = true
			}
		}
	}
	return keycodes
}

func (k *Keyboard) Keysymget(keycode xproto.Keycode, column byte) xproto.Keysym {
	i := (int(keycode)-int(k.Min))*int(k.KeysymsPerKeycode) + int(column)
	return k.Keysyms[i]
}

func (k *Keyboard) KeycodeToString(keycode xproto.Keycode) string {
	return strKeysyms[k.Keysyms[(int(keycode)-int(k.Min))*int(k.KeysymsPerKeycode)]]
}

func ByteToString(k *Keyboard, b1 byte, b2 byte) string {
	if b1 != 0 {
		return k.KeycodeToString(xproto.Keycode(b1))
	}
	if b2 != 0 {
		return fmt.Sprintf("button%d", b2)
	}
	return ""
}

func (k *Keyboard) StrToKeycode(str string) []xproto.Keycode {
	sym, ok := keysyms[str]
	if !ok {
		sym, ok = keysyms[strings.Title(str)]
	}
	if !ok {
		sym, ok = keysyms[strings.ToLower(str)]
	}
	if !ok {
		sym, ok = keysyms[strings.ToUpper(str)]
	}

	if !ok {
		return []xproto.Keycode{}
	}
	return k.Keycodes(sym)
}

func NewKeyboard(s *xproto.SetupInfo, c *xgb.Conn) (*Keyboard, error) {
	min, max := minMaxKeycodeGet(s)
	keymap, err := xproto.GetKeyboardMapping(c, min, byte(max-min+1)).Reply()
	if err != nil {
		return nil, UnMappable("keyboard", err)
	}
	return &Keyboard{
		Min: min,
		Max: max,
		GetKeyboardMappingReply: keymap,
	}, nil
}

var IgnoreMods []uint16 = []uint16{
	0,
	xproto.ModMaskLock,                   // Caps lock
	xproto.ModMask2,                      // Num lock
	xproto.ModMaskLock | xproto.ModMask2, // Caps and Num lock
}

func GrabKeyChecked(c *xgb.Conn, win xproto.Window, mods uint16, key xproto.Keycode) error {
	var err error
	for _, m := range IgnoreMods {
		err = xproto.GrabKeyChecked(c, true, win, mods|m, key, xproto.GrabModeAsync, xproto.GrabModeAsync).Check()
		if err != nil {
			return err
		}
	}
	return nil
}

func UngrabKeyChecked(c *xgb.Conn, win xproto.Window, mods uint16, key xproto.Keycode) {
	for _, m := range IgnoreMods {
		xproto.UngrabKeyChecked(c, key, win, mods|m).Check()
	}
}

var pointerMasks uint16 = xproto.EventMaskButtonRelease | xproto.EventMaskButtonPress

func MouseGrabChecked(c *xgb.Conn, win xproto.Window, mods uint16, button xproto.Button) error {
	var err error
	for _, m := range IgnoreMods {
		err = xproto.GrabButtonChecked(
			c,
			true,
			win,
			pointerMasks,
			xproto.GrabModeAsync,
			xproto.GrabModeAsync,
			0,
			0,
			byte(button),
			mods|m,
		).Check()
		if err != nil {
			return err
		}
	}
	return nil
}

func MouseUngrabChecked(c *xgb.Conn, win xproto.Window, mods uint16, button xproto.Button) {
	for _, m := range IgnoreMods {
		xproto.UngrabButtonChecked(c, byte(button), win, mods|m).Check()
	}
}
