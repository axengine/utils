package pprof

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/pkg/errors"
)

func MustStart(listen string) {
	go func() {
		if err := http.ListenAndServe(listen, nil); err != nil {
			panic(errors.Wrap(err, "pprof bind "+listen))
		}
	}()
}

func Start(listen string) error {
	var ch chan error
	go func() {
		if err := http.ListenAndServe(listen, nil); err != nil {
			ch <- errors.Wrap(err, "pprof bind "+listen)
			return
		}
		ch <- nil
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(time.Second):
		return nil
	}
}

// StartAsync run pprof http server async
// usage:go Start(ctx,wg,listen)
func StartAsync(ctx context.Context, wg *sync.WaitGroup, listen string) {
	server := &http.Server{
		Addr:    listen,
		Handler: http.DefaultServeMux,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			log.Println("start pprof server failed", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown pprof server failed", err)
	}
	log.Println("pprof server stopped")
}
