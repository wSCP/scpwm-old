package atomic

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Atomic interface {
	Atom(string) (xproto.Atom, error)
	AtomName(xproto.Atom) (string, error)
	GetProp(xproto.Window, string) (*xproto.GetPropertyReply, error)
	ChangeProp(xproto.Window, byte, string, string, []byte) error
	ChangeProp32(xproto.Window, string, string, ...uint) error
	Converter
	AtomVals
	WindowVals
	NumVals
	StringVals
}

type Converter interface {
	AtomToUint(ids []xproto.Atom) []uint
	StrToAtoms([]string) ([]uint, error)
	WindowToInt([]xproto.Window) []uint
}

type AtomVals interface {
	PropValAtom(*xproto.GetPropertyReply, error) (string, error)
	PropValAtoms(*xproto.GetPropertyReply, error) ([]string, error)
}

type WindowVals interface {
	PropValWindow(*xproto.GetPropertyReply, error) (xproto.Window, error)
	PropValWindows(*xproto.GetPropertyReply, error) ([]xproto.Window, error)
}

type NumVals interface {
	PropValNum(*xproto.GetPropertyReply, error) (uint, error)
	PropValNums(*xproto.GetPropertyReply, error) ([]uint, error)
	PropValNum64(*xproto.GetPropertyReply, error) (int64, error)
}

type StringVals interface {
	PropValStr(*xproto.GetPropertyReply, error) (string, error)
	PropValStrs(*xproto.GetPropertyReply, error) ([]string, error)
}

type atomic struct {
	c            *xgb.Conn
	Atoms        map[string]xproto.Atom
	AtomsLck     *sync.RWMutex
	AtomNames    map[xproto.Atom]string
	AtomNamesLck *sync.RWMutex
}

func New(c *xgb.Conn) Atomic {
	return &atomic{
		c:            c,
		Atoms:        make(map[string]xproto.Atom, 50),
		AtomsLck:     &sync.RWMutex{},
		AtomNames:    make(map[xproto.Atom]string, 50),
		AtomNamesLck: &sync.RWMutex{},
	}
}

var AtomIdError = Xrror("atom: '%s' returned an identifier of 0.").Out

//func (a *atomic) init() {
//	a.Atom("WM_DELETE_WINDOW")
//	a.Atom("WM_TAKE_FOCUS")
//	a.Atom("_BSPWM_FLOATING_WINDOW")
//}

func (a *atomic) Atom(name string) (xproto.Atom, error) {
	aid, err := a.atom(name, false)
	if err != nil {
		return 0, err
	}
	if aid == 0 {
		return 0, AtomIdError(name)
	}

	return aid, err
}

var AtomInternError = Xrror("Atom: Error interning atom '%s': %s").Out

func (a *atomic) atom(name string, onlyIfExists bool) (xproto.Atom, error) {
	if aid, ok := a.atomGet(name); ok {
		return aid, nil
	}

	reply, err := xproto.InternAtom(a.c, onlyIfExists, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, AtomInternError(name, err)
	}

	a.cacheAtom(name, reply.Atom)

	return reply.Atom, nil
}

var AtomNameError = Xrror("AtomName: Error fetching name for ATOM id '%d': %s").Out

func (a *atomic) AtomName(aid xproto.Atom) (string, error) {
	if atomName, ok := a.atomNameGet(aid); ok {
		return string(atomName), nil
	}

	reply, err := xproto.GetAtomName(a.c, aid).Reply()
	if err != nil {
		return "", AtomNameError(aid, err)
	}

	atomName := string(reply.Name)
	a.cacheAtom(atomName, aid)

	return atomName, nil
}

func (a *atomic) atomGet(name string) (xproto.Atom, bool) {
	a.AtomsLck.RLock()
	defer a.AtomsLck.RUnlock()

	aid, ok := a.Atoms[name]
	return aid, ok
}

func (a *atomic) atomNameGet(aid xproto.Atom) (string, bool) {
	name, ok := a.AtomNames[aid]

	a.AtomNamesLck.RLock()
	defer a.AtomNamesLck.RUnlock()

	return name, ok
}

func (a *atomic) cacheAtom(name string, aid xproto.Atom) {
	a.AtomsLck.Lock()
	a.AtomNamesLck.Lock()
	defer a.AtomsLck.Unlock()
	defer a.AtomNamesLck.Unlock()

	a.Atoms[name] = aid
	a.AtomNames[aid] = name
}

var GetPropertyError = Xrror("%s").Out

func (a *atomic) GetProp(w xproto.Window, s string) (*xproto.GetPropertyReply, error) {
	atomId, err := a.Atom(s)
	if err != nil {
		return nil, err
	}

	reply, err := xproto.GetProperty(a.c, false, w, atomId, xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()

	if err != nil {
		msg := fmt.Sprintf("Error retrieving property '%s' on window %x: %s", s, w, err)
		return nil, GetPropertyError(msg)
	}

	if reply.Format == 0 {
		msg := fmt.Sprintf("No such property '%s' on window %x", s, w)
		return nil, GetPropertyError(msg)
	}

	return reply, nil
}

func (a *atomic) ChangeProp(w xproto.Window, format byte, prop, typ string, data []byte) error {
	propAtom, err := a.Atom(prop)
	if err != nil {
		return err
	}

	typAtom, err := a.Atom(typ)
	if err != nil {
		return err
	}

	return xproto.ChangePropertyChecked(a.c, xproto.PropModeReplace, w, propAtom, typAtom, format, uint32(len(data)/(int(format)/8)), data).Check()
}

func (a *atomic) ChangeProp32(w xproto.Window, prop, typ string, data ...uint) error {
	buf := make([]byte, len(data)*4)
	for i, datum := range data {
		xgb.Put32(buf[(i*4):], uint32(datum))
	}

	return a.ChangeProp(w, 32, prop, typ, buf)
}

func (a *atomic) WindowToInt(ids []xproto.Window) []uint {
	ids32 := make([]uint, len(ids))
	for i, v := range ids {
		ids32[i] = uint(v)
	}
	return ids32
}

func (a *atomic) AtomToUint(ids []xproto.Atom) []uint {
	ids32 := make([]uint, len(ids))
	for i, v := range ids {
		ids32[i] = uint(v)
	}
	return ids32
}

func (a *atomic) StrToAtoms(atomNames []string) ([]uint, error) {
	var err error
	atoms := make([]uint, len(atomNames))
	for i, atomName := range atomNames {
		a, err := a.Atom(atomName)
		if err != nil {
			return nil, err
		}
		atoms[i] = uint(a)
	}
	return atoms, err
}

func (a *atomic) PropValAtom(reply *xproto.GetPropertyReply, err error) (string, error) {
	if err != nil {
		return "", err
	}

	if reply.Format != 32 {
		return "", fmt.Errorf("PropValAtom: Expected format 32 but got %d", reply.Format)
	}

	return a.AtomName(xproto.Atom(xgb.Get32(reply.Value)))
}

func (a *atomic) PropValAtoms(reply *xproto.GetPropertyReply, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, fmt.Errorf("PropValAtoms: Expected format 32 but got %d", reply.Format)
	}

	ids := make([]string, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		ids[i], err = a.AtomName(xproto.Atom(xgb.Get32(vals)))
		if err != nil {
			return nil, err
		}

		vals = vals[4:]
	}
	return ids, nil
}

func (a *atomic) PropValWindow(reply *xproto.GetPropertyReply, err error) (xproto.Window, error) {
	if err != nil {
		return 0, err
	}
	if reply.Format != 32 {
		return 0, fmt.Errorf("PropValId: Expected format 32 but got %d", reply.Format)
	}
	return xproto.Window(xgb.Get32(reply.Value)), nil
}

func (a *atomic) PropValWindows(reply *xproto.GetPropertyReply, err error) ([]xproto.Window, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, fmt.Errorf("PropValIds: Expected format 32 but got %d", reply.Format)
	}

	ids := make([]xproto.Window, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		ids[i] = xproto.Window(xgb.Get32(vals))
		vals = vals[4:]
	}
	return ids, nil
}

func (a *atomic) PropValNum(reply *xproto.GetPropertyReply, err error) (uint, error) {
	if err != nil {
		return 0, err
	}
	if reply.Format != 32 {
		return 0, fmt.Errorf("PropValNum: Expected format 32 but got %d",
			reply.Format)
	}
	return uint(xgb.Get32(reply.Value)), nil
}

func (a *atomic) PropValNums(reply *xproto.GetPropertyReply, err error) ([]uint, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, fmt.Errorf("PropValIds: Expected format 32 but got %d",
			reply.Format)
	}

	nums := make([]uint, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		nums[i] = uint(xgb.Get32(vals))
		vals = vals[4:]
	}
	return nums, nil
}

func (a *atomic) PropValNum64(reply *xproto.GetPropertyReply, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	if reply.Format != 32 {
		return 0, fmt.Errorf("PropValNum: Expected format 32 but got %d",
			reply.Format)
	}
	return int64(xgb.Get32(reply.Value)), nil
}

func (a *atomic) PropValStr(reply *xproto.GetPropertyReply, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if reply.Format != 8 {
		return "", fmt.Errorf("PropValStr: Expected format 8 but got %d", reply.Format)
	}
	return string(reply.Value), nil
}

func (a *atomic) PropValStrs(reply *xproto.GetPropertyReply, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 8 {
		return nil, fmt.Errorf("PropValStrs: Expected format 8 but got %d", reply.Format)
	}

	var strs []string
	sstart := 0
	for i, c := range reply.Value {
		if c == 0 {
			strs = append(strs, string(reply.Value[sstart:i]))
			sstart = i + 1
		}
	}
	if sstart < int(reply.ValueLen) {
		strs = append(strs, string(reply.Value[sstart:]))
	}
	return strs, nil
}
