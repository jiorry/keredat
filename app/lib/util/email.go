package util

import (
	"net/mail"
	"strings"

	"github.com/kere/gos"
)

func SendEmail(title, body string) error {
	conf := gos.Configuration.GetConf("mail")
	addr := conf.Get("addr")
	from := mail.Address{conf.Get("mail_user_name"), conf.Get("mail")}
	user := conf.Get("mail_user")
	password := conf.Get("mail_password")
	mailList := strings.Split(conf.Get("mail_to_list"), ",")

	client := gos.NewSmtpPlainMail(addr, from, user, password)
	var to = make([]*mail.Address, len(mailList))
	var arr []string
	for i, _ := range mailList {
		arr = strings.Split(mailList[i], ":")
		if len(arr) == 2 {
			to[i] = &mail.Address{arr[0], arr[1]}
		} else {
			to[i] = &mail.Address{"", mailList[i]}
		}
	}

	gos.Log.Info("Send Email", title)
	return client.Send(title, body, to)
}
