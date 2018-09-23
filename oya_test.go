package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/bilus/oya/build"
	"github.com/pkg/errors"
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

	lastCommandErr error
	stdout         *bytes.Buffer
	stderr         *bytes.Buffer
}

func (s *SuiteContext) MustSetUp() {
	projectDir, err := ioutil.TempDir("", "oya")
	if err != nil {
		panic(err)
	}
	s.projectDir = projectDir
	s.stdout = bytes.NewBuffer(nil)
	s.stderr = bytes.NewBuffer(nil)
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

func (c *SuiteContext) iRunOyaBuild(job string) error {
	c.lastCommandErr = build.Build(c.projectDir, job, c.stdout, c.stderr)
	return nil
}

func (c *SuiteContext) theCommandSucceeds() error {
	if c.lastCommandErr != nil {
		return errors.Wrap(c.lastCommandErr, "build failed")
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
	s.Step(`^file (.+) containing$`, c.fileContaining)
	s.Step(`^I run "oya build (.+)"$`, c.iRunOyaBuild)
	s.Step(`^file (.+) contains$`, c.fileContains)
	s.Step(`^file (.+) exists$`, c.fileExists)
	s.Step(`^the command succeeds$`, c.theCommandSucceeds)
	s.Step(`^the command outputs to (stdout|stderr)$`, c.theCommandOutputs)

	s.BeforeScenario(func(interface{}) { c.MustSetUp() })
	s.AfterScenario(func(interface{}, error) { c.MustTearDown() })
}
