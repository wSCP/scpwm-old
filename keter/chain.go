package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

type Chain interface {
	Details
	Activity
	Key
	AddCmd(*Cmd)
	Cmd() *Cmd
	Chainable
}

type Details interface {
	Mechanic() int
	Raw() string
	String() string
	Chained() string
}

type Activity interface {
	Clear()
	Active() (int64, bool)
	Activated() bool
}

type Chainable interface {
	Head() Chain
	Tail() Chain
	SetPrev(Chain)
	Prev() Chain
	SetNext(Chain)
	Next() Chain
	Append(Chain)
}

type chain struct {
	mec     int
	expires int64
	prev    Chain
	next    Chain
	khord   string
	chord   string
	cmd     *Cmd
}

func (c *chain) expired() bool {
	if time.Now().Unix() < c.expires {
		return false
	}
	return true
}

func (c *chain) touch() {
	c.expires = time.Now().Add(time.Duration(ChainExpiry) * time.Second).Unix()
}

func clear(c Chain) {
	head := c.Head()
	head.Clear()
	n := head.Next()
	for n != nil {
		n.Clear()
		n = n.Next()
	}
}

func (c *chain) Clear() {
	c.expires = 0
}

func (c *chain) Active() (int64, bool) {
	if !c.expired() {
		return c.expires, true
	}
	return c.expires, false
}

func (c *chain) Activated() bool {
	head := c.Head()
	if last, atv := head.Active(); atv {
		for head.Next() != nil {
			head = head.Next()
			curr, atv := head.Active()
			if !atv {
				return false
			}
			if curr < last {
				return false
			}
			last = curr
		}
		return true
	}
	return false
}

func (c *chain) Head() Chain {
	var head Chain
	head = c
	for head.Prev() != nil {
		head = head.Prev()
	}
	return head
}

func (c *chain) Tail() Chain {
	var tail Chain
	tail = c
	for tail.Next() != nil {
		tail = tail.Next()
	}
	return tail
}

func (c *chain) SetPrev(o Chain) {
	c.prev = o
}

func (c *chain) Prev() Chain {
	return c.prev
}

func (c *chain) SetNext(o Chain) {
	o.SetPrev(c)
	c.next = o
}

func (c *chain) Next() Chain {
	return c.next
}

func (c *chain) Append(o Chain) {
	c.Tail().SetNext(o)
}

func (c *chain) Cmd() *Cmd {
	return c.cmd
}

func (c *chain) AddCmd(cmd *Cmd) {
	c.cmd = cmd
}

func attach(c Chain, h XHandle, w xproto.Window) error {
	switch c.Mechanic() {
	case KeyPress, KeyRelease:
		mods, keycodes, err := ParseKeyInput(h.Keyboard(), c.String())
		if err != nil {
			return err
		}
		for _, kc := range keycodes {
			GrabKeyChecked(h.Conn(), w, mods, kc)
			ky := mkInput(c.Mechanic(), w, mods, byte(kc), 0)
			h.Put(ky, c)
		}
	case ButtonPress, ButtonRelease:
		mods, button, err := ParseMouseInput(c.String())
		if err != nil {
			return err
		}
		MouseGrabChecked(h.Conn(), w, mods, button)
		ky := mkInput(c.Mechanic(), w, mods, 0, byte(button))
		h.Put(ky, c)
	}
	return nil
}

func (c *chain) Attach(h XHandle, w xproto.Window) error {
	var err error
	err = attach(c, h, w)
	if err != nil {
		return err
	}
	if c.Next() != nil {
		n := c.Next()
		for n != nil {
			err = attach(n, h, w)
			if err != nil {
				return err
			}
			n = n.Next()
		}
	}
	return nil
}

func (c *chain) Run(h XHandle, param string) error {
	c.touch()
	if c.Activated() {
		//spew.Dump(c.Head().Chained())
		cmd := c.Cmd()
		return cmd.Exec(param).Run()
	}
	return nil
}

func (c *chain) Mechanic() int {
	return c.mec
}

func (c *chain) Raw() string {
	return c.khord
}

func (c *chain) String() string {
	return c.chord
}

func (c *chain) Chained() string {
	cmdstring := func(c Chain, b *bytes.Buffer) {
		if cmd := c.Cmd(); cmd != nil {
			b.Write(cmd.Bytes())
		}
	}
	ret := new(bytes.Buffer)
	ret.WriteString(c.khord)
	cmdstring(c, ret)
	n := c.Next()
	for n != nil {
		ret.WriteString("; ")
		ret.WriteString(n.String())
		cmdstring(n, ret)
		n = n.Next()
	}
	return ret.String()
}

func newChain(in []byte) *chain {
	chord, khord, mechanic := chordMechanic(in)
	return &chain{
		mec:   mechanic,
		khord: khord,
		chord: chord,
	}
}

func chordMechanic(in []byte) (string, string, int) {
	var chord, khord string
	var mec int
	var isRelease, isButton bool
	if spl := bytes.Split(in, []byte("+")); spl[len(spl)-1][0] == RELEASE {
		isRelease = true
	}
	if bytes.Contains(in, []byte("button")) {
		isButton = true
	}
	cut := func(r rune) rune {
		if r == '@' {
			return -1
		}
		return r
	}
	switch {
	case !isRelease && !isButton:
		s := string(in)
		chord = s
		khord = s
		mec = KeyPress
	case isRelease && !isButton:
		chord = string(bytes.Map(cut, in))
		khord = string(in)
		mec = KeyRelease
	case !isRelease && isButton:
		s := string(in)
		chord = s
		khord = s
		mec = ButtonPress
	case isRelease && isButton:
		chord = string(bytes.Map(cut, in))
		khord = string(in)
		mec = ButtonRelease
	default:
		s := string(in)
		chord = s
		khord = s
		mec = KeyPress
	}
	return chord, khord, mec
}

func mkChain(in [][]byte) Chain {
	head := newChain(in[0])
	curr := head
	for _, v := range in[1:] {
		n := newChain(v)
		n.prev = curr
		curr.next = n
		curr = n
	}
	return head
}

func mkChains(in [][][]byte) []Chain {
	var ret []Chain
	for _, v := range in {
		ret = append(ret, mkChain(v))
	}
	return ret
}

type Cmd struct {
	raw  string
	t    *template.Template
	s    []*state
	next *Cmd
}

func newCmd(in []byte) *Cmd {
	t, raw, s := parseCmd(in)
	return &Cmd{
		raw: raw,
		t:   t,
		s:   s,
	}
}

func parseCmd(in []byte) (*template.Template, string, []*state) {
	var sts [][]byte
	var raw string
	r := rxMulti.ReplaceAllFunc(
		in,
		func(i []byte) []byte {
			sts = append(sts, i)
			return []byte(fmt.Sprintf("{{.var%d}}", len(sts)))
		})
	r2 := rxRangeInt.ReplaceAllFunc(
		r,
		func(i []byte) []byte {
			return []byte("{{.p}}")
		})
	raw = string(r2)
	tpl, err := template.New("cmd").Parse(raw)
	if err != nil {
		raw = err.Error()
		tpl, _ = template.New("cmd").Parse(raw)
	}
	states := newStates(sts)
	return tpl, raw, states
}

func (c *Cmd) Bytes() []byte {
	ret := new(bytes.Buffer)
	ret.WriteString("; (")
	ret.WriteString(c.raw)
	n := c.next
	for n != nil {
		ret.WriteString("; ")
		ret.WriteString(n.raw)
		n = n.next
	}
	ret.WriteString(")")
	return ret.Bytes()
}

func (c *Cmd) String() string {
	return string(c.Bytes())
}

func (c *Cmd) data(param string) map[string]string {
	ret := make(map[string]string)
	ret["p"] = param
	for i, v := range c.s {
		ret[fmt.Sprintf("var%d", i+1)] = v.Next()
	}
	return ret
}

func (c *Cmd) write(param string) string {
	b := new(bytes.Buffer)
	err := c.t.Execute(b, c.data(param))
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func (c *Cmd) Exec(param string) *exec.Cmd {
	spl := strings.Split(c.write(param), " ")
	return exec.Command(spl[0], spl[1:]...)
}

type state struct {
	index  int
	values []string
}

func newStates(in [][]byte) []*state {
	var ret []*state
	for _, v := range in {
		pre := bytes.Trim(v, "}{")
		post := bytes.Split(pre, []byte(","))
		var s []string
		for _, o := range post {
			s = append(s, string(o))
		}
		ret = append(ret, newState(s))
	}
	return ret
}

func newState(v []string) *state {
	return &state{
		index:  -1,
		values: v,
	}
}

func (s *state) incr() {
	if (s.index + 1) == len(s.values) {
		s.index = 0
	} else {
		s.index++
	}
}

func (s *state) Next() string {
	s.incr()
	return s.values[s.index]
}
