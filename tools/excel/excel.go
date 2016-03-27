package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/tealeg/xlsx"
)

type ExcelRow struct {
	Seq      string
	Code     string
	Name     string
	EpsNext  float64
	RateENow float64
	Eps      float64
	RateE    float64
	E0       float64
	E1       float64
	E2       float64
	E3       float64

	RoeNext  float64
	RateRNow float64
	Roe      float64
	RateR    float64
	R0       float64
	R1       float64
	R2       float64
	R3       float64
}

func (e ExcelRow) calculateEPS() ExcelRow {
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
		e.RateENow = math.Pow(e.E0/e.E2, n) // 明年
		e.EpsNext = e.E0 * e.RateENow       // 明年
	}

	e.RateE = math.Pow(e.E1/e.E3, n) // 当年
	e.Eps = e.E1 * e.RateE           // 当年

	return e
}

func (e ExcelRow) calculateROE() ExcelRow {
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
		e.RateRNow = math.Pow(e.R0/e.R2, n) // 明年
		e.RoeNext = e.R0 * e.RateRNow       // 明年
	}

	e.RateR = math.Pow(e.R1/e.R3, n) // 当年
	e.Roe = e.R1 * e.RateR           // 当年
	return e
}

func readExcel() ([]ExcelRow, error) {
	xlFile, err := xlsx.OpenFile("a.xlsx")
	if err != nil {
		return nil, err
	}
	sheet := xlFile.Sheets[0]
	result := make([]ExcelRow, 0)

	fmt.Println("totle:", len(sheet.Rows))

	var item ExcelRow
	for _, row := range sheet.Rows {
		item = ExcelRow{}
		item.Seq, _ = row.Cells[0].String()
		item.Code, _ = row.Cells[1].String()
		item.Name, _ = row.Cells[2].String()

		item.E0, _ = row.Cells[3].Float()
		item.E1, _ = row.Cells[4].Float()
		item.E2, _ = row.Cells[5].Float()
		item.E3, _ = row.Cells[6].Float()
		item = item.calculateEPS()

		item.R0, _ = row.Cells[7].Float()
		item.R1, _ = row.Cells[8].Float()
		item.R2, _ = row.Cells[9].Float()
		item.R3, _ = row.Cells[10].Float()
		item = item.calculateROE()

		result = append(result, item)
	}

	return result, nil
}

var greyStyle *xlsx.Style
var defaultStyle *xlsx.Style
var headStyle *xlsx.Style
var rateLimit float64 = 1.1

func main() {
	border := xlsx.NewBorder("thin", "thin", "thin", "thin")
	border.BottomColor = "FFAAAAAA"
	border.TopColor = "FFAAAAAA"
	border.RightColor = "FFAAAAAA"
	border.LeftColor = "FFAAAAAA"

	defaultStyle = xlsx.NewStyle()
	defaultStyle.ApplyBorder = true
	defaultStyle.Border = *border

	greyStyle = xlsx.NewStyle()
	greyStyle.ApplyBorder = true
	greyStyle.Border = *border
	greyStyle.ApplyFill = true
	greyStyle.Fill = *xlsx.NewFill("solid", "FFEEEEEE", "")

	headStyle = xlsx.NewStyle()
	headStyle.ApplyBorder = true
	headStyle.Border = *border
	headStyle.ApplyFill = true
	headStyle.Fill = *xlsx.NewFill("solid", "FFD9D9D9", "")

	result, err := readExcel()
	if len(result) == 0 {
		fmt.Println("result is empty")
		return
	}

	file := xlsx.NewFile()
	// sheet EPS --------------------------------
	itemsEPS := make([]ExcelRow, 0)
	for _, item := range result {
		if math.IsNaN(item.RateE) || item.RateENow < rateLimit || item.RateE < rateLimit {
			continue
		}

		itemsEPS = append(itemsEPS, item)
	}

	funcRateE := func(p1, p2 *ExcelRow) bool {
		return p1.RateE > p2.RateE
	}
	SortBy(funcRateE).Sort(itemsEPS)
	sheet1, _ := file.AddSheet(fmt.Sprint("EPS-", len(itemsEPS)))
	fillEPS(sheet1, itemsEPS)

	// sheet ROE --------------------------------
	itemsROE := make([]ExcelRow, 0)
	for _, item := range result {
		if math.IsNaN(item.RateR) || item.RateRNow < rateLimit || item.RateR < rateLimit {
			continue
		}

		itemsROE = append(itemsROE, item)
	}

	funcRateR := func(p1, p2 *ExcelRow) bool {
		return p1.RateR > p2.RateR
	}
	SortBy(funcRateR).Sort(itemsROE)

	sheet2, _ := file.AddSheet(fmt.Sprint("ROE-", len(itemsROE)))
	fillROE(sheet2, itemsROE)

	// sheet ROE && EPS --------------------------------
	itemsAll := make([]ExcelRow, 0)
	for _, eps := range itemsEPS {
		for _, roe := range itemsROE {
			if eps.Code == roe.Code {
				itemsAll = append(itemsAll, eps)
				break
			}
		}
	}
	sheet3, _ := file.AddSheet(fmt.Sprint("BOTH-", len(itemsAll)))
	fillBoth(sheet3, itemsAll)
	// Save file ---------------------------
	err = file.Save(fmt.Sprint("a-finish.xlsx"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Save File.")
}

func fillEPS(sheet *xlsx.Sheet, result []ExcelRow) {
	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("序号")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("代码")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("名称")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("明年预估EPS")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("最新复合增长率")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("当年实际EPS")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("预估当年EPS")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("复合增长率")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("EPS-1")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("EPS-2")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("EPS-3")

	for _, item := range result {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Seq)
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Code)
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Name)

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.EpsNext) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.EpsNext, "0.0000") // 明年预估EPS
		}
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.RateENow) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.RateENow, "0.0000") // 明年复合增长率
		}
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.E0) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.E0, "0.0000") // 当年实际EPS
		}
		cell = row.AddCell()
		cell.SetStyle(greyStyle)
		cell.SetFloatWithFormat(item.Eps, "0.0000") // 当年估算EPS

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.RateE, "0.0000") // 当年复合增长率

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.E1, "0.0000")
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.E2, "0.0000")
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.E3, "0.0000")
	}

}

func fillROE(sheet *xlsx.Sheet, result []ExcelRow) {
	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("序号")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("代码")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("名称")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("明年预估ROE")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("最新复合增长率")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("当年实际ROE")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("预估当年ROE")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("复合增长率")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("ROE-1")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("ROE-2")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("ROE-3")

	for _, item := range result {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Seq)
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Code)
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Name)

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.RoeNext) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.RoeNext, "0.0000") // 明年预估ROE
		}
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.RateRNow) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.RateRNow, "0.0000") // 明年复合增长率
		}
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.R0) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.R0, "0.0000") // 当年实际ROE
		}

		cell = row.AddCell()
		cell.SetStyle(greyStyle)
		cell.SetFloatWithFormat(item.Roe, "0.0000") // 当年估算ROE

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.RateR, "0.0000") // 当年复合增长率

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.R1, "0.0000")
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.R2, "0.0000")
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.R3, "0.0000")
	}

}

type SortBy func(p1, p2 *ExcelRow) bool

func (by SortBy) Sort(rows []ExcelRow) {
	ps := &Sorter{
		rows: rows,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// Sorter joins a By function and a slice of Planets to be sorted.
type Sorter struct {
	rows []ExcelRow
	by   func(p1, p2 *ExcelRow) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *Sorter) Len() int {
	return len(s.rows)
}

// Swap is part of sort.Interface.
func (s *Sorter) Swap(i, j int) {
	s.rows[i], s.rows[j] = s.rows[j], s.rows[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *Sorter) Less(i, j int) bool {
	return s.by(&s.rows[i], &s.rows[j])
}

func fillBoth(sheet *xlsx.Sheet, result []ExcelRow) {
	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("序号")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("代码")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("名称")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("明年预估EPS")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("最新复合增长率")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("当年实际EPS")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("预估当年EPS")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("复合增长率")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("明年预估ROE")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("最新复合增长率")

	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("当年实际ROE")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("预估当年ROE")
	cell = row.AddCell()
	cell.SetStyle(headStyle)
	cell.SetString("复合增长率")

	for _, item := range result {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Seq)
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Code)
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetString(item.Name)

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.EpsNext) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.EpsNext, "0.0000") // 明年预估EPS
		}
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.RateENow) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.RateENow, "0.0000") // 明年复合增长率
		}
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.E0) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.E0, "0.0000") // 当年实际EPS
		}
		cell = row.AddCell()
		cell.SetStyle(greyStyle)
		cell.SetFloatWithFormat(item.Eps, "0.0000") // 当年估算EPS

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.RateE, "0.0000") // 当年复合增长率

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.RoeNext) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.RoeNext, "0.0000") // 明年预估ROE
		}
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.RateRNow) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.RateRNow, "0.0000") // 明年复合增长率
		}
		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		if math.IsNaN(item.R0) {
			cell.SetString("-")
		} else {
			cell.SetFloatWithFormat(item.R0, "0.0000") // 当年实际ROE
		}

		cell = row.AddCell()
		cell.SetStyle(greyStyle)
		cell.SetFloatWithFormat(item.Roe, "0.0000") // 当年估算ROE

		cell = row.AddCell()
		cell.SetStyle(defaultStyle)
		cell.SetFloatWithFormat(item.RateR, "0.0000") // 当年复合增长率
	}

}
