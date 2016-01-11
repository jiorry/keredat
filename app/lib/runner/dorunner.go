package runner

import (
	"github.com/jiorry/keredat/app/lib/tools/dfcf"
	"github.com/jiorry/keredat/app/lib/tools/sina/indexB"
)

// RunTimer
func RunTimer() {
	r := &Runner{}
	r.Open = []*AlertItem{NewAlertItem(dfcf.AlertAtHgtChanged), NewAlertItem(indexB.Alert)}
	runtimer(r)
}
