package hx50etf

// 华夏基金ETF

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/kere/gos/db"

	"github.com/jiorry/keredat/app/lib/util"
	"github.com/kere/gos"
)

var codeListCurentUp [][]byte
var codeListCurentDown [][]byte
var codeListNextUp [][]byte
var codeListNextDown [][]byte

var ajax *util.Ajax
var hx50ETFIndexData *hx50ETFIndex
var hx50ETFDatasetCurrent []*hx50ETF
var hx50ETFDatasetNext []*hx50ETF

func init() {
	ajax = util.NewAjax("")
	hx50ETFIndexData = &hx50ETFIndex{}
}

// StoreTodayETFData store etf data
func StoreTodayETFData() error {
	if hx50ETFIndexData != nil && hx50ETFIndexData.CreatedAt.Format("20060102") == gos.NowInLocation().Format("20060102") {
		return nil
	}
	gos.Log.Info("StoreTodayETFData")
	var err error
	hx50ETFIndexData, hx50ETFDatasetCurrent, err = fetchTodayETFData(etfDateString(0))
	if err != nil {
		return err
	}
	err = storeETF(hx50ETFDatasetCurrent)
	if err != nil {
		return err
	}
	_, hx50ETFDatasetNext, err = fetchTodayETFData(etfDateString(1))
	if err != nil {
		return err
	}
	return storeETF(hx50ETFDatasetNext)
}

func storeETF(dataset []*hx50ETF) error {
	for _, item := range dataset {
		if db.NewExistsBuilder(item.Table()).Where("name=? and date=?", item.Name, item.Date).Exists() {
			continue
		}

		item.Init(item)
		err := item.Create()
		if err != nil {
			return err
		}
	}
	return nil
}

// fetchTodayETFData store etf data
func fetchTodayETFData(monthStr string) (*hx50ETFIndex, []*hx50ETF, error) {
	indexData, listUp, listDown, err := prepareETF(monthStr)
	if err != nil {
		return nil, nil, err
	}

	dataset, err := fetchETFItems(indexData, listUp)
	if err != nil {
		return nil, nil, err
	}

	datasetDown, err := fetchETFItems(indexData, listDown)
	if err != nil {
		return nil, nil, err
	}
	dataset = append(dataset, datasetDown...)

	return indexData, dataset, nil
}

// fetchETFItems
func fetchETFItems(indexData *hx50ETFIndex, list [][]byte) ([]*hx50ETF, error) {
	// http://hq.sinajs.cn/list=
	// var hq_str_CON_OP_10000326="5,0.0233,0.0234,0.0237,1,19353,-24.27,2.3000,0.0309,0.0360,0.2565,0.0001,0.0244,10,0.0243,20,0.0239,10,0.0238,2,0.0237,1,0.0233,5,0.0231,10,0.0228,20,0.0226,20,0.0225,10,2015-12-14 11:33:59,1,T0,EBS,510050,50ETF沽12月2300,59.55,0.0388,0.0204,5323,1597707.00";
	//							   买量,买价,  最新 ,卖价 ,卖量,持仓量，涨幅  ，行权价，前结   ，今开  ，涨停  ，跌停 ，卖5   ， ，卖四,  量 ，卖三  ， ，卖二，  ，卖一  ， ,买1   ， ，买2  ，  ，买3   ， ， 买4 ， ，买5   ，  ，日期              ， ，  ，  ，     ,              ， 振幅 ，最高  ， 最低 ,总量 ，金额
	// 								0 , 1  ,  2   , 3  , 4 ,   5 ,  6   , 7   ,   8    ,  9   , 10   ,  11 , 12   ,13,14  ,15  , 16   , 17, 18, 19,  20  ,21, 22  ,23,24   , 25, 26  , 27,  28, 29, 30   , 31,32                ,33,34,35, 36    , 37           ,38    , 39   , 40   ,41  ,42

	// var hq_str_CON_SO_10000326="50ETF沽12月2300,,,,5323,-0.2876,3.4114,-0.6633,0.1259,0.3099,0.0388,0.0204,510050P1512M02300,2.3000,0.0234,0.0184";
	// 								                 ,总量,Delta  ,Gamma  ,Theta, Vega  ，波动率 ，最高  ，最低  ，                 ，     ，时间价值，理论价值
	// 							          0 ,1,2,3,   4 ,   5   ,  6   ,    7  ,  8    ,    9  ,  10  ,  11  ,  12             , 13   , 14    ,  15
	// var hq_str_CON_ZL_10000326="50ETF沽12月2300,,,,19353,0.14%,0.0234,-86.32%,0.0233,0.0237,0.0388,0.0204,5323,-0.2876,3.4114,-0.6633,0.1259,0.3099,510050P1512M02300,2.3000,0.0184";
	// 								                 ,持仓 ，占比 ，最新 ，跌幅    ，买价  ， 卖价 ， 最高 ， 最低 ，成交量，delta,gamma,theta   ,vega  ,波动率 ,                 ,      , 理论价值
	var dataset = make([]*hx50ETF, 0)
	var listStrOP = make([]string, 0)
	var listStrSO = make([]string, 0)

	for _, b := range list {
		listStrOP = append(listStrOP, fmt.Sprint("CON_OP_", string(b)))
	}
	if len(listStrOP) == 0 {
		return nil, gos.DoError("listStrOP is empty")
	}
	body, err := ajax.GetBody(fmt.Sprintf("http://hq.sinajs.cn/list=%s", strings.Join(listStrOP, ",")))
	if err != nil {
		return nil, gos.DoError(err)
	}

	for _, b := range list {
		listStrSO = append(listStrSO, fmt.Sprint("CON_SO_", string(b)))
	}
	bodySO, err := ajax.GetBody(fmt.Sprintf("http://hq.sinajs.cn/list=%s", strings.Join(listStrSO, ",")))
	if err != nil {
		return nil, gos.DoError(err)
	}

	lOP := bytes.Split(body, []byte("\";"))
	lSO := bytes.Split(bodySO, []byte("\";"))
	isFound := false
	for _, bCode := range list {
		data := &hx50ETF{}
		isFound = false
		for _, lop := range lOP {
			// 查找 var hq_str_CON_OP_10000326 的行
			if bytes.Index(lop, []byte(fmt.Sprint("var hq_str_CON_OP_", string(bCode)))) > 0 {
				// 截取 =" 之后的字符，并把它们按照逗号分割成数组
				arr := bytes.Split(lop[bytes.Index(lop, []byte("=\""))+2:], gos.B_COMMA)
				if len(arr) < 1 {
					return nil, gos.DoError("parse CON_OP string list error")
				}
				data.IndexID = indexData.ID
				data.Date = string(arr[32])
				data.Exec = util.ParseMoney(string(arr[7]))
				data.Price = util.ParseMoney(string(arr[2]))
				data.Prev = util.ParseMoney(string(arr[8]))
				data.Open = util.ParseMoney(string(arr[9]))
				data.OpenInt = int64(util.ParseMoney(string(arr[5])))
				data.ChangePercent = util.ParseMoney(string(arr[38]))
				data.High = util.ParseMoney(string(arr[39]))
				data.Low = util.ParseMoney(string(arr[40]))
				data.Amount = int64(util.ParseMoney(string(arr[41])))
				data.Total = int64(util.ParseMoney(string(arr[42])))

				isFound = true
				break
			}
		}

		for _, lso := range lSO {
			// 查找 var hq_str_CON_SO_10000326 的行
			if bytes.Index(lso, []byte(fmt.Sprint("var hq_str_CON_SO_", string(bCode)))) > 0 {
				// 截取 =" 之后的字符，并把它们按照逗号分割成数组
				arr := bytes.Split(lso[bytes.Index(lso, []byte("=\""))+2:], gos.B_COMMA)
				if len(arr) < 1 {
					return nil, gos.DoError("parse CON_OP string list error")
				}

				data.Name = string(arr[12])
				// 510050P1512M02300
				data.EndOn = etfEndOn(data.Name[7:11])
				data.Delta = util.ParseMoney(string(arr[5]))
				data.Gamma = util.ParseMoney(string(arr[6]))
				data.Vega = util.ParseMoney(string(arr[7]))
				data.Theta = util.ParseMoney(string(arr[8]))
				data.VIX = util.ParseMoney(string(arr[9]))
				data.TimeValue = util.ParseMoney(string(arr[14]))
				data.TheoValue = util.ParseMoney(string(arr[15]))

				isFound = true
				break
			}
		}

		if isFound {
			dataset = append(dataset, data)
		}
	}
	return dataset, nil
}

// prepareETF 抓取数据
func prepareETF(monthStr string) (*hx50ETFIndex, [][]byte, [][]byte, error) {
	// http://hq.sinajs.cn/list=OP_UP_5100501512,OP_DOWN_5100501512,s_sh510050,sh510050
	url := "http://hq.sinajs.cn/list=OP_UP_510050%s,OP_DOWN_510050%s,sh510050"

	body, err := ajax.GetBody(fmt.Sprintf(url, monthStr, monthStr))
	if err != nil {
		return nil, nil, nil, gos.DoError(err)
	}

	arr := bytes.Split(body, []byte(",\";"))
	if len(arr) < 0 {
		return nil, nil, nil, gos.DoError("arr length must > 0")
	}

	n := bytes.Index(arr[0], []byte("=\""))
	list := bytes.Split(arr[0][n+2:], gos.B_COMMA)

	listUp := make([][]byte, len(list))
	for i, _ := range list {
		// CON_OP_10000393
		listUp[i] = list[i][7:]
	}

	n = bytes.Index(arr[1], []byte("=\""))
	list = bytes.Split(arr[1][n+2:], gos.B_COMMA)
	listDown := make([][]byte, len(list))
	for i, _ := range list {
		// CON_OP_10000393
		listDown[i] = list[i][7:]
	}

	// var hq_str_s_sh510050="50ETF,2.344,-0.011,-0.47,2490440,58266";
	// 收盘价，涨跌，涨幅，总量，金额（万）
	// var hq_str_sh510050="50ETF,2.323,2.344,2.415,2.420,2.319,2.416,2.417,371982384,879489760,91750,2.416,362300,2.415,443600,2.414,4054000,2.413,129600,2.412,10200,2.417,135100,2.418,733100,2.419,1054400,2.420,35200,2.421,2015-12-14,15:03:11,00";
	// 						     , 今开 ，昨收 ，现价 ，最高 ，最低  ，买价，卖价 ，  总手   ，  金额   ，买1 ，     ， 买2  ，    ， 买3  ，    ，  买4  ，    ，  买5 ，    ，     ，    ，     ，     ，     ，    ，      ，     ，    ，    ，   日期    ，  时间  ，
	//                        0  ,   1 ,  2  ,  3   ,  4 ,  5   ,  6 ,  7  ,    8    ,    9    ,  10 ,  11 ,  12  ,  13 , 14   , 15  ,   16  ,  17 ,   18 ,  19 ,  20 ,  21 ,  22  ,  23  ,  24 ,  25 ,  26  ,   27 ,  28 ,  29 ,  30       ,  31    ,
	n = bytes.Index(arr[2], []byte("=\""))
	arr = bytes.Split(arr[2][n+2:], gos.B_COMMA)

	hx50ETFIndexData.Open = util.ParseMoney(string(arr[1]))
	hx50ETFIndexData.Prev = util.ParseMoney(string(arr[2]))
	hx50ETFIndexData.Price = util.ParseMoney(string(arr[3]))
	hx50ETFIndexData.High = util.ParseMoney(string(arr[4]))
	hx50ETFIndexData.Low = util.ParseMoney(string(arr[5]))
	hx50ETFIndexData.Amount = int64(util.ParseMoney(string(arr[8])))
	hx50ETFIndexData.Total = int64(util.ParseMoney(string(arr[9])))
	hx50ETFIndexData.Date = string(arr[30])
	hx50ETFIndexData.Time = string(arr[31])

	if !db.NewExistsBuilder(hx50ETFIndexData.Table()).Where("date=?", hx50ETFIndexData.Date).Exists() {
		hx50ETFIndexData.Init(hx50ETFIndexData)
		hx50ETFIndexData.Create()
	}

	vo, err := db.NewQueryBuilder(hx50ETFIndexData.Table()).Where("date=?", hx50ETFIndexData.Date).Struct(hx50ETFIndexData).FindOne()
	if err != nil {
		return nil, nil, nil, gos.DoError(err)
	}

	return vo.(*hx50ETFIndex), listUp, listDown, nil
}

func etfEndOn(s string) string {
	d, err := time.ParseInLocation("20060102", fmt.Sprint("20", s, "01"), gos.GetSite().Location)
	if err != nil {
		gos.DoError(err)
		return ""
	}

	d = d.AddDate(0, 1, 0)
	for {
		d = d.AddDate(0, 0, -1)
		if d.Weekday() == 5 {
			break
		}
	}
	// 最后一个礼拜的周三，为期权交割日
	return d.AddDate(0, 0, -2).Format("2006-01-02")
}

func etfDateString(itype int) string {
	l := []int{1, 3, 6, 9, 12}
	now := gos.NowInLocation()
	year := now.Year()
	month := int(now.Month())
	index := 0

	for i, v := range l {
		if month <= v {
			index = i
			break
		}
	}

	if index+itype >= len(l) {
		index = 0
		year++
	}

	return fmt.Sprintf("%d%02d", year, l[index])[2:]
}

type hx50ETFIndex struct {
	db.BaseVO
	ID        int64     `json:"id" skip:"all"`
	Price     float64   `json:"price"`
	Open      float64   `json:"open"`
	Prev      float64   `json:"prev"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Amount    int64     `json:"amount"`
	Total     int64     `json:"total"`
	Date      string    `json:"date"`
	Time      string    `json:"etime"`
	CreatedAt time.Time `json:"created_at" autotime:"true" skip:"update"`
}

func (a *hx50ETFIndex) Table() string {
	return "hx50etf_index"
}

type hx50ETF struct {
	db.BaseVO
	IndexID       int64   `json:"index_id" skip:"update"`
	Name          string  `json:"name"`
	EndOn         string  `json:"end_on"`
	Exec          float64 `json:"exec"`
	Price         float64 `json:"price"`
	Prev          float64 `json:"prev"`
	Open          float64 `json:"open"`
	OpenInt       int64   `json:"open_int"`
	ChangePercent float64 `json:"change_percent"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Amount        int64   `json:"amount"`
	Total         int64   `json:"total"`
	Date          string  `json:"date"`

	Delta     float64   `json:"delta"`
	Gamma     float64   `json:"gamma"`
	Vega      float64   `json:"vega"`
	Theta     float64   `json:"theta"`
	VIX       float64   `json:"vix"`
	TimeValue float64   `json:"time_value"`
	TheoValue float64   `json:"theo_value"`
	CreatedAt time.Time `json:"created_at" autotime:"true" skip:"update"`
}

func (a *hx50ETF) Table() string {
	return "hx50etf_items"
}
