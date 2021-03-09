package util

type Extension struct {
	Value [8]byte
}

func (extension *Extension) Slice() []byte {
	return extension.Value[:]
}
