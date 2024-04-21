package ssh

import (
	"os"
	"testing"
)

func createTmpConfig(t testing.TB, content []byte) (string, func()) {
	t.Helper()

	path, _ := os.Getwd()
	tmp, err := os.CreateTemp(path, "example_config")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	tmp.Write(content)

	removeFile := func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}

	return tmp.Name(), removeFile
}

func TestLoadConfig(t *testing.T) {
	checkValidValues := func(t testing.TB, value, expected any) {
		t.Helper()
		if value != expected {
			t.Errorf("got '%v' :: want '%v'", value, expected)
		}
	}

	t.Run("test valid config file", func(t *testing.T) {
		filename, removeTmpConfig := createTmpConfig(t, []byte(`Host example
User lorem
Hostname 127.0.0.1
Port 22
IdentityFile ~/.ssh/test`))
		defer removeTmpConfig()

		result := LoadSSHConfig(filename)

		if len(result) != 1 {
			t.Fatalf("LoadSSHConfig returned empty list")
		}

		checkValidValues(t, result[0].Host, "example")
		checkValidValues(t, result[0].Hostname, "127.0.0.1")
		checkValidValues(t, result[0].User, "lorem")
		checkValidValues(t, result[0].Port, uint16(22))
		checkValidValues(t, result[0].IdentityFile, "~/.ssh/test")
	})

	t.Run("test valid config file second format", func(t *testing.T) {
		filename, removeTmpConfig := createTmpConfig(
			t,
			[]byte(`Host example
User=lorem
Hostname=127.0.0.1
Port=22
IdentityFile=~/.ssh/test`))
		defer removeTmpConfig()

		result := LoadSSHConfig(filename)

		if len(result) != 1 {
			t.Fatalf("LoadSSHConfig returned empty list")
		}

		checkValidValues(t, result[0].Host, "example")
		checkValidValues(t, result[0].Hostname, "127.0.0.1")
		checkValidValues(t, result[0].User, "lorem")
		checkValidValues(t, result[0].Port, uint16(22))
		checkValidValues(t, result[0].IdentityFile, "~/.ssh/test")
	})
}
