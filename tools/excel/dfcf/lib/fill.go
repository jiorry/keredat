package lib

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/tealeg/xlsx"
)

// fill --------------------------------------
var greyStyle *xlsx.Style
var defaultStyle *xlsx.Style
var headStyle *xlsx.Style

func InitStyle() {
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
}

func FillExcel(sheet *xlsx.Sheet, result []ExcelRow, ctype string) {
	if len(result) == 0 {
		return
	}

	typ := reflect.TypeOf(result[0])

	var cell *xlsx.Cell
	var row *xlsx.Row
	var fieldName string
	var indexStr string
	var index int

	// fill header
	row = sheet.AddRow()
	num := typ.NumField()
	var fieldNum = 0
	for i := 0; i < num; i++ {
		if indexStr = typ.Field(i).Tag.Get(ctype); indexStr == "" {
			continue
		}
		fieldNum++
		row.AddCell()
	}

	for i := 0; i < num; i++ {
		if indexStr = typ.Field(i).Tag.Get(ctype); indexStr == "" {
			continue
		}

		if fieldName = typ.Field(i).Tag.Get("name"); fieldName == "" {
			continue
		}
		index, _ = strconv.Atoi(indexStr)
		cell = row.Cells[index]
		cell.SetStyle(headStyle)
		cell.SetString(fieldName)
	}

	var val reflect.Value
	var field reflect.StructField

	// fill item
	for _, item := range result {
		val = reflect.ValueOf(item)
		typ = val.Type()
		row = sheet.AddRow()
		for i := 0; i < fieldNum; i++ {
			row.AddCell()
		}

		for i := 0; i < num; i++ {
			field = typ.Field(i)

			if indexStr = field.Tag.Get(ctype); indexStr == "" {
				continue
			}
			index, _ = strconv.Atoi(indexStr)
			cell = row.Cells[index]

			switch field.Tag.Get("style") {
			case "grey":
				cell.SetStyle(greyStyle)
			default:
				cell.SetStyle(defaultStyle)
			}

			switch field.Tag.Get("type") {
			case "float":
				v := val.Field(i).Interface().(float64)
				if math.IsNaN(v) {
					cell.SetString("-")
				} else {
					// cell.SetString(fmt.Sprintf("%.4f", v))
					cell.SetFloat(v)
				}
			default:
				cell.SetString(fmt.Sprint(val.Field(i).Interface()))
			}

		}
	}
}
