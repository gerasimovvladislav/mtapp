package mtapp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func InitSignalHandler(
	cancelFunc context.CancelFunc,
) {
	osSigCh := make(chan os.Signal, 1)

	signal.Notify(
		osSigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		s := <-osSigCh
		switch s {
		case syscall.SIGHUP:
		case syscall.SIGINT:
			cancelFunc()
		case syscall.SIGTERM:
			cancelFunc()
		case syscall.SIGQUIT:
			cancelFunc()
		}
	}()
}

func FatalJsonLog(msg string, err error) string {
	escape := func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `"`, `\"`)
	}

	errString := ""
	if err != nil {
		errString = err.Error()
	}

	return fmt.Sprintf(
		`{"level":"fatal","ts":"%s","msg":"%s","error":"%s"}`,
		time.Now().Format(time.RFC3339),
		escape(msg),
		escape(errString),
	)
}
