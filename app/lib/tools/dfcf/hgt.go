package dfcf

// 沪港通
import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net/mail"
	"time"

	"github.com/jiorry/gotock/app/lib/tools/wget"
	"github.com/kere/gos"
	"github.com/kere/gos/lib/util"
)

type HgtAmount struct {
	Date    time.Time `json:"date"`
	AmountA float64   `json:"amount_a"`
	AmountB float64   `json:"amount_b"`
}

// http://datainterface.eastmoney.com/EM_DataCenter/JS.aspx?type=SHT&sty=SHTTMYE&rt=0.6341747129336
// curl "http://datainterface.eastmoney.com/EM_DataCenter/JS.aspx?type=SHT&sty=SHTTMYE&rt=0.6341747129336" -H "Pragma: no-cache" -H "Accept-Encoding: gzip, deflate, sdch" -H "Accept-Language: en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4,zh-TW;q=0.2" -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36" -H "Accept: */*" -H "Referer: http://data.eastmoney.com/bkzj/hgt.html" -H "X-Requested-With: ShockwaveFlash/18.0.0.232" -H "Cookie: HAList=f-0-000003-B"%"u80A1"%"u6307"%"u6570; pgv_pvi=3511567662; pgv_info=ssi=s7100762258" -H "Connection: keep-alive" -H "Cache-Control: no-cache" --compressed
// RzrqSumData 抓取数据
func GetHgtAmount() ([]*HgtAmount, error) {
	r := rand.New(rand.NewSource(99))
	formt := "http://datainterface.eastmoney.com/EM_DataCenter/JS.aspx?type=SHT&sty=SHTTMYE&rt=%v"

	body, err := wget.GetBody(fmt.Sprintf(formt, r.Float64()))
	if err != nil {
		return nil, gos.DoError(err)
	}

	now := gos.NowInLocation()
	arr := bytes.Split(body, []byte("\r\n"))
	var tmp [][]byte
	result := make([]*HgtAmount, len(arr))
	var amountA float64
	var amountB float64
	var date time.Time

	for i, v := range arr {
		tmp = bytes.Split(v, gos.B_SEMICOLON)
		amountA, err = util.Str2Float64(string(tmp[1]))
		if err != nil {
			amountA = 0
		}
		amountB, err = util.Str2Float64(string(tmp[2]))
		if err != nil {
			amountA = 0
		}
		date, err = time.ParseInLocation("2006/1/2 15:04", fmt.Sprint(now.Format("2006/1/2"), " ", string(tmp[0])), gos.GetSite().Location)
		if err != nil {
			date = now
		}
		result[i] = &HgtAmount{
			Date:    date,
			AmountA: amountA,
			AmountB: amountB,
		}
	}
	return result, nil
}

// func RunAlertHgt() {
// 	var errCh = make(chan error)
// 	c := time.Tick(1 * time.Minute)
// 	conf := gos.Configuration.GetConf("other")
// 	minute := conf.GetInt("hgt_check_minute")
// 	amount := conf.GetFloat("hgt_check_amount")
//
// 	go func() {
// 		for range c {
// 			errCh <- AlertAtHgtChanged(minute, amount)
// 		}
// 	}()
//
// 	go func() {
// 		for {
// 			select {
// 			case err := <-errCh:
// 				if err != nil {
// 					gos.DoError(err)
// 				}
// 			}
// 		}
// 	}()
// }

var isAlertAtHgtChanged = false
var countAlertAtHgtChanged = 0
var countAlertAtHgtChangedStep = 10

// AlertAtHgtChanged
// n range of minute
func AlertAtHgtChanged() error {
	conf := gos.Configuration.GetConf("other")
	n := conf.GetInt("hgt_check_minute")
	diff := conf.GetFloat("hgt_check_amount")
	now := gos.NowInLocation()

	switch now.Weekday() {
	case time.Sunday, time.Saturday:
		return nil
	}

	appConf := gos.Configuration.GetConf("other")
	if util.InStringSlice(appConf.GetStringSlice("holiday"), now.Format("20060102")) {
		return nil
	}

	//如果已经通知过，那么未来一段时间内不再通知
	if countAlertAtHgtChanged > countAlertAtHgtChangedStep {
		isAlertAtHgtChanged = false
	}

	if isAlertAtHgtChanged && countAlertAtHgtChanged < countAlertAtHgtChangedStep {
		countAlertAtHgtChanged++
		return nil
	}

	nowUnix := now.Unix()
	df := "2006-01-02 15:04"
	// t := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	t := now.Format("2006-01-02")

	begin, err := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 9, 0), gos.GetSite().Location)
	if err != nil {
		return err
	}
	beginUnix := begin.Unix()
	// end, err := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 15, 15), gos.GetSite().Location)
	// if err != nil {
	// 	return err
	// }

	// gos.Log.Info("AlertAtHgtChanged A", n, diff, now, begin, end)
	// if nowUnix < beginUnix || nowUnix > end.Unix() {
	// 	return nil
	// }

	minute := int((nowUnix - beginUnix) / 60)

	// 排除中午时间
	midA, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 12, 0), gos.GetSite().Location)
	midB, _ := time.ParseInLocation(df, fmt.Sprintf("%s %02d:%02d", t, 13, 0), gos.GetSite().Location)
	if nowUnix > midA.Unix() && nowUnix < midB.Unix() {
		return nil
	} else if nowUnix > midB.Unix() {
		minute -= 60
	}

	if minute < n+2 {
		return nil
	}

	items, err := GetHgtAmount()
	if err != nil {
		return err
	}
	// gos.Log.Info("AlertAtHgtChanged B", minute, n, len(items))

	amountCurrent := items[minute].AmountA
	// 如果当前时间的金额==0，那么退一步取值
	if amountCurrent == 0 {
		minute = minute - 1
		amountCurrent = items[minute].AmountA
	}
	// 如果任然==0，判断是数据有误
	if amountCurrent == 0 && minute > 10 {
		gos.Log.Info("amountCurrent==0 min=", minute, " ", items[minute], items[minute-1], items[minute-2], items[minute-3])
		return nil
	}

	amountBefore := items[minute-n].AmountA
	if amountBefore == 0 && minute > 10 {
		gos.Log.Info("amountBefore==0 min=", minute, " ", items[minute], items[minute-1], items[minute-2], items[minute-3])
		return nil
	}

	diffCurrent := amountCurrent - amountBefore
	// gos.Log.Info("AlertAtHgtChanged C", "cur", amountCurrent, "bef", amountBefore, "Diff:", diffCurrent, items[minute].Date, items[minute-n].Date)
	// 如果幅度小于预期，则退出检查
	if math.Abs(diffCurrent) < diff {
		return nil
	}

	conf = gos.Configuration.GetConf("mail")
	addr := conf.Get("addr")
	from := mail.Address{conf.Get("mail_user_name"), conf.Get("mail")}
	user := conf.Get("mail_user")
	password := conf.Get("mail_password")

	client := gos.NewSmtpPlainMail(addr, from, user, password)
	to := make([]*mail.Address, 1)
	to[0] = &mail.Address{"jiorry", "a@kere.me"}

	var title string
	var body string

	if diffCurrent > 0 {
		title = fmt.Sprintf("沪港通资金异动 %.2f", diffCurrent)
	} else {
		title = fmt.Sprintf("沪港通资金异动 %.2f", diffCurrent)
	}

	for i := 0; i < n; i++ {
		body += fmt.Sprintln(items[minute-i].Date.Format("15:04"), " ", items[minute-i].AmountA)
	}

	body += fmt.Sprintf("资金变动：%.2f\n", diffCurrent)
	body += fmt.Sprintln("http://data.eastmoney.com/bkzj/hgt.html")

	err = client.Send(title, body, to)
	if err != nil {
		return err
	}
	gos.Log.Info("Send Email", t)

	isAlertAtHgtChanged = true
	countAlertAtHgtChanged = 0

	return nil
}
