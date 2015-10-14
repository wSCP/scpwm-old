package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"unicode"
)

const (
	RELEASE      = byte('@')
	SPACE        = byte(' ')
	COMMENT      = byte('#')
	LINECONTINUE = byte('\\')
	SEMICOLON    = byte(';')
	COMMA        = byte(',')
	DASH         = byte('-')
	//SEQ_NONE  = byte('_')
)

func Configure(chains []Chain) (XHandle, error) {
	X, err := NewXHandle("")
	if err != nil {
		return nil, err
	}

	for _, c := range chains {
		err = c.Attach(X, X.Root())
	}
	if err != nil {
		return nil, err
	}

	return X, nil
}

func LoadConfig(f string) ([]Chain, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	chains, err := parseConfig(reader)
	if err != nil {
		return nil, err
	}
	return chains, nil
}

var ParseError = Krror("error parsing configuration: %s").Out

func parseConfig(r *bufio.Reader) ([]Chain, error) {
	var err error
	o, err := parseOrders(r)
	if err != nil {
		return nil, err
	}
	cs, err := parseChains(o)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func parseOrders(r *bufio.Reader) ([]*order, error) {
	lineno := 0
	var keys, actions []byte
	var o []*order
	var err error
	for err == nil {
		l, _, err := r.ReadLine()
		if err != nil {
			break
		}
		lineno++
		if len(l) == 0 || l[0] == COMMENT {
			continue
		}
		for l[len(l)-1] == LINECONTINUE {
			nl, _, err := r.ReadLine()
			if err != nil {
				return nil, ParseError(fmt.Sprintf("line %d -- %s", lineno, err.Error()))
			}
			l = append(l, bytes.TrimFunc(nl, unicode.IsSpace)...)
		}
		if l[0] != SPACE {
			keys = bytes.Join(bytes.Split(l, []byte{SPACE}), []byte(""))
		} else {
			actions = bytes.TrimFunc(l, unicode.IsSpace)
		}
		if keys != nil && actions != nil {
			o = append(o, neworder(keys, actions))
			keys = nil
			actions = nil
		}
	}
	return o, nil
}

type order struct {
	keys    *chunk
	actions *chunk
}

func neworder(k, a []byte) *order {
	return &order{
		keys:    extractChunks(k, SEMICOLON),
		actions: extractChunks(a, SEMICOLON),
	}
}

func (o *order) chainKeys() []Chain {
	var ll [][][]byte
	kys := o.keys
	if kys != nil {
		if kys.next == nil {
			ll = append(ll, kys.split)
		} else {
			for kys.next != nil {
				ll = append(ll, kys.split)
				kys = kys.next
			}
		}
	}
	return mkChains(combine(ll))
}

func (o *order) chainCmds() *Cmd {
	var ret *Cmd
	c := o.actions
	if c.nxtlen() == 0 {
		ret = newCmd(c.raw.Bytes())
	} else {
		head := newCmd(c.raw.Bytes())
		curr := head
		for c.next != nil {
			n := newCmd(c.raw.Bytes())
			curr.next = n
			curr = n
			c = c.next
		}
		ret = head
	}
	return ret
}

var MalformedOrder = Krror("malformed order:\n%+v\n, must contain both keys & action").Out

func parseChains(o []*order) ([]Chain, error) {
	var ret []Chain
	for _, od := range o {
		if od.keys == nil || od.actions == nil {
			return nil, MalformedOrder(od)
		}
		chordKeys(od.keys)
		cns := od.chainKeys()
		cc := od.chainCmds()
		for _, cn := range cns {
			cn.Tail().AddCmd(cc)
		}
		ret = append(ret, cns...)
	}
	return ret, nil
}

type chunk struct {
	raw   *bytes.Buffer
	split [][]byte
	next  *chunk
}

func (c *chunk) nxtlen() int {
	return c.next.raw.Len()
}

func newchunk() *chunk {
	return &chunk{raw: new(bytes.Buffer)}
}

func extractChunks(in []byte, spr byte) *chunk {
	c := newchunk()
	head := c
	for _, ck := range bytes.Split(in, []byte{spr}) {
		c.raw.Write(ck)
		c.next = newchunk()
		c = c.next
	}
	return head
}

func chordKeys(c *chunk) {
	c.split = sequenceKeys(c.raw.Bytes())
	n := c.next
	if n != nil {
		n.split = sequenceKeys(n.raw.Bytes())
		n = n.next
	}
}

func collapseChords(in [][][]byte) [][]byte {
	var ret [][]byte
	for _, c := range in {
		b := new(bytes.Buffer)
		for _, cc := range c {
			b.Write(cc)
		}
		ret = append(ret, b.Bytes())
	}
	return ret
}

func extractSeq(r rune) bool {
	if r == '{' || r == '}' {
		return true
	}
	return false
}

func expandSeq(in []byte) [][]byte {
	var ret [][]byte
	switch seqType(wrapUnknownSeq(in)) {
	case multiSequence:
		m, err := multiSeq(in)
		if err == nil {
			ret = m
		}
	case intrangeSequence:
		r, err := rangeSeq(in)
		if err == nil {
			ret = r
		}
	default:
		ret = append(ret, in)
	}
	return ret
}

var NotaSequence = Krror("is not a %s sequence: %s").Out

func multiSeq(in []byte) ([][]byte, error) {
	m := bytes.Split(in, []byte{COMMA})
	if len(m) > 1 {
		return m, nil
	}
	return nil, NotaSequence("multiple item (e.g. a,b,c,d) range", in)
}

func rangeSeq(in []byte) ([][]byte, error) {
	s := bytes.Split(in, []byte{DASH})

	if len(s) != 2 {
		return nil, NotaSequence("not an integer range", in)
	}

	var ret [][]byte

	initial, err := strconv.ParseInt(string(s[0]), 10, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseInt(string(s[1]), 10, 64)
	if err != nil {
		return nil, err
	}

	if initial < last {
		for i := initial; i <= last; i++ {
			ret = append(ret, strconv.AppendInt([]byte{}, i, 10))
		}
		return ret, nil
	}

	return nil, NotaSequence("not an integer range", in)
}

func expandSeqs(in [][]byte) [][][]byte {
	var ret [][][]byte
	for _, seq := range in {
		ret = append(ret, expandSeq(seq))
	}
	return ret
}

func sequenceKeys(in []byte) [][]byte {
	var ret [][]byte
	s := bytes.FieldsFunc(in, extractSeq)
	if len(s) > 1 {
		ret = collapseChords(combine(expandSeqs(s)))
	} else {
		ret = s
	}
	return ret
}

func nextIndex(ix []int, lens func(i int) int) {
	for j := len(ix) - 1; j >= 0; j-- {
		ix[j]++
		if j == 0 || ix[j] < lens(j) {
			return
		}
		ix[j] = 0
	}
}

func combine(in [][][]byte) [][][]byte {
	var ret [][][]byte
	lens := func(i int) int { return len(in[i]) }

	for ix := make([]int, len(in)); ix[0] < lens(0); nextIndex(ix, lens) {
		var r [][]byte
		for j, k := range ix {
			r = append(r, in[j][k])
		}
		ret = append(ret, r)
	}
	return ret
}

var (
	rxMulti    = regexp.MustCompile(`\B\{(.*[\w-]+)(.*,)(.*[\w-]+)(?:,[\w-]+=[\w-]+)*\}`) // {a,b,c,d}
	rxRangeInt = regexp.MustCompile(`\B\{([\d]?)-([\d]?)(?:,[\d]?=[\d]?)*\}`)             // {1-9}
)

type sequenceType int

const (
	noSequence sequenceType = iota
	multiSequence
	intrangeSequence
)

func wrapUnknownSeq(b []byte) []byte {
	id := new(bytes.Buffer)
	id.WriteByte('{')
	id.Write(b)
	id.WriteByte('}')
	return id.Bytes()
}

func seqType(b []byte) sequenceType {
	if rxMulti.Match(b) {
		return multiSequence
	}
	if rxRangeInt.Match(b) {
		return intrangeSequence
	}
	return noSequence
}
