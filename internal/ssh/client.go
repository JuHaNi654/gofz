package ssh

import (
	"fmt"
	"os"
	"strings"
)

type EventType int

const (
	List EventType = iota
	Wd
	Get
	Put
	Quit
)

type Event struct {
	Event   EventType
	Payload string
}

type RecvEvent struct {
	Event   EventType
	Payload any
}

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
func (c *SftpClient) Put(path string) {
	c.eventChan <- Event{
		Event: Put,
	}
}

// Download entry
func (c *SftpClient) Get(path string) {
	c.eventChan <- Event{
		Event: Get,
	}
}

// TODO close the channels if needed
func (c *SftpClient) Quit() {
	c.eventChan <- Event{
		Event: Quit,
	}
}

func loadPrivateKey(path string) []byte {
	wd, _ := os.UserHomeDir()
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", wd, 1)
	}

	key, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	return key
}
