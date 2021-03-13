package message

type Directive struct {
	Address string
	Message Message
}

func NewDirective(message Message, address string) *Directive {
	return &Directive{
		Address: address,
		Message: message,
	}
}
