package dfcf

// 股指期货

import (
	"encoding/json"
	"fmt"

	"github.com/kere/gos"
)

func FetchGzqh() (*gzqhStruct, error) {
	// r := rand.New(rand.NewSource(99))
	// http://datainterface.eastmoney.com/EM_DataCenter/JS.aspx?type=QHCC&sty=QHSYCC&stat=3&mkt=069001009&fd=2016-03-01&code=if1603&cb=callback&callback=callback&_=1456832661681
	str := "http://datainterface.eastmoney.com/EM_DataCenter/JS.aspx?type=QHCC&sty=QHSYCC&stat=3&mkt=069001009&fd=%s&code=if%s&cb=callback&callback=callback&_=%d"
	now := gos.NowInLocation()
	// now.Weekday() == time.Friday

	str = fmt.Sprintf(str, now.Format("2006-01-02"), now.Format("0601"), now.Unix())
	body, err := ajaxClient.GetBody(str)
	if err != nil {
		return nil, gos.DoError(err)
	}

	if len(body) < 21 {
		return nil, gos.DoError("length less 21:", string(body))
	}

	v := &gzqhStruct{}
	src := body[10 : len(body)-2]
	err = json.Unmarshal(src, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type gzqhStruct struct {
	series1 interface{}
	series2 interface{}
	series3 interface{}
	a1      []string // 成交量龙虎榜
	a2      []string // 多头持仓龙虎榜
	a3      []string // 多头增仓龙虎榜
	a4      []string // 多头减仓龙虎榜
	a5      []string // 净多头龙虎榜
	a6      []string // 空头持仓龙虎榜
	a7      []string // 空头增仓龙虎榜
	a8      []string // 空头减仓龙虎榜
	a9      []string // 净空头龙虎榜
}

func (i *gzqhStruct) ParseItem() (string, float64, float64) {
	return "", 0.0, 0.0
}
