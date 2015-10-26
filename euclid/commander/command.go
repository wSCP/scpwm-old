package commander

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thrisp/scpwm/euclid/clients"
	"github.com/thrisp/scpwm/euclid/desktops"
	"github.com/thrisp/scpwm/euclid/monitors"
	"github.com/thrisp/scpwm/euclid/settings"
)

var (
	FailureError = Xrror("command failure: %s").Out
	SyntaxError  = Xrror("received message `%s`, but syntax was incorrect: %s").Out
	UnknownError = Xrror("received unknown message of type `%s`: %s").Out
	LengthError  = Xrror("inappropriate length message: `%d : %s`, message must be longer than 1 word and shorter than 15").Out
)

type Command struct {
	pri string
	sub []string
	cmd []string
	raw []byte
}

func NewCommand(cmd []byte) Command {
	raw := cmd
	in := string(cmd)
	spl := strings.Split(in, " ")
	pri, sub, cmds := parseCmd(spl)
	return Command{
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

type (
	GetMonitors func() []monitors.Monitor
	GetDesktops func() []desktops.Desktop
	GetClients  func() []clients.Client
)

func (c Command) process(d Data) (Response, error) {
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

func (c Command) config(d Data) (Response, error) {
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
	return mUNKNOWN, nil
}

func (c Command) monitor(d Data) (Response, error) {
	//int cmd_monitor(char **args, int num);
	return mUNKNOWN, nil
}

func (c Command) desktop(d Data) (Response, error) {
	//int cmd_desktop(char **args, int num);
	return mUNKNOWN, nil
}

func (c Command) client(d Data) (Response, error) {
	//int cmd_window(char **args, int num);
	return mUNKNOWN, nil
}

func (c Command) query(d Data) (Response, error) {
	//int cmd_query(char **args, int num, FILE *rsp);
	return mUNKNOWN, nil
}

func (c Command) rule(d Data) (Response, error) {
	//int cmd_rule(char **args, int num, FILE *rsp);
	return mUNKNOWN, nil
}

func (c Command) pointer(d Data) (Response, error) {
	//int cmd_pointer(char **args, int num);
	return mUNKNOWN, nil
}

func (c Command) restore(d Data) (Response, error) {
	//int cmd_restore(char **args, int num);
	return mUNKNOWN, nil
}

func (c Command) control(d Data) (Response, error) {
	//int cmd_control(char **args, int num, FILE *rsp);
	return mUNKNOWN, nil
}

func (c Command) quit(d Data) (Response, error) {
	//int cmd_quit(char **args, int num);
	return mFAILURE, nil
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

func (c Command) setSetting(s settings.Settings) (Response, error) {
	return mUNKNOWN, errors.New("SET SETTING UNKNOWN")
}

func (c Command) getSetting(s settings.Settings) (Response, error) {
	return mUNKNOWN, errors.New("GET SETTING UNKNOWN")
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
