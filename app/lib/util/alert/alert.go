package alert

import (
	"fmt"
	"time"

	"github.com/kere/gos"
)

type AlertMessage struct {
	Title     string
	Message   []byte
	CreatedAt time.Time
	IType     int
}

func NewAlertMessage(title string, message []byte, itype int) *AlertMessage {
	return &AlertMessage{Title: title, Message: message, CreatedAt: gos.NowInLocation(), IType: itype}
}

func (a *AlertMessage) TitleString() string {
	if a.IType > 0 {
		return fmt.Sprint("** ", a.Title)
	} else if a.IType < 0 {
		return fmt.Sprint("-- ", a.Title)
	}

	return a.Title
}
