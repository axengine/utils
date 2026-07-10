package pprof

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

func TestMustStart(t *testing.T) {
	MustStart(":6060")
	select {}
}

func TestStart(t *testing.T) {
	if err := Start(":6060"); err != nil {
		t.Fatal(err)
	}
	select {}
}

func TestStartAsync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go StartAsync(ctx, &wg, ":6060")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
	wg.Wait()
}
