// Protocol the interface of protocol, include three method,
// you can define custom protocol through implement three method.
package protocols

type Protocol interface {
	// Encode encode package before send to client.
	Encode(data []byte) []byte
	// Decode decode package and emit.
	Decode(recvBuffer []byte) interface{}
	// Input check the integrity of package.
	Input(recvBuffer []byte) int
}
