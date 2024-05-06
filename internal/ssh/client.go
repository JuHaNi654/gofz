package ssh

type SftpClient struct {
	passphrase []byte
	eventChan  chan Event
	Recv       chan RecvEvent
}

func NewSftpClient() *SftpClient {
	return &SftpClient{
		eventChan: make(chan Event),
		Recv:      make(chan RecvEvent),
	}
}

func (c *SftpClient) SetPassphrase(val string) {
	c.passphrase = []byte(val)
}

func (c *SftpClient) Passphrase() []byte {
	return c.passphrase
}

func (c *SftpClient) Getwd() {
	c.eventChan <- Event{
		Event: Wd,
	}
}

// List entries
func (c *SftpClient) List(path string) {
	c.eventChan <- Event{
		Event:   List,
		Payload: path,
	}
}

// Upload entry
func (c *SftpClient) Put(target, dest string) {
	c.eventChan <- Event{
		Event: Put,
    Payload: []string{target, dest},
	}
}

// Download entry
func (c *SftpClient) Get(target, dest string) {
	c.eventChan <- Event{
		Event: Get,
    Payload: []string{target, dest}, 
	}
}

// TODO close the channels if needed
func (c *SftpClient) Quit() {
	c.eventChan <- Event{
		Event: Quit,
	}
}
