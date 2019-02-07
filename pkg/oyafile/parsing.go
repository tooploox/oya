package oyafile

import (
	"fmt"
	"strings"

	"github.com/bilus/oya/pkg/raw"
	"github.com/bilus/oya/pkg/task"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
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
		default:
			err := parseUserTask(name, value, oyafile)
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

func parseUserTask(name string, value interface{}, o *Oyafile) error {
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
