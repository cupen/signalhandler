package signalhub

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"

	"github.com/cupen/signalhub/signaldefs"
)

// Handler for `os/signal`
type Handler func(os.Signal)

// Hub is a group of signal handlers.
type SignalHub struct {
	handlers map[os.Signal][]Handler
	queue    chan os.Signal
}

// New a handlers hub
func New() *SignalHub {
	return &SignalHub{
		handlers: map[os.Signal][]Handler{},
		queue:    make(chan os.Signal),
	}
}

// Watch `os/signal`, which meaning handler will be triggered by `os/signal`.
// You can watch `os/signal` with multiple handlers, they will triggerred with
// ordering by first-watch-first-trigger.
func (h *SignalHub) Watch(sig os.Signal, handler Handler) {
	if sig == signaldefs.SIG_NONE {
		log.Printf("[warn] invalid signal: %s, skipped it", sig)
		return
	}

	handlers, ok := h.handlers[sig]
	if !ok {
		signal.Notify(h.queue, sig)
		h.handlers[sig] = []Handler{}
	}
	h.handlers[sig] = append(handlers, handler)
}

func (h *SignalHub) Unwatch(sig os.Signal) {
	if len(h.handlers) <= 0 {
		return
	}
	delete(h.handlers, sig)
}

// Touch will trigger `os/signal` manually.
func (h *SignalHub) Touch(sig os.Signal) {
	if sig == signaldefs.SIG_NONE {
		log.Printf("[warn] invalid signal received")
		return
	}
	h.queue <- sig
}

// Run in current goroutine. It will causes blocking.
func (h *SignalHub) Run(defaultHandler ...Handler) {
	for sig := range h.queue {
		handlers, ok := h.handlers[sig]
		if !ok || len(handlers) <= 0 {
			if len(defaultHandler) > 0 {
				handler := defaultHandler[0]
				h._touch(sig, handler)
			}
			continue
		}
		for _, handler := range handlers {
			h._touch(sig, handler)
		}
	}
}

// Start is similar as Run, but it's running asynchronously in separate goroutine.
func (h *SignalHub) Start(defaultHandler ...Handler) {
	go h.Run(defaultHandler...)
}

// Stop all
func (h *SignalHub) Stop() {
	close(h.queue)
}

func (h *SignalHub) _touch(sig os.Signal, handler Handler) {
	defer func() {
		if r := recover(); r != nil {
			fname := getFuncName(handler)
			log.Printf("panic:%v casued by handler '%s'.", r, fname)
		}
	}()
	handler(sig)
}

func getFuncName(f interface{}) string {
	if f == nil {
		return "null"
	}
	val := reflect.ValueOf(f)
	name := runtime.FuncForPC(val.Pointer()).Name()
	tmpArr := strings.Split(name, ".")
	return tmpArr[len(tmpArr)-1]
}
