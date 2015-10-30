package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"
)

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
