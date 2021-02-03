package barrelman

type EventEmitter interface {
	Send(Event)
}

type Event struct {
	Device *Device
	Key    string
	Value  string
}
