// Protocol the interface of protocol, include three method,
// you can define custom protocol through implement three method.
package protocols

type Protocol interface {
	// Encode encode package before send to client.
	// The type of return is different for each protocol.
	Encode(data []byte) interface{}
	// Decode decode package and emit.
	Decode(recvBuffer []byte) []byte
	// Input check the integrity of package.
	// if return the value of bool, close connection, indicates that the package is greater than MaxPackageSize,
	// else if return 0, the package is not a integrity package, continue to receive.
	Input(recvBuffer []byte) interface{}
}
