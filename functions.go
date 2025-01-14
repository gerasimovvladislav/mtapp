package mtapp

import (
	"context"
	"fmt"
	"log/slog"
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
			slog.Info("Received signal SIGHUP! Renew configs")
		case syscall.SIGINT:
			slog.Info("Received signal SIGINT! Process exited")
			cancelFunc()
		case syscall.SIGTERM:
			slog.Info("Received signal SIGTERM! Process exited")
			cancelFunc()
		case syscall.SIGQUIT:
			slog.Info("Received signal SIGQUIT! Process exited")
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
