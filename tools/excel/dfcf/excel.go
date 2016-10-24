package main

import (
	"fmt"
	"math"

	"github.com/jiorry/keredat/tools/excel/dfcf/lib"
	"github.com/kere/gos"
	"github.com/tealeg/xlsx"
)

func main() {
	gos.InitLog("")

	var rateLimit = float64(1.1)
	lib.InitStyle()

	result, err := lib.ReadExcel()
	if len(result) == 0 {
		fmt.Println("result is empty")
		return
	}

	file := xlsx.NewFile()
	// sheet EPS --------------------------------
	var itemsEPS = make([]lib.ExcelRow, 0)
	for _, item := range result {
		if math.IsNaN(item.RateE) || item.RateENow < rateLimit || item.RateE < rateLimit {
			continue
		}

		item = item.WgetStockData()
		itemsEPS = append(itemsEPS, item)
	}

	funcRateE := func(p1, p2 *lib.ExcelRow) bool {
		return p1.RateE > p2.RateE
	}
	lib.SortBy(funcRateE).Sort(itemsEPS)
	sheet1, _ := file.AddSheet(fmt.Sprint("EPS-", len(itemsEPS)))
	lib.FillExcel(sheet1, itemsEPS, "sheet1")

	// sheet ROE --------------------------------
	var itemsROE = make([]lib.ExcelRow, 0)
	for _, item := range result {
		if math.IsNaN(item.RateR) || item.RateRNow < rateLimit || item.RateR < rateLimit {
			continue
		}

		item = item.WgetStockData()
		itemsROE = append(itemsROE, item)
	}

	funcRateR := func(p1, p2 *lib.ExcelRow) bool {
		return p1.RateR > p2.RateR
	}
	lib.SortBy(funcRateR).Sort(itemsROE)

	sheet2, _ := file.AddSheet(fmt.Sprint("ROE-", len(itemsROE)))
	lib.FillExcel(sheet2, itemsROE, "sheet2")

	// sheet ROE && EPS --------------------------------
	var itemsAll = make([]lib.ExcelRow, 0)
	for _, eps := range itemsEPS {
		for _, roe := range itemsROE {
			if eps.Code == roe.Code {
				itemsAll = append(itemsAll, eps)
				break
			}
		}
	}
	sheet3, _ := file.AddSheet(fmt.Sprint("BOTH-", len(itemsAll)))
	lib.FillExcel(sheet3, itemsAll, "sheet3")
	// Save file ---------------------------
	err = file.Save(fmt.Sprint("a-finish.xlsx"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Save File.")
}
