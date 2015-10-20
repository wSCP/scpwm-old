package commander

import (
	"github.com/thrisp/scpwm/euclid/handler"
	"github.com/thrisp/scpwm/euclid/settings"
)

type Commander interface {
	Process([]byte, *settings.Settings, handler.Handler) cmdrResponse
}

type commander struct {
	comm chan string
}

func New(comm chan string) Commander {
	return &commander{comm: comm}
}

func (c *commander) Process(msg []byte, s *settings.Settings, h handler.Handler) cmdrResponse {
	cmd := NewCommand(msg)
	resp, err := cmd.process(s, h)
	if err != nil {
		c.comm <- err.Error()
	}
	return resp
}

type cmdrResponse []byte

var (
	mSUCCESS cmdrResponse = []byte("success")
	mSYNTAX  cmdrResponse = []byte("syntax")
	mUNKNOWN cmdrResponse = []byte("unknown")
	mLENGTH  cmdrResponse = []byte("length")
	mFAILURE cmdrResponse = []byte("failure")
)
