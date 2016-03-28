package lib

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jiorry/keredat/app/lib/util/ajax"
	"github.com/tealeg/xlsx"
)

var stockList map[string]*Stock

type Stock struct {
	Open float64
	Date string
}

var countWget = 0

func init() {
	stockList = make(map[string]*Stock, 0)
}

func WgetStock(code string) (*Stock, error) {
	if item, isOk := stockList[code]; isOk {
		fmt.Println("is found:", code)
		return item, nil
	}

	countWget++
	fmt.Println("countWget:", countWget)

	ajaxClient := ajax.NewAjax("")

	now := time.Now()
	ctype := "sz"
	if code[0:1] == "6" {
		ctype = "sh"
	}
	// http://hq.sinajs.cn/rn=1459142710888&list=sz002335
	str := fmt.Sprintf("http://hq.sinajs.cn/rn=%d&list=%s%s", now.Unix(), ctype, code)
	body, err := ajaxClient.GetBody(str)
	if err != nil {
		return nil, err
	}
	if len(body) < 23 {
		return nil, err
	}

	// var hq_str_sz002335="科华恒盛,39.58,39.18,40.71,42.00,39.58,40.71,40.72,3064511,126285436.73,4700,40.71,22300,40.70,800,40.68,1000,40.65,5200,40.64,700,40.72,1200,40.74,7825,40.75,100,40.90,1000,40.91,2016-03-28,14:03:11,00";
	//                         0      1      2    3     4     5     6     7      8         9         10   11    12    13   14   15   16   17     18   19   20   21    22   23   24    25   26   27   28    29     30         31    32

	l := bytes.Split(body[21:len(body)-1], []byte(","))
	if len(l) < 2 {
		return nil, err
	}
	stock := &Stock{}
	stock.Open, err = strconv.ParseFloat(string(l[1]), 64)
	if err != nil {
		return nil, err
	}
	stock.Date = string(l[30])
	stockList[code] = stock
	return stock, nil
}

func ReadExcel() ([]ExcelRow, error) {
	xlFile, err := xlsx.OpenFile("a.xlsx")
	if err != nil {
		return nil, err
	}
	sheet := xlFile.Sheets[0]
	var result = make([]ExcelRow, 0)

	fmt.Println("totle:", len(sheet.Rows))

	var item ExcelRow
	var prefix string
	for _, row := range sheet.Rows {
		item = ExcelRow{}
		item.Seq, _ = row.Cells[0].String()
		if item.Seq == "序号" {
			continue
		}

		item.Code, _ = row.Cells[1].String()
		prefix = item.Code[0:1]
		if prefix == "9" || prefix == "A" {
			continue
		}
		item.Name, _ = row.Cells[2].String()

		item.E0, _ = row.Cells[3].Float()
		item.E1, _ = row.Cells[4].Float()
		item.E2, _ = row.Cells[5].Float()
		item.E3, _ = row.Cells[6].Float()
		item = item.CalculateEPS()

		item.R0, _ = row.Cells[7].Float()
		item.R1, _ = row.Cells[8].Float()
		item.R2, _ = row.Cells[9].Float()
		item.R3, _ = row.Cells[10].Float()
		item = item.CalculateROE()

		item.PE, _ = row.Cells[11].Float()
		sv, _ := row.Cells[12].Float()
		if math.IsNaN(sv) {
			item.SV = "-"
		} else {
			item.SV = fmt.Sprintf("%.2f亿", sv/100000000)
		}

		result = append(result, item)
	}

	return result, nil
}
