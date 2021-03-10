package message

type Interested struct {
}

func (interest *Interested) GetMessageID() MsgID {
	return MsgInterested
}

func (interest *Interested) GetPayload() []byte {
	return nil
}

type NotInterested struct {
}

func (disinterest *NotInterested) GetMessageID() MsgID {
	return MsgNotInterested
}

func (disinterest *NotInterested) GetPayload() []byte {
	return nil
}
