package main

import (
	"fmt"

	"github.com/jiorry/keredat/app/api"
	"github.com/jiorry/keredat/app/lib/runner"
	"github.com/jiorry/keredat/app/lib/util"
	"github.com/jiorry/keredat/app/page/home"
	"github.com/jiorry/keredat/app/page/rzrq"
	"github.com/jiorry/keredat/app/page/user"
	"github.com/kere/gos"
	_ "github.com/lib/pq"
)

func main() {
	gos.Init()

	gos.Route("/", &home.Default{})
	gos.Route("/login", &user.Login{})
	gos.Route("/regist", &user.Regist{})

	gos.Route("/rzrq/sum", &rzrq.Sum{})
	gos.Route("/rzrq/stock/:code", &rzrq.Stock{})

	// open api router
	// gos.WebApiRoute("web", &api.Public{})

	// open api
	// api.RegistOpenApi()
	gos.WebApiRoute("open", &api.OpenApi{})

	// websocket router
	// gos.WebSocketRoute("conn", (*hiuser.UserWebSock)(nil))

	err := runner.RunTimer()
	if err != nil {
		util.SendEmail("keredat err", fmt.Sprint(err))
		return
	}

	gos.Start()
}
