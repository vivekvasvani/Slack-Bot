package server

import (
	"github.com/golang/glog"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)


// NewServer exposes api for admission control
func NewServer() {

	router := fasthttprouter.New()

	router.POST("/getavailable",func(ctx *fasthttp.RequestCtx) {
		postAdmissionResults(ctx)
	})




	router.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
		glog.V(0).Infof("Panic occurred %s", p,ctx.Request.URI().String())
	}
	log.Println("Service Started on port " + "5498")
	glog.Fatal(fasthttp.ListenAndServe(":" + "5498", fasthttp.CompressHandler(router.Handler)))

}

