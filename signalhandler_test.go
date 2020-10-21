package signalhandler

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func assertString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.FailNow()
	}
}

func TestWatch(t *testing.T) {
	var countor = make(chan string)
	Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit"
	})
	go Run(func(sig os.Signal) {
		t.Logf("default handler. signal = %v", sig)
		countor <- "miss"
	})
	Touch(syscall.SIGTERM)
	assertString(t, "miss", <-countor)

	Touch(syscall.SIGHUP)
	assertString(t, "hit", <-countor)
}

func TestWatch_MultiHandlers(t *testing.T) {
	var countor = make(chan string)
	Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit"
	})
	Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit2"
	})
	Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit3"
	})
	go Run(func(sig os.Signal) {
		t.Logf("default handler. signal = %v", sig)
		countor <- "miss"
	})
	Touch(syscall.SIGTERM)
	assertString(t, "miss", <-countor)

	Touch(syscall.SIGHUP)
	assertString(t, "hit", <-countor)
	assertString(t, "hit2", <-countor)
	assertString(t, "hit3", <-countor)
}

func TestWatch_Panic(t *testing.T) {
	var countor = make(chan string)
	Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit"
		panic(fmt.Errorf("panic"))
	})
	Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit2"
		panic(fmt.Errorf("panic2"))
	})
	Watch(syscall.SIGHUP, func(sig os.Signal) {
		countor <- "hit3"
		panic(fmt.Errorf("panic3"))
	})
	go Run(func(sig os.Signal) {
		t.Logf("default handler. signal = %v", sig)
		countor <- "miss"
	})
	Touch(syscall.SIGHUP)
	assertString(t, "hit", <-countor)
	assertString(t, "hit2", <-countor)
	assertString(t, "hit3", <-countor)
	time.Sleep(time.Millisecond)
}
