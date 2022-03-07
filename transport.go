package barrelman

// Transport is the interface met by an API transport for the system
type Transport interface {
	Serve(address string) error
}
