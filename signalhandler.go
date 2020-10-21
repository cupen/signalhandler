package signalhandler

import (
	"log"
	"os"
	"os/signal"
)

var (
	hub      = map[os.Signal][]Handler{}
	sigQueue = make(chan os.Signal)
)

// Handler for `os/signal`
type Handler func(sig os.Signal)

// Watch `os/signal`, which meaning handler will be triggered by `os/signal`.
// You can watch `os/signal` with multiple handlers, they will be ordering by first-watch-first-trigger.
func Watch(sig os.Signal, h Handler) {
	handlers, ok := hub[sig]
	if !ok {
		signal.Notify(sigQueue, sig)
		hub[sig] = []Handler{}
	}
	hub[sig] = append(handlers, h)
}

// Touch will trigger `os/signal` manually.
func Touch(sig os.Signal) {
	sigQueue <- sig
}

// Run in current goroutine. It will causes blocking.
func Run(defaultHandler ...Handler) {
	for sig := range sigQueue {
		handlers, ok := hub[sig]
		if !ok || len(handlers) <= 0 {
			if len(defaultHandler) > 0 {
				defaultHandler[0](sig)
			}
			continue
		}
		for _, h := range handlers {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic recovered. err=%#v", r)
					}
				}()
				h(sig)
			}()
		}
	}
}

// Start is similar as Run, but it's in a new goroutine.
func Start(defaultHandler ...Handler) {
	go Run(defaultHandler...)
}
