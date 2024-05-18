package ssh

type SftpChannel struct {
	Sender chan Event
	Recv   chan RecvEvent
}

func NewSftpChannel() *SftpChannel {
	return &SftpChannel{
		Recv: make(chan RecvEvent),
	}
}

func (c *SftpChannel) OpenSender() {
	c.Sender = make(chan Event)
}

func (c *SftpChannel) Getwd() {
	c.Sender <- Event{
		Event: Wd,
	}
}

// List entries
func (c *SftpChannel) List(path string) {
	c.Sender <- Event{
		Event:   List,
		Payload: path,
	}
}

// Upload entry
func (c *SftpChannel) Put(target, dest string) {
	c.Sender <- Event{
		Event:   Put,
		Payload: []string{target, dest},
	}
}

// Download entry
func (c *SftpChannel) Get(target, dest string) {
	c.Sender <- Event{
		Event:   Get,
		Payload: []string{target, dest},
	}
}

func (c *SftpChannel) Quit() {
	close(c.Sender)
}

func (c *SftpChannel) handleResponse(t EventType, payload any, err error) {
	if err != nil {
		c.Recv <- RecvEvent{
			Event:   Error,
			Payload: err,
		}
	} else {
		c.Recv <- RecvEvent{
			Event:   t,
			Payload: payload,
		}
	}
}
