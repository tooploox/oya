package internal

import (
	"fmt"
	"io"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
)

func Import(rootDir, uri string, stdout, stderr io.Writer) error {
	fmt.Println("Hello")
	fmt.Printf("%v\n", rootDir)
	fmt.Printf("%v\n", uri)
	err := oyafile.AddImport(rootDir, uri)
	if err != nil {
		return errors.Wrapf(err, "Error while importing pack %v", uri)
	}
	return nil
}
