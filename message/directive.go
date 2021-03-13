package message

import (
	"bittorrent-go/util"
)

type Directive struct {
	Address   string
	MessageID MsgID
	Message   Message
}

func NewDirective(message Message, address *util.Address) *Directive {
	return &Directive{
		Address:   address.String(),
		MessageID: message.GetMessageID(),
		Message:   message,
	}
}
