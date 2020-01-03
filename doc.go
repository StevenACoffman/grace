// Package Grace provides graceful shutdown made simple for zero to many goroutines
// Tiny library to gracefully shutdown your application by catching the OS signals using [`sync.errgroup`](https://godoc.org/golang.org/x/sync/errgroup).
//
// I often find I have invoked one or more persistent blocking methods, and some other method is needed be invoked in another goroutine to tell it to gracefully shut down when an interrupt is received.
//
// For instance, when [`ListenAndServe()`](https://golang.org/pkg/net/http/#ListenAndServe) is invoked, [`Shutdown`](https://godoc.org/net/http#Server.Shutdown) needs to be called.
//
// This library allows you to start zero or more concurrent goroutines, and trigger a graceful shutdown when an interrupt is received.
//
// • Go `net/http` package offers [`Shutdown`](https://godoc.org/net/http#Server.Shutdown) function to gracefully shutdown your http server.
// • Go `database/sql` package offers [`Close`](https://godoc.org/database/sql#DB.Close) function to gracefully close the connection to your SQL database.
// • Google `google.golang.org/grpc` package offers [`Server.GracefulStop`](https://godoc.org/google.golang.org/grpc#Server.GracefulStop), stops accepting new connections, and blocks until all the pending RPCs are finished
//
// Alternatively, this library allows you to invoke zero or more concurrent goroutines with an optional timeout.
//
// Example Usage with ListenAndServe and Shutdown
//  package main
//
//  import (
//  	"fmt"
//  	"log"
//  	"net/http"
//  	"time"
//
//  	"github.com/StevenACoffman/grace"
//  )
//
//  func main() {
//
//  	wait, ctx := grace.NewWait()
//  	var httpServer *http.Server
//
//  	err := wait.WaitWithFunc(
//  		func() error {
//  			http.HandleFunc("/", healthCheck)
//  			httpServer = newHTTPServer()
//
//  			if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
//  				return err
//  			}
//  			return nil
//  		},
//  		func() error {
//  			//cleanup: on interrupt, shutdown server
//  			<-ctx.Done()
//  			log.Printf("closing http goroutine\n")
//  			return httpServer.Shutdown(ctx)
//  		})
//
//  	if err != nil {
//  		log.Println("finished clean")
//  	} else {
//  		log.Printf("received error: %v", err)
//  	}
//  }
//
//  func healthCheck(w http.ResponseWriter, r *http.Request) {
//  	w.Header().Set("Content-Type", "text/plain")
//  	w.Header().Set("Content-Length", "0")
//  	w.WriteHeader(200)
//  }
//
//  func newHTTPServer() *http.Server {
//  	httpServer := &http.Server{
//  		Addr:         fmt.Sprintf(":8080"),
//  		ReadTimeout:  10 * time.Second,
//  		WriteTimeout: 10 * time.Second,
//  	}
//  	log.Printf("HTTP Metrics server serving at %s", ":8080")
//  	return httpServer
//  }
package grace


