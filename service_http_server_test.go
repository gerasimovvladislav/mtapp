package mtapp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"go.uber.org/goleak"
)

func TestServer_ListenAndServe(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()

	s := &ServiceHttpServer{
		httpServer: newFastHttpServer(withServe(func(server *fasthttp.Server, _ string) error {
			return server.Serve(ln)
		})),
	}

	const listen = ":64102"

	ctx, cancel := context.WithCancel(context.Background())
	serverFinishErrCh := s.ListenAndServe(ctx, listen, false)

	var err error
	cancel()
	err = <-serverFinishErrCh
	assert.NoError(t, err)
	goleak.VerifyNone(t, goleak.IgnoreTopFunction("time.Sleep"))
}

func TestServer_Ping(t *testing.T) {
	s := &ServiceHttpServer{
		httpServer: newFastHttpServer(),
	}

	ctx := &fasthttp.RequestCtx{}
	s.HandlePing(ctx)

	assert.Equal(t, []byte("text/plain; charset=utf-8"), ctx.Response.Header.ContentType())
	assert.Equal(t, []byte("PONG"), ctx.Response.Body())
}
func TestServer_ListenAndServe_ServiceRoutes(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()

	s := &ServiceHttpServer{
		httpServer: newFastHttpServer(withServe(func(server *fasthttp.Server, _ string) error {
			return server.Serve(ln)
		})),
	}

	const listen = ":64102"

	hostname, err := os.Hostname()
	assert.NoError(t, err)

	host := fmt.Sprintf("%s%s", hostname, listen)

	ctx, cancel := context.WithCancel(context.Background())
	serverFinishErrCh := s.ListenAndServe(ctx, listen, false)

	routes := []string{
		"/metrics",
		"/ping",
	}

	for _, r := range routes {
		r := r
		t.Run(fmt.Sprintf("route %s", r), func(t *testing.T) {
			client := http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
						return ln.Dial()
					},
				},
			}
			resp, err := client.Get(fmt.Sprintf("http://%s%s", host, r))

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.NoError(t, resp.Body.Close())
		})
	}
	cancel()
	err = <-serverFinishErrCh
	assert.NoError(t, err)
	goleak.VerifyNone(t, goleak.IgnoreTopFunction("time.Sleep"))
}
