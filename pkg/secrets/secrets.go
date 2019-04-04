package secrets

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const SecretsFileName = "secrets.oya"

func Decrypt(path string) ([]byte, bool, error) {
	if ok, err := isSopsFile(path); !ok || err != nil {
		return nil, false, err
	}

	var output []byte
	if _, err := os.Stat(path); err != nil {
		return output, false, ErrNoSecretsFile{Path: path}
	}
	decryptCmd := exec.Command("sops", "-d", path)
	decrypted, err := decryptCmd.CombinedOutput()
	if err != nil {
		return output, false,
			ErrSecretsFailure{Path: path, CmdError: string(decrypted)}
	}
	return decrypted, true, nil
}

func Encrypt(workDir string) error {
	file := filepath.Join(workDir, SecretsFileName)
	if alreadyEncrypted(file) {
		return ErrSecretsAlreadyEncrypted{Path: file}
	}
	cmd := exec.Command("sops", "-e", file)
	encoded, err := cmd.CombinedOutput()
	if err != nil {
		return ErrSecretsFailure{Path: file, CmdError: string(encoded)}
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

func isSopsFile(path string) (bool, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return false, err
	}
	dec := json.NewDecoder(jsonFile)
	var v map[string]interface{}
	if err := dec.Decode(&v); err != nil {
		return false, nil // Yes, ignoring error.
	}
	_, ok := v["sops"]
	return ok, nil
}

func alreadyEncrypted(file string) bool {
	// Trying to decrypt and check if succeed
	cmd := exec.Command("sops", "-d", file)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
