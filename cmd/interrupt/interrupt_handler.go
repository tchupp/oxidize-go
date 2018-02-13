package interrupt

import (
	"os"
	"os/signal"
)

type Handler interface {
	AddInterruptCallback(handler func())
	WaitForInterrupt()
}

type handler struct {
	signalChannel      chan os.Signal
	interruptReceived  chan struct{}
	addCallbackChannel chan func()
	interruptCallbacks []func()
}

func NewHandler() Handler {
	h := &handler{
		signalChannel:      make(chan os.Signal, 1),
		interruptReceived:  make(chan struct{}),
		addCallbackChannel: make(chan func(), 5),
	}
	go h.masterInterruptHandler()

	return h
}

func (h *handler) masterInterruptHandler() {
	signal.Notify(h.signalChannel, os.Interrupt)

	for {
		select {
		case sig := <-h.signalChannel:
			log.Infof("Received signal (%s).  Shutting down...", sig)
			h.invokeCallbacks()
			return

		case handler := <-h.addCallbackChannel:
			h.interruptCallbacks = append([]func(){handler}, h.interruptCallbacks...)
		}
	}
}

func (h *handler) WaitForInterrupt() {
	<-h.interruptReceived
}

func (h *handler) simulateSignal(sig os.Signal) {
	h.signalChannel <- sig
}

func (h *handler) AddInterruptCallback(handler func()) {
	h.addCallbackChannel <- handler
}

func (h *handler) invokeCallbacks() {
	for _, callback := range h.interruptCallbacks {
		callback()
	}
	close(h.interruptReceived)
}
