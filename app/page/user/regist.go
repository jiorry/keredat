package user

import (
	"github.com/jiorry/keredat/app/page/common"

	"github.com/kere/gos"
)

type Regist struct {
	gos.Page
}

// func (p *Regist) Befor() bool {
// 	p.View.Folder = "user"
// 	p.Cache.Type = gos.PAGE_CACHE_FILE
// 	return true
// }

func (p *Regist) Prepare() bool {
	p.Title = "用户注册"
	p.View.Folder = "user"
	common.SetupPage(&p.Page, "default")

	p.Layout.TopRenderList = nil

	return true
}
