package pprof

import (
	"github.com/pkg/errors"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func MustStart(addr string) {
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(errors.Wrap(err, "pprof bind "+addr))
		}
	}()
}

func Start(addr string) error {
	var ch chan error
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			ch <- errors.Wrap(err, "pprof bind "+addr)
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
