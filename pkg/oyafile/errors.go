package oyafile

import (
	"fmt"
	"strings"

	"github.com/tooploox/oya/pkg/task"
)

type ErrTaskFail struct {
	TaskName task.Name
	Args     []string
}

func (e ErrTaskFail) Error() string {
	var optArgs string
	if len(e.Args) > 0 {
		optArgs = fmt.Sprintf(" with the following arguments: %s", strings.Join(e.Args, ", "))
	}
	return fmt.Sprintf("task %q failed%v", e.TaskName, optArgs)
}
