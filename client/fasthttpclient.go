package client


import (
	"log"
	"time"
	"github.com/valyala/fasthttp"
	"github.com/fatih/color"
	"strconv"
	"crypto/tls"
)

var fasthttpClient = &fasthttp.Client{
	MaxConnsPerHost: 	       DefaultConnections,
	MaxIdleConnDuration:           DefaultConnectionTimeout,
	ReadTimeout:                   DefaultTimeout,
	WriteTimeout:                  DefaultTimeout,
	TLSConfig:		       &tls.Config{
		InsecureSkipVerify: true,
		ClientSessionCache: tls.NewLRUClientSessionCache(0),
	},
}


const (
	// DefaultTimeout is the default amount of time an Attacker waits for a request before it times out.
	DefaultTimeout = 60 * time.Second
	// DefaultConnections is the default amount of max open idle connections per target host.
	DefaultConnections = 15000
	DefaultConnectionTimeout = 30 * time.Second
)





func HitRequest(url string,method string,header map[string]string,payload string) []byte{
	req := buildRequest(url,method,header,payload)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttpClient.Do(req, resp)
	if err != nil {
		log.Printf("HTTP request failed: %s", err)
		return []byte("")
	}
	if(resp.StatusCode()!=200){
		reqUrl:=req.URI().String()
		color.Red(reqUrl+"------------>"+strconv.Itoa(resp.StatusCode()))
	}
	//reqUrl:=req.URI().String()
	response:=resp.Body()
	//color.Blue(reqUrl+"------------>"+strconv.Itoa(resp.StatusCode())+"------------>"+string(response))
	return response
}



func buildRequest(url string,method string,header map[string]string,payload string) *fasthttp.Request{
	var req = fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	for key, value := range header {
		req.Header.Add(key,value)
	}
	req.Header.SetMethod(method)
	if(payload!="") {
		req.SetBodyString(payload)
	}
	return req
}

