# Introduction
A easy way to hanlde `os/signals`. See https://gobyexample.com/signals

# Usage

```go
import (
	"log"
	"os"
	"syscall"
	"time"

	"github.com/cupen/signalhandler"
)

func main() {
	exitGracefully := func(sig os.Signal) {
		// do something before process exit
		countDown(10)
		log.Printf("exited by os/signal(%v)", sig)
		os.Exit(0)
	}

	signalhandler.Watch(syscall.SIGQUIT, exitGracefully)
	signalhandler.Watch(syscall.SIGTERM, exitGracefully)
	signalhandler.Watch(syscall.SIGINT, func(sig os.Signal) {
		exitGracefully(sig)
	})

	log.Printf("started. CTRL-C or kill -INT ${pid}")
	signalhandler.Run()
}

func countDown(secs int) {
	for i := secs; i > 0; i-- {
		log.Printf("exit after %d seconds.", i)
		time.Sleep(time.Second)
	}
}}

```