package main

import (
	"fmt"

	"github.com/jiorry/keredat/app/lib/tools/sina/hx50etf"
	"github.com/kere/gos"
	_ "github.com/lib/pq"
)

func main() {
	gos.ConfName = "../../app/app.conf"
	gos.Init()
	err := hx50etf.StoreTodayETFData()
	fmt.Println(err)

}
