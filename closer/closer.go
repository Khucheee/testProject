package closer

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var CloseFunctions []func()

func InitGracefulShutdown() *sync.WaitGroup {
	slog.Debug("InitGracefulShutdown started")
	var wg sync.WaitGroup
	slog.Debug("WaitGroup for graceful shutdown created")
	wg.Add(1)
	slog.Debug("WaitGroup count: 1")
	go func() {
		slog.Debug("Call CtrlC func in InitGracefulShutdown")
		CtrlC()
		wg.Done()
	}()
	return &wg
}

func CtrlC() {
	slog.Debug("func CtrlC started")
	sigChan := make(chan os.Signal, 1)
	slog.Debug("channel for catching ctrl+c created")
	signal.Notify(sigChan, syscall.SIGINT)
	slog.Debug("channel was set for receiving SIGINT signal, waiting for signal")
	<-sigChan
	slog.Debug("signal from sigChan received, starting graceful shutdown")
	for iterator := len(CloseFunctions) - 1; iterator >= 0; iterator-- {
		CloseFunctions[iterator]()
	}
	slog.Debug("graceful shutdown completed")
}
