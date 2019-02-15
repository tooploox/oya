package secrets

import (
	"os"
	"os/exec"
)

func DecryptSecrets(file string) ([]byte, error) {
	var output []byte
	if _, err := os.Stat(file); err != nil {
		return output, err
	}
	decryptCmd := exec.Command("sops", "-d", file)
	decrypted, err := decryptCmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return decrypted, nil
}
