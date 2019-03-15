package project_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/tooploox/oya/pkg/project"
	"github.com/tooploox/oya/pkg/template"
	tu "github.com/tooploox/oya/testutil"
)

var (
	noArgs  []string
	noFlags map[string]string
	noScope template.Scope
)

func TestProject_Detect_NoOya(t *testing.T) {
	workDir := "./fixtures/empty_project"
	installDir := "" // Unused
	_, err := project.Detect(workDir, installDir)
	tu.AssertErr(t, err, "Expected error trying to detect Oya project in empty dir")
}

func TestProject_Detect_InRootDir(t *testing.T) {
	workDir := "./fixtures/project"
	installDir := "" // Unused
	_, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
}

func TestProject_Detect_InSubDir(t *testing.T) {
	workDir := "./fixtures/project/subdir"
	installDir := "" // Unused
	_, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
}

func TestProject_Detect_InEmptySubDir(t *testing.T) {
	workDir := "./fixtures/project/empty_subdir"
	installDir := "" // Unused
	_, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
}

func TestProject_Run_NoTask(t *testing.T) {
	workDir := "./fixtures/project"
	installDir := "" // Unused
	project, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	err = project.Run(workDir, "noSuchTask", false, false, nil, noScope, ioutil.Discard, ioutil.Discard)
	tu.AssertErr(t, err, "Expected error when trying to run without matching task")
}
func TestProject_Run_NoTaskRecurse(t *testing.T) {
	workDir := "./fixtures/project"
	installDir := "" // Unused
	project, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	err = project.Run(workDir, "noSuchTask", true, false, nil, noScope, ioutil.Discard, ioutil.Discard)
	tu.AssertErr(t, err, "Expected error when trying to run without matching task")
}

func TestProject_Run_NoChanges(t *testing.T) {
	workDir := "./fixtures/empty_changeset_project"
	installDir := "" // Unused
	project, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err = project.Run(workDir, "build", false, true, nil, noScope, stdout, stderr)
	tu.AssertNoErr(t, err, "Expected no error running with empty changeset")
	tu.AssertEqual(t, 0, len(stdout.String()))
	tu.AssertEqual(t, 0, len(stderr.String()))
}

func TestProject_Run_NoChangesRecurse(t *testing.T) {
	workDir := "./fixtures/empty_changeset_project"
	installDir := "" // Unused
	project, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err = project.Run(workDir, "build", true, true, nil, noScope, stdout, stderr)
	tu.AssertNoErr(t, err, "Expected no error running with empty changeset")
	tu.AssertEqual(t, 0, len(stdout.String()))
	tu.AssertEqual(t, 0, len(stderr.String()))
}

func TestProject_Run_WithChanges(t *testing.T) {
	workDir := "./fixtures/project"
	installDir := "" // Unused
	project, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err = project.Run(workDir, "build", false, true, nil, template.Scope{"Args": noArgs, "Flags": noFlags}, stdout, stderr)
	tu.AssertNoErr(t, err, "Expected no error running non-empty changeset")
	tu.AssertEqual(t, "build run", stdout.String())
}

func TestProject_Run_WithChangesRecurse(t *testing.T) {
	workDir := "./fixtures/project"
	installDir := "" // Unused
	project, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err = project.Run(workDir, "build", true, true, nil, template.Scope{"Args": noArgs, "Flags": noFlags}, stdout, stderr)
	tu.AssertNoErr(t, err, "Expected no error running non-empty changeset")
	tu.AssertEqual(t, "build run", stdout.String())
}

func TestProject_Run_WithArgs(t *testing.T) {
	workDir := "./fixtures/project"
	installDir := "" // Unused
	project, err := project.Detect(workDir, installDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err = project.Run(workDir, "build", false, true, []string{"arg1", "arg2", "--flag1", "--flag2"},
		template.Scope{
			"Args": []string{"arg1", "arg2"},
			"Flags": map[string]string{
				"flag1": "flag1",
				"flag2": "flag2",
			},
		},
		stdout, stderr)
	tu.AssertNoErr(t, err, "Expected no error running non-empty changeset")
	tu.AssertEqual(t, "build run\nAll args: arg1 arg2 --flag1 --flag2\nArgs: arg1 arg2\nFlags: flag1 flag2", stdout.String())
}
