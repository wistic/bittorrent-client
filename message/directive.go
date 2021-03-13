package message

import (
	"bittorrent-go/util"
)

type Directive struct {
	Address   *util.Address
	MessageID MsgID
	Message   Message
}

func NewDirective(message Message, address *util.Address) *Directive {
	return &Directive{
		Address:   address,
		MessageID: message.GetMessageID(),
		Message:   message,
	}
}
