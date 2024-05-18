package ssh

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

func Connect(client *SftpChannel, srvConfig *Config) (bool, error) {
	var err error
	var signer ssh.Signer

	key, err := loadPrivateKey(srvConfig.IdentityFile)
	if err != nil {
		return false, err
	}

	if srvConfig.Protected {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, srvConfig.PassphraseAsBytes())
	} else {
		signer, err = ssh.ParsePrivateKey(key)
	}

	if err != nil {
		return false, PassphraseError("Invalid passphrase")
	}

	hostKeyCallback, err := kh.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return false, err
	}

	config := ssh.ClientConfig{
		User: srvConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
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
