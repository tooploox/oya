package oyafile

import (
	"fmt"
	"strings"

	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/raw"
	"github.com/bilus/oya/pkg/semver"
	"github.com/bilus/oya/pkg/task"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Parse(raw *raw.Oyafile) (*Oyafile, error) {
	of, err := raw.Decode()
	if err != nil {
		return nil, err
	}
	oyafile, err := New(raw.Path, raw.RootDir)
	if err != nil {
		return nil, err
	}
	for name, value := range of {
		switch name {
		case "Import":
			err := parseImports(value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		case "Values":
			err := parseValues(value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		case "Secrets":
			err := parseValues(value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		case "Project":
			err := parseProject(value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		case "Ignore":
			err := parseIgnore(value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		case "Changeset":
			err := parseTask(name, value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		case "Require":
			err := parseRequire(name, value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}

		default:
			taskName := task.Name(name)
			if taskName.IsBuiltIn() {
				log.Debugf("WARNING: Unrecognized built-in task or directive %q; skipping.", name)
				continue
			}

			err := parseTask(name, value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		}
	}

	err = oyafile.resolveImports()
	if err != nil {
		return nil, err
	}
	err = oyafile.addBuiltIns()
	if err != nil {
		return nil, err
	}
	return oyafile, nil
}

func parseMeta(metaName, key string) (task.Name, bool) {
	taskName := strings.TrimSuffix(key, "."+metaName)
	return task.Name(taskName), taskName != key
}

func parseImports(value interface{}, o *Oyafile) error {
	imports, ok := value.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("expected map of aliases to paths")
	}
	for alias, path := range imports {
		alias, ok := alias.(string)
		if !ok {
			return fmt.Errorf("expected import alias")
		}
		path, ok := path.(string)
		if !ok {
			return fmt.Errorf("expected import path")
		}
		o.Imports[types.Alias(alias)] = types.ImportPath(path)
	}
	return nil
}

func parseValues(value interface{}, o *Oyafile) error {
	values, ok := value.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("expected map of keys to values")
	}
	for k, v := range values {
		valueName, ok := k.(string)
		if !ok {
			return fmt.Errorf("expected map of keys to values")
		}
		o.Values[valueName] = v
	}
	return nil
}

func parseProject(value interface{}, o *Oyafile) error {
	projectName, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected project name, actual: %v", value)
	}
	o.Project = projectName
	return nil
}

func parseIgnore(value interface{}, o *Oyafile) error {
	rulesI, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("expected an array of ignore rules, actual: %v", value)
	}
	rules := make([]string, len(rulesI))
	for i, ri := range rulesI {
		rule, ok := ri.(string)
		if !ok {
			return fmt.Errorf("expected an array of ignore rules, actual: %v", ri)
		}
		rules[i] = rule
	}
	o.Ignore = rules
	return nil
}

func parseTask(name string, value interface{}, o *Oyafile) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected a script, actual: %v", name)
	}
	if taskName, ok := parseMeta("Doc", name); ok {
		o.Tasks.AddDoc(taskName, s)
	} else {
		o.Tasks.AddTask(task.Name(name), task.Script{
			Script: s,
			Shell:  o.Shell,
			Scope:  &o.Values,
		})
	}
	return nil
}

func parseRequire(name string, value interface{}, o *Oyafile) error {
	defaultErr := fmt.Errorf("expected entries mapping pack import paths to their version, example: \"github.com/tooploox/oya-packs/docker: v1.0.0\"")

	requires, ok := value.(map[interface{}]interface{})
	if !ok {
		return defaultErr
	}

	packs := make([]pack.Pack, 0, len(requires))
	for importPathI, versionI := range requires {
		importPath, ok := importPathI.(string)
		if !ok {
			return defaultErr
		}
		version, ok := versionI.(string)
		if !ok {
			return defaultErr
		}
		l, err := pack.OpenLibrary(importPath)
		if err != nil {
			return err
		}

		ver, err := semver.Parse(version)
		if err != nil {
			return err
		}
		pack, err := l.Version(ver)
		if err != nil {
			return err
		}
		packs = append(packs, pack)
	}

	o.Require = packs
	return nil
}
