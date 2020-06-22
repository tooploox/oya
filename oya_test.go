package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/cucumber/godog"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tooploox/oya/cmd"
	"github.com/tooploox/oya/pkg/secrets"
	"github.com/tooploox/oya/testutil/gherkin"
)

const sopsPgpKey = "317D 6971 DD80 4501 A6B8  65B9 0F1F D46E 2E8C 7202"

type SuiteContext struct {
	projectDir string

	lastCommandErr      error
	lastCommandExitCode int
	stdin               io.WriteCloser
	stdout              *bytes.Buffer
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
}

func (c *SuiteContext) MustTearDown() {
	if err := removePGPKeys(c.projectDir); err != nil {
		panic(err)
	}

	if err := os.RemoveAll(c.projectDir); err != nil {
		panic(err)
	}
}

func setEnv(projectDir string) {
	err := os.Setenv("OYA_HOME", projectDir)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("SOPS_PGP_FP", sopsPgpKey)
	if err != nil {
		panic(err)
	}
}

// removePGPKeys removes PGP keys based on fingerprints(s) in .sops.yaml, NOT sopsPgpKey ^.
func removePGPKeys(projectDir string) error {
	if err := os.Chdir(projectDir); err != nil {
		return err
	}
	sops, err := secrets.LoadPGPSopsYaml()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}
	fingerprints := make([]string, 0)

	for _, rule := range sops.CreationRules {
		fingerprints = append(fingerprints, strings.Split(rule.PGP, ",")...)
	}
	return secrets.RemovePGPKeypairs(fingerprints)
}

// overrideOyaCmd overrides `oya` command used by $Tasks in templates
// to run oya tasks.
// It builds oya to a temporary directory and use it to launch Oya in scripts.
func overrideOyaCmd(projectDir string) {
	executablePath := filepath.Join(projectDir, "_bin/oya")
	oyaCmdOverride := fmt.Sprintf(
		"function oya() { (cd %v && go build -o %v oya.go) && %v $@; }",
		sourceFileDirectory(), executablePath, executablePath)
	os.Setenv("OYA_CMD", oyaCmdOverride)
}

func (c *SuiteContext) writeFile(path, contents string) error {
	targetPath := c.resolvePath(path)
	dir := filepath.Dir(targetPath)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(targetPath, []byte(contents), 0600)
}

func (c *SuiteContext) readFile(path string) (string, error) {
	sourcePath := c.resolvePath(path)
	contents, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return "", err
	}
	return string(contents), err
}

func (c *SuiteContext) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(c.projectDir, path)
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
	_, err := os.Stat(c.resolvePath(path))
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
	r, w := io.Pipe()
	c.stdin = w
	c.stdout.Reset()
	cmd.ResetFlags()

	oldArgs := os.Args
	os.Args = parseCommand(command)
	defer func() {
		os.Args = oldArgs
	}()
	cmd.SetInput(r)
	cmd.SetOutput(c.stdout)
	c.lastCommandExitCode, c.lastCommandErr = cmd.ExecuteE()
	return nil
}

func parseCommand(command string) []string {
	argv := make([]string, 0)
	r := regexp.MustCompile(`([^\s"']+)|"([^"]*)"|'([^']*)'`)
	matches := r.FindAllStringSubmatch(command, -1)
	for _, match := range matches {
		for _, group := range match[1:] {
			if group != "" {
				argv = append(argv, group)
			}
		}
	}
	return argv
}

func (c *SuiteContext) iRunOya(command string) error {
	return c.execute("oya " + command)
}

func (c *SuiteContext) iRunOyaInteractively(command string) error {
	go func() {
		err := c.execute("oya " + command)
		if err != nil {
			log.Fatalf("Unexpected error from interactive oya command: %v", err)
		}
	}()

	return nil
}

func (c *SuiteContext) modifyFileToContain(path string, contents *gherkin.DocString) error {
	return c.writeFile(path, contents.Content)
}

func (c *SuiteContext) iSendToRepl(line string) error {
	_, err := c.stdin.Write([]byte(line + "\n"))
	return err
}

func (c *SuiteContext) theCommandSucceeds() error {
	if c.lastCommandErr != nil {
		log.Println(c.stdout.String())
		return errors.Wrap(c.lastCommandErr, "command unexpectedly failed")
	}
	return nil
}

func (c *SuiteContext) theCommandFails() error {
	if c.lastCommandErr == nil {
		return errors.Errorf("last command unexpectedly succeeded")
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

func (c *SuiteContext) theCommandOutputs(expected *gherkin.DocString) error {
	actual := c.stdout.String()
	if actual != expected.Content {
		return fmt.Errorf("unexpected %v; expected: %q", actual, expected.Content)
	}
	return nil
}

func (c *SuiteContext) theCommandOutputsTextMatching(expected *gherkin.DocString) error {
	actual := c.stdout.String()
	rx := regexp.MustCompile(expected.Content)
	if !rx.MatchString(actual) {
		return fmt.Errorf("unexpected %v; expected to match: %q", actual, expected.Content)
	}
	return nil
}

func (c *SuiteContext) theCommandExitCodeIs(expectedExitCode int) error {
	if c.lastCommandExitCode != expectedExitCode {
		return errors.Errorf("unexpected exit code from the last command: %v; expected: %v", c.lastCommandExitCode, expectedExitCode)
	}
	return nil
}

func (c *SuiteContext) oyafileIsEncryptedUsingKeyInSopsyaml(oyafilePath string) error {
	sops, err := secrets.LoadPGPSopsYaml()
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadFile(oyafilePath)
	if err != nil {
		return err
	}
	for _, rule := range sops.CreationRules {
		fingerprint := rule.PGP
		if strings.Contains(string(contents), fingerprint) {
			return nil
		}
	}

	return errors.Errorf("%q not encrypted using key is .sops.yaml", oyafilePath)
}

func FeatureContext(s *godog.Suite) {
	c := SuiteContext{}
	s.Step(`^I'm in project dir$`, c.iAmInProjectDir)
	s.Step(`^I\'m in the (.+) dir$`, c.imInDir)
	s.Step(`^file (.+) containing$`, c.fileContaining)
	s.Step(`^I run "oya (.+)"$`, c.iRunOya)
	s.Step(`^I run "oya (.+)" interactively$`, c.iRunOyaInteractively)
	s.Step(`^I modify file (.+) to contain$`, c.modifyFileToContain)
	s.Step(`^I send "([^"]*)" to repl$`, c.iSendToRepl)
	s.Step(`^file (.+) contains$`, c.fileContains)
	s.Step(`^file (.+) does not contain$`, c.fileDoesNotContain)
	s.Step(`^file (.+) exists$`, c.fileExists)
	s.Step(`^file (.+) does not exist$`, c.fileDoesNotExist)
	s.Step(`^the command succeeds$`, c.theCommandSucceeds)
	s.Step(`^the command fails$`, c.theCommandFails)
	s.Step(`^the command fails with error$`, c.theCommandFailsWithError)
	s.Step(`^the command fails with error matching$`, c.theCommandFailsWithErrorMatching)
	s.Step(`^the command outputs$`, c.theCommandOutputs)
	s.Step(`^the command outputs text matching$`, c.theCommandOutputsTextMatching)
	s.Step(`^the command exit code is (.+)$`, c.theCommandExitCodeIs)
	s.Step(`^the ([^ ]+) environment variable set to "([^"]*)"$`, c.environmentVariableSet)
	s.Step(`^([^ ]+) is encrypted using PGP key in .sops.yaml$`, c.oyafileIsEncryptedUsingKeyInSopsyaml)
	s.BeforeScenario(func(*gherkin.Scenario) { c.MustSetUp() })
	s.AfterScenario(func(*gherkin.Scenario, error) { c.MustTearDown() })
}

// sourceFileDirectory returns the current .go source file directory.
func sourceFileDirectory() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
