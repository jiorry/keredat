package main

import (
	"flag"

	"github.com/jiorry/keredat/app/api"
	"github.com/jiorry/keredat/app/lib/runner"
	"github.com/jiorry/keredat/app/page/home"
	"github.com/jiorry/keredat/app/page/rzrq"
	"github.com/jiorry/keredat/app/page/user"
	"github.com/kere/gos"
	_ "github.com/lib/pq"
)

func main() {
	flag.Parse()
	gos.Init()

	gos.Route("/", &home.Default{})
	gos.Route("/data", &home.Data{})
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

	runner.RunTimer()

	gos.Start()
}
