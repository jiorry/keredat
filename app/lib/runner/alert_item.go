package runner

import (
	"github.com/jiorry/keredat/app/lib/util/alert"

	"github.com/kere/gos"
	"github.com/kere/gos/lib/util"
)

type AlertItem struct {
	Func         func() (*alert.AlertMessage, error)
	CheckHoliday bool
	isRun        bool
	count        int
}

func NewAlertItem(f func() (*alert.AlertMessage, error)) *AlertItem {
	return &AlertItem{Func: f, CheckHoliday: true}
}

func (a *AlertItem) ResetStatus() {
	a.isRun = false
	a.count = 0
}

func (a *AlertItem) Run() {
	if a.CheckHoliday && util.InStringSlice(holiday, gos.NowInLocation().Format("20060102")) {
		return
	}

	if a.count > 5 {
		a.count = 0
	} else if a.count > 0 {
		a.count++
	}

	go a.doAlert()
}

func (a *AlertItem) RunOnce() {
	if a.isRun {
		return
	}

	a.Run()
}

func (a *AlertItem) doAlert() {
	if a.count > 0 {
		return
	}

	alertItem, err := a.Func()
	if err != nil {
		errCh <- err
	}

	if alertItem != nil {
		alertCh <- alertItem
	}

	a.isRun = true
	a.count++
}
