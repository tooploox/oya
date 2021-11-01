package printers

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/types"
)

type taskInfo struct {
	taskName       task.Name
	alias          types.Alias
	bareTaskName   string
	meta           task.Meta
	relOyafilePath string
}

type TaskList struct {
	workDir string
	tasks   []taskInfo
}

func NewTaskList(workDir string) *TaskList {
	return &TaskList{
		workDir: workDir,
	}
}

func (p *TaskList) AddTask(taskName task.Name, meta task.Meta, oyafilePath string) error {
	alias, bareTaskName := taskName.Split()
	relPath, err := filepath.Rel(p.workDir, oyafilePath)
	if err != nil {
		return err
	}

	p.tasks = append(p.tasks, taskInfo{taskName, alias, bareTaskName, meta, relPath})
	return nil
}

func (p *TaskList) Print(w io.Writer) {
	sortTasks(p.tasks)

	printTask := p.taskPrinter()

	lastRelPath := ""
	first := true
	for _, t := range p.tasks {
		if t.relOyafilePath != lastRelPath {
			if !first {
				fmt.Fprintln(w)
			} else {
				first = false
			}

			fmt.Fprintf(w, "# in ./%s\n", t.relOyafilePath)
		}
		printTask(w, t.taskName, t.meta)
		lastRelPath = t.relOyafilePath
	}
}

func (p *TaskList) taskPrinter() func(io.Writer, task.Name, task.Meta) {
	docOffset := maxTaskWidth(p.tasks)
	return func(w io.Writer, taskName task.Name, meta task.Meta) {
		fmt.Fprintf(w, "oya run %s", taskName)
		exposed := meta.IsTaskExposed(taskName)
		if len(meta.Doc) > 0 || exposed {
			padding := strings.Repeat(" ", docOffset-len(taskName))
			fmt.Fprintf(w, "%s #", padding)
			if len(meta.Doc) > 0 {
				fmt.Fprintf(w, " %s", meta.Doc)
			}

			if exposed {
				fmt.Fprintf(w, " (%s)", string(meta.OriginalTaskName))
			}
		}
		fmt.Fprintln(w)
	}
}

func maxTaskWidth(tasks []taskInfo) int {
	w := 0
	for _, t := range tasks {
		l := len(string(t.taskName))
		if l > w {
			w = l
		}
	}
	return w
}

func isParentPath(p1, p2 string) bool {
	relPath, _ := filepath.Rel(p2, p1)
	return strings.Contains(relPath, "../")
}

func sortTasks(tasks []taskInfo) {
	sort.SliceStable(tasks, func(i, j int) bool {
		lt := tasks[i]
		rt := tasks[j]

		ldir := path.Dir(lt.relOyafilePath)
		rdir := path.Dir(rt.relOyafilePath)

		// Top-level tasks go before tasks in subdirectories.
		if isParentPath(ldir, rdir) {
			return true
		}
		if isParentPath(rdir, ldir) {
			return false
		}

		if rdir == ldir {
			if len(lt.alias) == 0 && len(rt.alias) != 0 {
				return true
			}
			if len(lt.alias) != 0 && len(rt.alias) == 0 {
				return false
			}
			// Tasks w/o alias before tasks with alias,
			// sort aliases alphabetically.
			if lt.alias < rt.alias {
				return true
			}

			// Sort tasks alphabetically.
			if lt.bareTaskName < rt.bareTaskName {
				return true
			}
		}

		return false
	})
}
