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

type Handler func()
type DefaultHandler func(sig os.Signal)

func Watch(sig os.Signal, h Handler) {
	handlers, ok := hub[sig]
	if !ok {
		signal.Notify(sigQueue, sig)
		hub[sig] = []Handler{}
	}
	hub[sig] = append(handlers, h)
}

func Touch(sig os.Signal) {
	sigQueue <- sig
}

func Run(defaultHandler ...DefaultHandler) {
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
				h()
			}()
		}
	}
}

func Start(defaultHandler ...DefaultHandler) {
	go Run(defaultHandler...)
}
