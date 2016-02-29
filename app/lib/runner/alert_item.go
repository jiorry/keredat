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

	// 每一次都运行
	if a.RunMode == 0 {
		am, err = a.Func()
		if err != nil {
			errCh <- err
		}

		return
	}

	if a.isDoBreak {
		a.breakCount++

		if a.breakCount < 5 {
			return
		} else {
			a.breakCount = 0
			a.isDoBreak = false
		}
	}

	// 运行成功后休息一段周期
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
		if am.Title == "alert test" {
			println("Title:alert test")
		} else {
			alertCh <- am
		}
		a.isDoBreak = true
	}

}

func AlertTest() (*alert.AlertMessage, error) {
	return alert.NewAlertMessage("alert test", []byte("alert test"), 1), nil
}
