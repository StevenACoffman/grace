# Grace - graceful shutdown made simple

Tiny library to gracefully shutdown your application by catching the OS signals using [`sync.errgroup`](https://godoc.org/golang.org/x/sync/errgroup).

I often find I have invoked one or more persistent blocking methods, and some other method is needed be invoked in another goroutine to tell it to gracefully shut down when an interrupt is received.

For instance, when [`ListenAndServe()`](https://golang.org/pkg/net/http/#ListenAndServe) is invoked, [`Shutdown`](https://godoc.org/net/http#Server.Shutdown) needs to be called.

This library allows you to start zero or more concurrent goroutines, and trigger a graceful shutdown when an interrupt is received.

+ Go `net/http` package offers [`Shutdown`](https://godoc.org/net/http#Server.Shutdown) function to gracefully shutdown your http server.
+ Go `database/sql` package offers [`Close`](https://godoc.org/database/sql#DB.Close) function to gracefully close the connection to your SQL database. 
+ Google `google.golang.org/grpc` package offers [`Server.GracefulStop`](https://godoc.org/google.golang.org/grpc#Server.GracefulStop), stops accepting new connections, and blocks until all the pending RPCs are finished

Alternatively, this library allows you to invoke zero or more concurrent goroutines with an optional timeout.

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

### Prior Art and Alternatives
This uses errgroup, but I found a number of other libraries that use other mechanisms:
+ [death](https://github.com/vrecan/death) (sync.WaitGroup)
+ [graceful](https://github.com/TV4/graceful) (context cancellation)
+ [finish](https://github.com/pseidemann/finish/) (context + mutexes)
+ [waitabit](https://github.com/heartwilltell/waitabit) (sync.WaitGroup)
+ [Gist using errgroup](https://gist.github.com/pteich/c0bb58b0b7c8af7cc6a689dd0d3d26ef)
+ [Gist using GRPC GracefulStop](https://gist.github.com/akhenakh/38dbfea70dc36964e23acc19777f3869)

Comparing them is pretty instructive. I wish I'd used some of their testing techniques!