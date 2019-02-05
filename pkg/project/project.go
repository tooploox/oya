package project

import (
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/template"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TODO: Duplicated in oyafile module.
const VendorDir = ".oya/vendor"

type Project struct {
	RootDir string
}

func Load(rootDir string) (Project, error) {
	prj, err := Detect(rootDir)
	if err != nil {
		return prj, err
	}

	rel, err := filepath.Rel(rootDir, prj.RootDir)
	if err != nil {
		return prj, errors.Wrapf(err, "%v is not the Oya project root directory (it's %v)", rootDir, prj.RootDir)
	}
	if rel != "." {
		return prj, errors.Errorf("%v is not an Oya project root directory", rootDir)
	}

	return prj, nil
}

func Detect(workDir string) (Project, error) {
	detectedRootDir, found, err := detectRoot(workDir)
	if err != nil {
		return Project{}, err
	}
	if !found {
		return Project{}, ErrNoProject{Path: workDir}
	}
	return Project{
		RootDir: detectedRootDir,
	}, nil
}

func (p Project) Run(workDir, taskName string, positionalArgs []string, flags map[string]string, stdout, stderr io.Writer) error {
	log.Debugf("Task %q at %v", taskName, workDir)

	changes, err := p.Changeset(workDir)
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		return nil
	}

	foundAtLeastOneTask := false
	for _, o := range changes {
		found, err := o.RunTask(taskName, toScope(positionalArgs, flags), stdout, stderr)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
		}
		if found {
			foundAtLeastOneTask = found
		}
	}

	if !foundAtLeastOneTask {
		return ErrNoTask{
			Task: taskName,
		}
	}
	return nil
}

func (p Project) Oyafile(oyafilePath string) (*oyafile.Oyafile, bool, error) {
	return oyafile.Load(oyafilePath, p.RootDir)
}

func (p Project) Vendor(pack pack.Pack) error {
	return pack.Vendor(filepath.Join(p.RootDir, VendorDir))
}

func isRoot(raw oyafile.RawModifier) (bool, error) {
	return raw.HasKey("Project")
}

// detectRoot attempts to detect the root project directory marked by
// root Oyafile, i.e. one containing Project: directive.
// It walks the directory tree, starting from startDir, going upwards,
// looking for root.
func detectRoot(startDir string) (string, bool, error) {
	path := startDir
	maxParts := 256
	for i := 0; i < maxParts; i++ {
		raw, found, err := oyafile.LoadRawFromDir(path)
		if err == nil && found {
			isRoot, err := isRoot(raw)
			if err != nil {
				return "", false, err
			}
			if isRoot {
				return path, true, nil
			}
		}

		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}

	return "", false, nil
}

func toScope(positionalArgs []string, flags map[string]string) template.Scope {
	return template.Scope{
		"Args":  positionalArgs,
		"Flags": camelizeFlags(flags),
	}
}

func camelizeFlags(flags map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range flags {
		result[camelize(k)] = v
	}
	return result
}

var sepRx = regexp.MustCompile("(-|_).")

// camelize turns - or _ separated identifiers into camel case.
// Example: "aa-bb" becomes "aaBb".
func camelize(s string) string {
	return sepRx.ReplaceAllStringFunc(s, func(match string) string {
		return strings.ToUpper(match[1:])
	})

}
