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
	o, found, err := detectRoot(rootDir)
	if err != nil {
		return Project{}, errors.Wrapf(err, "error when attempting to load Oya project from %v", rootDir)
	}
	if !found {
		return Project{}, errors.Errorf("missing Oyafile at %v", rootDir)
	}

	rel, err := filepath.Rel(rootDir, o.RootDir)
	if err != nil {
		return Project{}, errors.Wrapf(err, "%v is not the Oya project root directory (it's %v)", rootDir, o.RootDir)
	}
	if rel != "." {
		return Project{}, errors.Errorf("%v is not an Oya project root directory", rootDir)
	}
	return Project{
		RootDir: o.RootDir,
	}, nil
}

func Detect(workDir string) (Project, error) {
	o, found, err := detectRoot(workDir)
	if err != nil {
		return Project{}, err
	}
	if !found {
		return Project{}, ErrNoProject{Path: workDir}
	}
	return Project{
		RootDir: o.RootDir,
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

func isRoot(o *oyafile.Oyafile) bool {
	return len(o.Project) > 0
}

// detectRoot attempts to detect the root project directory marked by
// root Oyafile, i.e. one containing Project: directive.
// It walks the directory tree, starting from startDir, going upwards,
// loading Oyafiles in each dir. It stops when it encounters the root
// Oyafile.
func detectRoot(startDir string) (*oyafile.Oyafile, bool, error) {
	path := startDir
	maxParts := 256
	var lastErr error
	for i := 0; i < maxParts; i++ {
		o, found, err := oyafile.LoadFromDir(path, path)
		if err != nil {
			lastErr = err
		}
		if err == nil && found && isRoot(o) {
			return o, true, nil
		}
		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}

	return nil, false, lastErr
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
