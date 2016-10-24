package lib

import "math"

type ExcelRow struct {
	Seq      string  `name:"序号" sheet1:"0" sheet2:"0" sheet3:"0"`
	Code     string  `name:"代码" sheet1:"1" sheet2:"1" sheet3:"1"`
	Name     string  `name:"名称" sheet1:"2" sheet2:"2" sheet3:"2"`
	EpsNext  float64 `name:"明年预估EPS" sheet1:"3" sheet3:"3" type:"float"`
	RateENow float64 `name:"明年复合增长率" sheet1:"4" sheet3:"4" type:"float"`
	Eps      float64 `name:"当年预估EPS" sheet1:"6" sheet3:"6" style:"grey" type:"float"`
	RateE    float64 `name:"复合增长率" sheet1:"7" sheet3:"7" type:"float"`
	E0       float64 `name:"实际EPS" sheet1:"5" sheet3:"5" type:"float"`
	E1       float64 `name:"EPS-1" sheet1:"8" type:"float"`
	E2       float64 `name:"EPS-2" sheet1:"9" type:"float"`
	E3       float64 `name:"EPS-3" sheet1:"10" type:"float"`

	RoeNext  float64 `name:"明年预估ROE" sheet2:"3" sheet3:"8" type:"float"`
	RateRNow float64 `name:"明年复合增长率" sheet2:"4" sheet3:"9" type:"float"`
	Roe      float64 `name:"当年预估ROE" sheet2:"6" sheet3:"11" style:"grey" type:"float"`
	RateR    float64 `name:"复合增长率" sheet2:"7" sheet3:"12" type:"float"`
	R0       float64 `name:"实际ROE" sheet2:"5" sheet3:"10" type:"float"`
	R1       float64 `name:"ROE-1" sheet2:"8" type:"float"`
	R2       float64 `name:"ROE-2" sheet2:"9" type:"float"`
	R3       float64 `name:"ROE-3" sheet2:"10" type:"float"`

	PE float64 `name:"动态市盈率" sheet1:"11" sheet2:"11" sheet3:"11" type:"float"`
	SV string  `name:"流通市值" sheet1:"12" sheet2:"12" sheet3:"12"`

	Open  float64 `name:"开盘价" sheet1:"13" sheet2:"13" sheet3:"13" type:"float"`
	Date  string  `name:"价格时间" sheet1:"14" sheet2:"14" sheet3:"14"`
	YYYGp float64 `name:"业绩预告变动%" sheet1:"15" sheet2:"15" sheet3:"15" type:"float"`
	YYYG  string  `name:"业绩预告" sheet1:"16" sheet2:"16" sheet3:"16"`
}

func (e ExcelRow) WgetStockData() ExcelRow {
	if e.Code == "" || e.Open > 0 {
		return e
	}

	stock, err := WgetStock(e.Code)
	if err != nil || stock == nil {
		return e
	}

	e.Open = stock.Open
	e.Date = stock.Date
	return e
}

func (e ExcelRow) CalculateEPS() ExcelRow {
	n := 1.00 / 3.00
	if math.IsNaN(e.E2) || math.IsNaN(e.E3) || e.E1 <= 0 || e.E2 <= 0 || e.E3 <= 0 ||
		e.E1 < e.E2 || e.E2 < e.E3 || e.E1 < e.E3 ||
		!math.IsNaN(e.E0) && e.E0 < e.E1 {
		e.RateENow = math.NaN()
		e.EpsNext = math.NaN()
		e.RateE = math.NaN()
		e.Eps = math.NaN()
		return e
	}

	// 复合增长率
	if math.IsNaN(e.E0) {
		e.RateENow = math.NaN()
		e.EpsNext = math.NaN()
	} else {
		e.RateENow = round(math.Pow(e.E0/e.E2, n)) // 明年
		e.EpsNext = round(e.E0 * e.RateENow)       // 明年
	}

	e.RateE = round(math.Pow(e.E1/e.E3, n)) // 当年
	e.Eps = round(e.E1 * e.RateE)           // 当年

	return e
}

func (e ExcelRow) CalculateROE() ExcelRow {
	n := 1.00 / 3.00
	if math.IsNaN(e.R2) || math.IsNaN(e.R3) || e.R1 <= 0 || e.R2 <= 0 || e.R3 <= 0 ||
		e.R1 < e.R2 || e.R2 < e.R3 || e.R1 < e.R3 ||
		!math.IsNaN(e.R0) && e.R0 < e.R1 {
		e.RateRNow = math.NaN()
		e.RoeNext = math.NaN()
		e.RateR = math.NaN()
		e.Roe = math.NaN()
		return e
	}

	// 复合增长率
	if math.IsNaN(e.R0) {
		e.RateRNow = math.NaN()
		e.RoeNext = math.NaN()
	} else {
		e.RateRNow = round(math.Pow(e.R0/e.R2, n)) // 明年
		e.RoeNext = round(e.R0 * e.RateRNow)       // 明年
	}

	e.RateR = round(math.Pow(e.R1/e.R3, n)) // 当年
	e.Roe = round(e.R1 * e.RateR)           // 当年
	return e
}
func round(v float64) float64 {
	return math.Floor(v*10000+0.5) / 10000.0
	// return v
}
