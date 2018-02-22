package blockchain

import (
	"testing"
)

func TestFeed_Send(t *testing.T) {
	tests := []struct {
		name  string
		event Event
	}{
		{
			name:  "send block saved event",
			event: BlockSaved,
		},
		{
			name:  "send header saved event",
			event: HeaderSaved,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				channel = make(chan Event, 1)
				f       = NewFeed()
				sub1    = f.Subscribe(channel)
			)

			f.Send(tt.event)

			if event, ok := <-channel; !ok || event != tt.event {
				t.Errorf("Unexpected value received. wanted - %d, got - %d", tt.event, event)
			}

			sub1.Unsubscribe()
		})
	}
}

func TestSubscriber_Unsubscribe(t *testing.T) {
	var (
		channel = make(chan Event, 1)
		f       = NewFeed()
		sub1    = f.Subscribe(channel)
	)

	sub1.Unsubscribe()
	f.Send(BlockSaved)

	if _, ok := <-channel; ok {
		t.Error("Expected channel to be closed")
	}
}
