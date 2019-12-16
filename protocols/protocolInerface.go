package protocols

type Protocol interface {
	// Encode encode package before send to client
	Encode()
	// Decode decode package and emit
	Decode()
	// Input check the integrity of package
	Input() int
}
