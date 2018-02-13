package interrupt

import (
	"os"
	"testing"
	"time"
)

func Test_handler_WaitForInterrupt_Waits(t *testing.T) {
	h := NewHandler().(*handler)

	sleepDuration := 250 * time.Millisecond

	start := time.Now()
	go func() {
		time.Sleep(sleepDuration)
		h.simulateSignal(os.Interrupt)
	}()

	h.WaitForInterrupt()

	duration := time.Since(start)
	if duration < sleepDuration {
		t.Errorf("did not wait for interrupt. wanted - %dµs, got - %s", sleepDuration, duration)
	}
}

func Test_handler_WaitForInterrupt_InvokesCallback(t *testing.T) {
	h := NewHandler().(*handler)

	invoked := false
	h.AddInterruptCallback(func() {
		invoked = true
	})

	go func() {
		time.Sleep(250 * time.Millisecond)
		h.simulateSignal(os.Interrupt)
	}()

	h.WaitForInterrupt()

	if invoked == false {
		t.Error("interrupt handler was not invoked")
	}
}
