package message

type MsgID uint8

type Message interface {
	GetMessageID() MsgID
	GetPayload() []byte
}

const (
	// MsgChoke chokes the receiver
	MsgChoke MsgID = 0
	// MsgUnchoke unchokes the receiver
	MsgUnchoke MsgID = 1
	// MsgInterested expresses interest in receiving data
	MsgInterested MsgID = 2
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested MsgID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave MsgID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield MsgID = 5
	// MsgRequest requests a block of data from the receiver
	MsgRequest MsgID = 6
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece MsgID = 7
	// MsgCancel cancels a request
	MsgCancel MsgID = 8
)
