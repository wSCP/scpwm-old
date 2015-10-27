package commander

import (
	"errors"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
	Raw() []byte
	Process(Data) (Response, error)
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
	spl := strings.Split(in, " ")
	pri, sel, opts := parseCmd(spl)
	return command{
		pri: pri,
		sel: sel,
		opt: opts,
		raw: raw,
	}
}

//var subcmdTranslate map[string]string = map[string]string{
//	"-m":        "monitor",
//	"--monitor": "monitor",
//	"-d":        "desktop",
//	"--desktop": "desktop",
//	"-w":        "window",
//	"--window":  "window",
//}

func splitOpts(r rune) bool {
	if r == '-' {
		return true
	}
	return false
}

func parseCmd(s []string) (string, string, [][]string) {
	//var subcmd []string
	pri := s[0]
	s = s[:0+copy(s[0:], s[0+1:])]
	fds := strings.FieldsFunc(strings.Join(s, " "), splitOpts)
	spew.Dump(fds)
	//if sc, ok := subcmdTranslate[s[1]]; ok {
	//	subcmd = append(subcmd, sc, s[2])
	//	s = s[:1+copy(s[1:], s[2+1:])]
	//}
	//return pri, subcmd, s[1:]
	return pri, "", nil
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

func (c command) Raw() []byte {
	return c.raw
}

func (c command) Process(d Data) (Response, error) {
	switch c.pri {
	case "config":
		return c.config(d)
	case "monitor":
		return c.monitor(d)
	case "desktop":
		return c.desktop(d)
	case "client":
		return c.client(d)
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
	return mUNKNOWN, UnknownError(c.pri, string(c.raw))
}

func (c command) config(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) monitor(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) desktop(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) client(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) query(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) rule(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) pointer(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) restore(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) control(d Data) (Response, error) {
	return mUNKNOWN, nil
}

func (c command) quit(d Data) (Response, error) {
	return mFAILURE, nil
}

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

func (c command) setSetting(s settings.Settings) (Response, error) {
	return mUNKNOWN, errors.New("SET SETTING UNKNOWN")
}

func (c command) getSetting(s settings.Settings) (Response, error) {
	return mUNKNOWN, errors.New("GET SETTING UNKNOWN")
}
