package clients

import (
	"github.com/thrisp/scpwm/euclid/branch"
	"github.com/thrisp/scpwm/euclid/rules"
)

func New() *branch.Branch {
	return branch.New("clients")
}

func Add(cs *branch.Branch, c Client, f Client, csq *rules.Consequence) {
	if c != nil {
		cs.PushFront(c)
	}
}

func Remove(cs *branch.Branch, c Client) error {
	return nil
}

/*
func IsHead(cs *branch.Branch, c Client) bool {
	f := cs.Front()
	fr := f.Value.(Client)
	if fr == c {
		return true
	}
	return false
}

func IsTail(cs *branch.Branch, c Client) bool {
	b := cs.Back()
	ba := b.Value.(Client)
	if ba == c {
		return true
	}
	return false
}
*/
