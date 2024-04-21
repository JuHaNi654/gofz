package ssh

import (
	"strconv"
)

type Config struct {
	Host         string
	User         string
	Hostname     string
	Port         uint16
	IdentityFile string
}

// TODO some kind of validation?
func newConfig(data map[string]string) *Config {
	port, _ := toPort(data["Port"])

	return &Config{
		Host:         data["Host"],
		User:         data["User"],
		Hostname:     data["Hostname"],
		Port:         port,
		IdentityFile: data["IdentityFile"],
	}
}

func toPort(port string) (uint16, error) {
	value, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(value), nil
}
