package message

type ErrorMessage struct {
	Address string
	Value   error
}
