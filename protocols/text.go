// text protocol.
package protocols

import (
	"bytes"
)

// Text text struct.
type Text struct {
}

// Input check the integrity of the package.
func (t *Text) Input(buffer []byte, maxPackageSize int) interface{} {
	if len(buffer) >= maxPackageSize {
		return false
	}

	// find the position of "\n", if not found, continue receive.
	position := bytes.IndexByte(buffer, '\n')
	if position == -1 {
		return 0
	}

	return position + 1
}

// Decode decode the buffer.
func (t *Text) Decode(buffer []byte) []byte {
	return bytes.TrimRight(buffer, "\n")
}

// Encode encode the buffer, the type of return value is []byte.
func (t *Text) Encode(buffer []byte) interface{} {
	return bytes.Join([][]byte{buffer, []byte("\n")}, []byte(""))
}

// NewTextProtocol.
func NewTextProtocol() *Text {
	return &Text{}
}
