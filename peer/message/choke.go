package message

type Choke struct {
}

func (choke *Choke) GetMessageID() MsgID {
	return MsgChoke
}

func (choke *Choke) GetPayload() []byte {
	return nil
}

type Unchoke struct {
}

func (unchoke *Unchoke) GetMessageID() MsgID {
	return MsgUnchoke
}

func (unchoke *Unchoke) GetPayload() []byte {
	return nil
}
