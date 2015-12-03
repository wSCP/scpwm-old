package icccm

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"

	"github.com/thrisp/scpwm/euclid/atomic"
)

func New(c *xgb.Conn, r xproto.Window, a atomic.Atomic) ICCCM {
	return &icccm{
		conn: c,
		root: r,
		atom: a,
		EWMH: newEwmh(c, r, a),
	}
}

type ICCCM interface {
	WmClassGet(xproto.Window) (*WmClass, error)
	WmClassSet(xproto.Window, *WmClass) error
	WmHintsGet(xproto.Window) (*Hints, error)
	WmNormalHintsGet(xproto.Window) (*NormalHints, error)
	EWMH
}

type icccm struct {
	conn *xgb.Conn
	root xproto.Window
	atom atomic.Atomic
	EWMH
}

var WmClassGetError = Xrror("Two strings make up WM_CLASS -- found %d in '%v'.").Out

type WmClass struct {
	Instance, Class string
}

// WmClassGet returns the class and instance name of a window, as well as any error.
func (i *icccm) WmClassGet(w xproto.Window) (*WmClass, error) {
	raw, err := i.atom.PropValStrs(i.atom.GetProp(w, "WM_CLASS"))
	if err != nil {
		return nil, err
	}
	if len(raw) != 2 {
		return nil, WmClassGetError(len(raw), raw)
	}

	return &WmClass{raw[0], raw[1]}, nil
}

func (i *icccm) WmClassSet(w xproto.Window, class *WmClass) error {
	raw := make([]byte, len(class.Instance)+len(class.Class)+2)
	copy(raw, class.Instance)
	copy(raw[(len(class.Instance)+1):], class.Class)

	return i.atom.ChangeProp(w, 8, "WM_CLASS", "STRING", raw)
}

const (
	HintInput = (1 << iota)
	HintState
	HintIconPixmap
	HintIconWindow
	HintIconPosition
	HintIconMask
	HintWindowGroup
	HintMessage
	HintUrgency
)

const (
	SizeHintUSPosition = (1 << iota)
	SizeHintUSSize
	SizeHintPPosition
	SizeHintPSize
	SizeHintPMinSize
	SizeHintPMaxSize
	SizeHintPResizeInc
	SizeHintPAspect
	SizeHintPBaseSize
	SizeHintPWinGravity
)

type Hints struct {
	Flags                   uint
	Input, InitialState     uint
	IconX, IconY            int
	IconPixmap, IconMask    xproto.Pixmap
	WindowGroup, IconWindow xproto.Window
}

var HintsError = Xrror("There are %d fields in %s, but %d were expected.").Out

func (i *icccm) WmHintsGet(win xproto.Window) (*Hints, error) {
	lenExpect := 9
	raw, err := i.atom.PropValNums(i.atom.GetProp(win, "WM_HINTS"))
	if err != nil {
		return nil, err
	}
	l := len(raw)
	if l != lenExpect {
		return nil, HintsError(l, "WM_HINTS", lenExpect)
	}

	hints := &Hints{}
	hints.Flags = raw[0]
	hints.Input = raw[1]
	hints.InitialState = raw[2]
	hints.IconPixmap = xproto.Pixmap(raw[3])
	hints.IconWindow = xproto.Window(raw[4])
	hints.IconX = int(raw[5])
	hints.IconY = int(raw[6])
	hints.IconMask = xproto.Pixmap(raw[7])
	hints.WindowGroup = xproto.Window(raw[8])

	return hints, nil
}

type NormalHints struct {
	Flags                                                   uint
	X, Y                                                    int
	Width, Height, MinWidth, MinHeight, MaxWidth, MaxHeight uint
	WidthInc, HeightInc                                     uint
	MinAspectNum, MinAspectDen, MaxAspectNum, MaxAspectDen  uint
	BaseWidth, BaseHeight, WinGravity                       uint
}

func (i *icccm) WmNormalHintsGet(win xproto.Window) (*NormalHints, error) {
	lenExpect := 18
	hints, err := i.atom.PropValNums(i.atom.GetProp(win, "WM_NORMAL_HINTS"))
	if err != nil {
		return nil, err
	}
	l := len(hints)
	if l != lenExpect {
		return nil, HintsError(l, "WM_NORMAL_HINTS", lenExpect)
	}

	nh := &NormalHints{}
	nh.Flags = hints[0]
	nh.X = int(hints[1])
	nh.Y = int(hints[2])
	nh.Width = hints[3]
	nh.Height = hints[4]
	nh.MinWidth = hints[5]
	nh.MinHeight = hints[6]
	nh.MaxWidth = hints[7]
	nh.MaxHeight = hints[8]
	nh.WidthInc = hints[9]
	nh.HeightInc = hints[10]
	nh.MinAspectNum = hints[11]
	nh.MinAspectDen = hints[12]
	nh.MaxAspectNum = hints[13]
	nh.MaxAspectDen = hints[14]
	nh.BaseWidth = hints[15]
	nh.BaseHeight = hints[16]
	nh.WinGravity = hints[17]

	if nh.WinGravity <= 0 {
		nh.WinGravity = xproto.GravityNorthWest
	}

	return nh, nil
}
