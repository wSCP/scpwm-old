package main

import (
	"errors"
	"fmt"
	"strings"
)

type msgResponse []byte

var (
	mSUCCESS msgResponse = []byte("success")
	mSYNTAX  msgResponse = []byte("syntax")
	mUNKNOWN msgResponse = []byte("unknown")
	mLENGTH  msgResponse = []byte("length")
	mFAILURE msgResponse = []byte("failure")
)

var (
	FailureError = Xrror("command failure: %s").Out
	SyntaxError  = Xrror("received message `%s`, but syntax was incorrect: %s").Out
	UnknownError = Xrror("received unknown message of type `%s`: %s").Out
	LengthError  = Xrror("inappropriate length message: `%d : %s`, message must be longer than 1 word and shorter than 15").Out
)

func (e *Euclid) processMsg(msg []byte, comm chan string) msgResponse {
	cmd := NewCommand(e, comm, msg)
	return cmd.process(e)
}

type Command struct {
	e   *Euclid
	com chan string
	pri string
	sub []string
	cmd []string
	raw []byte
}

func NewCommand(e *Euclid, com chan string, cmd []byte) Command {
	raw := cmd
	in := string(cmd)
	spl := strings.Split(in, " ")
	pri, sub, cmds := parseCmd(spl)
	return Command{
		e:   e,
		com: com,
		pri: pri,
		sub: sub,
		cmd: cmds,
		raw: raw,
	}
}

var subcmdTranslate map[string]string = map[string]string{
	"-m":        "monitor",
	"--monitor": "monitor",
	"-d":        "desktop",
	"--desktop": "desktop",
	"-w":        "window",
	"--window":  "window",
}

func parseCmd(s []string) (string, []string, []string) {
	var pri string
	var subcmd []string
	pri = s[0]
	if sc, ok := subcmdTranslate[s[1]]; ok {
		subcmd = append(subcmd, sc, s[2])
		s = s[:1+copy(s[1:], s[2+1:])]
	}
	return pri, subcmd, s[1:]
}

func (c Command) process(e *Euclid) msgResponse {
	switch c.pri {
	case "config":
		return c.config(e)
	case "monitor":
		return c.monitor(e)
	case "desktop":
		return c.desktop(e)
	case "client":
		return c.client(e)
	case "query":
		return c.query(e)
	case "rule":
		return c.rule(e)
	case "pointer":
		return c.pointer(e)
	case "restore":
		return c.restore(e)
	case "control":
		return c.control(e)
	case "quit":
		return c.quit(e)
	default:
		return mUNKNOWN
		c.com <- UnknownError(c.pri, string(c.raw)).Error()
	}
	return nil
}

func (c Command) config(e *Euclid) msgResponse {
	/*mon := e.Monitors.Focused()
	desk := mon.Desktops.Current()
	ref := Coordinates(e, mon, desk, nil)
	//ref := Coordinates(e, mon, desk, desk.focus)
	trg := Coordinates(e, nil, nil, nil)
	if c.subCmd() {
		switch c.sub[0] {
		case "monitor":
			if !MonitorFromDescription(c.cmd, &ref, &trg) {
				return c.respond(mFAILURE, FailureError(c.subCmdErrorString()))
			}
		case "desktop":
			if !DesktopFromDescription(c.cmd, &ref, &trg) {
				return c.respond(mFAILURE, FailureError(c.subCmdErrorString()))
			}
		case "window":
			if !NodeFromDescription(c.cmd, &ref, &trg) {
				return c.respond(mFAILURE, FailureError(c.subCmdErrorString()))
			}
		}
	}
	switch c.length() {
	case 1:
		return c.getSetting(e)
	case 2:
		return c.setSetting(e)
	default:
		return c.respond(mSYNTAX, SyntaxError(c.pri, c.cmdString()))
	}
	return c.respond(mFAILURE, FailureError(c.raw))*/
	return mUNKNOWN
}

func (c Command) monitor(e *Euclid) msgResponse {
	//int cmd_monitor(char **args, int num);
	return mUNKNOWN
}

func (c Command) desktop(e *Euclid) msgResponse {
	//int cmd_desktop(char **args, int num);
	return mUNKNOWN
}

func (c Command) client(e *Euclid) msgResponse {
	//int cmd_window(char **args, int num);
	return mUNKNOWN
}

func (c Command) query(e *Euclid) msgResponse {
	//int cmd_query(char **args, int num, FILE *rsp);
	return mUNKNOWN
}

func (c Command) rule(e *Euclid) msgResponse {
	//int cmd_rule(char **args, int num, FILE *rsp);
	return mUNKNOWN
}

func (c Command) pointer(e *Euclid) msgResponse {
	//int cmd_pointer(char **args, int num);
	return mUNKNOWN
}

func (c Command) restore(e *Euclid) msgResponse {
	//int cmd_restore(char **args, int num);
	return mUNKNOWN
}

func (c Command) control(e *Euclid) msgResponse {
	//int cmd_control(char **args, int num, FILE *rsp);
	return mUNKNOWN
}

func (c Command) quit(e *Euclid) msgResponse {
	//int cmd_quit(char **args, int num);
	return mFAILURE
}

func (c Command) respond(mr msgResponse, err error) msgResponse {
	c.com <- err.Error()
	return mr
}

func (c Command) cmdString() string {
	return strings.Join(c.cmd, " ")
}

func (c Command) subCmd() bool {
	if len(c.sub) > 0 {
		return true
	}
	return false
}

func (c Command) subCmdErrorString() string {
	return fmt.Sprintf("impossible %s, `%s` does not exist", c.pri, c.sub)
}

func (c Command) length() int {
	return len(c.cmd)
}

func (c Command) setSetting(e *Euclid) msgResponse {
	return c.respond(mUNKNOWN, errors.New("SET SETTING UNKNOWN"))
}

func (c Command) getSetting(e *Euclid) msgResponse {
	return c.respond(mUNKNOWN, errors.New("GET SETTING UNKNOWN"))
}

/*
int set_setting(coordinates_t loc, char *name, char *value);
int get_setting(coordinates_t loc, char *name, FILE* rsp);
bool parse_subscriber_mask(char *s, subscriber_mask_t *mask);
bool parse_bool(char *value, bool *b);
bool parse_layout(char *s, layout_t *l);
bool parse_direction(char *s, direction_t *d);
bool parse_cycle_direction(char *s, cycle_dir_t *d);
bool parse_circulate_direction(char *s, circulate_dir_t *d);
bool parse_history_direction(char *s, history_dir_t *d);
bool parse_flip(char *s, flip_t *f);
bool parse_pointer_action(char *s, pointer_action_t *a);
bool parse_child_polarity(char *s, child_polarity_t *p);
bool parse_degree(char *s, int *d);
bool parse_window_id(char *s, long int *i);
bool parse_bool_declaration(char *s, char **key, bool *value, alter_state_t *state);
bool parse_index(char *s, int *i);
*/
