package runner

import (
	"github.com/jiorry/keredat/app/lib/tools/dfcf"
	"github.com/jiorry/keredat/app/lib/tools/sina/indexB"
)

// RunTimer
func RunTimer() {
	r := &Runner{}
	r.Open = []*AlertRunner{NewAlertRunner(dfcf.AlertAtHgtChanged, 1), NewAlertRunner(indexB.Alert, 1)}
	runtimer(r)
}
