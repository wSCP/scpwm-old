package commander

import (
	"bytes"
	"net"

	"github.com/thrisp/scpwm/euclid/clients"
	"github.com/thrisp/scpwm/euclid/desktops"
	"github.com/thrisp/scpwm/euclid/handler"
	"github.com/thrisp/scpwm/euclid/monitors"
	"github.com/thrisp/scpwm/euclid/ruler"
	"github.com/thrisp/scpwm/euclid/settings"
)

type Data interface {
	settings.Settings
	handler.Handler
	ruler.Ruler
	Monitors() []monitors.Monitor
	Desktops() []desktops.Desktop
	Clients() []clients.Client
}

type Commander interface {
	Listen(*net.UnixListener, Data)
	Process([]byte, Data) Response
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

func (c *commander) Process(msg []byte, d Data) Response {
	cmd := NewCommand(msg)
	resp, err := cmd.process(d)
	if err != nil {
		c.comm <- err.Error()
	}
	return resp
}

type Response []byte

var (
	mSUCCESS Response = []byte("success")
	mSYNTAX  Response = []byte("syntax")
	mUNKNOWN Response = []byte("unknown")
	mLENGTH  Response = []byte("length")
	mFAILURE Response = []byte("failure")
)
