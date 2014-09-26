package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/html"
	//"github.com/elazarl/goproxy/ext/image"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	verbose := flag.Bool("verbose", true, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8888", "proxy listen address")
	flag.Parse()

	fmt.Printf("verbose: %t\naddress: %s\n\n", *verbose, *addr)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	// changes request headers
	/*proxy.OnRequest().DoFunc(
	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		r.Header.Set("X-GoProxy", "being proxied")
		return r, nil
	})*/

	// changes response headers
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		r.Header.Set("X-GoProxy", "being proxied")
		//fmt.Println(ctx.Req.Host, "->", r.Header.Get("Content-Type"))
		return r
	})

	// returns different content on a criteria
	/*proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("sapo.pt"))).DoFunc(
	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		fmt.Printf("trolling request %s\n", r.RequestURI)
		return r, goproxy.NewResponse(
			r,
			goproxy.ContentTypeText,
			http.StatusForbidden,
			"Don't waste your time!")
	})*/

	/*proxy.OnRequest(goproxy.ReqCondition(
	func(r *http.Request, ctx *goproxy.ProxyCtx) bool {

	}))*/

	// serve local file
	/*kitten, err := ioutil.ReadFile("/tmp/kitten.jpg")
	if err != nil {
		panic(err)
	}
	proxy.OnResponse(goproxy_image.RespIsImage).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		log.Println("PROXYING KITTEN! %s", ctx.Req.URL.String())
		return goproxy.NewResponse(
			ctx.Req,
			"image/jpeg",
			http.StatusOK,
			string(kitten[:]))
	})*/

	// patches HTML
	rgx := regexp.MustCompile("josepedrodias\\.com")
	proxy.OnResponse(goproxy_html.IsHtml).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		//fmt.Println("HTML: %s", ctx.Req.URL.String())

		// if ctx.Req.URL.String() == "http://josepedrodias.com/" {
		if rgx.MatchString(ctx.Req.URL.String()) {
			fmt.Println("PATCHING HTML %s", ctx.Req.URL.String())

			bodyBytes, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				panic(err)
			}

			bodyString := string(bodyBytes[:])

			//fmt.Println("BODY: %s\n\n\n", bodyString)

			i := strings.LastIndex(bodyString, "</body>")

			patch := "<script>console.log(\"pwned!\");</script>"

			bodyString2 := bodyString[:i] + patch + bodyString[i:]

			return goproxy.NewResponse(
				ctx.Req,
				ctx.Req.Header.Get("Content-Type"),
				http.StatusOK,
				bodyString2)
		}

		return resp
	})

	log.Fatal(http.ListenAndServe(*addr, proxy))
}
