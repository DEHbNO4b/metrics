package interfaces

// Pinger it is an interface for ping relation database.
type Pinger interface {
	Ping() bool
}
