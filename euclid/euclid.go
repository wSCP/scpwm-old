package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
)

type Euclid struct {
	XHandle
	*Settings
	*loop
	monitors Monitors
	pending  []Rule
	//history  *History
	//rules   *Rule
	//pending *Pending
}

func New() *Euclid {
	e := &Euclid{
		Settings: DefaultSettings(),
	}
	hndl, err := NewXHandle("", ewmhSupported)
	if err != nil {
		panic(err)
	}
	e.XHandle = hndl
	e.monitors = NewMonitors(e)
	e.pending = make([]Rule, 0)
	return e
}

func defaultSocketPath() string {
	var socketPath = os.Getenv(socketEnv)
	if socketPath == "" {
		socketPath = fmt.Sprintf(socketPathTpl, "", 0, 0)
	}
	return socketPath
}

func (e *Euclid) lstn(l *net.UnixListener, loop *loop) {
	for {
		conn, err := l.AcceptUnix()
		if err != nil {
			panic(err)
		}
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		r := bytes.Trim(buf[:n], " ")
		resp := e.processMsg(r, loop.comm)
		conn.Write(resp)
		conn.Close()
	}
}

func (e *Euclid) Loop(l *net.UnixListener) *loop {
	pre := make(chan struct{}, 0)
	post := make(chan struct{}, 0)
	quit := make(chan struct{}, 0)
	com := make(chan string, 0)
	e.loop = &loop{pre, post, quit, com}
	go func() {
		e.lstn(l, e.loop)
	}()
	go func() {
		e.Evt(e.loop.pre, e.loop.post, e.loop.quit)
	}()
	return e.loop
}

type loop struct {
	pre  chan struct{}
	post chan struct{}
	quit chan struct{}
	comm chan string
}

func (e *Euclid) ToggleVisibility() {
	//void toggle_visibility(void)
	//visible = !visible;
	//if (!visible)
	//	clear_input_focus();
	//for (monitor_t *m = mon_head; m != NULL; m = m->next)
	//	for (node_t *n = first_extrema(m->desk->root); n != NULL; n = next_leaf(n, m->desk->root))
	//		window_set_visibility(n->client->window, visible);
	//if (visible)
	//	update_input_focus();
}

func main() {
	e := New()
	sckt, err := net.ListenUnix("unix", &net.UnixAddr{socketPath, "unix"})
	if err != nil {
		panic(err)
	}
	defer os.Remove(socketPath)

	defer e.Conn().Close()

	l := e.Loop(sckt)

	err = e.LoadConfig(ConfigPath)

	if err != nil {
		e.Println(err.Error())
	}

EVT:
	for {
		select {
		case <-l.pre:
			<-l.post
		case msg := <-l.comm:
			if verbose {
				e.Println(msg)
			}
		case <-l.quit:
			break EVT
		}
	}

	e.Println("EXITING......\n")
}

func init() {
	flag.StringVar(&ConfigPath, "config", defaultConfigPath(), "Reads the main configuration from the given file. The default location is 'XDG_CONFIG_HOME/euclid/euclidrc'")
	flag.StringVar(&socketEnv, "e", socketEnv, "Reads the socket from the given env variable, scpwm will attempt to read from 'SCPWM_SOCKET' by default")
	flag.StringVar(&socketPath, "p", defaultSocketPath(), "Reads the socket from the given path. Default is '/tmp/scpwm_0_0-socket'")
	flag.BoolVar(&verbose, "v", verbose, "Verbose logging messages, default is false.")
	flag.Parse()
}
