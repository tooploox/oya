package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/bilus/oya/build"
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
}

func (s *SuiteContext) MustSetUp() {
	projectDir, err := ioutil.TempDir("", "oya")
	if err != nil {
		panic(err)
	}
	err = os.Chdir(projectDir)
	if err != nil {
		panic(err)
	}
	s.projectDir = projectDir
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

func (c *SuiteContext) fileProjectToopfileContaining(path string, contents *gherkin.DocString) error {
	return c.writeFile(path, contents.Content)
}

func (c *SuiteContext) fileProjectToopfileContains(path string, contents *gherkin.DocString) error {
	actual, err := c.readFile(path)
	if err != nil {
		return err
	}
	if actual != contents.Content {
		return fmt.Errorf("unexpected file %v contents: %v expected: %v", path, actual, contents.Content)
	}
	return nil
}

func (c *SuiteContext) oyaBuildIsRun(job string) error {
	return build.Build(c.projectDir, job)
}

func FeatureContext(s *godog.Suite) {
	c := SuiteContext{}
	s.Step(`^file (.+) containing$`, c.fileProjectToopfileContaining)
	s.Step(`^file (.+) contains$`, c.fileProjectToopfileContains)
	s.Step(`^"oya build (.+)" is run$`, c.oyaBuildIsRun)

	s.BeforeScenario(func(interface{}) { c.MustSetUp() })
	s.AfterScenario(func(interface{}, error) { c.MustTearDown() })
}
