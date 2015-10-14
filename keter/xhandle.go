package main

import (
	"errors"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Event struct {
	evt xgb.Event
	err xgb.Error
}

type Eventr interface {
	Enqueue(evt xgb.Event, err xgb.Error)
	Dequeue() (xgb.Event, xgb.Error)
	Empty() bool
	Quit()
	Quitting() bool
}

type Storer interface {
	Put(Input, Key)
	Get(interface{}, int, xproto.Window) ([]Key, string, error)
}

type XHandle interface {
	Conn() *xgb.Conn
	Root() xproto.Window
	Keyboard() *Keyboard
	Storer
	Eventr
}

type xhandle struct {
	conn     *xgb.Conn
	root     xproto.Window
	keyboard *Keyboard
	Events   []Event
	EvtsLck  *sync.RWMutex
	Keys     map[int]map[xproto.Window]map[Input][]Key
	KeysLck  *sync.RWMutex
	quit     bool
}

func (h *xhandle) Conn() *xgb.Conn {
	return h.conn
}

func (h *xhandle) Root() xproto.Window {
	return h.root
}

func (h *xhandle) Keyboard() *Keyboard {
	return h.keyboard
}

func (h *xhandle) Enqueue(evt xgb.Event, err xgb.Error) {
	h.EvtsLck.Lock()
	defer h.EvtsLck.Unlock()

	h.Events = append(h.Events, Event{
		evt: evt,
		err: err,
	})
}

func (h *xhandle) Dequeue() (xgb.Event, xgb.Error) {
	h.EvtsLck.Lock()
	defer h.EvtsLck.Unlock()

	e := h.Events[0]
	h.Events = h.Events[1:]
	return e.evt, e.err
}

func (h *xhandle) Empty() bool {
	h.EvtsLck.Lock()
	defer h.EvtsLck.Unlock()

	return len(h.Events) == 0
}

func (h *xhandle) Put(k Input, fn Key) {
	h.KeysLck.Lock()
	defer h.KeysLck.Unlock()

	if _, ok := h.Keys[k.evt]; !ok {
		h.Keys[k.evt] = make(map[xproto.Window]map[Input][]Key)
	}
	if _, ok := h.Keys[k.evt][k.win]; !ok {
		h.Keys[k.evt][k.win] = make(map[Input][]Key)
	}

	h.Keys[k.evt][k.win][k] = append(h.Keys[k.evt][k.win][k], fn)
}

func parseEvent(e interface{}) (uint16, byte, byte, error) {
	var s uint16
	var d, b byte
	var err error
	switch evt := e.(type) {
	case xproto.KeyPressEvent:
		s = evt.State
		d = byte(evt.Detail)
		err = nil
	case xproto.KeyReleaseEvent:
		s = evt.State
		d = byte(evt.Detail)
		err = nil
	case xproto.ButtonPressEvent:
		s = evt.State
		b = byte(evt.Detail)
		err = nil
	case xproto.ButtonReleaseEvent:
		s = evt.State
		b = byte(evt.Detail)
		err = nil
	default:
		err = errors.New("event %+v should not have been passed to parseEvent")
	}
	return s, d, b, err
}

var NoKey = Krror("there is no corresponding key for %+v").Out

func (h *xhandle) Get(e interface{}, evtype int, win xproto.Window) ([]Key, string, error) {
	h.KeysLck.RLock()
	defer h.KeysLck.RUnlock()

	mods, key, button, err := parseEvent(e)
	if err != nil {
		return nil, "", err
	}

	if k, ok := h.Keys[evtype][win][Input{evtype, win, mods, key, button}]; ok {
		return k, ByteToString(h.Keyboard(), key, button), nil
	}

	return nil, "", NoKey(e)
}

func (h *xhandle) Quit() {
	h.quit = true
}

func (h *xhandle) Quitting() bool {
	return h.quit
}

func NewXHandle(display string) (*xhandle, error) {
	c, err := xgb.NewConnDisplay(display)
	if err != nil {
		return nil, err
	}

	s := xproto.Setup(c)
	screen := s.DefaultScreen(c)

	kb, err := NewKeyboard(s, c)
	if err != nil {
		return nil, err
	}

	h := &xhandle{
		conn:     c,
		root:     screen.Root,
		Events:   make([]Event, 0, 1000),
		EvtsLck:  &sync.RWMutex{},
		keyboard: kb,
		Keys:     make(map[int]map[xproto.Window]map[Input][]Key),
		KeysLck:  &sync.RWMutex{},
	}

	return h, nil
}

func Main(h XHandle) (chan struct{}, chan struct{}, chan struct{}) {
	pingBefore := make(chan struct{}, 0)
	pingAfter := make(chan struct{}, 0)
	pingQuit := make(chan struct{}, 0)
	go func() {
		evtLoop(h, pingBefore, pingAfter, pingQuit)
	}()
	return pingBefore, pingAfter, pingQuit
}

func evtLoop(h XHandle, pingBefore, pingAfter, pingQuit chan struct{}) {
	for {
		if h.Quitting() {
			if pingQuit != nil {
				pingQuit <- struct{}{}
			}
			break
		}

		Read(h, true)

		processEvts(h, pingBefore, pingAfter)
	}
}

var EventError = Krror("event error: %s").Out

func runKey(h XHandle, k Key, param string) {
	go k.Run(h, param)
}

func processEvts(h XHandle, pingBefore, pingAfter chan struct{}) {
	for !h.Empty() {
		if h.Quitting() {
			return
		}

		if pingBefore != nil && pingAfter != nil {
			pingBefore <- struct{}{}
		}

		ev, err := h.Dequeue()

		if err != nil {
			Logger.Println(EventError(err.Error()))
			if pingBefore != nil && pingAfter != nil {
				pingAfter <- struct{}{}
			}
			continue
		}

		if ev == nil {
			Logger.Fatal("BUG: Expected an event but got nil.")
		}

		switch event := ev.(type) {
		case xproto.KeyPressEvent:
			if keys, param, err := h.Get(event, KeyPress, event.Event); err == nil {
				for _, k := range keys {
					runKey(h, k, param)
				}
			}
		case xproto.KeyReleaseEvent:
			if keys, param, err := h.Get(event, KeyRelease, event.Event); err == nil {
				for _, k := range keys {
					runKey(h, k, param)
				}
			}
		case xproto.ButtonPressEvent:
			if keys, param, err := h.Get(event, ButtonPress, event.Event); err == nil {
				for _, k := range keys {
					runKey(h, k, param)
				}
			}
		case xproto.ButtonReleaseEvent:
			if keys, param, err := h.Get(event, ButtonPress, event.Event); err == nil {
				for _, k := range keys {
					runKey(h, k, param)
				}
			}
		}

		if pingBefore != nil && pingAfter != nil {
			pingAfter <- struct{}{}
		}
	}
}

func Read(h XHandle, block bool) {
	if block {
		ev, err := h.Conn().WaitForEvent()
		if ev == nil && err == nil {
			Logger.Fatal("BUG: Could not read an event or an error.")
		}
		h.Enqueue(ev, err)
	}

	for {
		ev, err := h.Conn().PollForEvent()

		if ev == nil && err == nil {
			break
		}

		h.Enqueue(ev, err)
	}
}

const (
	NoMechanic    = 1
	KeyPress      = xproto.KeyPress
	KeyRelease    = xproto.KeyRelease
	ButtonPress   = xproto.ButtonPress
	ButtonRelease = xproto.ButtonRelease
)

type Input struct {
	evt int
	win xproto.Window
	mod uint16
	cod byte
	but byte
}

func mkInput(e int, w xproto.Window, m uint16, c byte, b byte) Input {
	return Input{e, w, m, c, b}
}

type Key interface {
	Attach(XHandle, xproto.Window) error
	Run(XHandle, string) error
}
