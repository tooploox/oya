package oyafile

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tooploox/oya/pkg/raw"
	"github.com/tooploox/oya/pkg/semver"
	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/types"
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

	for nameI, value := range of {
		name, ok := nameI.(string)
		if !ok {
			return nil, errors.Errorf("Incorrect value name: %v", name)
		}
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
			if err := ensureProject(raw); err != nil {
				return nil, errors.Wrapf(err, "unexpected Require directive")
			}
			err = parseRequire(name, value, oyafile)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		case "Replace":
			if err := ensureProject(raw); err != nil {
				return nil, errors.Wrapf(err, "unexpected Replace directive")
			}
			if err := parseReplace(name, value, oyafile); err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}

		default:
			taskName := task.Name(name)
			if taskName.IsBuiltIn() {
				log.Debugf("WARNING: Unrecognized built-in task or directive %q; skipping.", name)
				continue
			}

			if err := parseTask(name, value, oyafile); err != nil {
				return nil, errors.Wrapf(err, "error parsing key %q", name)
			}
		}
	}

	err = oyafile.resolveReplacements()
	if err != nil {
		return nil, err
	}

	err = oyafile.addBuiltIns()
	if err != nil {
		return nil, err
	}
	return oyafile, nil
}

// resolveReplacements replaces Requires paths based on Requires directives, if any.
func (oyafile *Oyafile) resolveReplacements() error {
	for i, ref := range oyafile.Requires {
		replPath, ok := oyafile.Replacements[ref.ImportPath]
		if ok {
			oyafile.Requires[i].ReplacementPath = replPath
		}
	}
	return nil
}

func parseMeta(metaName, key string) (task.Name, bool) {
	taskName := strings.TrimSuffix(key, "."+metaName)
	return task.Name(taskName), taskName != key
}

func parseImports(value interface{}, o *Oyafile) error {
	if value == nil {
		return nil
	}
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
	if value == nil {
		return nil
	}
	values, ok := value.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("expected map of keys to values; got %T", value)
	}
	for k, v := range values {
		key, ok := k.(string)
		if !ok {
			return fmt.Errorf("expected map of keys to values")
		}
		o.Values[key] = v
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
	if value == nil {
		return nil
	}

	defaultErr := fmt.Errorf("expected entries mapping pack import paths to their version, example: \"github.com/tooploox/oya-packReferences/docker: v1.0.0\"")

	requires, ok := value.(map[interface{}]interface{})
	if !ok {
		return defaultErr
	}

	packReferences := make([]PackReference, 0, len(requires))
	for importPathI, versionI := range requires {
		importPath, ok := importPathI.(string)
		if !ok {
			return defaultErr
		}
		version, ok := versionI.(string)
		if !ok {
			return defaultErr
		}

		ver, err := semver.Parse(version)
		if err != nil {
			return err
		}
		packReferences = append(packReferences,
			PackReference{
				ImportPath: types.ImportPath(importPath),
				Version:    ver,
			})
	}

	o.Requires = packReferences
	return nil
}

func parseReplace(name string, value interface{}, o *Oyafile) error {
	defaultErr := fmt.Errorf("expected entries mapping pack import paths to paths relative to the project root directory, example: \"github.com/tooploox/oya-pack/docker: /packs/docker\"")

	replacements, ok := value.(map[interface{}]interface{})
	if !ok {
		return defaultErr
	}

	for importPathI, pathI := range replacements {
		importPath, ok := importPathI.(string)
		if !ok {
			return defaultErr
		}
		path, ok := pathI.(string)
		if !ok {
			return defaultErr
		}

		o.Replacements[types.ImportPath(importPath)] = path
	}

	return nil

}

func ensureProject(raw *raw.Oyafile) error {
	_, hasProject, err := raw.Project()
	if err != nil {
		return err
	}
	if hasProject {
		return nil
	}
	return errors.Errorf("must be in file with a Project directive")
}
