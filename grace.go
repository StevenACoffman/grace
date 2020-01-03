// Package grace provides functionality for handling system interrupts.
package grace

import (
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func getStopSignalsChannel() chan os.Signal {

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel,
		os.Interrupt,    // interrupt is syscall.SIGINT, Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)
	return signalChannel

}

// Wait struct represent wait functionality.
type Wait struct {
	g       *errgroup.Group
	done    context.CancelFunc
	gctx    *context.Context
	signals *chan os.Signal
}

// NewWait
func NewWait() *Wait {
	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)
	signalChannel := getStopSignalsChannel()
	wait := &Wait{
		g:       g,
		done:    done,
		signals: &signalChannel,
		gctx:    &gctx,
	}
	// goroutine to check for signals to gracefully finish all functions
	wait.g.Go(wait.listen)
	return wait
}

func (w *Wait) Done() <-chan struct{} {
	gctx := *w.gctx
	return gctx.Done()
}

func (w *Wait) Err() error {
	gctx := *w.gctx
	return gctx.Err()
}

// Wait - simply wait for interruption.
func (w *Wait) Wait() error {
	return w.WaitWithFunc()
}

// WaitWithTimeout wait a given amount for time and then exit.
func (w *Wait) WaitWithTimeout(timeout time.Duration) error {
	return w.WaitWithTimeoutAndFunc(timeout)
}

// WaitWithFunc receive function(s) as argument and call Wait() to wait interrupt signal.
// The function(s) will be called before Wait().
// Returns the first non-nil error (if any) from the functions.
func (w *Wait) WaitWithFunc(functions ...func() error) error {
	for i := range functions {
		w.g.Go(functions[i])
	}
	// wait for all errgroup goroutines
	if err := w.g.Wait(); err != nil {
		return err
	}
	return nil
}

// WaitWithTimeoutAndFunc call function, wait a given amount of time and then exit.
// Returns the first non-nil error (if any) from the functions.
func (w *Wait) WaitWithTimeoutAndFunc(timeout time.Duration, functions ...func() error) error {
	for i := range functions {
		w.g.Go(functions[i])
	}
	// force a stop after timeout
	time.AfterFunc(timeout, func() {
		stdlog.Printf("force finished after %fs", timeout.Seconds())
		w.done()
	})

	// wait for all errgroup goroutines
	if err := w.g.Wait(); err != nil {
		return err
	}
	return nil
}

// listen - a simple function that listens to the signals channel for interruption signals and then call Done() of the errgroup.
func (w *Wait) listen() error {
	signalChannel := getStopSignalsChannel()
	select {
	case sig := <-signalChannel:
		stdlog.Printf("Received signal: %s\n", sig)
		w.done()
	case <-w.Done():
		stdlog.Printf("closing signal goroutine\n")
		return w.Err()
	}

	return nil
}
