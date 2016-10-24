package dfcf

import (
	"github.com/jiorry/keredat/app/lib/util/ajax"
)

var (
	ajaxClient *ajax.Ajax
)

func init() {
	ajaxClient = ajax.NewAjax("")
}
