package ssh

type PassphraseError string

func (err PassphraseError) Error() string {
	return string(err)
}

type SSHClientError string

func (err SSHClientError) Error() string {
	return string(err)
}
