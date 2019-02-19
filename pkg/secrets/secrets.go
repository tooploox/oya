package secrets

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const SecretsFileName = "secrets.oya"

type ErrNoSecretsFile struct {
	FileName    string
	OyafilePath string
}

type ErrSecretsAlreadyEncrypted struct {
	FileName    string
	OyafilePath string
}

type ErrSecretsFailure struct {
	FileName    string
	OyafilePath string
	CmdError    string
}

func (err ErrNoSecretsFile) Error() string {
	return fmt.Sprintf("Oya secrets file \"%v\" not found in %v", err.FileName, err.OyafilePath)
}

func (err ErrSecretsAlreadyEncrypted) Error() string {
	return fmt.Sprintf("Oya secrets file \"%v\" already encrypted in %v", err.FileName, err.OyafilePath)
}

func (err ErrSecretsFailure) Error() string {
	return fmt.Sprintf("Oya secrets failure on file \"%v\" in %v, with error: %v", err.FileName, err.OyafilePath, err.CmdError)
}

func Decrypt(workDir string) ([]byte, error) {
	var output []byte
	file := workDir + "/" + SecretsFileName
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
	file := workDir + "/" + SecretsFileName
	if alreadyEncrypted(file) {
		return ErrSecretsAlreadyEncrypted{FileName: SecretsFileName, OyafilePath: workDir}
	}
	cmd := exec.Command("sops", "-e", file)
	encoded, err := cmd.CombinedOutput()
	if err != nil {
		return ErrSecretsFailure{FileName: SecretsFileName, OyafilePath: workDir, CmdError: string(encoded)}
	}
	err = ioutil.WriteFile(file, encoded, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ViewCmd(workDir string) *exec.Cmd {
	return exec.Command("sops", SecretsFileName)
}

func alreadyEncrypted(file string) bool {
	// Trying to decrypt and check if succeed
	cmd := exec.Command("sops", "-d", file)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
