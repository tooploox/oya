package debug

import (
	log "github.com/sirupsen/logrus"
	"github.com/tooploox/oya/pkg/oyafile"
)

func LogOyafiles(msg string, oyafiles []*oyafile.Oyafile) {
	if log.GetLevel() == log.DebugLevel {
		log.Debug(msg)
		for _, o := range oyafiles {
			log.Debugf("  %v", o.Dir)
		}
	}
}
