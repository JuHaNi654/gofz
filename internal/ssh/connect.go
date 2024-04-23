package ssh

import (
	"fmt"
	"gofz/internal/debug"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Connect(client *SftpClient, srvConfig *Config) (bool, error) {
	key := loadPrivateKey(srvConfig.IdentityFile)
	var err error
	var signer ssh.Signer
	signer, err = ssh.ParsePrivateKeyWithPassphrase(key, client.Passphrase())

	if err != nil {
		return false, PassphraseError("Invalid passphrase")
	}

	config := ssh.ClientConfig{
		User: srvConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", srvConfig.Hostname, srvConfig.Port)
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect to [%s]\n", addr)
		return false, SSHClientError(msg)
	}

	go handleConnection(conn, client)
	return true, nil
}

func handleConnection(conn *ssh.Client, client *SftpClient) {
	debug.Write("Open", "Connection")
	defer conn.Close()
	defer func() {
		debug.Write("Closed", "Connection")
	}()

	sc, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
		os.Exit(1)
	}

	defer sc.Close()
	// TODO handle errors
	for {
		event := <-client.eventChan
		switch event.Event {
		case List:
			debug.Write("List", "Incoming event type")
			files, _ := listFiles(*sc, event.Payload)
			client.Recv <- RecvEvent{
				Event:   List,
				Payload: files,
			}
		case Wd:
			debug.Write("Wd", "Incoming event type")
			path, _ := getWd(*sc)
			client.Recv <- RecvEvent{
				Event:   Wd,
				Payload: path,
			}
		case Get:
			debug.Write("Get", "Incoming event type")
		case Put:
			debug.Write("Put", "Incoming event type")
		case Quit:
			debug.Write("Quit", "Incoming event type")
			client.Recv <- RecvEvent{
				Event:   Quit,
				Payload: true,
			}
			return
		}
	}
}

func getWd(sc sftp.Client) (string, error) {
	return sc.Getwd()
}

func listFiles(sc sftp.Client, path string) ([]os.FileInfo, error) {
	return sc.ReadDir(path)
}
