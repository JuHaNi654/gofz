package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Config struct {
	Host         string
	User         string
	Hostname     string
	Port         uint16
	IdentityFile string

  Passphrase string
  Protected bool
}
func (c *Config) PassphraseAsBytes() []byte {
  return []byte(c.Passphrase)
}

func newConfig(data map[string]string) *Config {
	port, _ := toPort(data["Port"])
  protected, _ := isKeyProtected(data["IdentityFile"])

  return &Config{
		Host:         data["Host"],
		User:         data["User"],
		Hostname:     data["Hostname"],
		Port:         port,
		IdentityFile: data["IdentityFile"],
    Protected:    protected,	
  }
}

func toPort(port string) (uint16, error) {
	value, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(value), nil
}

func isKeyProtected(keyFile string) (bool, error) {
  path, _ := os.UserHomeDir()
	
  if strings.HasPrefix(keyFile, "~") {
		keyFile = strings.Replace(keyFile, "~", path, 1)
	} else {
    keyFile = fmt.Sprintf("%s/%s", path, keyFile)
  }
 
  cmd := exec.Command("ssh-keygen", "-f", keyFile, "-y")
  _, err := cmd.CombinedOutput()
 

  if err != nil {
    return true, nil
  }

  return false, nil
}
