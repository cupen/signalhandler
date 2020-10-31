package signalhub

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func assertString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Logf("\"%s\" != \"%s\"", expected, actual)
		t.FailNow()
	}
}

func TestWatch(t *testing.T) {
	var countor = make(chan string)
	hub := New()
	hub.Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit"
	})
	go hub.Run(func(sig os.Signal) {
		t.Logf("default handler. signal = %v", sig)
		countor <- "miss"
	})
	hub.Touch(syscall.SIGTERM)
	assertString(t, "miss", <-countor)

	hub.Touch(syscall.SIGHUP)
	assertString(t, "hit", <-countor)
}

func TestWatch_MultiHandlers(t *testing.T) {
	var countor = make(chan string)
	hub := New()
	hub.Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit"
	})
	hub.Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit2"
	})
	hub.Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit3"
	})
	go hub.Run(func(sig os.Signal) {
		t.Logf("default handler. signal = %v", sig)
		countor <- "miss"
	})
	hub.Touch(syscall.SIGTERM)
	assertString(t, "miss", <-countor)

	hub.Touch(syscall.SIGHUP)
	assertString(t, "hit", <-countor)
	assertString(t, "hit2", <-countor)
	assertString(t, "hit3", <-countor)
}

func TestWatch_Panic(t *testing.T) {
	var countor = make(chan string)
	hub := New()
	hub.Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit"
		panic(fmt.Errorf("panic"))
	})
	hub.Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit2"
		panic(fmt.Errorf("panic2"))
	})
	hub.Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit3"
		panic(fmt.Errorf("panic3"))
	})
	go hub.Run(func(sig os.Signal) {
		t.Logf("default handler. signal = %v", sig)
		countor <- "miss"
	})
	hub.Touch(syscall.SIGHUP)
	assertString(t, "hit", <-countor)
	assertString(t, "hit2", <-countor)
	assertString(t, "hit3", <-countor)
	time.Sleep(time.Millisecond)
}

func yes() {}
func Test_getFuncName(t *testing.T) {
	// yes := func() {}
	assertString(t, "yes", getFuncName(yes))
}
