package home

import (
	"fmt"

	"github.com/kere/gos"
)

type Data struct {
	gos.Page
}

// func (p *Data) Befor() bool {
// 	p.View.Folder = "home"
// 	p.Cache.Type = gos.PAGE_CACHE_FILE
// 	return true
// }

func (p *Data) Prepare() bool {
	fmt.Println(p.Ctx.Request.Form)
	return false
}
