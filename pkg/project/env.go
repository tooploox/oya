package project

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
)

func InstallDir() (string, error) {
	homeDir, found := os.LookupEnv("OYA_HOME")
	if !found {
		user, err := user.Current()
		if err != nil {
			return "", err
		}

		if len(user.HomeDir) == 0 {
			return "", errors.Errorf("Could not detect user home directory")
		}
		homeDir = user.HomeDir
	}

	return filepath.Join(homeDir, ".oya", "packs"), nil
}

func LookupOyaScope() (string, bool) {
	return os.LookupEnv("OYA_SCOPE")
}

func SetOyaScope(scope string) error {
	return os.Setenv("OYA_SCOPE", scope)
}
