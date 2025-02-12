package router

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Route struct {
}

func NewRoute() *Route {
	return &Route{}
}

func (r *Route) Register() *router.Router {

	route := router.New()
	route.GET("/", func(ctx *fasthttp.RequestCtx) {
		return
	})

	return route
}
