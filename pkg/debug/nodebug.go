// +build !debug

package debug

import (
	"github.com/tooploox/oya/pkg/oyafile"
)

func FP(bool) bool {
	return false
}

func LogOyafiles(msg string, oyafiles []*oyafile.Oyafile) {
}
