package ssh

import (
	"fmt"
	"gofz/internal/debug"
	"io"
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
		  path, _ := event.Payload.(string)	
      files, _ := listFiles(sc, path)
			client.Recv <- RecvEvent{
				Event:   List,
				Payload: files,
			}
		case Wd:
			debug.Write("Wd", "Incoming event type")
			path, _ := getWd(sc)
			client.Recv <- RecvEvent{
				Event:   Wd,
				Payload: path,
			}
		case Get:
			debug.Write("Get", "Incoming event type")
      items, _ := event.Payload.([]string)
      // TODO maybe make better solution than array of items
      msg, err := getFile(sc, items[0], items[1])
   
      if err != nil {
        client.Recv <- RecvEvent{
          Event: Error,
          Payload: err,
        }
      }

			client.Recv <- RecvEvent{
				Event: Get,
				Payload: msg,
			}
    case Put:
			debug.Write("Put", "Incoming event type")
      items, _ := event.Payload.([]string)
      msg, err := putFile(sc, items[0], items[1])

      if err != nil {
        client.Recv <- RecvEvent{
          Event: Error,
          Payload: err,
        }
      }

			client.Recv <- RecvEvent{
				Event: Put,
				Payload: msg,
			}
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

func putFile(sc *sftp.Client, target, dest string) (string, error)  {  
  localFile, err := os.Open(target)
  if err != nil {
    return "", err
  }
  defer localFile.Close()

  remoteFile, err := sc.OpenFile(dest, (os.O_WRONLY|os.O_CREATE|os.O_TRUNC))
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
    debug.Write(err.Error(), "Error opening remote file")
    return "", err
  }
  
  defer remoteFile.Close()

  localFile, err := os.Create(dest)
  if err != nil {
    debug.Write(err.Error(), "Error creating or opening local file")
    return "", err
  }

  defer localFile.Close()

  bytes, err := io.Copy(localFile, remoteFile)

  if err != nil {
    debug.Write(err.Error(), "Error while copying bytes")
    return "", err
  }

  msg := fmt.Sprintf("bytes copied: %d", bytes)
  return msg, nil
}

func getWd(sc *sftp.Client) (string, error) {
	return sc.Getwd()
}

func listFiles(sc *sftp.Client, path string) ([]os.FileInfo, error) {
	return sc.ReadDir(path)
}
