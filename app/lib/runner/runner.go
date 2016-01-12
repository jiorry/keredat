package runner

import (
	"fmt"
	"time"

	"github.com/jiorry/keredat/app/lib/util/alert"
	"github.com/jiorry/keredat/app/lib/util/email"

	"github.com/kere/gos"
)

var errCh chan error
var alertCh chan *alert.AlertMessage
var holiday []string

func init() {
	errCh = make(chan error)
	alertCh = make(chan *alert.AlertMessage)

	oConf := gos.Configuration.GetConf("other")
	holiday = oConf.GetStringSlice("holiday")

}

func runtimer(r *Runner) {
	go r.do1MinuteAction()
	go handlerError()
}

// Runner moniter
type Runner struct {
	Before          []*AlertRunner
	MorningPrepare  []*AlertRunner
	Noon            []*AlertRunner
	Open            []*AlertRunner
	ClosedAfterNoon []*AlertRunner
	Night           []*AlertRunner
}

func (a *Runner) do1MinuteAction() {
	// every 1 minute
	c := time.Tick(1 * time.Minute)

	for range c {

		now := gos.NowInLocation()
		switch now.Weekday() {
		case time.Sunday, time.Saturday:
			continue
		}

		df := "2006-01-02 15:04"
		t := now.Format("2006-01-02")

		prepare, err := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 7, 0), gos.GetSite().Location)
		if err != nil {
			panic(err)
		}

		beginA, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 9, 0), gos.GetSite().Location)

		endA, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 11, 30), gos.GetSite().Location)
		beginB, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 13, 0), gos.GetSite().Location)
		endB, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 15, 15), gos.GetSite().Location)
		night, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 19, 0), gos.GetSite().Location)

		if now.Before(prepare) {
			for _, item := range a.Before {
				item.Run()
			}

		} else if now.Before(beginA) {
			// 7-9点
			for _, item := range a.MorningPrepare {
				item.Run()
			}

		} else if now.Before(endA) {
			// 早场
			for _, item := range a.Open {
				item.Run()
			}

		} else if now.Before(beginB) {
			// 中午休息
			for _, item := range a.Noon {
				item.Run()
			}

		} else if now.Before(endB) {
			// 下半场 PartB
			for _, item := range a.Open {
				item.Run()
			}

		} else if now.Before(night) {
			// 下午
			for _, item := range a.ClosedAfterNoon {
				item.Run()
			}

		} else if night.Before(now) {
			// 晚上
			for _, item := range a.Night {
				item.Run()
			}

		}

	}
}

func handlerError() {
	for {
		select {
		case a := <-alertCh:
			gos.Log.Info("email:", a.TitleString())
			err := email.SendPlainEmail(a.TitleString(), a.Message)
			if err != nil {
				gos.DoError(err)
			}
		case err := <-errCh:
			gos.DoError(err)
		}
	}
}
