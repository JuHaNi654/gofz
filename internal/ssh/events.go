package ssh

type EventType int

const (
	List EventType = iota
  Wd
	Get
	Put
  Error	
  Quit
  Connected
)

type Event struct {
	Event   EventType
	Payload any
}

type RecvEvent struct {
	Event   EventType
	Payload any
}
