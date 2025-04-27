package mtapp

import (
	"context"
	"net/http/pprof"

	"github.com/buaazp/fasthttprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

const pprofUrlPrefix = "/debug/pprof"

func NewServiceHttpServer() *ServiceHttpServer {
	return &ServiceHttpServer{
		httpServer: newFastHttpServer(),
	}
}

type ServiceHttpServer struct {
	httpServer *FastHttpServer
}

func (s *ServiceHttpServer) ListenAndServe(ctx context.Context, listen string, enablePprof bool) <-chan error {
	router := fasthttprouter.New()

	router.GET("/metrics", s.HandleMetrics)
	router.GET("/ping", s.HandlePing)

	if enablePprof {
		for _, path := range []string{"/", "/allocs", "/block", "/goroutine", "/heap", "/mutex", "/threadcreate"} {
			router.GET(pprofUrlPrefix+path, fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index))
		}
		router.GET(pprofUrlPrefix+"/cmdline", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline))
		router.GET(pprofUrlPrefix+"/profile", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile))
		router.GET(pprofUrlPrefix+"/symbol", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol))
		router.GET(pprofUrlPrefix+"/trace", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace))
	}

	return s.httpServer.run(ctx, router.Handler, listen)
}

func (s *ServiceHttpServer) HandleMetrics(ctx *fasthttp.RequestCtx) {
	prometheusHandler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	prometheusHandler(ctx)
}

func (s *ServiceHttpServer) HandlePing(ctx *fasthttp.RequestCtx) {
	ctx.SuccessString("text/plain; charset=utf-8", "PONG")
}
