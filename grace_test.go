package grace

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWait_Wait(t *testing.T) {

	t.Run("Wait with timeout and func", func(t *testing.T) {
		start := time.Now().Unix()
		var result int
		wait, _ := NewWait()
		syscall.Kill(os.Getegid(), syscall.SIGTERM)
		wait.WaitWithTimeoutAndFunc(time.Second, func() error {
			result = 1
			return nil
		})
		end := time.Now().Unix()
		if end-start != 1 {
			t.Error("Something went wrong")
		}
		if result != 1 {
			t.Error("Result is not equal 1")
		}
	})

	t.Run("Wait", func(t *testing.T) {
		wait, _ := NewWait()
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
		wait.Wait()
	})

	t.Run("Wait with timeout", func(t *testing.T) {
		start := time.Now().Unix()
		wait, _ := NewWait()
		wait.WaitWithTimeout(time.Second)
		end := time.Now().Unix()
		if end-start != 1 {
			t.Error("Something went wrong")
		}
	})

	t.Run("Wait with func", func(t *testing.T) {
		var result int
		wait, _ := NewWait()
		syscall.Kill(os.Getpid(), syscall.SIGTERM)

		wait.WaitWithFunc(func() error {
			result = 1
			return nil
		})

		if result != 1 {
			t.Error("Result is not equal 1")
		}
	})
}
