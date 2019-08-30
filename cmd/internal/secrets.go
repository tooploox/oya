package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/secrets"
)

var ErrUnsupportedType = errors.New("Unsupported type")

func SecretsInit(typ, email, name, desc, format string, stdout, stderr io.Writer) error {
	if typ != "pgp" {
		return ErrUnsupportedType
	}

	keyPair, err := secrets.Init(email, name, desc)
	if err != nil {
		return err
	}

	if err = secrets.GeneratePGPSopsYaml(keyPair); err != nil {
		return err
	}

	if err = secrets.ImportPGPKeypair(keyPair); err != nil {
		return err
	}

	if format == "json" {
		b, err := json.MarshalIndent(keyPair, "", "  ")
		if err != nil {
			return err
		}
		stdout.Write(b)
	} else {
		fmt.Fprintf(stdout, "Generated a new PGP key (%q).\n", email)
		fmt.Fprintf(stdout, "Fingerprint: %v\n", keyPair.Fingerprint)
		fmt.Fprintf(stdout, "Imported the generated PGP key into GPG.\n")
		fmt.Fprintf(stdout, "Generated .sops.yaml referencing the new key.\n")
	}

	return nil
}

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
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func SecretsEncrypt(path string, stdout, stderr io.Writer) error {
	if err := secrets.Encrypt(path, path); err != nil {
		return err
	}
	return nil
}
