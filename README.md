# Grace - graceful shutdown made simple

Tiny library for manage you application shutdown in graceful way by catching the OS signals using errgroup instead of WaitGroup.

For a waitgroup implementation see the excellent [waitabit](https://github.com/heartwilltell/waitabit/) which is what I based this on. 

## Documentation

[![](https://goreportcard.com/badge/github.com/StevenACoffman/grace)](https://goreportcard.com/report/github.com/StevenACoffman/grace)
[![](https://godoc.org/github.com/StevenACoffman/grace?status.svg)](https://godoc.org/github.com/StevenACoffman/grace)

## Installation

```bash
go get -u github.com/StevenACoffman/grace
```

## Usage
#### Simple Run until Interrupt signal received 
```go
package main

import (
	"log"
	"time"

	"github.com/StevenACoffman/grace"
)

func main() {

	wait, ctx := grace.NewWait()

	err := wait.WaitWithFunc(func() error {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				log.Printf("ticker 2s ticked\n")
				// testcase what happens if an error occured
				//return fmt.Errorf("test error ticker 2s")
			case <-ctx.Done():
				log.Printf("closing ticker 2s goroutine\n")
				return nil
			}
		}
	})

	if err != nil {
		log.Println("finished clean")
	} else {
		log.Printf("received error: %v", err)
	}
}

```

#### Usage with a default timeout:

```go
package main

import (
	"log"
	"time"

	"github.com/StevenACoffman/grace"
)

func main() {

	wait, ctx := grace.NewWait()

	err := wait.WaitWithTimeoutAndFunc(15*time.Second, func() error {
		ticker := time.NewTicker(2 * time.Second)
		for {
         			select {
         			case <-ticker.C:
         				log.Printf("ticker 2s ticked\n")
         				// testcase what happens if an error occured
         				//return fmt.Errorf("test error ticker 2s")
         			case <-ctx.Done():
         				log.Printf("closing ticker 2s goroutine\n")
         				return nil
         			}
         		}
	})

	if err != nil {
		log.Println("finished clean")
	} else {
		log.Printf("received error: %v", err)
	}
}

```

#### Usage with cleanup on shutdown
Bring your own cleanup function!
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/StevenACoffman/grace"
)

func main() {

	wait, ctx := grace.NewWait()
	var httpServer *http.Server

	err := wait.WaitWithFunc(
		func() error {
			http.HandleFunc("/", healthCheck)
			httpServer = newHTTPServer()

			if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
				return err
			}
			return nil
		},
		func() error { 
			//cleanup: on interrupt, shutdown server
			<-ctx.Done()
			log.Printf("closing http goroutine\n")
			return httpServer.Shutdown(ctx)
		})

	if err != nil {
		log.Println("finished clean")
	} else {
		log.Printf("received error: %v", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(200)
}

func newHTTPServer() *http.Server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":8080"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("HTTP Metrics server serving at %s", ":8080")
	return httpServer
}

```

