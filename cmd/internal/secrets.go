package internal

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/secrets"
)

func SecretsView(path string, stdout, stderr io.Writer) error {
	output, found, err := secrets.Decrypt(path)
	if err != nil {
		return err
	}
	if !found {
		return errors.Errorf("secret file %q not found", path)
	}
	stdout.Write(output)
	return nil
}

func SecretsEdit(filename string, stdout, stderr io.Writer) error {
	cmd := secrets.ViewCmd(filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout // os.Stdout
	cmd.Stderr = stderr // os.Stderr
	return cmd.Run()
}

func SecretsEncrypt(path string, stdout, stderr io.Writer) error {
	if err := secrets.Encrypt(path, path); err != nil {
		return err
	}
	return nil
}

func terminalRun(cmd *exec.Cmd) error {
	return cmd.Run()
}
