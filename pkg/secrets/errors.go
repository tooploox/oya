package secrets

import "fmt"

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
	return fmt.Sprintf("secrets file \"%v\" not found in %v", err.FileName, err.OyafilePath)
}

func (err ErrSecretsAlreadyEncrypted) Error() string {
	return fmt.Sprintf("secrets file \"%v\" already encrypted in %v", err.FileName, err.OyafilePath)
}

func (err ErrSecretsFailure) Error() string {
	return fmt.Sprintf("secrets failure on file \"%v\" in %v, with error: %v", err.FileName, err.OyafilePath, err.CmdError)
}
