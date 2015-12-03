package icccm

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func clientEvent(c *xgb.Conn, root, w xproto.Window, a xproto.Atom, data ...interface{}) error {
	evMask := (xproto.EventMaskSubstructureNotify | xproto.EventMaskSubstructureRedirect)
	cm, err := mkClientMessage(32, w, a, data...)
	if err != nil {
		return err
	}

	return xproto.SendEventChecked(c, false, root, uint32(evMask), string(cm.Bytes())).Check()
}

func mkClientMessage(format byte, w xproto.Window, t xproto.Atom, data ...interface{}) (*xproto.ClientMessageEvent, error) {
	// Create the client data list first
	var clientData xproto.ClientMessageDataUnion

	// Don't support formats 8 or 16 yet.
	switch format {
	case 8:
		buf := make([]byte, 20)
		for i := 0; i < 20; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = data[i].(byte)
		}
		clientData = xproto.ClientMessageDataUnionData8New(buf)
	case 16:
		buf := make([]uint16, 10)
		for i := 0; i < 10; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = uint16(data[i].(int16))
		}
		clientData = xproto.ClientMessageDataUnionData16New(buf)
	case 32:
		buf := make([]uint32, 5)
		for i := 0; i < 5; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = uint32(data[i].(int))
		}
		clientData = xproto.ClientMessageDataUnionData32New(buf)
	default:
		return nil, fmt.Errorf("mkClientMessage: Unsupported format '%d'.", format)
	}

	return &xproto.ClientMessageEvent{
		Format: format,
		Window: w,
		Type:   t,
		Data:   clientData,
	}, nil
}
