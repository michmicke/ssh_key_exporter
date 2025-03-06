package ssh

import (
	"bufio"
	"os"

	"golang.org/x/crypto/ssh"
)

type SshKey struct {
	Keytype     string
	Fingerprint string
	Comment     string
}

func ParseAuthorizedKeysFile(path string) ([]SshKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keys []SshKey
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		key, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
		if err != nil {
			continue
		}

		sshKey := SshKey{
			Keytype:     key.Type(),
			Fingerprint: ssh.FingerprintSHA256(key),
			Comment:     comment,
		}
		keys = append(keys, sshKey)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
