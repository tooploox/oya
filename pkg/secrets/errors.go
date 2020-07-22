package secrets

import "fmt"

type ErrNoSecretsFile struct {
	Path string
}

type ErrSecretsAlreadyEncrypted struct {
	Path string
}

type ErrSecretsFailure struct {
	Path   string
	Err    error
	Stderr []byte
}

func (err ErrNoSecretsFile) Error() string {
	return fmt.Sprintf("secret file %q not found", err.Path)
}

func (err ErrSecretsAlreadyEncrypted) Error() string {
	return fmt.Sprintf("secret file %q already encrypted", err.Path)
}

func (err ErrSecretsFailure) Error() string {
	if err.Stderr != nil {
		return fmt.Sprintf("error procesing secret file %q: %v\n%s",
			err.Path, err.Err, string(err.Stderr))

	}
	return fmt.Sprintf("error procesing secret file %q: %v",
		err.Path, err.Err)
}
