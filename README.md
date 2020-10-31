[![GoDoc][doc]][doc-to] 
[![License][license]][license-to] 
![Go Report][status]
[![Coverage][cover]][cover-to]

[status]: https://goreportcard.com/badge/github.com/cupen/signalhub?style=flat-square
[doc]:    https://godoc.org/github.com/cupen/signalhub?status.svg
[doc-to]: https://pkg.go.dev/github.com/cupen/signalhub
[license]:  https://img.shields.io/badge/license-WTFPL-blue.svg
[license-to]: LICENSE
[cover]:    https://codecov.io/gh/cupen/signalhub/branch/master/graph/badge.svg?token=HQODXQHLK3
[cover-to]: https://codecov.io/gh/cupen/signalhub

# Introduction
Hanlde `os/signals` with function style. See https://gobyexample.com/signals

# Usage

```go
import (
	"log"
	"os"
	"syscall"
	"time"

	"github.com/cupen/signalhub"
)

func main() {
	exitHook := func(sig os.Signal) {
		// do something before process exit
		countDown(10)
		log.Printf("exited by os/signal(%v)", sig)
		os.Exit(0)
	}
	hub := signalhub.New()
	hub.Watch(syscall.SIGQUIT, exitHook)
	hub.Watch(syscall.SIGTERM, exitHook)
	hub.Watch(syscall.SIGINT, func(sig os.Signal) {
		exitHook(sig)
	})

	log.Printf("started. CTRL-C or kill -INT ${pid}")
	h.Run()
}

func countDown(secs int) {
	for i := secs; i > 0; i-- {
		log.Printf("exit after %d seconds.", i)
		time.Sleep(time.Second)
	}
}

```

# License
```
        DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE 
                    Version 2, December 2004 

 Copyright (C) 2020 cupen <xcupen@gmail.com> 

 Everyone is permitted to copy and distribute verbatim or modified 
 copies of this license document, and changing it is allowed as long 
 as the name is changed. 

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE 
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION 

  0. You just DO WHAT THE FUCK YOU WANT TO.
```