package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/bilus/oya/cmd"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const SOPS_PGP_KEY = "317D 6971 DD80 4501 A6B8  65B9 0F1F D46E 2E8C 7202"

type SuiteContext struct {
	projectDir string

	lastCommandErr error
	stdout         *bytes.Buffer
	stderr         *bytes.Buffer
}

func (c *SuiteContext) MustSetUp() {
	projectDir, err := ioutil.TempDir("", "oya")
	if err != nil {
		panic(err)
	}

	overrideOyaCmd(projectDir)
	setEnv(projectDir)

	log.SetLevel(log.DebugLevel)
	c.projectDir = projectDir
	c.stdout = bytes.NewBuffer(nil)
	c.stderr = bytes.NewBuffer(nil)
}

func (c *SuiteContext) MustTearDown() {
	err := os.RemoveAll(c.projectDir)
	if err != nil {
		panic(err)
	}
}

func setEnv(projectDir string) {
	err := os.Setenv("OYA_HOME", projectDir)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("SOPS_PGP_FP", SOPS_PGP_KEY)
	if err != nil {
		panic(err)
	}
}

// overrideOyaCmd overrides `oya` command used by $Tasks in templates
// to run oya tasks.
// It builds oya to a temporary directory and use it to launch Oya in scripts.
func overrideOyaCmd(projectDir string) {
	executablePath := filepath.Join(projectDir, "_bin/oya")
	oyaCmdOverride := fmt.Sprintf(
		"(cd %v && go build -o %v oya.go); %v",
		sourceFileDirectory(), executablePath, executablePath)
	oyafile.OyaCmdOverride = &oyaCmdOverride
}

func (c *SuiteContext) writeFile(relPath, contents string) error {
	sourceFileDirectory := path.Join(c.projectDir, relPath)
	dir := path.Dir(sourceFileDirectory)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(sourceFileDirectory, []byte(contents), 0600)
}

func (c *SuiteContext) readFile(relPath string) (string, error) {
	sourceFileDirectory := path.Join(c.projectDir, relPath)
	contents, err := ioutil.ReadFile(sourceFileDirectory)
	if err != nil {
		return "", err
	}
	return string(contents), err
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

func (c *SuiteContext) environmentVariableSet(name, value string) error {
	return os.Setenv(name, value)
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

func (c *SuiteContext) fileDoesNotContain(path string, contents *gherkin.DocString) error {
	actual, err := c.readFile(path)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(".*" + contents.Content + ".*")
	if len(re.FindString(actual)) > 0 {
		return fmt.Errorf("unexpected file %v contents: %q NOT expected: %q", path, actual, contents.Content)
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
	c.stdout.Reset()
	c.stderr.Reset()

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

func (c *SuiteContext) modifyFileToContain(path string, contents *gherkin.DocString) error {
	return c.writeFile(path, contents.Content)
}

func (c *SuiteContext) theCommandSucceeds() error {
	if c.lastCommandErr != nil {
		log.Println(c.stdout.String())
		log.Println(c.stderr.String())
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

func (c *SuiteContext) theCommandOutputsTextMatching(target string, expected *gherkin.DocString) error {
	var actual string
	switch target {
	case "stdout":
		actual = c.stdout.String()
	case "stderr":
		actual = c.stderr.String()
	default:
		return fmt.Errorf("Unexpected command output target: %v", target)
	}
	rx := regexp.MustCompile(expected.Content)
	if !rx.MatchString(actual) {
		return fmt.Errorf("unexpected %v output: %q expected to match: %q", target, actual, expected.Content)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	c := SuiteContext{}
	s.Step(`^I'm in project dir$`, c.iAmInProjectDir)
	s.Step(`^I\'m in the (.+) dir$`, c.imInDir)
	s.Step(`^file (.+) containing$`, c.fileContaining)
	s.Step(`^I run "oya (.+)"$`, c.iRunOya)
	s.Step(`^I modify file (.+) to contain$`, c.modifyFileToContain)
	s.Step(`^file (.+) contains$`, c.fileContains)
	s.Step(`^file (.+) does not contain$`, c.fileDoesNotContain)
	s.Step(`^file (.+) exists$`, c.fileExists)
	s.Step(`^file (.+) does not exist$`, c.fileDoesNotExist)
	s.Step(`^the command succeeds$`, c.theCommandSucceeds)
	s.Step(`^the command fails with error$`, c.theCommandFailsWithError)
	s.Step(`^the command fails with error matching$`, c.theCommandFailsWithErrorMatching)
	s.Step(`^the command outputs to (stdout|stderr)$`, c.theCommandOutputs)
	s.Step(`^the command outputs to (stdout|stderr) text matching$`, c.theCommandOutputsTextMatching)
	s.Step(`^the ([^ ]+) environment variable set to "([^"]*)"$`, c.environmentVariableSet)

	s.BeforeScenario(func(interface{}) { c.MustSetUp() })
	s.AfterScenario(func(interface{}, error) { c.MustTearDown() })
}

// sourceFileDirectory returns the current .go source file directory.
func sourceFileDirectory() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
