package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/bilus/oya/cmd"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("oya", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

type SuiteContext struct {
	projectDir string
	vendorDir  string

	lastCommandErr error
	stdout         *bytes.Buffer
	stderr         *bytes.Buffer
}

func (c *SuiteContext) MustSetUp() {
	projectDir, err := ioutil.TempDir("", "oya")
	if err != nil {
		panic(err)
	}
	vendorDir := filepath.Join(projectDir, "oya/vendor")
	err = os.MkdirAll(vendorDir, 0700)
	if err != nil {
		panic(err)
	}
	log.SetLevel(log.DebugLevel)
	c.projectDir = projectDir
	c.vendorDir = vendorDir
	c.stdout = bytes.NewBuffer(nil)
	c.stderr = bytes.NewBuffer(nil)
}

func (c *SuiteContext) writeFile(relPath, contents string) error {
	fullPath := path.Join(c.projectDir, relPath)
	dir := path.Dir(fullPath)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, []byte(contents), 0600)
}

func (c *SuiteContext) readFile(relPath string) (string, error) {
	fullPath := path.Join(c.projectDir, relPath)
	contents, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(contents), err
}

func (c *SuiteContext) MustTearDown() {
	err := os.RemoveAll(c.projectDir)
	if err != nil {
		panic(err)
	}
}

func (c *SuiteContext) iAmInProjectDir() error {
	return os.Chdir(c.projectDir)
}

func (c *SuiteContext) imInDir(subdir string) error {
	return os.Chdir(subdir)
}

func (c *SuiteContext) fileContaining(path string, contents *gherkin.DocString) error {
	return c.writeFile(path, contents.Content)
}

func (c *SuiteContext) fileContains(path string, contents *gherkin.DocString) error {
	actual, err := c.readFile(path)
	if err != nil {
		return err
	}
	if actual != contents.Content {
		return fmt.Errorf("unexpected file %v contents: %q expected: %q", path, actual, contents.Content)
	}
	return nil
}

func (c *SuiteContext) fileExists(path string) error {
	_, err := os.Stat(path)
	return err
}

func (c *SuiteContext) fileDoesNotExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.Errorf("expected %v to not exist", path)
}

func (c *SuiteContext) execute(command string) error {
	oldArgs := os.Args
	os.Args = strings.Fields(command)
	defer func() {
		os.Args = oldArgs
	}()
	cmd.SetOutput(c.stdout)
	c.lastCommandErr = cmd.ExecuteE()
	return nil
}

func (c *SuiteContext) iRunOya(command string) error {
	return c.execute("oya " + command)
}

func (c *SuiteContext) theCommandSucceeds() error {
	if c.lastCommandErr != nil {
		return errors.Wrap(c.lastCommandErr, "command unexpectedly failed")
	}
	return nil
}

func (c *SuiteContext) theCommandFailsWithError(errMsg *gherkin.DocString) error {
	errMsg.Content = fmt.Sprintf("^%v$", errMsg.Content)
	return c.theCommandFailsWithErrorMatching(errMsg)
}

func (c *SuiteContext) theCommandFailsWithErrorMatching(errMsg *gherkin.DocString) error {
	if c.lastCommandErr == nil {
		return errors.Errorf("last command unexpectedly succeeded")
	}

	rx := regexp.MustCompile(errMsg.Content)
	if !rx.MatchString(c.lastCommandErr.Error()) {
		return errors.Wrap(c.lastCommandErr,
			fmt.Sprintf("unexpected error %q; expected to match %q", c.lastCommandErr, errMsg.Content))
	}
	return nil
}

func (c *SuiteContext) theCommandOutputs(target string, expected *gherkin.DocString) error {
	var actual string
	switch target {
	case "stdout":
		actual = c.stdout.String()
	case "stderr":
		actual = c.stderr.String()
	default:
		return fmt.Errorf("Unexpected command output target: %v", target)
	}
	if actual != expected.Content {
		return fmt.Errorf("unexpected %v output: %q expected: %q", target, actual, expected.Content)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	c := SuiteContext{}
	s.Step(`^I'm in project dir$`, c.iAmInProjectDir)
	s.Step(`^I\'m in the (.+) dir$`, c.imInDir)
	s.Step(`^file (.+) containing$`, c.fileContaining)
	s.Step(`^I run "oya (.+)"$`, c.iRunOya)
	s.Step(`^file (.+) contains$`, c.fileContains)
	s.Step(`^file (.+) exists$`, c.fileExists)
	s.Step(`^file (.+) does not exist$`, c.fileDoesNotExist)
	s.Step(`^the command succeeds$`, c.theCommandSucceeds)
	s.Step(`^the command fails with error$`, c.theCommandFailsWithError)
	s.Step(`^the command fails with error matching$`, c.theCommandFailsWithErrorMatching)
	s.Step(`^the command outputs to (stdout|stderr)$`, c.theCommandOutputs)

	s.BeforeScenario(func(interface{}) { c.MustSetUp() })
	s.AfterScenario(func(interface{}, error) { c.MustTearDown() })
}
