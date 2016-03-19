package runner

import (
	"github.com/jiorry/keredat/app/lib/util/alert"

	"github.com/kere/gos"
	"github.com/kere/gos/lib/util"
)

type AlertRunner struct {
	Func         func() (*alert.AlertMessage, error)
	breakCount   int
	isDoBreak    bool
	CheckHoliday bool
	RunMode      int
}

// NewAlertRunner
func NewAlertRunner(f func() (*alert.AlertMessage, error), runMode int) *AlertRunner {
	return &AlertRunner{Func: f, CheckHoliday: true, RunMode: runMode}
}

// Run function
func (a *AlertRunner) Run() bool {
	if a.CheckHoliday && util.InStringSlice(holiday, gos.NowInLocation().Format("20060102")) {
		return false
	}
	go a.doAlert()
	return true
}

func (a *AlertRunner) doAlert() {
	var am *alert.AlertMessage
	var err error
	am, err = a.Func()

	switch a.RunMode {
	case 0: // 每一次都运行
		if err != nil {
			errCh <- err
			return
		}
		if am != nil {
			alertCh <- am
		}

	case 1: // 运行成功后休息一段周期
		if a.isDoBreak {
			a.breakCount++

			if a.breakCount < 5 {
				return
			} else {
				a.breakCount = 0
				a.isDoBreak = false
			}
		}

		am, err = a.Func()
		if err != nil {
			errCh <- err
			return
		}

		if am == nil {
			a.isDoBreak = false
		} else {
			alertCh <- am
			a.isDoBreak = true
		}
	}

}

func AlertTest() (*alert.AlertMessage, error) {
	return alert.NewAlertMessage("alert test", []byte("alert test"), 1), nil
}
