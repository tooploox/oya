package secrets

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const SecretsFileName = "secrets.oya"

func Decrypt(workDir string) ([]byte, error) {
	var output []byte
	file := filepath.Join(workDir, SecretsFileName)
	if _, err := os.Stat(file); err != nil {
		return output, ErrNoSecretsFile{FileName: SecretsFileName, OyafilePath: workDir}
	}
	decryptCmd := exec.Command("sops", "-d", file)
	decrypted, err := decryptCmd.CombinedOutput()
	if err != nil {
		return output, ErrSecretsFailure{FileName: SecretsFileName, OyafilePath: workDir, CmdError: string(decrypted)}
	}
	return decrypted, nil
}

func Encrypt(workDir string) error {
	file := filepath.Join(workDir, SecretsFileName)
	if alreadyEncrypted(file) {
		return ErrSecretsAlreadyEncrypted{FileName: SecretsFileName, OyafilePath: workDir}
	}
	cmd := exec.Command("sops", "-e", file)
	encoded, err := cmd.CombinedOutput()
	if err != nil {
		return ErrSecretsFailure{FileName: SecretsFileName, OyafilePath: workDir, CmdError: string(encoded)}
	}
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, encoded, fi.Mode())
	if err != nil {
		return err
	}
	return nil
}

func ViewCmd(workDir string) *exec.Cmd {
	file := filepath.Join(workDir, SecretsFileName)
	return exec.Command("sops", file)
}

func alreadyEncrypted(file string) bool {
	// Trying to decrypt and check if succeed
	cmd := exec.Command("sops", "-d", file)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
