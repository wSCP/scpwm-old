package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Handlr interface {
	Conn() *xgb.Conn
	Root() xproto.Window
	Keyboard() *Keyboard
	Storr
	Eventr
}

type Storr interface {
	Put(Input, Keyable)
	Get(interface{}, int, xproto.Window) ([]Keyable, string, error)
}

type Eventr interface {
	Enqueue(evt xgb.Event, err xgb.Error)
	Dequeue() (xgb.Event, xgb.Error)
	Empty() bool
	Quit()
	Quitting() bool
}

type event struct {
	evt xgb.Event
	err xgb.Error
}

type handlr struct {
	conn     *xgb.Conn
	root     xproto.Window
	keyboard *Keyboard
	events   []event
	evtsLck  *sync.RWMutex
	keys     map[int]map[xproto.Window]map[Input][]Keyable
	keysLck  *sync.RWMutex
	quit     bool
}

func (h *handlr) Conn() *xgb.Conn {
	return h.conn
}

func (h *handlr) Root() xproto.Window {
	return h.root
}

func (h *handlr) Keyboard() *Keyboard {
	return h.keyboard
}

func (h *handlr) Enqueue(evt xgb.Event, err xgb.Error) {
	h.evtsLck.Lock()
	defer h.evtsLck.Unlock()

	h.events = append(h.events, event{
		evt: evt,
		err: err,
	})
}

func (h *handlr) Dequeue() (xgb.Event, xgb.Error) {
	h.evtsLck.Lock()
	defer h.evtsLck.Unlock()

	e := h.events[0]
	h.events = h.events[1:]
	return e.evt, e.err
}

func (h *handlr) Empty() bool {
	h.evtsLck.Lock()
	defer h.evtsLck.Unlock()

	return len(h.events) == 0
}

func (h *handlr) Put(k Input, fn Keyable) {
	h.keysLck.Lock()
	defer h.keysLck.Unlock()

	if _, ok := h.keys[k.evt]; !ok {
		h.keys[k.evt] = make(map[xproto.Window]map[Input][]Keyable)
	}
	if _, ok := h.keys[k.evt][k.win]; !ok {
		h.keys[k.evt][k.win] = make(map[Input][]Keyable)
	}

	h.keys[k.evt][k.win][k] = append(h.keys[k.evt][k.win][k], fn)
}

var ParseEventError = Krror("event %+v should not have been passed to parseEvent").Out

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
		err = ParseEventError(e)
	}
	return s, d, b, err
}

var NoKey = Krror("there is no corresponding key for %+v").Out

func (h *handlr) Get(e interface{}, evtype int, win xproto.Window) ([]Keyable, string, error) {
	h.keysLck.RLock()
	defer h.keysLck.RUnlock()

	mods, key, button, err := parseEvent(e)
	if err != nil {
		return nil, "", err
	}

	if k, ok := h.keys[evtype][win][Input{evtype, win, mods, key, button}]; ok {
		return k, ByteToString(h.Keyboard(), key, button), nil
	}

	return nil, "", NoKey(e)
}

func (h *handlr) Quit() {
	h.quit = true
}

func (h *handlr) Quitting() bool {
	return h.quit
}

func NewHandlr(display string) (Handlr, error) {
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

	h := &handlr{
		conn:     c,
		root:     screen.Root,
		events:   make([]event, 0, 1000),
		evtsLck:  &sync.RWMutex{},
		keyboard: kb,
		keys:     make(map[int]map[xproto.Window]map[Input][]Keyable),
		keysLck:  &sync.RWMutex{},
	}

	return h, nil
}

func signals() chan os.Signal {
	s := make(chan os.Signal, 0)
	signal.Notify(
		s,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGCHLD,
		syscall.SIGPIPE,
	)
	return s
}

func Loop(h Handlr) (chan struct{}, chan struct{}, chan struct{}, chan os.Signal) {
	before := make(chan struct{}, 0)
	after := make(chan struct{}, 0)
	quit := make(chan struct{}, 0)
	signals := signals()
	go func() {
		loop(h, before, after, quit)
	}()
	return before, after, quit, signals
}

func loop(h Handlr, before, after, quit chan struct{}) {
	for {
		if h.Quitting() {
			if quit != nil {
				quit <- struct{}{}
			}
			break
		}

		read(h, true)

		process(h, before, after)
	}
}

func read(h Handlr, block bool) {
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

var EventError = Krror("event error: %s").Out

func runKey(h Handlr, k Keyable, param string) {
	go k.Run(h, param)
}

func process(h Handlr, before, after chan struct{}) {
	for !h.Empty() {
		if h.Quitting() {
			return
		}

		if before != nil && after != nil {
			before <- struct{}{}
		}

		ev, err := h.Dequeue()

		if err != nil {
			Logger.Println(EventError(err.Error()))
			if before != nil && after != nil {
				after <- struct{}{}
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

		if before != nil && after != nil {
			after <- struct{}{}
		}
	}
}

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
