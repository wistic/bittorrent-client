package message

type Done struct {
}

func (done *Done) GetMessageID() MsgID {
	return MsgDone
}
func (done *Done) GetPayload() []byte {
	return nil
}
