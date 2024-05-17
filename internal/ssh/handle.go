package ssh

import (
	"fmt"
	"gofz/internal/assert"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func handleConnection(conn *ssh.Client, client *SftpClient) {
	defer conn.Close()

	sc, err := sftp.NewClient(conn)
	assert.Assert("Unable to start SFTP subsystem", err)

	defer sc.Close()

	for {
		event := <-client.eventChan
		switch event.Event {
		case List:
			path, _ := event.Payload.(string)
			files, err := sc.ReadDir(path)

			if err != nil {
				client.handleResponse(Error, err)
			}

			client.handleResponse(List, files)
		case Wd:
			path, err := sc.Getwd()

			if err != nil {
				client.handleResponse(Error, err)
			}

			client.handleResponse(Wd, path)
		case Get:
			items, _ := event.Payload.([]string)
			// TODO maybe make better solution than array of items
			msg, err := getFile(sc, items[0], items[1])

			if err != nil {
				client.handleResponse(Error, err)
			}

			client.handleResponse(Get, msg)
		case Put:
			items, _ := event.Payload.([]string)
			msg, err := putFile(sc, items[0], items[1])

			if err != nil {
				client.handleResponse(Error, err)
			}

			client.handleResponse(Put, msg)
		case Quit:
			client.handleResponse(Quit, true)
			return
		}
	}
}

func putFile(sc *sftp.Client, target, dest string) (string, error) {
	localFile, err := os.Open(target)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	remoteFile, err := sc.OpenFile(dest, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		return "", err
	}
	defer remoteFile.Close()

	bytes, err := io.Copy(remoteFile, localFile)

	if err != nil {
		return "", err
	}

	msg := fmt.Sprintf("bytes copied: %d", bytes)
	return msg, nil
}

func getFile(sc *sftp.Client, target, dest string) (string, error) {
	remoteFile, err := sc.OpenFile(target, (os.O_RDONLY))
	if err != nil {
		return "", err
	}

	defer remoteFile.Close()

	localFile, err := os.Create(dest)
	if err != nil {
		return "", err
	}

	defer localFile.Close()

	bytes, err := io.Copy(localFile, remoteFile)

	if err != nil {
		return "", err
	}

	msg := fmt.Sprintf("bytes copied: %d", bytes)
	return msg, nil
}
