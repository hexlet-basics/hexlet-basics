package main

import (
	"syscall"
	"testing"
	"time"
)

// TestMainConsumesUntilTheProcessIsAskedToStop runs main itself, as the
// platform does: the worker starts its queue and router, and SIGTERM makes it
// drain and return. A failed start would call os.Exit and take the test binary
// down.
func TestMainConsumesUntilTheProcessIsAskedToStop(t *testing.T) {
	workerEnv(t)

	returned := make(chan struct{})
	go func() {
		main()
		close(returned)
	}()

	// The worker has no endpoint to probe; give River and the event router the
	// time to start before the signal.
	time.Sleep(2 * time.Second)
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	select {
	case <-returned:
	case <-time.After(60 * time.Second):
		t.Fatal("main did not return after SIGTERM")
	}
}
