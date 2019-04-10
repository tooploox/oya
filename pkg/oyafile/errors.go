package oyafile

import (
	"fmt"

	"github.com/tooploox/oya/pkg/task"
)

type ErrTaskFail struct {
	Cause       error
	OyafilePath string
	TaskName    task.Name
	Args        []string
}

func (e ErrTaskFail) Error() string {
	return fmt.Sprintf("task %v failed in %v: %v", string(e.TaskName), e.OyafilePath, e.Cause.Error())
}
