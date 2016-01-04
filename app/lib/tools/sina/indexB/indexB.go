package indexB

// 成分B指数

import (
	"bytes"
	"container/list"
	"fmt"
	"math"
	"time"

	"github.com/jiorry/keredat/app/lib/util"
	"github.com/jiorry/keredat/app/lib/util/ajax"
	"github.com/jiorry/keredat/app/lib/util/alert"
	"github.com/kere/gos"
)

var (
	ajaxClient *ajax.Ajax
	indexBList list.List
)

func init() {
	ajaxClient = ajax.NewAjax("")
}

func Alert() (*alert.AlertMessage, error) {
	gos.Log.Info("indexB Alert")
	conf := gos.Configuration.GetConf("other")
	n := conf.GetInt("indexb_check_minute")
	diff := conf.GetFloat("indexb_check_value")

	data, err := fetchIndexB()
	if err != nil {
		return nil, err
	}

	l := list.New()
	l.PushFront(data)
	if l.Len() < n {
		return nil, nil
	}
	if l.Len() > n {
		l.Remove(l.Back())
	}

	front := l.Front().Value.(*indexBStruct)
	back := l.Back().Value.(*indexBStruct)

	val := 100 * (front.Price - back.Price) / back.Price
	if math.Abs(val) < diff {
		return nil, nil
	}

	msg := fmt.Sprint(front.String(), "\n", back.String())
	if val > 0 {
		return alert.NewAlertMessage(fmt.Sprintf("指数B异常增长 %.4f", val), []byte(msg), 1), nil
	} else {
		return alert.NewAlertMessage(fmt.Sprintf("指数B异常下跌 %.4f", val), []byte(msg), -1), nil
	}
}

func fetchIndexB() (*indexBStruct, error) {
	// r := rand.New(rand.NewSource(99))
	// http://hq.sinajs.cn/?_=0.4559384104795754&list=sh000003
	body, err := ajaxClient.GetBody("http://hq.sinajs.cn/list=sh000003")
	if err != nil {
		return nil, gos.DoError(err)
	}

	// var hq_str_sh000003="Ｂ股指数,444.4195,442.7577,407.7723,453.2994,407.5514,0,0,4127282,3448536298,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2015-12-28,15:03:55,00";
	//                             ,今  开   ,昨  收  , 现   价 , 最  高 ,  最  低 ， ，，成交量 ，  成交额   ，
	//                          1  ，  2    ，   3   ，   4    ，  5    ，   6   ，7，8，  9   ，   10   ， ，， ，，， ，， ，，，， ，，， ， ，，，，， 31       ，32
	if len(body) < 21 {
		return nil, gos.DoError("length less 21:", string(body))
	}
	n := bytes.Index(body, []byte("=\""))
	src := body[n : len(body)-2]
	arr := bytes.Split(src, gos.B_COMMA)

	data := &indexBStruct{}
	data.Price = util.ParseMoney(string(arr[3]))
	data.High = util.ParseMoney(string(arr[4]))
	data.Low = util.ParseMoney(string(arr[5]))
	data.Amount = int64(util.ParseMoney(string(arr[8])))
	data.Total = int64(util.ParseMoney(string(arr[9])))

	date, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprint(string(arr[30]), " ", string(arr[31])), gos.GetSite().Location)
	if err != nil {
		return nil, err
	}

	data.Date = date

	return data, nil
}

type indexBStruct struct {
	Price  float64
	High   float64
	Low    float64
	Amount int64
	Total  int64
	Date   time.Time
}

func (i *indexBStruct) String() string {
	return fmt.Sprint(i.Date.Format("2006-01-02 15:04:05"), " 现价:", i.Price, " 交易额:", i.Amount, " 最高:", i.High, " 最低:", i.Low)
}
