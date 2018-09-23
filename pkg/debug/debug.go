package debug

import (
	"github.com/bilus/oya/pkg/oyafile"
	log "github.com/sirupsen/logrus"
)

func LogOyafiles(msg string, oyafiles []*oyafile.Oyafile) {
	if log.GetLevel() == log.DebugLevel {
		log.Debug(msg)
		for _, o := range oyafiles {
			log.Debugf("  %v", o.Dir)
		}
	}
}
