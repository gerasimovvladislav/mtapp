package mtapp

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/buaazp/fasthttprouter"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"go.uber.org/goleak"
	"golang.org/x/net/context"
)

func TestRunFastHttpServer(t *testing.T) {
	router := fasthttprouter.New()
	router.GET("/ping", func(ctx *fasthttp.RequestCtx) {
		ctx.SuccessString("text/plain; charset=utf-8", "PONG")
	})

	ctx, cancel := context.WithCancel(context.Background())
	const listen = ":64102"

	ln := fasthttputil.NewInmemoryListener()

	server := FastHttpServer{
		serve: func(server *fasthttp.Server, _ string) error {
			return server.Serve(ln)
		},
	}
	errCh := server.run(ctx, router.Handler, listen)

	hostname, err := os.Hostname()
	assert.NoError(t, err)

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}
	resp, err := client.Get(fmt.Sprintf("http://%s%s/ping", hostname, listen))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, resp.Body.Close())

	cancel()
	err = <-errCh
	assert.NoError(t, err)

	// в workerPool.Start у fastHttp остается подвисшая горутина, игнорируем ее.
	goleak.VerifyNone(t, goleak.IgnoreTopFunction("time.Sleep"))
}
