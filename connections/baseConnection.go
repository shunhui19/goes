// base struct of connection.
package connections

type BaseConnection struct {
	// ConnectionCount count of connection.
	ConnectionCount int
	// TotalRequest count of request.
	TotalRequest int
	// ThrowException count of error.
	ThrowException int
	// SendFail count of send fail.
	SendFail int
}
