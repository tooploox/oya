package testutil

import (
	"io"
	"os"
	"testing"
)

func CopyFile(fromPath, toPath string) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}

func MustCopyFile(t *testing.T, fromPath, toPath string) {
	err := CopyFile(fromPath, toPath)
	AssertNoErr(t, err, "Error copying file from %v to %v", fromPath, toPath)
}
