package ssh

import (
	"bufio"
	"gofz/internal/assert"
	"os"
	"strings"
)

func LoadSSHConfig(filename string) []*Config {
	var lines []string

	f, err := os.Open(filename)
	assert.Assert("could not open file", err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		lines = append(lines, strings.Trim(line, " "))
	}

	assert.Assert("Scanner error", scanner.Err())

	return parseConfig(lines)
}

func parseConfig(lines []string) []*Config {
	var items []*Config
	tmp := map[string]string{}

	split := func(r rune) bool {
		return r == 32 || r == 61 // Split by whitespace or '='
	}

	for i, line := range lines {
		pair := strings.FieldsFunc(line, split)
		if pair[0] == "Host" && len(tmp) != 0 {
			config := newConfig(tmp)
			items = append(items, config)
			tmp = map[string]string{}
		}

		tmp[pair[0]] = pair[1]

		if len(lines)-1 == i {
			config := newConfig(tmp)
			items = append(items, config)
		}
	}

	return items
}

func loadPrivateKey(path string) ([]byte, error) {
	wd, _ := os.UserHomeDir()
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", wd, 1)
	}

	key, err := os.ReadFile(path)
	if err != nil {
    return nil, err
	}

	return key, nil
}
