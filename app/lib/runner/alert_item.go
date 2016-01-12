package runner

import (
	"github.com/jiorry/keredat/app/lib/util/alert"

	"github.com/kere/gos"
	"github.com/kere/gos/lib/util"
)

type AlertRunner struct {
	Func         func() (*alert.AlertMessage, error)
	breakCount   int
	isBreak      bool
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
	// 每一次都运行
	if a.RunMode == 0 {
		am, err = a.Func()
		if err != nil {
			errCh <- err
			return
		}
	}

	if a.isBreak {
		if a.breakCount < 5 {
			return
		} else {
			a.breakCount = 0
			a.isBreak = false
		}
	}

	// 运行成功后停止一段时间
	if a.RunMode == 1 {
		am, err = a.Func()
		if err != nil {
			errCh <- err
			return
		}
	}

	if am == nil {
		return
	} else {
		alertCh <- am
		a.isBreak = true
	}

	if a.isBreak {
		a.breakCount++
	}
}
