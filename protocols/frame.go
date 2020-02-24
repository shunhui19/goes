// frame protocol
package protocols

import (
	"encoding/binary"
)

// Frame frame struct.
type Frame struct {
}

// Input check the integrity of the package.
func (t *Frame) Input(buffer []byte, maxPackageSize int) interface{} {
	if len(buffer) < 4 {
		return 0
	}

	len := int(binary.BigEndian.Uint32(buffer[:4]))
	if len > maxPackageSize {
		return false
	} else {
		return len
	}
}

// Decode decode the buffer.
func (t *Frame) Decode(buffer []byte) []byte {
	r := make([]byte, len(buffer)-4)
	copy(r, buffer[4:])

	return r
}

// Encode encode the buffer, the type of return value is []byte.
func (t *Frame) Encode(buffer []byte) interface{} {
	encodeBuffer := make([]byte, 4+len(buffer))
	copy(encodeBuffer[4:], buffer)
	binary.BigEndian.PutUint32(encodeBuffer, 4+uint32(len(buffer)))

	return encodeBuffer
}

// NewFrameProtocol.
func NewFrameProtocol() *Frame {
	return &Frame{}
}
