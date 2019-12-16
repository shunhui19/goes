// Protocol the interface of protocol, include three method,
// you can define custom protocol through implement three method.
package protocols

type Protocol interface {
	// Encode encode package before send to client.
	Encode()
	// Decode decode package and emit.
	Decode()
	// Input check the integrity of package.
	Input() int
}
