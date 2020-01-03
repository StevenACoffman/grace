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
```go
package main

import (
    "log"
    "os"
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
				return ctx.Err()
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

##### Usage with a default timeout:

```go
package main

import (
    "log"
    "os"
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
				return ctx.Err()
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

