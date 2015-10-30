package commander

import (
	"errors"
	"strings"

	"github.com/thrisp/scpwm/euclid/settings"
)

var (
	FailureError = Xrror("command failure: %s").Out
	SyntaxError  = Xrror("received message `%s`, but syntax was incorrect: %s").Out
	UnknownError = Xrror("received unknown message of type `%s`: %s").Out
	LengthError  = Xrror("inappropriate length message: `%d : %s`, message must be longer than 1 word and shorter than 15").Out
)

type Command interface {
	Primary() string
	Selector() string
	Opts() [][]string
	Raw() string
	Process(Data) (Response, string, error)
}

type command struct {
	pri string
	sel string
	opt [][]string
	raw []byte
}

func NewCommand(cmd []byte) Command {
	raw := cmd
	in := string(cmd)
	pri, sel, opts := parseCmd(in)
	return command{
		pri: pri,
		sel: sel,
		opt: opts,
		raw: raw,
	}
}

func splitByHyphen(r rune) bool {
	if r == '-' {
		return true
	}
	return false
}

func parseCmd(s string) (string, string, [][]string) {
	fds := strings.FieldsFunc(s, splitByHyphen)
	front := fds[0]
	options := fds[:0+copy(fds[0:], fds[0+1:])]
	pri := strings.Split(front, " ")
	primary := pri[0]
	sel := pri[:0+copy(pri[0:], pri[0+1:])]
	selector := strings.Join(sel, " ")
	var opts [][]string
	for _, v := range options {
		opts = append(opts, strings.Split(v, " "))
	}
	return primary, selector, opts
}

func (c command) Primary() string {
	return c.pri
}

func (c command) Selector() string {
	return c.sel
}

func (c command) Opts() [][]string {
	return c.opt
}

func (c command) Raw() string {
	return string(c.raw)
}

func (c command) Process(d Data) (Response, string, error) {
	switch c.pri {
	case "config":
		return c.config(d)
	case "monitor":
		return c.monitor(d)
	case "desktop":
		return c.desktop(d)
	case "window":
		return c.window(d)
	case "query":
		return c.query(d)
	case "rule":
		return c.rule(d)
	case "pointer":
		return c.pointer(d)
	case "restore":
		return c.restore(d)
	case "control":
		return c.control(d)
	case "quit":
		return c.quit(d)
	}
	return mUNKNOWN, "", UnknownError(c.pri, string(c.raw))
}

func (c command) config(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) monitor(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) desktop(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) window(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) query(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) rule(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) pointer(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) restore(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) control(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) quit(d Data) (Response, string, error) {
	return mUNKNOWN, "", nil
}

func (c command) setSetting(s settings.Settings) (Response, string, error) {
	return mUNKNOWN, "", errors.New("SET SETTING UNKNOWN")
}

func (c command) getSetting(s settings.Settings) (Response, string, error) {
	return mUNKNOWN, "", errors.New("GET SETTING UNKNOWN")
}

//var subcmdTranslate map[string]string = map[string]string{
//	"-m":        "monitor",
//	"--monitor": "monitor",
//	"-d":        "desktop",
//	"--desktop": "desktop",
//	"-w":        "window",
//	"--window":  "window",
//}

//func (c Command) cmdString() string {
//return strings.Join(c.cmd, " ")
//}

//func (c Command) subCmd() bool {
//	if len(c.sub) > 0 {
//		return true
//	}
//	return false
//}

//func (c Command) subCmdErrorString() string {
//	return fmt.Sprintf("impossible %s, `%s` does not exist", c.pri, c.sub)
//}

//func (c Command) length() int {
//	return len(c.cmd)
//}
