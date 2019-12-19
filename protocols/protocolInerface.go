// Protocol the interface of protocol, include three method,
// you can define custom protocol through implement three method.
package protocols

import "goes/connections"

type Protocol interface {
	// Encode encode package before send to client.
	Encode(data string, connectionInterface connections.ConnectionInterface)
	// Decode decode package and emit.
	Decode(recvBuffer string, connectionInterface connections.ConnectionInterface)
	// Input check the integrity of package.
	Input(recvBuffer string, connectionInterface connections.ConnectionInterface) int
}
