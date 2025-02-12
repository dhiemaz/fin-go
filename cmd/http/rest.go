package http

import (
	"github.com/dhiemaz/fin-go/infrastructure/server"
	"github.com/dhiemaz/fin-go/infrastructure/server/router"
)

func Start() {
	router := router.NewRoute()
	server.Start(router.Register())
}
