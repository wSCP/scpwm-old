package commander

import (
	"bytes"
	"net"

	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/clients"
	"github.com/thrisp/scpwm/euclid/desktops"
	"github.com/thrisp/scpwm/euclid/handler"
	"github.com/thrisp/scpwm/euclid/monitors"
	"github.com/thrisp/scpwm/euclid/rules"
	"github.com/thrisp/scpwm/euclid/settings"
)

type Commander interface {
	Listen(*net.UnixListener, Data)
	Process([]byte, Data) Response
}

type Response []byte

type Data interface {
	settings.Settings
	handler.Handler
	rules.Ruler
	Tree() *branch.Branch
	Monitors() []monitors.Monitor
	Desktops() []desktops.Desktop
	Clients() []clients.Client
}

type commander struct {
	comm chan string
}

func New(comm chan string) Commander {
	return &commander{comm: comm}
}

func (c *commander) Listen(l *net.UnixListener, d Data) {
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
		resp := c.Process(r, d)
		conn.Write(resp)
		conn.Close()
	}
}

var (
	mSUCCESS Response = []byte("success")
	mSYNTAX  Response = []byte("syntax")
	mUNKNOWN Response = []byte("unknown")
	mLENGTH  Response = []byte("length")
	mFAILURE Response = []byte("failure")
)

func (c *commander) Process(msg []byte, d Data) Response {
	cmd := NewCommand(msg)
	resp, result, err := cmd.Process(d)
	if result != "" {
		c.comm <- result
	}
	if err != nil {
		c.comm <- err.Error()
	}
	return resp
}
