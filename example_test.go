package grace

import (
	"fmt"
	"os"
	"syscall"
	"time"
)


func ExampleWait_Wait() {
	wait, _ := NewWait()
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
	wait.Wait()
	fmt.Println("Bye")
	// Output: Bye
}

func ExampleWait_WaitWithFunc() {
	wait, _ := NewWait()
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
	wait.WaitWithFunc(func() error {
		// your logic here
		fmt.Println("Bye")
		return nil
	})
	// Output: Bye
}

func ExampleWait_WaitWithFuncErr() {
	wait, _ := NewWait()
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
	if err := wait.WaitWithFunc(func() error {
		// your logic here
		return fmt.Errorf("something went wrong")
	}); err != nil {
		fmt.Println("something bad happened:", err)
	}
	fmt.Println("Bye")
	// Output: something bad happened: something went wrong
	// Bye
}

func ExampleWait_WaitWithTimeout() {
	wait, _ := NewWait()
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
	wait.WaitWithTimeout(1 * time.Second)
	fmt.Println("Bye")
	// Output: Bye
}

func ExampleWait_WaitWithTimeoutAndFunc() {
	wait, _ := NewWait()
	wait.WaitWithTimeoutAndFunc(1*time.Second, func() error {
		// your logic here
		fmt.Println("Bye")
		return nil
	})
	// Output: Bye
}

func ExampleWait_WaitWithTimeoutAndFuncErr() {
	wait, _ := NewWait()
	if err := wait.WaitWithTimeoutAndFunc(1*time.Second, func() error {
		// your logic here
		return fmt.Errorf("something went wrong")
	}); err != nil {
		fmt.Println("something bad happened:", err)
	}
	fmt.Println("Bye")
	// Output: something bad happened: something went wrong
	// Bye
}

