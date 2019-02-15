package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

const SecretsFile = "Oyafile.secrets"

type ErrNoSecretsFile struct {
	FileName    string
	OyafilePath string
}

func (err ErrNoSecretsFile) Error() string {
	return fmt.Sprintf("Oya secrets file \"%v\" not found in %v", err.FileName, err.OyafilePath)
}

type ErrSecretsAlreadyEncrypted struct {
	FileName    string
	OyafilePath string
}

func (err ErrSecretsAlreadyEncrypted) Error() string {
	return fmt.Sprintf("Oya secrets file \"%v\" already encrypted in %v", err.FileName, err.OyafilePath)
}

func SecretsView(workDir string, stdout, stderr io.Writer) error {
	if _, err := os.Stat(SecretsFile); err != nil {
		return ErrNoSecretsFile{FileName: SecretsFile, OyafilePath: workDir}
	}
	cmd := exec.Command("sops", "-d", SecretsFile)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func SecretsEdit(workDir string, stdout, stderr io.Writer) error {
	cmd := exec.Command("sops", SecretsFile)
	if err := terminalRun(cmd); err != nil {
		return err
	}
	return nil
}

func SecretsEncrypt(workDir string, stdout, stderr io.Writer) error {
	if _, err := os.Stat(SecretsFile); err != nil {
		return ErrNoSecretsFile{FileName: SecretsFile, OyafilePath: workDir}
	}
	// r, _, err := raw.Load(SecretsFile, workDir)
	// if err != nil {
	// 	fmt.Println("Error loading raw secrets file")
	// }
	if alreadyEncrypted() {
		return ErrSecretsAlreadyEncrypted{FileName: SecretsFile, OyafilePath: workDir}
	}
	cmd := exec.Command("sops", "-e", SecretsFile)
	encoded, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(SecretsFile, encoded, 0644)
	if err != nil {
		return err
	}
	return nil
}

func alreadyEncrypted() bool {
	// Trying to decrypt and check if succeed
	cmd := exec.Command("sops", "-d", SecretsFile)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func terminalRun(cmd *exec.Cmd) error {
	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}
