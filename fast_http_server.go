package mtapp

import (
	"context"

	"github.com/valyala/fasthttp"
)

type fastHttpServerOption interface {
	apply(*FastHttpServer)
}

type serveOption func(server *fasthttp.Server, addr string) error

func (o serveOption) apply(server *FastHttpServer) {
	server.serve = o
}

func withServe(s func(server *fasthttp.Server, addr string) error) fastHttpServerOption {
	return serveOption(s)
}

type FastHttpServer struct {
	serve func(server *fasthttp.Server, addr string) error
}

func newFastHttpServer(opts ...fastHttpServerOption) *FastHttpServer {
	server := &FastHttpServer{
		serve: defaultServe,
	}

	for _, o := range opts {
		o.apply(server)
	}

	return server
}

// run запускает fasthttp сервер, и возвращает канал с ошибкой при завершении работы сервера.
func (s *FastHttpServer) run(
	ctx context.Context,
	handler fasthttp.RequestHandler,
	listen string,
) <-chan error {
	ctx, cancel := context.WithCancel(ctx)
	server := &fasthttp.Server{
		Handler: handler,
	}

	result := make(chan error)
	go func() {
		defer close(result)
		<-ctx.Done()
		result <- server.Shutdown()
	}()

	go func() {
		_ = s.serve(server, listen)
		cancel()
	}()

	return result
}

func defaultServe(server *fasthttp.Server, addr string) error {
	return server.ListenAndServe(addr)
}
