package pprof

import "testing"

func TestMustStart(t *testing.T) {
	MustStart(":8080")
	select {}
}

func TestStart(t *testing.T) {
	if err := Start(":8080"); err != nil {
		t.Fatal(err)
	}
	select {}
}
