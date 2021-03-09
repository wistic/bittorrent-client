package util

type Extension struct {
	Value [8]byte
}

func NewExtension(value [8]byte) *Extension {
	return &Extension{Value: value}
}

func DefaultExtension() *Extension {
	return NewExtension([8]byte{})
}
