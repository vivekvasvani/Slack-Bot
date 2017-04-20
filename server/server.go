package server

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/golang/glog"
	"github.com/valyala/fasthttp"
)

func NewServer() {

	router := fasthttprouter.New()

	router.POST("/", func(ctx *fasthttp.RequestCtx) {
		getSTGStatus(ctx)
	})

	router.POST("/slack/sendoptions", func(ctx *fasthttp.RequestCtx) {
		sendOptions(ctx)
	})

	router.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
		glog.V(0).Infof("Panic occurred %s", p, ctx.Request.URI().String())
	}
	log.Println("Service Started on port " + "5498")
	glog.Fatal(fasthttp.ListenAndServe(":"+"5498", fasthttp.CompressHandler(router.Handler)))

}
