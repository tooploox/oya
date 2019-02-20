package internal

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/bilus/oya/pkg/secrets"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

func SecretsView(workDir string, stdout, stderr io.Writer) error {
	output, err := secrets.Decrypt(workDir)
	if err != nil {
		return err
	}
	stdout.Write(output)
	return nil
}

func SecretsEdit(workDir string, stdout, stderr io.Writer) error {
	viewCmd := secrets.ViewCmd(workDir)
	if err := terminalRun(viewCmd); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}

func SecretsEncrypt(workDir string, stdout, stderr io.Writer) error {
	if err := secrets.Encrypt(workDir); err != nil {
		return err
	}
	return nil
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
	defer close(ch)
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
		return err
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() {
		_, err = io.Copy(ptmx, os.Stdin)
		if err != nil {
			log.Printf("error %s", err)
		}
	}()

	_, err = io.Copy(os.Stdout, ptmx)
	if err != nil {
		return err
	}

	return nil
}
