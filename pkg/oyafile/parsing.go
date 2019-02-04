package oyafile

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func parseOyafile(path, rootDir string, of OyafileFormat) (*Oyafile, error) {
	oyafile, err := New(path, rootDir)
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

	return oyafile, nil
}

func parseMeta(metaName, key string) (string, bool) {
	taskName := strings.TrimSuffix(key, "."+metaName)
	return taskName, taskName != key
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
		o.Imports[Alias(alias)] = ImportPath(path)
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
		o.Tasks.AddTask(name, ScriptedTask{
			Name:   name,
			Script: Script(s),
			Shell:  o.Shell,
			Scope:  &o.Values,
		})
	}
	return nil
}
