package runner

import (
	"fmt"
	"time"

	"github.com/jiorry/keredat/app/lib/tools/dfcf"
	"github.com/jiorry/keredat/app/lib/util/alert"
	"github.com/jiorry/keredat/app/lib/util/email"

	"github.com/kere/gos"
	"github.com/kere/gos/lib/util"
)

var errCh chan error
var alertCh chan *alert.AlertMessage

func init() {
	errCh = make(chan error)
}

// RunTimer
func RunTimer() error {
	go run1MinuteAction()
	go handlerError()

	return nil
}

func run1MinuteAction() {
	// every 1 minute
	c := time.Tick(1 * time.Minute)
	oConf := gos.Configuration.GetConf("other")
	for range c {

		now := gos.NowInLocation()
		switch now.Weekday() {
		case time.Sunday, time.Saturday:
			continue
		}

		if util.InStringSlice(oConf.GetStringSlice("holiday"), now.Format("20060102")) {
			continue
		}

		df := "2006-01-02 15:04"
		t := now.Format("2006-01-02")

		beginA, err := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 9, 0), gos.GetSite().Location)
		if err != nil {
			panic(err)
		}

		endA, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 11, 30), gos.GetSite().Location)
		beginB, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 13, 0), gos.GetSite().Location)
		endB, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 15, 15), gos.GetSite().Location)
		night, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 19, 0), gos.GetSite().Location)

		if now.Before(beginA) {
			// 9点以前
		} else if beginA.Before(now) && now.Before(endA) {
			// 早场
			go func() {
				alertItem, err := dfcf.AlertAtHgtChanged()
				if err != nil {
					errCh <- err
				} else {
					alertCh <- alertItem
				}
			}()
		} else if endA.Before(now) && now.Before(beginB) {
			// 中午休息
		} else if beginB.Before(now) && now.Before(endB) {
			// 下半场
			go func() {
				alertItem, err := dfcf.AlertAtHgtChanged()
				if err != nil {
					errCh <- err
				} else {
					alertCh <- alertItem
				}
			}()
		} else if endB.Before(now) && now.Before(night) {
			// 下午
		} else if night.Before(now) {
			// 晚上
			// go func() {
			// 	err := hx50etf.StoreTodayETFData()
			// 	errCh <- err
			// }()
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
			if err != nil {
				gos.DoError(err)
			}
		}
	}
}
