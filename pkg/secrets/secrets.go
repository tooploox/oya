package secrets

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

type KeyPair struct {
	Public      string `json:"public_key"`
	Private     string `json:"private_key"`
	Fingerprint string `json:"fingerprint"`
}

func Init(email, name, desc string) (KeyPair, error) {
	return generatePGPKeyPair(email, name, desc)
}

func Decrypt(path string) ([]byte, bool, error) {
	if ok, err := isSopsFile(path); !ok || err != nil {
		return nil, false, err
	}

	var output []byte
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return output, false, ErrNoSecretsFile{Path: path}
		} else {
			return output, false, err
		}
	}
	decryptCmd := exec.Command("sops", "-d", path)
	decrypted, err := decryptCmd.CombinedOutput()
	if err != nil {
		return output, false,
			ErrSecretsFailure{Path: path, Err: err}
	}
	return decrypted, true, nil
}

func Encrypt(inputPath, outputPath string) error {
	if alreadyEncrypted(inputPath) {
		return ErrSecretsAlreadyEncrypted{Path: inputPath}
	}
	cmd := exec.Command("sops", "-e", inputPath)
	encoded, err := cmd.CombinedOutput()
	if err != nil {
		return ErrSecretsFailure{Path: inputPath, Err: err}
	}
	fi, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outputPath, encoded, fi.Mode())
	if err != nil {
		return err
	}
	return nil
}

func ViewCmd(filename string) *exec.Cmd {
	return exec.Command("sops", filename)
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

func alreadyEncrypted(path string) bool {
	// Trying to decrypt and check if succeed
	cmd := exec.Command("sops", "-d", path)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
