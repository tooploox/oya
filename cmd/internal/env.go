package internal

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
)

func installDir() (string, error) {
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

func lookupOyaScope() (string, bool) {
	return os.LookupEnv("OYA_SCOPE")
}

func setOyaScope(scope string) error {
	return os.Setenv("OYA_SCOPE", scope)
}
