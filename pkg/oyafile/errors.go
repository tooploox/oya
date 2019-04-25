package oyafile

import (
	"fmt"

	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/types"
)

type ErrTaskFail struct {
	OyafilePath string
	TaskName    task.Name
	Args        []string
	ImportPath  *types.ImportPath
}

func (e ErrTaskFail) Error() string {
	return fmt.Sprintf("task failed")
}
